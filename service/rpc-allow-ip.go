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
	account, ip, err := s.authorize(ctx)
	if err != nil {
		return nil, unAuthenticatedError(err)
	}
	if s.isKnownIP(ip, ctx) {
		return &rpc.AllowIPAddressResponse{Message: "ip address is already in whitelist"}, nil
	}

	violations := validateAllowIPAddressRequest(req)
	if violations != nil {
		return nil, invalidArgumentsError(violations)
	}
	tokenPayload, err := s.tokens.VerifyToken(ip.Token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "invalid allow ip address verification code")
	}

	if tokenPayload.Payload.AccountEmail != account.Email {
		return nil, status.Errorf(codes.Internal, "invalid allow ip address verification code")
	}

	if tokenPayload.Payload.AccountID.String() != account.ID.String() {
		return nil, status.Errorf(codes.Internal, "invalid allow ip address verification code")
	}

	if tokenPayload.Payload.AllowIPCode != req.GetCode() {
		return nil, status.Errorf(codes.Internal, "invalid allow ip address verification code")
	}

	meta := s.extractMetadata(ctx)
	ip.AllowedIps = append(ip.AllowedIps, meta.ClientIp)

	_, err = s.store.UpdateIP(ctx, sqlc.UpdateIPParams{
		ID:         ip.ID,
		AllowedIps: ip.AllowedIps,
		BlockedIps: ip.BlockedIps,
		Token:      "null",
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update ip record: %s", err.Error())
	}
	return &rpc.AllowIPAddressResponse{Message: fmt.Sprintf("your ip address %s has been successfully added to whitelist", meta.ClientIp)}, nil
}

func validateAllowIPAddressRequest(req *rpc.AllowIPAddressRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if len(req.GetCode()) != 6 {
		violations = append(violations, fieldViolation("code", fmt.Errorf("invalid code")))
	}
	return
}
