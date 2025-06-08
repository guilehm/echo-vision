//go:build !prod
// +build !prod

package postgres

// EncodeCursor encodes a cursor for tests only, should NOT be used in the app.
var EncodeCursorForTest = encodeCursor
