package postgres

import "github.com/rotisserie/eris"

var ErrCouldNotAcquireLock = eris.New("could not acquire lock")
