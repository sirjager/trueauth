package server

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

func unAuthorizedError(err error) error {
	return status.Errorf(_unauthenticated, err.Error())
}

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func invalidArgumentsError(violations []*errdetails.BadRequest_FieldViolation) error {
	invalidReq := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(_invalidArgument, "invalid parameters")
	statusDetails, err := statusInvalid.WithDetails(invalidReq)
	if err != nil {
		return statusInvalid.Err()
	}
	return statusDetails.Err()
}
