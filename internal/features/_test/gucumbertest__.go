
package main

import (
	"github.com/gucumber/gucumber"
	_i0 "github.com/ernestio/definition-mapper/internal/features/aws"
	
)

var (
	_ci0 = _i0.IMPORT_MARKER
	
)

func main() {
	
	gucumber.GlobalContext.RunDir("internal/features")
}
