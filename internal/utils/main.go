package utils

import "github.com/google/uuid"

func IsValidUUID(u string) bool {
	err := uuid.Validate(u)
	return err == nil
}
