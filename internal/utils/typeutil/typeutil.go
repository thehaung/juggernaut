package typeutil

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"strings"
	"sync"
)

var (
	syncOnce      sync.Once
	validatorUtil *validator.Validate
)

func GetValidator() *validator.Validate {
	syncOnce.Do(func() {
		validatorUtil = validator.New()
	})

	return validatorUtil
}

func UUIDWithLimit(length int16) string {
	uuidStr := uuid.NewString()[:length]
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	return uuidStr
}
