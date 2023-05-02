package sqlc

import (
	"context"
)

type UpdateUserVerifyTokenTxParams struct {
	UpdateUserVerifyTokenParams
	BeforeUpdate func() error
}

func (store *SQLStore) UpdateUserVerifyTokenTx(ctx context.Context, arg UpdateUserVerifyTokenTxParams) error {
	return store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}
		if err = q.UpdateUserVerifyToken(ctx, arg.UpdateUserVerifyTokenParams); err != nil {
			return err
		}
		return
	})
}

type UpdateUserEmailVerifiedTxParams struct {
	UpdateUserEmailVerifiedParams
	AfterUpdate func() error
}

func (store *SQLStore) UpdateUserEmailVerifiedTx(ctx context.Context, arg UpdateUserEmailVerifiedTxParams) error {
	return store.execTx(ctx, func(q *Queries) (err error) {
		if err = q.UpdateUserEmailVerified(ctx, arg.UpdateUserEmailVerifiedParams); err != nil {
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

type UpdateUserRecoveryTokenTxParams struct {
	UpdateUserRecoveryTokenParams
	BeforeUpdate func() error
}

func (store *SQLStore) UpdateUserRecoveryTokenTx(ctx context.Context, arg UpdateUserRecoveryTokenTxParams) error {
	return store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}
		if err = q.UpdateUserRecoveryToken(ctx, arg.UpdateUserRecoveryTokenParams); err != nil {
			return err
		}
		return
	})
}

type UpdateUserResetPasswordTxParams struct {
	UpdateUserResetPasswordParams
	BeforeUpdate func() error
}

func (store *SQLStore) UpdateUserResetPasswordTx(ctx context.Context, arg UpdateUserResetPasswordTxParams) error {
	return store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}
		if err = q.UpdateUserResetPassword(ctx, arg.UpdateUserResetPasswordParams); err != nil {
			return err
		}
		return
	})
}

type UpdateUserDeleteTokenTxParams struct {
	UpdateUserDeleteTokenParams
	BeforeUpdate func() error
}

func (store *SQLStore) UpdateUserDeleteTokenTx(ctx context.Context, arg UpdateUserDeleteTokenTxParams) error {
	return store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}
		if err = q.UpdateUserDeleteToken(ctx, arg.UpdateUserDeleteTokenParams); err != nil {
			return err
		}
		return
	})
}

type UpdateUserAllowIPTokenTxParams struct {
	UpdateUserAllowIPTokenParams
	BeforeUpdate func() error
}

func (store *SQLStore) UpdateUserAllowIPTokenTx(ctx context.Context, arg UpdateUserAllowIPTokenTxParams) error {
	return store.execTx(ctx, func(q *Queries) (err error) {
		if arg.BeforeUpdate != nil {
			if err = arg.BeforeUpdate(); err != nil {
				return err
			}
		}
		if err = q.UpdateUserAllowIPToken(ctx, arg.UpdateUserAllowIPTokenParams); err != nil {
			return err
		}
		return
	})
}
