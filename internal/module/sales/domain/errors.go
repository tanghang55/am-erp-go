package domain

import "errors"

var (
	ErrOrderNotFound       = errors.New("sales order not found")
	ErrInvalidTransition   = errors.New("invalid order status transition")
	ErrInvalidQuantity     = errors.New("invalid quantity")
	ErrItemNotFound        = errors.New("sales order item not found")
	ErrImportDuplicateFile = errors.New("import file already exists")
	ErrImportNotFound      = errors.New("import batch not found")
	ErrImportInvalidFile   = errors.New("invalid import file")
	ErrProductNotFound     = errors.New("product not found")
)
