package service

import (
	"context"

	rpc "github.com/sirjager/rpcs/trueauth/go"
)

func (s *TrueAuthService) DeleteAccount(ctx context.Context, req *rpc.DeleteAccountRequest) (*rpc.DeleteAccountResponse, error) {

	return &rpc.DeleteAccountResponse{}, nil
}
