package utils

import (
	"github.com/rs/xid"
)

func XIDNew() xid.ID {
	return xid.New()
}

func XIDFromString(value string) (xid.ID, error) {
	return xid.FromString(value)
}

func XIDFromBytes(value []byte) (xid.ID, error) {
	return xid.FromBytes(value)
}

