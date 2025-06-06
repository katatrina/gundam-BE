package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	
	"github.com/gin-gonic/gin"
	db "github.com/katatrina/gundam-BE/internal/db/sqlc"
	"github.com/katatrina/gundam-BE/internal/token"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
	authorizationPayloadKey = "authPayload"
	
	sellerPayloadKey    = "sellerPayload"
	moderatorPayloadKey = "moderatorPayload"
	adminPayloadKey     = "adminPayload"
)

// authMiddleware authenticates the user.
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if authorizationHeader == "" {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		
		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		
		authorizationHeaderType := fields[0]
		if authorizationHeaderType != authorizationTypeBearer {
			err := errors.New("unsupported authorization header type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func optionalAuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		
		// Nếu không có header xác thực, vẫn cho phép tiếp tục
		if authorizationHeader == "" {
			ctx.Set(authorizationPayloadKey, nil) // Set payload là nil để biết là chưa xác thực
			ctx.Next()
			return
		}
		
		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			// Định dạng không đúng nhưng vẫn cho phép tiếp tục
			ctx.Set(authorizationPayloadKey, nil)
			ctx.Next()
			return
		}
		
		authorizationHeaderType := fields[0]
		if authorizationHeaderType != authorizationTypeBearer {
			// Loại header không được hỗ trợ nhưng vẫn cho phép tiếp tục
			ctx.Set(authorizationPayloadKey, nil)
			ctx.Next()
			return
		}
		
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			// Token không hợp lệ nhưng vẫn cho phép tiếp tục
			ctx.Set(authorizationPayloadKey, nil)
			ctx.Next()
			return
		}
		
		// Nếu token hợp lệ, lưu payload vào context
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func requiredSellerRole(dbStore db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		authenticatedUserID := authPayload.Subject
		sellerID := ctx.Param("sellerID")
		
		seller, err := dbStore.GetUserByID(ctx, sellerID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				err = fmt.Errorf("user ID %s not found", sellerID)
				ctx.AbortWithStatusJSON(http.StatusNotFound, errorResponse(err))
				return
			}
			
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		if seller.Role != db.UserRoleSeller {
			err = fmt.Errorf("user ID %s is not a seller", sellerID)
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
		
		if authenticatedUserID != seller.ID {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(ErrSellerIDMismatch))
			return
		}
		
		ctx.Set(sellerPayloadKey, &seller)
		ctx.Next()
	}
}

func requiredModeratorRole(dbStore db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		authenticatedUserID := authPayload.Subject
		
		user, err := dbStore.GetUserByID(ctx, authenticatedUserID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				err = fmt.Errorf("user ID %s not found", authenticatedUserID)
				ctx.AbortWithStatusJSON(http.StatusNotFound, errorResponse(err))
				return
			}
			
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		if user.Role != db.UserRoleModerator {
			err = fmt.Errorf("user ID %s is not a moderator", authenticatedUserID)
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
		
		ctx.Set(moderatorPayloadKey, &user)
		ctx.Next()
	}
}

func requiredAdminRole(dbStore db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		authenticatedUserID := authPayload.Subject
		
		user, err := dbStore.GetUserByID(ctx, authenticatedUserID)
		if err != nil {
			if errors.Is(err, db.ErrRecordNotFound) {
				err = fmt.Errorf("user ID %s not found", authenticatedUserID)
				ctx.AbortWithStatusJSON(http.StatusNotFound, errorResponse(err))
				return
			}
			
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		
		if user.Role != db.UserRoleAdmin {
			err = fmt.Errorf("user ID %s is not an admin", authenticatedUserID)
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
			return
		}
		
		ctx.Set(adminPayloadKey, &user)
		ctx.Next()
	}
}
