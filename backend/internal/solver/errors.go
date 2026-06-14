package solver

import "errors"

var (
	ErrNotImplemented         = errors.New("constraint solver not implemented")
	ErrInvalidFurniture       = errors.New("invalid furniture definition")
	ErrInvalidRoom            = errors.New("invalid room geometry")
	ErrInstallZoneUndefined   = errors.New("installation zone is undefined")
	ErrReferenceWallNotFound  = errors.New("reference wall not found")
	ErrInstallSpaceTooSmall   = errors.New("installation space too small after clearances")
	ErrDimensionExceedsSpace  = errors.New("dimension exceeds available space")
	ErrChildrenWithoutSplit   = errors.New("volume node has children without split")
	ErrSplitChildMismatch     = errors.New("split definition does not match children")
	ErrSplitUndefined         = errors.New("split has no ratios or fixed sizes")
	ErrInvalidSplitAxis       = errors.New("invalid split axis")
	ErrUnknownConstraintMode  = errors.New("unknown constraint mode")
)
