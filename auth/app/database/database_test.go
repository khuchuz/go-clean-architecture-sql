package database

import (
	"testing"

	"github.com/khuchuz/go-clean-architecture-sql/auth/app/config"
	"github.com/stretchr/testify/assert"
)

func TestSetupDatabase(t *testing.T) {
	dbreal := SetupDatabase()
	assert.NoError(t, dbreal.Error)

	func() {
		config.DBName = "seharusnya_gaada_sih"
		defer func() {
			if r := recover(); r == nil {
				t.Error("SetupDatabase should have panicked!")
			}
		}()
		// This function should cause a panic
		SetupDatabase()
	}()
}
