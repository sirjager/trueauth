package sqlc

import "context"

type QuerierExtended interface {
	Create_UserTx(ctx context.Context, arg Create_UserTxParams) (User, error)

	Update_User_VerifyTokenTx(ctx context.Context, arg Update_User_VerifyTokenTxParams) error
	Update_User_EmailVerifiedTx(ctx context.Context, arg Update_User_EmailVerifiedTxParams) error

	Update_User_RecoveryTokenTx(ctx context.Context, arg Update_User_RecoveryTokenTxParams) error
	Update_User_ResetPasswordTx(ctx context.Context, arg Update_User_ResetPasswordTxParams) error

	Update_User_DeleteTokenTx(ctx context.Context, arg Update_User_DeleteTokenTxParams) error
}
