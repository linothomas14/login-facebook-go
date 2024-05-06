package util

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateHexUUID(t *testing.T) {
	// Generate UUID
	result := GenerateHexUUID()

	// Check if the generated string is not empty
	assert.NotEmpty(t, result, "Generated UUID should not be empty")

	// Check if the generated string is a valid UUID (hex)
	_, err := uuid.Parse(strings.Replace(result, "-", "", -1))
	assert.NoError(t, err, "Generated UUID should be a valid UUID")
}

func TestGenerateReqID(t *testing.T) {
	// Generate Request ID
	result := GenerateReqID()

	// Check if the generated string is not empty
	assert.NotEmpty(t, result, "Generated Request ID should not be empty")

	// Check if the generated string is a valid UUID
	_, err := uuid.Parse(result)
	assert.NoError(t, err, "Generated Request ID should be a valid UUID")
}
