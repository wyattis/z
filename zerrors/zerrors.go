package zerrors

// Returns the first non-nil error in an slice of errors
func AnyError(errors ...error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}
