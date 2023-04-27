package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirjager/trueauth/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rpc "github.com/sirjager/rpcs/trueauth/go"
)

func (s *TrueAuthService) Logout(ctx context.Context, req *rpc.LogoutRequest) (*rpc.LogoutResponse, error) {
	if violations := validateLogoutRequest(req); violations != nil {
		return nil, invalidArgumentsError(violations)
	}

	authorized, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthenticatedError(err)
	}

	if len(req.GetSessions()) > 0 {
		for _, sess := range req.GetSessions() {
			sessionID, err := uuid.Parse(sess)
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error())
			}
			err = s.store.DeleteSession(ctx, sessionID)
			if err != nil {
				return nil, status.Errorf(codes.Internal, err.Error())
			}
		}
	} else {
		err = s.store.DeleteSessionByAccount(ctx, authorized.Account.ID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}

	return &rpc.LogoutResponse{Message: "sessions deleted"}, nil
}

func validateLogoutRequest(req *rpc.LogoutRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	for index, id := range req.GetSessions() {
		if err := validator.ValidateUUID(id); err != nil {
			violations = append(violations, fieldViolation(fmt.Sprintf("index %d", index+1), err))
		}
	}
	return
}
