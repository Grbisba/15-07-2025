package production

import (
	"github.com/pkg/errors"
)

var (
	ErrNilConfig    = errors.New("config unexpectedly nil, but it's required")
	ErrNilProvider  = errors.New("provider unexpectedly nil, but it's required")
	ErrNilStore     = errors.New("store unexpectedly nil, but it's required")
	ErrNilCache     = errors.New("cache unexpectedly nil, but it's required")
	ErrNilEncryptor = errors.New("encryptor unexpectedly nil, but it's required")
	ErrNilManager   = errors.New("manager unexpectedly nil, but it's required")
)
