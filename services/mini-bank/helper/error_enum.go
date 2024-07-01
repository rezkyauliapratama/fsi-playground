package error_enum

import "github.com/rezkyauliapratama/fsi-playground/libs/error"

var (
	ErrTransactionTypeNotFound error.Error = error.NewError(error.TypeError, "TRANSACTION_TYPE_NOT_FOUND").WithMessage("transaction type is not exist")
)
