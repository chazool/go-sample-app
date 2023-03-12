package mock

import "errors"

func getError(err string) error {
	if len(err) != 0 {
		return errors.New(err)
	}
	return nil
}
