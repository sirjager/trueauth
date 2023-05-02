package service

import (
	"context"
	"fmt"

	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/pkg/validator"
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
	if authorized.User.ID.String() == req.GetIdentity() || authorized.User.Email == req.GetIdentity() || authorized.User.Username == req.GetIdentity() {
		return &rpc.UserResponse{User: publicProfile(authorized.User)}, nil
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
