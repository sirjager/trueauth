package server

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/hibiken/asynq"
	"google.golang.org/grpc/status"

	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/pkg/utils"
	rpc "github.com/sirjager/trueauth/rpc"
	"github.com/sirjager/trueauth/worker"
)

const (
	deleteCodeDigitsCount = 10
)

func (s *Server) Delete(
	ctx context.Context,
	req *rpc.DeleteRequest,
) (*rpc.DeleteResponse, error) {
	auth, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthorizedError(err)
	}

	meta := s.extractMetadata(ctx)

	// NOTE: If the code is not provided, that means user does not have deletion code,
	// we will create a new one and send it to user via email
	if len(req.GetCode()) == 0 {

		// check if user has requested deletion code recently, if yes, then return error
		if time.Since(auth.user.LastUserDeletion) < s.config.Auth.DeleteTokenCooldown {
			tryAfter := s.config.Auth.DeleteTokenCooldown - time.Since(
				auth.user.LastUserDeletion,
			)
			errMessage := "account deletion has been requested recently, please try again after %s"
			return nil, status.Errorf(_aborted, errMessage, tryAfter)
		}

		// generate a new random deletion code
		code := utils.RandomNumberAsString(deleteCodeDigitsCount)
		// later we will check if token is valid
		params := tokens.PayloadData{
			Code:      code,
			UserID:    auth.user.ID,
			UserEmail: auth.user.Email,
			ClientIP:  meta.clientIP,
			UserAgent: meta.userAgent,
		}
		token, _, tokenErr := s.tokens.CreateToken(params, s.config.Auth.DeleteTokenExpDur)
		if tokenErr != nil {
			return nil, status.Errorf(_internal, "failed to create token, %s", tokenErr.Error())
		}

		// generate task params and options
		taskParams := worker.PayloadUserDeletionCode{Token: token}
		randomDelay := time.Millisecond * time.Duration(utils.RandomInt(100, 600))
		taskOptions := []asynq.Option{
			asynq.MaxRetry(5),            // max retries if any error occurs
			asynq.Group(worker.QueueLow), // queue task in low priority
			asynq.ProcessIn(randomDelay), // random delay before processing
		}

		// now we will distribute task send email deletion code
		if err = s.tasks.DistributeTaskSendUserDeletionCode(ctx, taskParams, taskOptions...); err != nil {
			errMsg := "failed to initiate account deletion, %s"
			return nil, status.Errorf(_internal, errMsg, err.Error())
		}

		// after task is successfully distributed then we will save deletion token in database
		// so that mutltiple requests can be minimized
		updateParam := db.UpdateUserDeletionTokenParams{
			ID:                auth.user.ID,
			TokenUserDeletion: token,
			LastUserDeletion:  time.Now(),
		}
		if uErr := s.store.UpdateUserDeletionToken(ctx, updateParam); uErr != nil {
			return nil, status.Errorf(_internal, uErr.Error())
		}

		// lastly we ask user to check their email, and hit same api route with code
		return &rpc.DeleteResponse{
			Message: "check your inbox for further instructions",
		}, nil
	}

	// NOTE: If the code is not empty we will validate and continue with deletion
	//
	// this will validate if token is invalid or expired  and what not...
	tokenPayoad, err := s.tokens.VerifyToken(auth.user.TokenUserDeletion)
	if err != nil {
		return nil, status.Errorf(_unauthenticated, err.Error())
	}
	if tokenPayoad.Payload.Code != req.GetCode() {
		return nil, status.Errorf(_unauthenticated, "invalid code")
	}
	if tokenPayoad.Payload.UserEmail != auth.user.Email {
		return nil, status.Errorf(_unauthenticated, "invalid code")
	}
	if !bytes.Equal(tokenPayoad.Payload.UserID, auth.user.ID) {
		return nil, status.Errorf(_unauthenticated, "invalid code")
	}

	// following 2 checks are optional, but it makes more secure
	// there is no need to enforce same ip and useragent
	if tokenPayoad.Payload.UserAgent != meta.userAgent {
		return nil, status.Errorf(_unauthenticated, "invalid code")
	}
	if tokenPayoad.Payload.ClientIP != meta.clientIP {
		return nil, status.Errorf(_unauthenticated, "invalid code")
	}

	// delete user document
	if err = s.store.DeleteUser(ctx, auth.user.ID); err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	if err = s.cache.DeleteWithPrefix(ctx, userSessionsKey(auth.user.ID)); err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	// clear cookies
	if err = s.sendCookies(ctx, []http.Cookie{
		{Name: "sessionId", Value: "", Path: "/", Expires: time.Now(), HttpOnly: true},
		{Name: "accessToken", Value: "", Path: "/", Expires: time.Now(), HttpOnly: true},
		{Name: "refreshToken", Value: "", Path: "/", Expires: time.Now(), HttpOnly: true},
	}); err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	return &rpc.DeleteResponse{Message: "account deleted successfully"}, nil
}
