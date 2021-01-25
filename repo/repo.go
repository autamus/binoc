package repo

import (
	"github.com/autamus/go-parspack/pkg"
)

// Result is a reported package and its
// parsed location from the parsed library.
type Result struct {
	Data pkg.Package
	Path string
}
