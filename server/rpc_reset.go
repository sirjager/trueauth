package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/sirjager/gopkg/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"

	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/pkg/hash"
	"github.com/sirjager/trueauth/pkg/tokens"
	rpc "github.com/sirjager/trueauth/rpc"
	"github.com/sirjager/trueauth/worker"
)

const passwordResetCodeDigitsCount = 7

func (s *Server) Reset(
	ctx context.Context,
	req *rpc.ResetRequest,
) (*rpc.ResetResponse, error) {
	//  this will validate email new pasword and everything
	if violations := validateResetRequest(req); violations != nil {
		return nil, invalidArgumentsError(violations)
	}

	// first we fetch user by email
	user, err := s.store.ReadUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			// we don't want to disclose if user exists or not, so we simply return email sent.
			return &rpc.ResetResponse{Message: "you will receive an email shortly"}, nil
		}
		// if something else goes wrong, we return error
		return nil, fmt.Errorf("something went wrong: %w", err)
	}

	meta := s.extractMetadata(ctx)

	// If code is not provided, that means user does not have password reset code,
	// we will create a new one and send it to user via email
	if req.GetCode() == "" {
		// check if user has requested deletion code recently, if yes, then return error
		if time.Since(user.LastPasswordReset) < s.config.Auth.ResetTokenCooldown {
			tryAfter := s.config.Auth.ResetTokenCooldown - time.Since(
				user.LastPasswordReset,
			)
			errMessage := "already requested, please try again after %s"
			return nil, status.Errorf(_aborted, errMessage, tryAfter)
		}

		// generate a new random deletion code
		code := utils.RandomNumberAsString(passwordResetCodeDigitsCount)
		// later we will check if token is valid
		params := tokens.PayloadData{
			Code:      code,
			UserID:    user.ID,
			UserEmail: user.Email,
			ClientIP:  meta.clientIP(),
			UserAgent: meta.userAgent(),
		}
		token, _, tokenErr := s.tokens.CreateToken(params, s.config.Auth.ResetTokenExpDur)
		if tokenErr != nil {
			return nil, status.Errorf(_internal, "failed to create token, %s", tokenErr.Error())
		}

		// generate task params and options
		taskParams := worker.PayloadPasswordResetCode{Token: token}
		randomDelay := time.Millisecond * time.Duration(utils.RandomInt(100, 600))
		taskOptions := []asynq.Option{
			asynq.MaxRetry(5),            // max retries if any error occurs
			asynq.Group(worker.QueueLow), // queue task in low priority
			asynq.ProcessIn(randomDelay), // random delay before processing
		}

		// now we will distribute task send email deletion code
		if err = s.tasks.DistributeTaskSendPasswordResetCode(ctx, taskParams, taskOptions...); err != nil {
			errMsg := "failed to initiate password reset process, %s"
			return nil, status.Errorf(_internal, errMsg, err.Error())
		}

		updateParams := db.UpdateUserPasswordResetTokenParams{
			ID:                 user.ID,
			TokenPasswordReset: token,
			LastPasswordReset:  time.Now(),
		}
		if err = s.store.UpdateUserPasswordResetToken(ctx, updateParams); err != nil {
			return nil, status.Errorf(_internal, err.Error())
		}

		return &rpc.ResetResponse{Message: "check your inbox for further instructions"}, nil
	}

	// If code is not empty, we will check if it is valid and process it
	//
	// this will validate if token is invalid or expired  and what not...
	tokenPayoad, err := s.tokens.VerifyToken(user.TokenPasswordReset)
	if err != nil {
		return nil, status.Errorf(_unauthenticated, err.Error())
	}

	// we can also return a normal errors like: "invalid code" instead of detailed error
	if tokenPayoad.Payload.Code != req.GetCode() {
		return nil, status.Errorf(_unauthenticated, "invalid code, code does not match")
	}
	if tokenPayoad.Payload.UserEmail != user.Email {
		return nil, status.Errorf(_unauthenticated, "invalid code, user email does not match")
	}
	if !bytes.Equal(tokenPayoad.Payload.UserID, user.ID) {
		return nil, status.Errorf(_unauthenticated, "invalid code, user id does not match")
	}

	// following 2 checks are optional, there is no need to enforce same ip and useragent
	// but it makes whole process more secure
	if tokenPayoad.Payload.UserAgent != meta.userAgent() {
		return nil, status.Errorf(_unauthenticated, "invalid code, user agent does not match")
	}
	if tokenPayoad.Payload.ClientIP != meta.clientIP() {
		return nil, status.Errorf(_unauthenticated, "invalid code, client ip does not match")
	}

	hashingSalt := hash.RandomSalt(len(req.GetNewPassword()))
	hashedPassword, err := s.hasher.Hash(hashingSalt, req.GetNewPassword())
	if err != nil {
		return nil, status.Errorf(_internal, "failed to hash password: %s", err.Error())
	}

	if err = s.store.UpdateUserPasswordTx(ctx, db.UpdateUserPasswordTx{
		UpdateUserPasswordParams: db.UpdateUserPasswordParams{
			ID:       user.ID,
			HashSalt: hashingSalt,
			HashPass: hashedPassword,
		},
	}); err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	return &rpc.ResetResponse{Message: "successfully updated password"}, nil
}

func validateResetRequest(
	req *rpc.ResetRequest,
) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if req.GetCode() != "" {
		if len(req.GetCode()) != passwordResetCodeDigitsCount {
			violations = append(violations, fieldViolation("code", fmt.Errorf("invalid code")))
		}
		if err := validatePassword(req.GetNewPassword()); err != nil {
			violations = append(violations, fieldViolation("newPassword", err))
		}
	}

	return
}
