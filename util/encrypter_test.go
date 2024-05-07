package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestEncryptPassword(t *testing.T) {
	tests := []struct {
		password string
		cost     int
		err      bool
	}{
		{"password123", bcrypt.DefaultCost, false},
		{"hello123", 14, false},
		{"abcde", bcrypt.MinCost, false},
		{"askdnasadjkasjkbfhgweyfgabaesadhusdugfsjhdchxzvcghasvcaschjxbchjxzcbhdbchdsbajhcjzbhdsgfgsdfjbcgeyufsbj", bcrypt.MinCost, true},
	}

	for _, test := range tests {
		t.Run(test.password, func(t *testing.T) {
			hash, err := EncryptPassword(test.password, test.cost)
			if err != nil {
				if test.err {
					assert.Equal(t, bcrypt.ErrPasswordTooLong, err)
					return
				}
				t.Errorf("Error encrypting password: %v", err)
			}

			assert.Nil(t, VerifyPassword(test.password, hash))
		})
	}
}
