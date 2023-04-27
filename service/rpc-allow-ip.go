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

func (s *TrueAuthService) AllowIPAddress(ctx context.Context, req *rpc.AllowIPAddressRequest) (*rpc.AllowIPAddressResponse, error) {
	authorized, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthenticatedError(err)
	}
	if !s.isUnKnownIP(ctx, authorized.Account) {
		return &rpc.AllowIPAddressResponse{Message: "ip address is already in whitelist"}, nil
	}

	violations := validateAllowIPAddressRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}
	tokenPayload, err := s.tokens.VerifyToken(authorized.Account.AllowIpToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid code: %s", err.Error())
	}

	if tokenPayload.Payload.AccountEmail != authorized.Account.Email {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched email")
	}

	if tokenPayload.Payload.AccountID.String() != authorized.Account.ID.String() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched account id")
	}

	if tokenPayload.Payload.EmailVerificationCode != req.GetCode() {
		return nil, status.Errorf(codes.Internal, "invalid code: mismatched code")
	}

	meta := s.extractMetadata(ctx)
	authorized.Account.AllowedIps = append(authorized.Account.AllowedIps, meta.ClientIp)

	err = s.store.UpdateAccountAllowIP(ctx, sqlc.UpdateAccountAllowIPParams{
		ID:           authorized.Account.ID,
		AllowIpToken: "null",
		AllowedIps:   authorized.Account.AllowedIps,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update allow ip list: %s", err.Error())
	}

	return &rpc.AllowIPAddressResponse{Message: fmt.Sprintf("your ip address %s has been successfully added to whitelist", meta.ClientIp)}, nil
}

func validateAllowIPAddressRequest(req *rpc.AllowIPAddressRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if len(req.GetCode()) != 6 {
		violations = append(violations, fieldViolation("code", fmt.Errorf("invalid code")))
	}
	return
}
