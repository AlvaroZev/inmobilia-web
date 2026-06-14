package manufacturing

import "errors"

var (
	ErrNotImplemented            = errors.New("manufacturing compiler not implemented")
	ErrInvalidResolvedFurniture  = errors.New("invalid resolved furniture")
	ErrInvalidManufacturingModel = errors.New("invalid manufacturing model")
)
