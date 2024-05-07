package util

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateHexUUID() string {
	id := uuid.New()
	return strings.Replace(id.String(), "-", "", -1)
}

func GenerateReqID() string {
	return uuid.New().String()
}
