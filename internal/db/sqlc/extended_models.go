package db

import (
	"time"
	
	"github.com/google/uuid"
)

type GundamDetails struct {
	ID                   int64                `json:"gundam_id"`
	OwnerID              string               `json:"owner_id"`
	Name                 string               `json:"name"`
	Slug                 string               `json:"slug"`
	Grade                string               `json:"grade"`
	Series               string               `json:"series"`
	PartsTotal           int64                `json:"parts_total"`
	Material             string               `json:"material"`
	Version              string               `json:"version"`
	Quantity             int64                `json:"quantity"`
	Condition            string               `json:"condition"`
	ConditionDescription *string              `json:"condition_description"`
	Manufacturer         string               `json:"manufacturer"`
	Weight               int64                `json:"weight"`
	Scale                string               `json:"scale"`
	Description          string               `json:"description"`
	Price                *int64               `json:"price"`
	ReleaseYear          *int64               `json:"release_year"`
	Status               string               `json:"status"`
	Accessories          []GundamAccessoryDTO `json:"accessories"`
	PrimaryImageURL      string               `json:"primary_image_url"`
	SecondaryImageURLs   []string             `json:"secondary_image_urls"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

type GundamAccessoryDTO struct {
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
}

func ConvertGundamAccessoryToDTO(accessory GundamAccessory) GundamAccessoryDTO {
	return GundamAccessoryDTO{
		Name:     accessory.Name,
		Quantity: accessory.Quantity,
	}
}

type Sender struct {
	User
	ShopName *string `json:"shop_name,omitempty"` // Tên shop (nếu có)
}

type MemberOrderInfo struct {
	Order      Order       `json:"order"`
	OrderItems []OrderItem `json:"order_items"`
}

type MemberOrderDetails struct {
	Sender                  Sender              `json:"sender"`                    // Thông tin người gửi hàng (null nếu là người gửi)
	IsSender                bool                `json:"is_sender"`                 // Có phải người gửi không
	Receiver                User                `json:"receiver"`                  // Thông tin người nhận hàng (null nếu là người nhận)
	IsReceiver              bool                `json:"is_receiver"`               // Có phải người nhận không
	Order                   Order               `json:"order"`                     // Thông tin đơn hàng
	OrderItems              []OrderItem         `json:"order_items"`               // Danh sách sản phẩm trong đơn hàng
	OrderDelivery           OrderDelivery       `json:"order_delivery"`            // Thông tin vận chuyển
	FromDeliveryInformation DeliveryInformation `json:"from_delivery_information"` // Địa chỉ gửi hàng
	ToDeliveryInformation   DeliveryInformation `json:"to_delivery_information"`   // Địa chỉ nhận hàng
}

type SalesOrderInfo struct {
	Order      Order       `json:"order"`
	OrderItems []OrderItem `json:"order_items"`
}

type SalesOrderDetails struct {
	Receiver                User                `json:"receiver"`                  // Thông tin người nhận hàng
	Order                   Order               `json:"order"`                     // Thông tin đơn hàng
	OrderItems              []OrderItem         `json:"order_items"`               // Danh sách sản phẩm trong đơn hàng
	OrderDelivery           OrderDelivery       `json:"order_delivery"`            // Thông tin vận chuyển
	FromDeliveryInformation DeliveryInformation `json:"from_delivery_information"` // Địa chỉ gửi hàng
	ToDeliveryInformation   DeliveryInformation `json:"to_delivery_information"`   // Địa chỉ nhận hàng
}

type OpenExchangePostInfo struct {
	ExchangePost      ExchangePost    `json:"exchange_post"`       // Thông tin bài đăng
	ExchangePostItems []GundamDetails `json:"exchange_post_items"` // Danh sách Gundam mà Người đăng bài cho phép trao đổi
	Poster            User            `json:"poster"`              // Thông tin Người đăng bài
	OfferCount        int64           `json:"offer_count"`         // Số lượng offer của bài đăng
	// AuthenticatedUserOffer       *ExchangeOffer  `json:"authenticated_user_offer"`        // Offer của người dùng đã đăng nhập (nếu có)
	// AuthenticatedUserOfferItems  []GundamDetails `json:"authenticated_user_offer_items"`  // Danh sách Gundam trong offer của người dùng đã đăng nhập (nếu có)
	// AuthenticatedUserWantedItems []GundamDetails `json:"authenticated_user_wanted_items"` // Danh sách Gundam mà người dùng đã đăng nhập muốn nhận (nếu có)
}

type UserExchangePostDetails struct {
	ExchangePost      ExchangePost        `json:"exchange_post"`       // Thông tin bài đăng
	ExchangePostItems []GundamDetails     `json:"exchange_post_items"` // Danh sách Gundam mà Người đăng bài cho phép trao đổi
	OfferCount        int64               `json:"offer_count"`         // Số lượng offer của bài đăng
	Offers            []ExchangeOfferInfo `json:"offers"`              // Danh sách các offer của bài đăng
}

type ExchangeOfferInfo struct {
	ID      uuid.UUID `json:"id"`      // ID đề xuất
	PostID  uuid.UUID `json:"post_id"` // ID bài đăng trao đổi
	Offerer User      `json:"offerer"` // Thông tin người đề xuất
	
	PayerID            *string `json:"payer_id"`            // ID người bù tiền
	CompensationAmount *int64  `json:"compensation_amount"` // Số tiền bù
	Note               *string `json:"note"`                // Ghi chú của đề xuất
	
	OffererExchangeItems []GundamDetails `json:"offerer_exchange_items"` // Danh sách Gundam của người đề xuất
	PosterExchangeItems  []GundamDetails `json:"poster_exchange_items"`  // Danh sách Gundam của người đăng bài mà người đề xuất muốn trao đổi
	
	NegotiationsCount    int64               `json:"negotiations_count"`    // Số lần đã thương lượng
	MaxNegotiations      int64               `json:"max_negotiations"`      // Số lần thương lượng tối đa
	NegotiationRequested bool                `json:"negotiation_requested"` // Đã yêu cầu thương lượng chưa
	LastNegotiationAt    *time.Time          `json:"last_negotiation_at"`   // Thời gian thương lượng gần nhất
	NegotiationNotes     []ExchangeOfferNote `json:"negotiation_notes"`     // Các ghi chú/tin nhắn thương lượng
	
	CreatedAt time.Time `json:"created_at"` // Thời gian tạo đề xuất
	UpdatedAt time.Time `json:"updated_at"` // Thời gian cập nhật đề xuất gần nhất
}

type UserExchangeOfferDetails struct {
	ExchangePost      ExchangePost      `json:"exchange_post"`       // Thông tin bài đăng
	Poster            User              `json:"poster"`              // Thông tin Người đăng bài
	ExchangePostItems []GundamDetails   `json:"exchange_post_items"` // Danh sách Gundam mà Người đăng bài cho phép trao đổi
	Offer             ExchangeOfferInfo `json:"offer"`               // Chi tiết đề xuất
}

type UserExchangeDetails struct {
	ID uuid.UUID `json:"id"` // ID của bài đăng trao đổi
	
	// Thông tin gốc về cuộc trao đổi (Ai đăng, ai đề xuất)
	PosterID  string `json:"poster_id"`  // ID người đăng bài
	OffererID string `json:"offerer_id"` // ID người đề xuất
	
	// Thông tin bù tiền
	PayerID            *string `json:"payer_id"`            // ID người trả tiền bù (nếu có)
	CompensationAmount *int64  `json:"compensation_amount"` // Số tiền bù (nếu có)
	
	// Thông tin cơ bản về cuộc trao đổi
	Status      string     `json:"status"`       // Trạng thái cuộc trao đổi
	CreatedAt   time.Time  `json:"created_at"`   // Thời gian tạo
	UpdatedAt   time.Time  `json:"updated_at"`   // Thời gian cập nhật
	CompletedAt *time.Time `json:"completed_at"` // Thời gian hoàn thành
	
	// Thông tin hủy (nếu có)
	CanceledBy     *string `json:"canceled_by"`     // ID người hủy
	CanceledReason *string `json:"canceled_reason"` // Lý do hủy
	
	// Thông tin về người tham gia
	CurrentUser ExchangeUserInfo `json:"current_user"` // Thông tin người dùng hiện tại
	Partner     ExchangeUserInfo `json:"partner"`      // Thông tin người còn lại
}

// ExchangeUserInfo chứa thông tin một bên tham gia trao đổi
type ExchangeUserInfo struct {
	// Thông tin cơ bản
	ID        string  `json:"id"`         // ID người dùng
	FullName  string  `json:"full_name"`  // Tên người dùng
	AvatarURL *string `json:"avatar_url"` // URL ảnh đại diện người dùng
	
	// Thông tin đơn hàng và vận chuyển
	Order                *Order               `json:"order"`                  // Thông tin đơn hàng (đơn hàng mà người này là người nhận)
	FromDelivery         *DeliveryInformation `json:"from_address"`           // Địa chỉ gửi hàng
	ToDelivery           *DeliveryInformation `json:"to_address"`             // Địa chỉ nhận hàng
	DeliveryFee          *int64               `json:"delivery_fee"`           // Phí vận chuyển
	DeliveryFeePaid      bool                 `json:"delivery_fee_paid"`      // Đã thanh toán phí vận chuyển chưa
	ExpectedDeliveryTime *time.Time           `json:"expected_delivery_time"` // Thời gian giao hàng dự kiến
	
	Note *string `json:"note"` // Ghi chú của người dùng
	
	Items []ExchangeItem `json:"items"` // Danh sách Gundam
}

type GundamSnapshot struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Grade    string `json:"grade"`
	Scale    string `json:"scale"`
	Quantity int64  `json:"quantity"`
	Weight   int64  `json:"weight"`
	ImageURL string `json:"image_url"`
}

type AuctionDetails struct {
	Auction             Auction              `json:"auction"`              // Thông tin phiên đấu giá
	AuctionParticipants []AuctionParticipant `json:"auction_participants"` // Danh sách người tham gia đấu giá
	AuctionBids         []AuctionBid         `json:"auction_bids"`         // Danh sách giá đấu
}

func NewWithdrawalRequestDetails(request WithdrawalRequest, bankAccount UserBankAccount) WithdrawalRequestDetails {
	return WithdrawalRequestDetails{
		ID:                   request.ID,
		UserID:               request.UserID,
		BankAccount:          bankAccount,
		Amount:               request.Amount,
		Status:               request.Status,
		ProcessedBy:          request.ProcessedBy,
		ProcessedAt:          request.ProcessedAt,
		RejectedReason:       request.RejectedReason,
		TransactionReference: request.TransactionReference,
		WalletEntryID:        request.WalletEntryID,
		CreatedAt:            request.CreatedAt,
		UpdatedAt:            request.UpdatedAt,
		CompletedAt:          request.CompletedAt,
	}
}

type WithdrawalRequestDetails struct {
	ID                   uuid.UUID               `json:"id"`
	UserID               string                  `json:"user_id"`
	BankAccount          UserBankAccount         `json:"bank_account"`
	Amount               int64                   `json:"amount"`
	Status               WithdrawalRequestStatus `json:"status"`
	ProcessedBy          *string                 `json:"processed_by"`
	ProcessedAt          *time.Time              `json:"processed_at"`
	RejectedReason       *string                 `json:"rejected_reason"`
	TransactionReference *string                 `json:"transaction_reference"`
	WalletEntryID        *int64                  `json:"wallet_entry_id"`
	CreatedAt            time.Time               `json:"created_at"`
	UpdatedAt            time.Time               `json:"updated_at"`
	CompletedAt          *time.Time              `json:"completed_at"`
}

type SellerDashboard struct {
	PublishedGundamsCount       int64 `json:"published_gundams_count"`        // Số lượng gundam đã đăng bán
	TotalIncome                 int64 `json:"total_income"`                   // Tổng thu nhập từ việc bán + đấu giá gundam
	CompletedOrdersCount        int64 `json:"completed_orders_count"`         // Số lượng đơn hàng đã hoàn thành
	ProcessingOrdersCount       int64 `json:"processing_orders_count"`        // Số lượng đơn hàng đang xử lý
	IncomeThisMonth             int64 `json:"income_this_month"`              // Thu nhập trong tháng hiện tại
	ActiveAuctionsCount         int64 `json:"active_auctions_count"`          // Số lượng phiên đấu giá đang diễn ra
	PendingAuctionRequestsCount int64 `json:"pending_auction_requests_count"` // Số lượng yêu cầu đấu giá đang chờ xử lý
}

type ModeratorDashboard struct {
	PendingAuctionRequestsCount    int64 `json:"pending_auction_requests_count"`    // Số lượng yêu cầu đấu giá đang chờ xử lý
	PendingWithdrawalRequestsCount int64 `json:"pending_withdrawal_requests_count"` // Số lượng yêu cầu rút tiền đang chờ xử lý
	TotalExchangesThisWeek         int64 `json:"total_exchanges_this_week"`         // Tổng số cuộc trao đổi trong tuần này
	TotalOrdersThisWeek            int64 `json:"total_orders_this_week"`            // Tổng số đơn hàng trong tuần này
}

type AdminDashboard struct {
	TotalBusinessUsers           int64 `json:"total_business_users"`             // Tổng số người dùng hoạt động trên nền tảng (chỉ tính role member và seller)
	TotalRegularOrdersThisMonth  int64 `json:"total_regular_orders_this_month"`  // Tổng số đơn hàng thường trong tháng này
	TotalExchangeOrdersThisMonth int64 `json:"total_exchange_orders_this_month"` // Tổng số đơn hàng trao đổi trong tháng này
	TotalAuctionOrdersThisMonth  int64 `json:"total_auction_orders_this_month"`  // Tổng số đơn hàng đấu giá trong tháng này
	TotalRevenueThisMonth        int64 `json:"total_revenue_this_month"`         // Tổng doanh thu trong tháng này
	CompletedExchangesThisMonth  int64 `json:"completed_exchanges_this_month"`   // Tổng số cuộc trao đổi đã hoàn thành trong tháng này
	CompletedAuctionsThisWeek    int64 `json:"completed_auctions_this_week"`     // Tổng số phiên đấu giá đã hoàn thành trong tuần này
	TotalWalletVolumeThisWeek    int64 `json:"total_wallet_volume_this_week"`    // Tổng khối lượng giao dịch ví trong tuần này
	TotalPublishedGundams        int64 `json:"total_published_gundams"`          // Tổng số gundam đã được đăng bán
	NewUsersThisWeek             int64 `json:"new_users_this_week"`              // Số lượng người dùng mới trong tuần này
}
