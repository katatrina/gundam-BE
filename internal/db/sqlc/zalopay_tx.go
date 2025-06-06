package db

import (
	"context"
	
	"github.com/katatrina/gundam-BE/internal/util"
)

type HandleZalopayCallbackTxParams struct {
	AppTransID string
	AppUser    string
}

func (store *SQLStore) HandleZalopayCallbackTx(ctx context.Context, arg HandleZalopayCallbackTxParams) error {
	err := store.ExecTx(ctx, func(qTx *Queries) error {
		// Kiểm tra giao dịch đã tồn tại trong DB chưa
		transaction, err := qTx.GetPaymentTransactionByProviderID(ctx, GetPaymentTransactionByProviderIDParams{
			Provider:              PaymentTransactionProviderZalopay,
			ProviderTransactionID: arg.AppTransID,
			UserID:                arg.AppUser,
		})
		if err != nil {
			return err
		}
		
		// Nếu đã xử lý rồi thì không làm gì nữa
		if transaction.Status == PaymentTransactionStatusCompleted {
			return nil
		}
		
		switch transaction.TransactionType {
		case PaymentTransactionTypeWalletDeposit:
			// Lấy thông tin ví người dùng
			wallet, err := qTx.GetWalletForUpdate(ctx, arg.AppUser)
			if err != nil {
				return err
			}
			
			// Cộng tiền vào ví người dùng
			_, err = qTx.AddWalletBalance(ctx, AddWalletBalanceParams{
				UserID: wallet.UserID,
				Amount: transaction.Amount,
			})
			if err != nil {
				return err
			}
			
			// Tạo bút toán nạp tiền vào ví
			_, err = qTx.CreateWalletEntry(ctx, CreateWalletEntryParams{
				WalletID:      wallet.UserID,
				ReferenceID:   &transaction.ProviderTransactionID,
				ReferenceType: WalletReferenceTypeDepositRequest,
				EntryType:     WalletEntryTypeDeposit,
				AffectedField: WalletAffectedFieldBalance,
				Amount:        transaction.Amount, // Số tiền nạp vào ví
				Status:        WalletEntryStatusCompleted,
				CompletedAt:   util.TimePointer(transaction.CreatedAt),
			})
			if err != nil {
				return err
			}
			
			// Cập nhật trạng thái giao dịch thanh toán thành công
			err = qTx.UpdatePaymentTransactionStatus(ctx, UpdatePaymentTransactionStatusParams{
				Status:                PaymentTransactionStatusCompleted,
				ProviderTransactionID: transaction.ProviderTransactionID,
				Provider:              PaymentTransactionProviderZalopay,
				UserID:                transaction.UserID,
			})
			if err != nil {
				return err
			}
		default:
			return nil
		}
		
		return nil
	})
	
	return err
}
