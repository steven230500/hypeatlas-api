package db

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IsDuplicateEntry(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "duplicate key") {
		return true
	}
	return errors.Is(err, gorm.ErrDuplicatedKey)
}

func IsForeignKeyConstraint(err error) bool {
	return errors.Is(err, gorm.ErrForeignKeyViolated)
}

func HandleCommonDBErrors(err error) error {
	// Aquí puedes mapear a mensajes de negocio si quieres.
	return err
}
