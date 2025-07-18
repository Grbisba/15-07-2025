package fiber

import (
	"github.com/pkg/errors"
)

var (
	errNilRef = errors.New("provided nil dependency")
)
