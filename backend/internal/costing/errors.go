package costing

import "errors"

var (
	ErrNotImplemented           = errors.New("cost engine not implemented")
	ErrInvalidManufacturingModel = errors.New("invalid manufacturing model")
	ErrInvalidCostResult        = errors.New("invalid cost result")
)
