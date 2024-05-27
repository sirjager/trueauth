package server

import (
	"bytes"
	"context"
	"time"

	"github.com/hibiken/asynq"
	"google.golang.org/grpc/status"

	"github.com/sirjager/trueauth/db/db"
	"github.com/sirjager/trueauth/pkg/tokens"
	"github.com/sirjager/trueauth/pkg/utils"
	rpc "github.com/sirjager/trueauth/stubs"
	"github.com/sirjager/trueauth/worker"
)

const (
	deleteCodeDigitsCount = 10
)

func (s *Server) Delete(ctx context.Context, req *rpc.DeleteRequest) (*rpc.DeleteResponse, error) {
	authorized, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthorizedError(err)
	}

	// NOTE: If the code is not provided, that means user does not have deletion code,
	// we will create a new one and send it to user via email
	if len(req.GetCode()) == 0 {

		// check if user has requested deletion code recently, if yes, then return error
		if time.Since(authorized.User.LastUserDeletion) < s.config.Auth.DeleteTokenCooldown {
			tryAfter := s.config.Auth.DeleteTokenCooldown - time.Since(
				authorized.User.LastUserDeletion,
			)
			errMessage := "account deletion has been requested recently, please try again after %s"
			return nil, status.Errorf(_aborted, errMessage, tryAfter)
		}

		// generate a new random deletion code
		code := utils.RandomNumberAsString(deleteCodeDigitsCount)
		// later we will check if token is valid
		params := tokens.PayloadData{
			Code: code, Type: "deletion",
			UserID:    authorized.User.ID,
			UserEmail: authorized.User.Email,
			ClientIP:  authorized.Meta.ClientIP,
			UserAgent: authorized.Meta.UserAgent,
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
			ID:                authorized.User.ID,
			TokenUserDeletion: token,
			LastUserDeletion:  time.Now(),
		}
		if uErr := s.store.UpdateUserDeletionToken(ctx, updateParam); uErr != nil {
			return nil, status.Errorf(_internal, uErr.Error())
		}

		// lastly we ask user to check their email, and hit same api route with code
		return &rpc.DeleteResponse{Message: "check your inbox for further instructions"}, nil
	}

	// NOTE: If the code is not empty we will validate and continue with deletion
	//
	// this will validate if token is invalid or expired  and what not...
	tokenPayoad, err := s.tokens.VerifyToken(authorized.User.TokenUserDeletion)
	if err != nil {
		return nil, status.Errorf(_unauthenticated, err.Error())
	}
	if tokenPayoad.Payload.Code != req.GetCode() {
		return nil, status.Errorf(_unauthenticated, "invalid deletion code")
	}
	if tokenPayoad.Payload.UserEmail != authorized.User.Email {
		return nil, status.Errorf(_unauthenticated, "invalid deletion code")
	}
	if !bytes.Equal(tokenPayoad.Payload.UserID, authorized.User.ID) {
		return nil, status.Errorf(_unauthenticated, "invalid deletion code")
	}

	// following 2 checks are optional, but it makes more secure
	// there is no need to enforce same ip and useragent
	if tokenPayoad.Payload.UserAgent != authorized.Meta.UserAgent {
		return nil, status.Errorf(_unauthenticated, "invalid deletion code")
	}
	if tokenPayoad.Payload.ClientIP != authorized.Meta.ClientIP {
		return nil, status.Errorf(_unauthenticated, "invalid deletion code")
	}

	// NOTE: now we can delete all the user related data
	//
	// delete all user's sessions
	if err = s.store.DeleteSessionByUserID(ctx, authorized.User.ID); err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}
	// delete user document
	if err = s.store.DeleteUser(ctx, authorized.User.ID); err != nil {
		return nil, status.Errorf(_internal, err.Error())
	}

	return &rpc.DeleteResponse{Message: "account deleted successfully"}, nil
}
