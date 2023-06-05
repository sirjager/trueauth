package utils

import (
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
)

// Features
//  1. Size: 12 bytes (96 bits), smaller than UUID, larger than snowflake
//  2. Base32 hex encoded by default (20 chars when transported as printable string, still sortable)
//  3. Non configured, you don't need set a unique machine and/or data center id
//  4. K-ordered
//  5. Embedded time with 1 second precision
//  6. Unicity guaranteed for 16,777,216 (24 bits) unique ids per second and per host/process
//  7. Lock-free (i.e.: unlike UUIDv1 and v2)
func UUID_XID() string {
	return xid.New().String()
}

// There are numerous methods for generating unique identifiers, so why KSUID?
//  1. Naturally ordered by generation time
//  2. Collision-free, coordination-free, dependency-free
//  3. Highly portable representations
//  4. The text representation is always 27 characters, encoded in alphanumeric base62 that will lexicographically sort by timestamp.
func UUID_KSUID() string {
	return ksuid.New().String()
}

// Universally Unique Lexicographically Sortable Identifier
//  1. Is compatible with UUID/GUID's
//  2. 1.21e+24 unique ULIDs per millisecond (1,208,925,819,614,629,174,706,176 to be exact)
//  3. Lexicographically sortable
//  4. Canonically encoded as a 26 character string, as opposed to the 36 character UUID
//  5. Uses Crockford's base32 for better efficiency and readability (5 bits per character)
//  6. Case insensitive
//  7. No special characters (URL safe)
//  8. Monotonic sort order (correctly detects and handles the same millisecond)
func UUID_ULID() string {
	return ulid.Make().String()
}

func UUID_GOOGLE() string {
	return uuid.NewString()
}
