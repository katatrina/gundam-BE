package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (s *GHNService) CreateOrder(ctx context.Context, arg CreateOrderRequest) (*CreateOrderResponse, error) {
	// API endpoint để tạo đơn hàng
	url := GHNBaseURL + "/shipping-order/create"
	
	totalWeight := int64(0)
	for _, item := range arg.OrderItems {
		totalWeight += item.Weight * item.Quantity
	}
	
	// Thông tin đơn hàng
	orderData := map[string]interface{}{
		"from_name":            arg.SenderAddress.FullName,
		"from_phone":           arg.SenderAddress.PhoneNumber,
		"from_address":         arg.SenderAddress.Detail,
		"from_ward_name":       arg.SenderAddress.WardName,
		"from_district_name":   arg.SenderAddress.DistrictName,
		"from_province_name":   arg.SenderAddress.ProvinceName,
		"to_name":              arg.ReceiverAddress.FullName,
		"to_phone":             arg.ReceiverAddress.PhoneNumber,
		"to_address":           arg.ReceiverAddress.Detail,
		"to_ward_name":         arg.ReceiverAddress.WardName,
		"to_district_name":     arg.ReceiverAddress.DistrictName,
		"to_province_name":     arg.ReceiverAddress.ProvinceName,
		"return_phone":         arg.SenderAddress.PhoneNumber,
		"return_address":       arg.SenderAddress.Detail,
		"return_district_name": arg.SenderAddress.DistrictName,
		"return_ward_name":     arg.SenderAddress.WardName,
		"return_province_name": arg.SenderAddress.ProvinceName,
		"client_order_code":    arg.Order.Code,
		"cod_amount":           int64(0), // Đã thanh toán bằng ví
		"content":              "Mô hình Gundam",
		"weight":               totalWeight,
		// Sử dụng giá trị mặc định cho toàn bộ đơn hàng
		"length":          int64(40), // cm
		"width":           int64(30),
		"height":          int64(20),
		"service_type_id": int64(2), // Chọn loại dịch vụ "Hàng nhẹ" cho đơn giản
		"payment_type_id": int64(2), // Người mua thanh toán phí dịch vụ
		"required_note":   "CHOXEMHANGKHONGTHU",
		"insurance_value": int64(0), // Không thêm phí bảo hiểm cho môi trường test
		// TODO: Thêm các thông tin khác nếu cần thiết
	}
	
	// Chuyển đổi dữ liệu thành JSON
	jsonData, err := json.Marshal(orderData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GHN order data: %w", err)
	}
	
	// Tạo request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create GHN order request: %w", err)
	}
	
	// Thiết lập header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", s.Token)
	req.Header.Set("ShopId", s.ShopID)
	
	// Gửi request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// Đọc response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// Kiểm tra mã trạng thái HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GHN API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}
	
	// Parse response
	var response CreateOrderResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse GHN order response: %w", err)
	}
	
	// Kiểm tra code trong response body
	if response.Code != int64(http.StatusOK) {
		return nil, fmt.Errorf("GHN API returned business error: code=%d, message=%s",
			response.Code, response.Message)
	}
	
	return &response, nil
}
