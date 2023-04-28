package service

import (
	"context"
	"fmt"

	rpc "github.com/sirjager/rpcs/trueauth/go"
	"github.com/sirjager/trueauth/db/sqlc"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TrueAuthService) AllowIP(ctx context.Context, req *rpc.AllowIPRequest) (*rpc.AllowIPResponse, error) {
	authorized, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthenticatedError(err)
	}
	if !s.isUnKnownIP(ctx, authorized.User) {
		return &rpc.AllowIPResponse{Message: "ip address is already in whitelist"}, nil
	}

	violations := validateAllowIPRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}
	tokenPayload, err := s.tokens.VerifyToken(authorized.User.AllowipToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid code: %s", err.Error())
	}

	if tokenPayload.Payload.UserEmail != authorized.User.Email {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched email")
	}

	if tokenPayload.Payload.UserID.String() != authorized.User.ID.String() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched account id")
	}

	if tokenPayload.Payload.AllowIPCode != req.GetCode() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched code")
	}

	meta := s.extractMetadata(ctx)
	authorized.User.AllowedIps = append(authorized.User.AllowedIps, meta.ClientIp)

	err = s.store.Update_User_AllowIP(ctx, sqlc.Update_User_AllowIPParams{
		AllowipToken: "null",
		ID:           authorized.User.ID,
		AllowedIps:   authorized.User.AllowedIps,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update allow ip list: %s", err.Error())
	}

	return &rpc.AllowIPResponse{Message: fmt.Sprintf("your ip address %s has been successfully added to whitelist", meta.ClientIp)}, nil
}

func validateAllowIPRequest(req *rpc.AllowIPRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if len(req.GetCode()) != 6 {
		violations = append(violations, fieldViolation("code", fmt.Errorf("invalid code")))
	}
	return
}
