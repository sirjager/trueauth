package server

import (
	"google.golang.org/grpc/codes"
)

// all status codes are defined here
// so that we donn't have to import google.golang.org/grpc/codes
// everywhere and use it as codes.Internal etc.
// now we can use directly by _internal etc
const (
	_internal           = codes.Internal
	_permissionDenied   = codes.PermissionDenied
	_invalidArgument    = codes.InvalidArgument
	_unauthenticated    = codes.Unauthenticated
	_conflict           = codes.AlreadyExists
	_notFound           = codes.NotFound
	_failedPreCondition = codes.FailedPrecondition
	_aborted            = codes.Aborted
	_deadlineExceeded   = codes.DeadlineExceeded
	_unimplemented      = codes.Unimplemented
	_unavailable        = codes.Unavailable
	_dataLoss           = codes.DataLoss
	_canceled           = codes.Canceled
	_unknown            = codes.Unknown
	_resourceExhausted  = codes.ResourceExhausted
	_outOfRange         = codes.OutOfRange
)

const (
	errUserDoesNotExist     = "user does not exists: %s"
	errUnauthorized         = "unauthorized, invalid or expired token"
	errMissingAuthorization = "unauthorized, missing or invalid token"
	errFailedToRetrieveUser = "failed to retrieve user: %s"
	errEmailNotVerified     = "email not verified, please verify your email first"
	errEmailNotRegistered   = "email not registered"
)
