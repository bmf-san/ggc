package config

import "errors"

// WarningError wraps a non-fatal configuration error. Callers can use
// errors.As to detect a warning versus a fatal error and decide whether to
// continue startup with default values or to abort.
type WarningError struct {
	Err error
}

// Error implements the error interface.
func (w *WarningError) Error() string {
	if w == nil || w.Err == nil {
		return "config warning"
	}
	return w.Err.Error()
}

// Unwrap exposes the underlying error for errors.Is/As.
func (w *WarningError) Unwrap() error {
	if w == nil {
		return nil
	}
	return w.Err
}

// IsWarning reports whether err (or any error it wraps) is a non-fatal
// config warning.
func IsWarning(err error) bool {
	if err == nil {
		return false
	}
	var w *WarningError
	return errors.As(err, &w)
}
