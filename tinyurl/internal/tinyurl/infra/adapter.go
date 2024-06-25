package sql

import (
	"database/sql"

	tinyError "github.com/christapa/tinyurl/pkg/error"
	"github.com/lib/pq"
)

func sqlToDomainError(err error) error {
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return tinyError.New(tinyError.NotFound, "not found")
	}

	// PRint type
	pqError, ok := err.(*pq.Error)
	if !ok {
		return tinyError.New(tinyError.Internal, err.Error())
	}

	// TODO Need to handle more error sql -> pkg error
	switch {
	case pqError.Code == "23505":
		return tinyError.New(tinyError.AlreadyExists, err.Error())
	default:
		return tinyError.New(tinyError.Internal, err.Error())
	}

}
