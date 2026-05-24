package plugins

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ObjectIDOrNil parses a hex string into a primitive.ObjectID.
// Returns primitive.NilObjectID if the string is empty or invalid.
func ObjectIDOrNil(id string) primitive.ObjectID {
	if strings.TrimSpace(id) == "" {
		return primitive.NilObjectID
	}
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID
	}
	return oid
}

// IsEmptyStringOrNil returns true if the pointer is nil or points to an empty/whitespace string.
func IsEmptyStringOrNil(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}
