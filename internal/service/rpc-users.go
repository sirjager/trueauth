package service

import (
	"context"
	"fmt"

	"github.com/sirjager/trueauth/pkg/validator"
	rpc "github.com/sirjager/trueauth/stubs/go"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoreService) User(ctx context.Context, req *rpc.UserRequest) (*rpc.UserResponse, error) {
	_, violations := validateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}
	authorized, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthenticatedError(err)
	}

	// authorized user can only access their own profile
	if string(authorized.user.ID) == req.GetIdentity() || authorized.user.Email == req.GetIdentity() || authorized.user.Username == req.GetIdentity() {
		return &rpc.UserResponse{User: publicProfile(authorized.user)}, nil
	}

	//TODO: if authorized user have permissions like admin then we will return requested user.
	//* For now we will just return user not found
	return nil, status.Errorf(codes.NotFound, "user not found")
}

func validateUserRequest(req *rpc.UserRequest) (identity string, violations []*errdetails.BadRequest_FieldViolation) {
	identity = ""
	if err := validator.ValidateUsername(req.GetIdentity()); err == nil {
		identity = "username"
	}
	if err := validator.ValidateEmail(req.GetIdentity()); err == nil {
		identity = "email"
	}
	if err := validator.ValidateUUID(req.GetIdentity()); err == nil {
		identity = "id"
	}
	if identity == "" {
		violations = append(violations, fieldViolation("identity", fmt.Errorf(err_invalid_identity)))
	}
	return identity, violations
}
