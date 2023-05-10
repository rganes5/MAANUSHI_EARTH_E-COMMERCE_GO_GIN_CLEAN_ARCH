package support

import (
	"errors"
	"regexp"
)

func Email_validator(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if re.MatchString(email) {
		return nil
	} else {
		return errors.New("invalid email")
	}
}
