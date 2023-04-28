package sqlc

import (
	"context"
)

type Create_UserTxParams struct {
	Create_UserParams
	AfterCreate func(user User) error
}

func (store *SQLStore) Create_UserTx(ctx context.Context, arg Create_UserTxParams) (User, error) {
	var user User
	err := store.execTx(ctx, func(q *Queries) (err error) {
		user, err = q.Create_User(ctx, arg.Create_UserParams)
		if err != nil {
			return err
		}
		if arg.AfterCreate != nil {
			if err = arg.AfterCreate(user); err != nil {
				return err
			}
		}
		return nil
	})

	return user, err
}

type Update_User_VerifyTokenTxParams struct {
	Update_User_VerifyTokenParams
	BeforeUpdate func() error
}

func (store *SQLStore) Update_User_VerifyTokenTx(ctx context.Context, arg Update_User_VerifyTokenTxParams) error {
	return store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}
		if err = q.Update_User_VerifyToken(ctx, arg.Update_User_VerifyTokenParams); err != nil {
			return err
		}
		return
	})
}

type Update_User_EmailVerifiedTxParams struct {
	Update_User_EmailVerifiedParams
	AfterUpdate func() error
}

func (store *SQLStore) Update_User_EmailVerifiedTx(ctx context.Context, arg Update_User_EmailVerifiedTxParams) error {
	return store.execTx(ctx, func(q *Queries) (err error) {
		if err = q.Update_User_EmailVerified(ctx, arg.Update_User_EmailVerifiedParams); err != nil {
			return err
		}
		if arg.AfterUpdate != nil {
			if err = arg.AfterUpdate(); err != nil {
				return err
			}
		}
		return
	})
}

type Update_User_RecoveryTokenTxParams struct {
	Update_User_RecoveryTokenParams
	BeforeUpdate func() error
}

func (store *SQLStore) Update_User_RecoveryTokenTx(ctx context.Context, arg Update_User_RecoveryTokenTxParams) error {
	return store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}
		if err = q.Update_User_RecoveryToken(ctx, arg.Update_User_RecoveryTokenParams); err != nil {
			return err
		}
		return
	})
}

type Update_User_ResetPasswordTxParams struct {
	Update_User_ResetPasswordParams
	BeforeUpdate func() error
}

func (store *SQLStore) Update_User_ResetPasswordTx(ctx context.Context, arg Update_User_ResetPasswordTxParams) error {
	return store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}
		if err = q.Update_User_ResetPassword(ctx, arg.Update_User_ResetPasswordParams); err != nil {
			return err
		}
		return
	})
}

type Update_User_DeleteTokenTxParams struct {
	Update_User_DeleteTokenParams
	BeforeUpdate func() error
}

func (store *SQLStore) Update_User_DeleteTokenTx(ctx context.Context, arg Update_User_DeleteTokenTxParams) error {
	return store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}
		if err = q.Update_User_DeleteToken(ctx, arg.Update_User_DeleteTokenParams); err != nil {
			return err
		}
		return
	})
}
