package api

import (
	"errors"
	"fmt"
	"net/http"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	db "github.com/katatrina/gundam-BE/internal/db/sqlc"
	"github.com/katatrina/gundam-BE/internal/token"
	"github.com/katatrina/gundam-BE/internal/util"
	"github.com/katatrina/gundam-BE/internal/worker"
	"github.com/rs/zerolog/log"
)

type createExchangeOfferRequest struct {
	// UUID của bài đăng trao đổi mà bạn muốn tạo offer
	ExchangePostID string `json:"exchange_post_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	
	// Danh sách ID các Gundam của chủ bài đăng mà bạn muốn nhận về (phải thuộc bài đăng này và có status "for exchange")
	PosterGundamIDs []int64 `json:"poster_gundam_ids" binding:"required" example:"123,456"`
	
	// Danh sách ID các Gundam của bạn để đưa ra trao đổi (phải thuộc về bạn và có status "in store")
	OffererGundamIDs []int64 `json:"offerer_gundam_ids" binding:"required" example:"789,321"`
	
	// ID người phải trả tiền bù (chỉ có thể là bạn hoặc chủ bài đăng, để null nếu không có ai bù)
	PayerID *string `json:"payer_id" example:"user_abc123"`
	
	// Số tiền bù theo VND (bắt buộc nếu có payer_id, phải > 0, chỉ trừ tiền khi offer được chấp nhận)
	CompensationAmount *int64 `json:"compensation_amount" example:"50000"`
	
	// Tin nhắn gửi kèm cho chủ bài đăng (tùy chọn)
	Note *string `json:"note" example:"Tôi rất thích RG Unicorn của bạn!"`
}

//	@Summary		Create an exchange offer
//	@Description	Create a new exchange offer for trading multiple Gundams between users with optional compensation.
//	@Description
//	@Description	**Business Rules:**
//	@Description	- Không thể tạo offer cho bài đăng của chính mình
//	@Description	- Mỗi user chỉ có 1 offer cho mỗi bài đăng
//	@Description	- Gundam của chủ bài đăng phải có status "for exchange"
//	@Description	- Gundam của người đề xuất phải có status "in store" (sẽ được chuyển thành "for exchange" sau khi tạo offer)
//	@Description	- Nếu có compensation, người trả phải có đủ số dư (chỉ kiểm tra nếu người đề xuất là người trả)
//	@Description	- Compensation chỉ được trừ tiền khi offer được chấp nhận, không trừ ngay
//	@Tags			exchanges
//	@Accept			json
//	@Produce		json
//	@Security		accessToken
//	@Param			request	body		createExchangeOfferRequest		true	"Create exchange offer request"
//	@Success		201		{object}	db.CreateExchangeOfferTxResult	"Exchange offer created successfully"
//	@Router			/users/me/exchange-offers [post]
func (server *Server) createExchangeOffer(c *gin.Context) {
	// Lấy thông tin người dùng đã đăng nhập
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	offererID := authPayload.Subject
	
	var req createExchangeOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	// Validate input arrays are not empty
	if len(req.PosterGundamIDs) == 0 {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("poster_gundam_ids cannot be empty")))
		return
	}
	
	if len(req.OffererGundamIDs) == 0 {
		c.JSON(http.StatusBadRequest, errorResponse(errors.New("offerer_gundam_ids cannot be empty")))
		return
	}
	
	postID, err := uuid.Parse(req.ExchangePostID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid post ID: %s", req.ExchangePostID)))
		return
	}
	
	// Kiểm tra bài đăng có tồn tại và đang mở không
	post, err := server.dbStore.GetExchangePost(c.Request.Context(), postID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			err = fmt.Errorf("exchange post ID %s not found", req.ExchangePostID)
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	if post.Status != db.ExchangePostStatusOpen {
		err = fmt.Errorf("exchange post ID %s is not open for offers", req.ExchangePostID)
		c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}
	
	// Kiểm tra người dùng không tự đề xuất cho bài đăng của mình
	if post.UserID == offererID {
		err = fmt.Errorf("you cannot make an offer to your own exchange post")
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	
	// Kiểm tra số tiền bù và người bù tiền
	if req.PayerID != nil && req.CompensationAmount == nil {
		c.JSON(http.StatusUnprocessableEntity, errorResponse(errors.New("compensation amount is required when payer is specified")))
		return
	}
	
	if req.PayerID == nil && req.CompensationAmount != nil {
		c.JSON(http.StatusUnprocessableEntity, errorResponse(errors.New("payer is required when compensation amount is specified")))
		return
	}
	
	if req.CompensationAmount != nil && *req.CompensationAmount <= 0 {
		c.JSON(http.StatusForbidden, errorResponse(errors.New("compensation amount must be positive")))
		return
	}
	
	// Kiểm tra người bù tiền phải là người đề xuất hoặc người đăng bài
	if req.PayerID != nil && *req.PayerID != offererID && *req.PayerID != post.UserID {
		c.JSON(http.StatusForbidden, errorResponse(errors.New("payer must be either the poster or the offerer")))
		return
	}
	
	// Chỉ kiểm tra số dư nếu người bù là người đề xuất.
	// Không trừ tiền ngay. Tiền sẽ được trừ khi đề xuất được chấp nhận.
	if req.PayerID != nil && *req.PayerID == offererID && req.CompensationAmount != nil {
		wallet, err := server.dbStore.GetWalletByUserID(c.Request.Context(), offererID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				err = fmt.Errorf("wallet not found for user ID %s", offererID)
				c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		if wallet.Balance < *req.CompensationAmount {
			err = fmt.Errorf("insufficient balance for compensation: needed %d, available %d", *req.CompensationAmount, wallet.Balance)
			c.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
	}
	
	// Kiểm tra các Gundam từ bài đăng
	for _, gundamID := range req.PosterGundamIDs {
		_, err = server.dbStore.GetExchangePostItemByGundamID(c.Request.Context(), db.GetExchangePostItemByGundamIDParams{
			PostID:   postID,
			GundamID: gundamID,
		})
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				err = fmt.Errorf("poster gundam ID %d is not part of the exchange post %s", gundamID, postID)
				c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		// Kiểm tra Gundam có tồn tại và có status phù hợp không
		posterGundam, err := server.dbStore.GetGundamByID(c.Request.Context(), gundamID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				err = fmt.Errorf("poster gundam ID %d not found", gundamID)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		// Kiểm tra Gundam có thuộc về người đăng bài không
		if posterGundam.OwnerID != post.UserID {
			err = fmt.Errorf("gundam ID %d does not belong to poster ID %s", gundamID, post.UserID)
			c.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		
		// Kiểm tra Gundam có được phép trao đổi không
		if posterGundam.Status != db.GundamStatusForexchange {
			err = fmt.Errorf("poster gundam ID %d is not available for exchange, current status: %s", gundamID, posterGundam.Status)
			c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
			return
		}
	}
	
	// Kiểm tra các Gundam của người đề xuất
	for _, gundamID := range req.OffererGundamIDs {
		// Kiểm tra Gundam có tồn tại không
		offererGundam, err := server.dbStore.GetGundamByID(c.Request.Context(), gundamID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				err = fmt.Errorf("gundam ID %d not found", gundamID)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		// Kiểm tra Gundam có thuộc về người đề xuất không
		if offererGundam.OwnerID != offererID {
			err = fmt.Errorf("offerer ID %s does not own gundam ID %d", offererID, gundamID)
			c.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		
		// Kiểm tra Gundam có được phép trao đổi không
		if offererGundam.Status != db.GundamStatusInstore {
			err = fmt.Errorf("offerer gundam ID %d is not available for exchange, current status: %s", gundamID, offererGundam.Status)
			c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
			return
		}
	}
	
	// TODO: Có thể kiểm tra Gundam đã tham gia các đề xuất khác chưa? Nếu có thì không cho phép tạo đề xuất mới.
	// Nhưng hiện tại chỉ cần kiểm tra trạng thái của Gundam là được rồi.
	
	// Tạo đề xuất trao đổi
	result, err := server.dbStore.CreateExchangeOfferTx(c.Request.Context(), db.CreateExchangeOfferTxParams{
		PostID:             postID,
		OffererID:          offererID,
		PosterGundamIDs:    req.PosterGundamIDs,
		OffererGundamIDs:   req.OffererGundamIDs,
		CompensationAmount: req.CompensationAmount,
		PayerID:            req.PayerID,
		Note:               req.Note,
	})
	if err != nil {
		if errors.Is(err, db.ErrExchangeOfferUnique) {
			err = fmt.Errorf("user ID %s already has an offer for exchange post ID %s", offererID, req.ExchangePostID)
			c.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Queue(worker.QueueCritical),
	}
	
	// Tạo thông báo ngắn gọn
	var message string
	if len(req.PosterGundamIDs) == 1 {
		// Lấy tên Gundam đầu tiên để hiển thị
		gundam, err := server.dbStore.GetGundamByID(c.Request.Context(), req.PosterGundamIDs[0])
		if err == nil {
			message = fmt.Sprintf("Có đề xuất trao đổi mới cho %s của bạn", gundam.Name)
		} else {
			message = "Có đề xuất trao đổi mới cho Gundam của bạn"
		}
	} else {
		message = fmt.Sprintf("Có đề xuất trao đổi mới cho %d Gundam của bạn", len(req.PosterGundamIDs))
	}
	
	// Gửi thông báo cho người đăng bài về đề xuất trao đổi mới
	err = server.taskDistributor.DistributeTaskSendNotification(c.Request.Context(), &worker.PayloadSendNotification{
		RecipientID: post.UserID,
		Title:       "Đề xuất trao đổi mới",
		Message:     message,
		Type:        "exchange",
		ReferenceID: result.Offer.ID.String(),
	}, opts...)
	if err != nil {
		log.Err(err).Msgf("failed to send notification to user ID %s", post.UserID)
	}
	
	c.JSON(http.StatusCreated, result)
}

// PostOfferURIParams định nghĩa tham số trên URI
type PostOfferURIParams struct {
	PostID  string `uri:"postID" binding:"required,uuid"`
	OfferID string `uri:"offerID" binding:"required,uuid"`
}

// requestNegotiationForOfferRequest là cấu trúc yêu cầu thương lượng
type requestNegotiationForOfferRequest struct {
	Note *string `json:"note"` // Ghi chú từ người yêu cầu thương lượng, không bắt buộc
}

//	@Summary		Request negotiation for an exchange offer
//	@Description	As a post owner, request negotiation with an offerer.
//	@Tags			exchanges
//	@Accept			json
//	@Produce		json
//	@Security		accessToken
//	@Param			postID	path		string									true	"Exchange Post ID"
//	@Param			offerID	path		string									true	"Exchange Offer ID"
//	@Param			request	body		requestNegotiationForOfferRequest		false	"Negotiation request"
//	@Success		200		{object}	db.RequestNegotiationForOfferTxResult	"Negotiation request response"
//	@Router			/users/me/exchange-posts/{postID}/offers/{offerID}/negotiate [patch]
func (server *Server) requestNegotiationForOffer(c *gin.Context) {
	// Lấy thông tin người dùng đã đăng nhập
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	userID := authPayload.Subject
	
	// Bind các tham số từ URI
	var uriParams PostOfferURIParams
	if err := c.ShouldBindUri(&uriParams); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	// Parse UUID từ string
	postID, err := uuid.Parse(uriParams.PostID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid post ID: %s", uriParams.PostID)))
		return
	}
	
	offerID, err := uuid.Parse(uriParams.OfferID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid offer ID: %s", uriParams.OfferID)))
		return
	}
	
	// Đọc request body
	var req requestNegotiationForOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	// -------------------
	// PHẦN 1: Kiểm tra business rules
	// -------------------
	
	// 1. Kiểm tra bài đăng tồn tại và người dùng là chủ sở hữu
	post, err := server.dbStore.GetExchangePost(c, postID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			err = fmt.Errorf("exchange post ID %s not found", uriParams.PostID)
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	if post.UserID != userID {
		err = fmt.Errorf("user ID %s is not the owner of exchange post ID %s", userID, uriParams.PostID)
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	
	if post.Status != db.ExchangePostStatusOpen {
		err = fmt.Errorf("exchange post ID %s is not open for negotiation", uriParams.PostID)
		c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}
	
	// 2. Kiểm tra đề xuất tồn tại và thuộc về bài đăng này
	offer, err := server.dbStore.GetExchangeOffer(c, offerID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			err = fmt.Errorf("exchange offer ID %s not found", uriParams.OfferID)
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	if offer.PostID != postID {
		err = fmt.Errorf("exchange offer ID %s does not belong to exchange post ID %s", uriParams.OfferID, uriParams.PostID)
		c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}
	
	// 3. Kiểm tra số lần thương lượng đã sử dụng
	if offer.NegotiationsCount >= offer.MaxNegotiations {
		err = fmt.Errorf("maximum number of negotiations reached for exchange offer ID %s, current count: %d, max: %d", uriParams.OfferID, offer.NegotiationsCount, offer.MaxNegotiations)
		c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}
	
	// 4. Kiểm tra xem hiện tại có đang yêu cầu thương lượng không
	if offer.NegotiationRequested {
		err = fmt.Errorf("negotiation already requested for exchange offer ID %s", uriParams.OfferID)
		c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}
	
	// -------------------
	// PHẦN 2: Xử lý transaction để cập nhật dữ liệu
	// -------------------
	
	// Thực hiện transaction
	result, err := server.dbStore.RequestNegotiationForOfferTx(c, db.RequestNegotiationForOfferTxParams{
		OfferID: offerID,
		UserID:  userID,
		Note:    req.Note,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	// Gửi thông báo cho người đề xuất về yêu cầu thương lượng
	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Queue(worker.QueueCritical),
	}
	
	// Tạo thông báo ngắn gọn nhưng đủ thông tin
	notificationMessage := fmt.Sprintf(
		"Chủ bài đăng '%s' đã yêu cầu thương lượng cho đề xuất trao đổi Gundam của bạn.",
		util.TruncateString(post.Content, 20), // Hàm rút gọn tiêu đề nếu quá dài
	)
	
	err = server.taskDistributor.DistributeTaskSendNotification(c.Request.Context(), &worker.PayloadSendNotification{
		RecipientID: offer.OffererID,
		Title:       "Yêu cầu thương lượng Gundam",
		Message:     notificationMessage,
		Type:        "exchange",
		ReferenceID: result.Offer.ID.String(),
	}, opts...)
	if err != nil {
		log.Err(err).
			Str("offerID", result.Offer.ID.String()).
			Str("postID", post.ID.String()).
			Msgf("failed to send notification to user ID %s", offer.OffererID)
	}
	
	// Trả về kết quả
	c.JSON(http.StatusOK, result)
}

type updateExchangeOfferRequest struct {
	RequireCompensation bool    `json:"require_compensation" binding:"required"` // true = yêu cầu bù tiền, false = không yêu cầu bù tiền
	CompensationAmount  *int64  `json:"compensation_amount"`                     // Bắt buộc khi require_compensation=true
	PayerID             *string `json:"payer_id"`                                // ID người trả tiền bù, bắt buộc khi require_compensation=true
	Note                *string `json:"note"`                                    // Tin nhắn thương lượng, không bắt buộc
}

//	@Summary		Update an exchange offer
//	@Description	As an offerer, update exchange offer details. Only allowed when a negotiation is requested by the post owner.
//	@Tags			exchanges
//	@Accept			json
//	@Produce		json
//	@Security		accessToken
//	@Param			offerID	path		string							true	"Exchange Offer OfferID"
//	@Param			request	body		updateExchangeOfferRequest		true	"Update offer request"
//	@Success		200		{object}	db.UpdateExchangeOfferTxResult	"Updated offer response"
//	@Router			/users/me/exchange-offers/{offerID} [patch]
func (server *Server) updateExchangeOffer(c *gin.Context) {
	// Lấy thông tin người dùng đã đăng nhập
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	userID := authPayload.Subject
	
	// Lấy OfferID của đề xuất từ URI
	offerIDStr := c.Param("offerID")
	offerID, err := uuid.Parse(offerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid offer ID: %s", offerIDStr)))
		return
	}
	
	// Đọc request body
	var req updateExchangeOfferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	// Kiểm tra thông tin bù tiền khi yêu cầu bù tiền
	if req.RequireCompensation {
		if req.CompensationAmount == nil {
			err = errors.New("compensation_amount is required when require_compensation is true")
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		
		if req.PayerID == nil {
			err = errors.New("payer_id is required when require_compensation is true")
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		
		if *req.CompensationAmount <= 0 {
			err = errors.New("compensation_amount must be positive when require_compensation is true")
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	
	// -------------------
	// PHẦN 1: Kiểm tra business rules
	// -------------------
	
	// 1. Kiểm tra đề xuất tồn tại và người dùng là người đề xuất
	offer, err := server.dbStore.GetExchangeOffer(c, offerID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			err = fmt.Errorf("exchange offer ID %s not found", offerIDStr)
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	if offer.OffererID != userID {
		err = fmt.Errorf("user ID %s is not the offerer of exchange offer ID %s", userID, offerIDStr)
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	
	// 2. Kiểm tra xem có yêu cầu thương lượng không
	if !offer.NegotiationRequested {
		err = fmt.Errorf("cannot update exchange offer ID %s as there is no negotiation requested", offerIDStr)
		c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}
	
	// 3. Kiểm tra số lần thương lượng chưa vượt quá giới hạn
	if offer.NegotiationsCount >= offer.MaxNegotiations {
		err = fmt.Errorf("maximum number of negotiations reached for exchange offer ID %s, current count: %d, max: %d", offerIDStr, offer.NegotiationsCount, offer.MaxNegotiations)
		c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}
	
	// Lấy thông tin bài đăng (cần thiết cho nhiều phần)
	post, err := server.dbStore.GetExchangePost(c, offer.PostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	// 4. Kiểm tra PayerID hợp lệ
	if req.RequireCompensation && req.PayerID != nil {
		validPayerID := *req.PayerID == offer.OffererID || *req.PayerID == post.UserID
		if !validPayerID {
			err = errors.New("payer_id must be either the offerer or the post owner")
			c.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	
	// 5. Kiểm tra số dư của người đề xuất nếu họ là người bù tiền
	if req.RequireCompensation && req.PayerID != nil && req.CompensationAmount != nil && *req.PayerID == userID {
		// Lấy thông tin ví và kiểm tra số dư
		wallet, err := server.dbStore.GetWalletByUserID(c, userID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				err = fmt.Errorf("wallet not found for user ID %s", userID)
				c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		// Kiểm tra số dư có đủ không
		if wallet.Balance < *req.CompensationAmount {
			err = fmt.Errorf("insufficient balance for compensation: needed %d, available %d", *req.CompensationAmount, wallet.Balance)
			c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
			return
		}
	}
	
	// -------------------
	// PHẦN 2: Xử lý transaction để cập nhật dữ liệu
	// -------------------
	
	arg := db.UpdateExchangeOfferTxParams{
		OfferID:              offerID,
		UserID:               userID,
		Note:                 req.Note,
		NegotiationRequested: util.BoolPointer(false),
		NegotiationsCount:    util.Int64Pointer(offer.NegotiationsCount + 1),
	}
	
	// Xử lý thông tin bù tiền
	if req.RequireCompensation {
		arg.CompensationAmount = req.CompensationAmount
		arg.PayerID = req.PayerID
	}
	
	// Thực hiện transaction
	result, err := server.dbStore.UpdateExchangeOfferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	// Gửi thông báo cho người đăng bài về việc đề xuất đã được cập nhật
	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Queue(worker.QueueCritical),
	}
	
	// Lấy thông tin người đề xuất
	user, err := server.dbStore.GetUserByID(c, userID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			err = fmt.Errorf("user OfferID %s not found", userID)
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	// Tạo thông báo đơn giản hơn
	notificationMessage := fmt.Sprintf(
		"Người đề xuất %s đã cập nhật lại đề xuất của họ.",
		user.FullName,
	)
	
	err = server.taskDistributor.DistributeTaskSendNotification(c.Request.Context(), &worker.PayloadSendNotification{
		RecipientID: post.UserID,
		Title:       "Cập nhật đề xuất trao đổi",
		Message:     notificationMessage,
		Type:        "exchange",
		ReferenceID: result.Offer.ID.String(),
	}, opts...)
	if err != nil {
		log.Err(err).Msgf("failed to send notification to user ID %s", post.UserID)
	}
	
	// Trả về kết quả
	c.JSON(http.StatusOK, result)
}

//	@Summary		Accept an exchange offer
//	@Description	As a post owner, accept an exchange offer. This will create an exchange transaction and related orders.
//	@Tags			exchanges
//	@Produce		json
//	@Security		accessToken
//	@Param			postID	path		string							true	"Exchange Post ID"
//	@Param			offerID	path		string							true	"Exchange Offer ID"
//	@Success		200		{object}	db.AcceptExchangeOfferTxResult	"Accepted offer response"
//	@Router			/users/me/exchange-posts/{postID}/offers/{offerID}/accept [patch]
func (server *Server) acceptExchangeOffer(c *gin.Context) {
	// Lấy thông tin người dùng đã đăng nhập
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	userID := authPayload.Subject
	
	var uri PostOfferURIParams
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	// Parse UUID từ string
	postID, err := uuid.Parse(uri.PostID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid postID: %s", uri.PostID)))
		return
	}
	
	offerID, err := uuid.Parse(uri.OfferID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid offerID: %s", uri.OfferID)))
		return
	}
	
	// -------------------
	// PHẦN 1: Kiểm tra business rules
	// -------------------
	
	// 1. Kiểm tra bài đăng tồn tại và người dùng là chủ bài đăng
	post, err := server.dbStore.GetExchangePost(c, postID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			err = fmt.Errorf("exchange post ID %s not found", postID.String())
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	if post.UserID != userID {
		err = fmt.Errorf("user ID %s is not the owner of exchange post ID %s", userID, post.ID.String())
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	
	// 2. Kiểm tra trạng thái bài đăng
	if post.Status != db.ExchangePostStatusOpen {
		err = fmt.Errorf("exchange post ID %s is not open, current status: %s", post.ID.String(), post.Status)
		c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return
	}
	
	// 3. Kiểm tra đề xuất tồn tại và thuộc về bài đăng
	offer, err := server.dbStore.GetExchangeOffer(c, offerID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			err = fmt.Errorf("exchange offer ID %s not found", uri.OfferID)
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	if offer.PostID != postID {
		err = fmt.Errorf("exchange offer ID %s does not belong to post ID %s", offer.ID.String(), post.ID.String())
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	
	// 4. Kiểm tra số dư của người trả tiền bù (nếu có)
	if offer.PayerID != nil && offer.CompensationAmount != nil && *offer.CompensationAmount > 0 {
		compensationAmount := *offer.CompensationAmount
		payerID := *offer.PayerID
		isPayerPoster := payerID == post.UserID
		isPayerOfferer := payerID == offer.OffererID
		
		// Lấy thông tin ví và kiểm tra số dư
		wallet, err := server.dbStore.GetWalletByUserID(c, payerID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				err = fmt.Errorf("wallet not found for user ID %s", payerID)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		// Kiểm tra số dư có đủ không
		if wallet.Balance < compensationAmount {
			switch {
			case isPayerPoster:
				err = fmt.Errorf("poster ID %s has insufficient balance for compensation: needed %d, available %d", payerID, compensationAmount, wallet.Balance)
			case isPayerOfferer:
				err = fmt.Errorf("offerer ID %s has insufficient balance for compensation: needed %d, available %d", payerID, compensationAmount, wallet.Balance)
			default:
				err = fmt.Errorf("payer ID %s has insufficient balance for compensation: needed %d, available %d", payerID, compensationAmount, wallet.Balance)
			}
			
			c.JSON(http.StatusUnprocessableEntity, errorResponse(err))
			return
		}
	}
	
	// -------------------
	// PHẦN 2: Xử lý transaction để chấp nhận đề xuất
	// -------------------
	
	arg := db.AcceptExchangeOfferTxParams{
		PostID:    postID,
		OfferID:   offerID,
		PosterID:  post.UserID,
		OffererID: offer.OffererID,
	}
	
	// Thêm thông tin bù tiền nếu có
	if offer.PayerID != nil && offer.CompensationAmount != nil && *offer.CompensationAmount > 0 {
		arg.CompensationAmount = offer.CompensationAmount
		arg.PayerID = offer.PayerID
	}
	
	// Thực hiện transaction
	result, err := server.dbStore.AcceptExchangeOfferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	opts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Queue(worker.QueueCritical),
	}
	
	// Gửi thông báo cho người bị trừ tiền bù (nếu có)
	if result.Exchange.PayerID != nil && result.Exchange.CompensationAmount != nil && *result.Exchange.CompensationAmount > 0 {
		err = server.taskDistributor.DistributeTaskSendNotification(c.Request.Context(), &worker.PayloadSendNotification{
			RecipientID: *result.Exchange.PayerID,
			Title:       "Thanh toán tiền bù cho giao dịch trao đổi",
			Message:     fmt.Sprintf("Số tiền %s đã được trừ từ ví của bạn để bù tiền cho giao dịch trao đổi Gundam.", util.FormatVND(*result.Exchange.CompensationAmount)),
			Type:        "exchange",
			ReferenceID: result.Exchange.ID.String(),
		}, opts...)
		if err != nil {
			log.Err(err).Msgf("failed to send notification to user ID %s", *result.Exchange.PayerID)
		}
		
	}
	
	// Gửi thông báo cho người được cộng tiền bù (nếu có)
	if result.Exchange.PayerID != nil && result.Exchange.CompensationAmount != nil && *result.Exchange.CompensationAmount > 0 {
		// Xác định người nhận tiền bù (người còn lại)
		var compensationReceiverID string
		if *result.Exchange.PayerID == result.Exchange.PosterID {
			compensationReceiverID = result.Exchange.OffererID
		} else {
			compensationReceiverID = result.Exchange.PosterID
		}
		
		err = server.taskDistributor.DistributeTaskSendNotification(c.Request.Context(), &worker.PayloadSendNotification{
			RecipientID: compensationReceiverID,
			Title:       "Nhận tiền bù cho giao dịch trao đổi",
			Message:     fmt.Sprintf("Bạn đã nhận được %s tiền bù cho giao dịch trao đổi Gundam. Số tiền này sẽ được cộng vào số dư tạm thời cho đến khi cuộc trao đổi hoàn tất.", util.FormatVND(*result.Exchange.CompensationAmount)),
			Type:        "exchange",
			ReferenceID: result.Exchange.ID.String(),
		}, opts...)
		if err != nil {
			log.Err(err).Msgf("failed to send notification to user ID %s", compensationReceiverID)
		}
	}
	
	// Gửi thông báo cho người có đề xuất được chấp nhận.
	err = server.taskDistributor.DistributeTaskSendNotification(c.Request.Context(), &worker.PayloadSendNotification{
		RecipientID: offer.OffererID,
		Title:       "Đề xuất trao đổi đã được chấp nhận",
		Message:     fmt.Sprintf("Đề xuất trao đổi của bạn cho bài đăng \"%s\" đã được chấp nhận. Vui lòng cung cấp thêm thông tin vận chuyển để hệ thống tạo đơn hàng cho bạn.", util.TruncateString(post.Content, 20)),
		Type:        "exchange",
		ReferenceID: result.Exchange.ID.String(),
	}, opts...)
	if err != nil {
		log.Err(err).Msgf("failed to send notification to user ID %s", offer.OffererID)
	}
	
	// Gửi thông báo cho những người khác có đề xuất không được chấp nhận.
	for _, rejectedOffer := range result.RejectedOffers {
		err = server.taskDistributor.DistributeTaskSendNotification(c.Request.Context(), &worker.PayloadSendNotification{
			RecipientID: rejectedOffer.OffererID,
			Title:       "Đề xuất trao đổi không được chấp nhận",
			Message:     fmt.Sprintf("Đề xuất trao đổi của bạn cho bài đăng \"%s\" đã không được chấp nhận.", util.TruncateString(post.Content, 20)),
			Type:        "exchange",
			ReferenceID: result.Exchange.ID.String(),
		}, opts...)
		if err != nil {
			log.Err(err).Msgf("failed to send notification to user ID %s", rejectedOffer.OffererID)
		}
	}
	
	c.JSON(http.StatusOK, result.Exchange)
}

//	@Summary		List user's exchange offers
//	@Description	Get a list of all exchange offers created by the authenticated user, including details about the exchange posts, items, and negotiation notes.
//	@Tags			exchanges
//	@Produce		json
//	@Security		accessToken
//	@Success		200	{array}	db.UserExchangeOfferDetails	"List of user's exchange offers"
//	@Router			/users/me/exchange-offers [get]
func (server *Server) listUserExchangeOffers(c *gin.Context) {
	// Lấy thông tin người dùng đã đăng nhập
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	userID := authPayload.Subject
	
	offers, err := server.dbStore.ListExchangeOffersByOfferer(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	var result []db.UserExchangeOfferDetails
	
	for _, offer := range offers {
		var offerDetails db.UserExchangeOfferDetails
		
		offerInfo := db.ExchangeOfferInfo{
			ID:     offer.ID,
			PostID: offer.PostID,
			// Offerer:              db.User{},
			PayerID:            offer.PayerID,
			CompensationAmount: offer.CompensationAmount,
			Note:               offer.Note,
			// OffererExchangeItems: nil,
			// PosterExchangeItems:  nil,
			NegotiationsCount:    offer.NegotiationsCount,
			MaxNegotiations:      offer.MaxNegotiations,
			NegotiationRequested: offer.NegotiationRequested,
			LastNegotiationAt:    offer.LastNegotiationAt,
			// NegotiationNotes:                nil,
			CreatedAt: offer.CreatedAt,
			UpdatedAt: offer.UpdatedAt,
		}
		
		// Lấy thông tin bài đăng của offer
		post, err := server.dbStore.GetExchangePost(c.Request.Context(), offer.PostID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				err = fmt.Errorf("exchange post ID %s not found", offer.PostID)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		offerDetails.ExchangePost = post
		
		// Lấy thông tin chủ bài đăng
		poster, err := server.dbStore.GetUserByID(c.Request.Context(), post.UserID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				err = fmt.Errorf("user ID %s not found", post.UserID)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		offerDetails.Poster = poster
		
		// Lấy thông tin các item của bài đăng
		postItems, err := server.dbStore.ListExchangePostItems(c.Request.Context(), post.ID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				err = fmt.Errorf("no postItems found for exchange post ID %s", post.ID)
				c.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		// Lấy thông tin chi tiết của từng item
		postGundams := make([]db.GundamDetails, 0, len(postItems))
		for _, item := range postItems {
			gundam, err := server.dbStore.GetGundamDetailsByID(c.Request.Context(), nil, item.GundamID)
			if err != nil {
				if errors.Is(err, db.ErrRecordNotFound) {
					err = fmt.Errorf("gundam ID %d not found", item.GundamID)
					c.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			
			postGundams = append(postGundams, gundam)
		}
		offerDetails.ExchangePostItems = postGundams
		
		// Lấy các item từ bài đăng mà người đề xuất muốn trao đổi
		posterItems, err := server.dbStore.ListExchangeOfferItems(c.Request.Context(), db.ListExchangeOfferItemsParams{
			OfferID:      offer.ID,
			IsFromPoster: util.BoolPointer(true),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		// Lấy thông tin chi tiết của từng posterItems
		posterGundams := make([]db.GundamDetails, 0, len(posterItems))
		for _, item := range posterItems {
			gundam, err := server.dbStore.GetGundamDetailsByID(c.Request.Context(), nil, item.GundamID)
			if err != nil {
				if errors.Is(err, db.ErrRecordNotFound) {
					err = fmt.Errorf("gundam ID %d not found", item.GundamID)
					c.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			
			posterGundams = append(posterGundams, gundam)
		}
		offerInfo.PosterExchangeItems = posterGundams
		
		// Lấy các offer item từ người đề xuất
		offererItems, err := server.dbStore.ListExchangeOfferItems(c.Request.Context(), db.ListExchangeOfferItemsParams{
			OfferID:      offer.ID,
			IsFromPoster: util.BoolPointer(false),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		// Lấy thông tin chi tiết của từng offererItems
		offererGundams := make([]db.GundamDetails, 0, len(offererItems))
		for _, item := range offererItems {
			gundam, err := server.dbStore.GetGundamDetailsByID(c.Request.Context(), nil, item.GundamID)
			if err != nil {
				if errors.Is(err, db.ErrRecordNotFound) {
					err = fmt.Errorf("gundam ID %d not found", item.GundamID)
					c.JSON(http.StatusNotFound, errorResponse(err))
					return
				}
				
				c.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			
			offererGundams = append(offererGundams, gundam)
		}
		offerInfo.OffererExchangeItems = offererGundams
		
		// Lấy thông các ghi chú thương lượng (nếu có)
		negotiationNotes, err := server.dbStore.ListExchangeOfferNotes(c.Request.Context(), offer.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		offerInfo.NegotiationNotes = negotiationNotes
		
		// Lấy thông tin người đề xuất (bỏ qua)
		
		offerDetails.Offer = offerInfo
		
		result = append(result, offerDetails)
	}
	
	c.JSON(http.StatusOK, result)
}

//	@Summary		Get user's exchange offer details
//	@Description	Retrieves detailed information about a specific exchange offer created by the authenticated user.
//	@Tags			exchanges
//	@Produce		json
//	@Security		accessToken
//	@Param			offerID	path		string	true	"Exchange Offer ID"
//	@Success		200		{object}	db.UserExchangeOfferDetails
//	@Failure		400		{object}	error	"Invalid offer ID"
//	@Failure		404		{object}	error	"Offer not found"
//	@Failure		403		{object}	error	"Unauthorized access"
//	@Failure		500		{object}	error	"Internal server error"
//	@Router			/users/me/exchange-offers/{offerID} [get]
func (server *Server) getUserExchangeOffer(c *gin.Context) {
	// Lấy thông tin người dùng đã đăng nhập
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	userID := authPayload.Subject
	
	// Lấy ID của đề xuất từ URL
	offerIDStr := c.Param("offerID")
	offerID, err := uuid.Parse(offerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid offer ID: %s", offerIDStr)))
		return
	}
	
	// Lấy thông tin đề xuất
	offer, err := server.dbStore.GetExchangeOffer(c.Request.Context(), offerID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("exchange offer ID %s not found", offerIDStr)))
			return
		}
		
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	// Kiểm tra quyền truy cập - chỉ chủ đề xuất mới có thể xem chi tiết
	if offer.OffererID != userID {
		err = fmt.Errorf("offer ID %v does not belong to user ID %v", offerID, userID)
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	
	var offerDetails db.UserExchangeOfferDetails
	
	offerInfo := db.ExchangeOfferInfo{
		ID:                   offer.ID,
		PostID:               offer.PostID,
		PayerID:              offer.PayerID,
		CompensationAmount:   offer.CompensationAmount,
		Note:                 offer.Note,
		NegotiationsCount:    offer.NegotiationsCount,
		MaxNegotiations:      offer.MaxNegotiations,
		NegotiationRequested: offer.NegotiationRequested,
		LastNegotiationAt:    offer.LastNegotiationAt,
		CreatedAt:            offer.CreatedAt,
		UpdatedAt:            offer.UpdatedAt,
	}
	
	// Lấy thông tin bài đăng của offer
	post, err := server.dbStore.GetExchangePost(c.Request.Context(), offer.PostID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("exchange post ID %s not found", offer.PostID)))
			return
		}
		
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	offerDetails.ExchangePost = post
	
	// Lấy thông tin chủ bài đăng
	poster, err := server.dbStore.GetUserByID(c.Request.Context(), post.UserID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user ID %s not found", post.UserID)))
			return
		}
		
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	offerDetails.Poster = poster
	
	// Lấy thông tin các item của bài đăng
	postItems, err := server.dbStore.ListExchangePostItems(c.Request.Context(), post.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	// Lấy thông tin chi tiết của từng item
	postGundams := make([]db.GundamDetails, 0, len(postItems))
	for _, item := range postItems {
		gundam, err := server.dbStore.GetGundamDetailsByID(c.Request.Context(), nil, item.GundamID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("gundam ID %d not found", item.GundamID)))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		postGundams = append(postGundams, gundam)
	}
	offerDetails.ExchangePostItems = postGundams
	
	// Lấy các item từ bài đăng mà người đề xuất muốn trao đổi
	posterItems, err := server.dbStore.ListExchangeOfferItems(c.Request.Context(), db.ListExchangeOfferItemsParams{
		OfferID:      offer.ID,
		IsFromPoster: util.BoolPointer(true),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	// Lấy thông tin chi tiết của từng posterItems
	posterGundams := make([]db.GundamDetails, 0, len(posterItems))
	for _, item := range posterItems {
		gundam, err := server.dbStore.GetGundamDetailsByID(c.Request.Context(), nil, item.GundamID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("gundam ID %d not found", item.GundamID)))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		posterGundams = append(posterGundams, gundam)
	}
	offerInfo.PosterExchangeItems = posterGundams
	
	// Lấy các offer item từ người đề xuất
	offererItems, err := server.dbStore.ListExchangeOfferItems(c.Request.Context(), db.ListExchangeOfferItemsParams{
		OfferID:      offer.ID,
		IsFromPoster: util.BoolPointer(false),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	// Lấy thông tin chi tiết của từng offererItems
	offererGundams := make([]db.GundamDetails, 0, len(offererItems))
	for _, item := range offererItems {
		gundam, err := server.dbStore.GetGundamDetailsByID(c.Request.Context(), nil, item.GundamID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("gundam ID %d not found", item.GundamID)))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		offererGundams = append(offererGundams, gundam)
	}
	offerInfo.OffererExchangeItems = offererGundams
	
	// Lấy thông tin các ghi chú thương lượng (nếu có)
	negotiationNotes, err := server.dbStore.ListExchangeOfferNotes(c.Request.Context(), offer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	offerInfo.NegotiationNotes = negotiationNotes
	
	offerDetails.Offer = offerInfo
	
	c.JSON(http.StatusOK, offerDetails)
}

//	@Summary		Delete an exchange offer
//	@Description	Delete an exchange offer created by the authenticated user.
//	@Tags			exchanges
//	@Produce		json
//	@Security		accessToken
//	@Param			offerID	path		string					true	"Exchange Offer ID"
//	@Success		200		{object}	db.ExchangeOfferInfo	"Deleted offer response"
//	@Router			/users/me/exchange-offers/{offerID} [delete]
func (server *Server) deleteExchangeOffer(c *gin.Context) {
	// Lấy thông tin người dùng đã đăng nhập
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	userID := authPayload.Subject
	
	// Lấy ID của đề xuất từ URL
	offerIDStr := c.Param("offerID")
	offerID, err := uuid.Parse(offerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid offer ID: %s", offerIDStr)))
		return
	}
	
	// Lấy thông tin đề xuất
	offer, err := server.dbStore.GetExchangeOffer(c.Request.Context(), offerID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("exchange offer ID %s not found", offerIDStr)))
			return
		}
		
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	// Kiểm tra quyền truy cập - chỉ chủ đề xuất mới có thể xóa
	if offer.OffererID != userID {
		err = fmt.Errorf("offer ID %v does not belong to user ID %v", offerID, userID)
		c.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	
	// Thực hiện xóa cứng đề xuất
	deletedOffer, err := server.dbStore.DeleteExchangeOffer(c.Request.Context(), offerID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("exchange offer ID %s not found", offerIDStr)))
			return
		}
		
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	
	// Gửi thống báo cho người đăng bài về việc đề xuất đã bị xóa nếu đề xuất đang được thương lượng
	if deletedOffer.NegotiationRequested {
		// Lấy thông tin người đề xuất
		offerer, err := server.dbStore.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user ID %s not found", userID)))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		// Lấy bài đăng của đề xuất
		post, err := server.dbStore.GetExchangePost(c.Request.Context(), deletedOffer.PostID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("exchange post ID %s not found", deletedOffer.PostID)))
				return
			}
			
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		opts := []asynq.Option{
			asynq.MaxRetry(3),
			asynq.Queue(worker.QueueCritical),
		}
		
		message := fmt.Sprintf("%s đã xóa đề xuất trao đổi của họ cho bài đăng \"%s\".", offerer.FullName, util.TruncateString(post.Content, 20))
		
		err = server.taskDistributor.DistributeTaskSendNotification(c.Request.Context(), &worker.PayloadSendNotification{
			RecipientID: post.UserID,
			Title:       "Đề xuất trao đổi đã bị xóa",
			Message:     message,
			Type:        "exchange",
			ReferenceID: deletedOffer.ID.String(),
		}, opts...)
		if err != nil {
			log.Err(err).Msgf("failed to send notification to user ID %s", post.UserID)
		}
	}
	
	// Trả về thông tin đề xuất đã xóa
	c.JSON(http.StatusOK, deletedOffer)
}
