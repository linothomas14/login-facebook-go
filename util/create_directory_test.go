package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDirectory(t *testing.T) {
	// Test case 1: Create a single directory
	testDir1 := "testdir1"
	err := CreateDirectory(testDir1)
	defer cleanup(testDir1)
	assert.NoError(t, err, "Error should be nil when creating a single directory")
	assert.DirExists(t, testDir1, "Directory should exist")

	// Test case 2: Create multiple directories
	testDir2 := "testdir2"
	testDir3 := "testdir3"
	err = CreateDirectory(testDir2, testDir3)
	defer cleanup(testDir2, testDir3)
	assert.NoError(t, err, "Error should be nil when creating multiple directories")
	assert.DirExists(t, testDir2, "Directory 2 should exist")
	assert.DirExists(t, testDir3, "Directory 3 should exist")

	// Test case 3: Attempt to create an already existing directory
	err = CreateDirectory(testDir1)
	assert.NoError(t, err, "Error should be nil when attempting to create an existing directory")
}

func cleanup(paths ...string) {
	for _, path := range paths {
		os.RemoveAll(path)
	}
}
