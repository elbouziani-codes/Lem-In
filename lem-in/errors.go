package main

import (
	"errors"
	"fmt"
)

// Common errors for lem-in.
// The project spec requires printing ONLY "ERROR: invalid data format" for invalid input format,
// but internally we can have more descriptive errors for debugging, then format them later.
var (
	ErrInvalidNumAnts        = errors.New("invalid number of ants")
	ErrMissingStart          = errors.New("missing start room")
	ErrMissingEnd            = errors.New("missing end room")
	ErrDuplicateRoom         = errors.New("duplicate room names")
	ErrDuplicateLink         = errors.New("duplicate links")
	ErrInvalidCoordinates    = errors.New("invalid coordinates")
	ErrLinkToUnknownRoom     = errors.New("link to unknown room")
	ErrRoomLinkToItself      = errors.New("room linking to itself")
	ErrNoPath                = errors.New("no path from start to end")
	ErrInvalidFormat         = errors.New("invalid data format")
)

// PrintError prints the standard error message and exits.
func PrintError() {
	fmt.Println("ERROR: invalid data format")
}
