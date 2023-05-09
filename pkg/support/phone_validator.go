package support

import (
	"errors"
	"regexp"
)

// func MobileNum_validator(number string) bool {
// 	re := regexp.MustCompile(`^[0-9]{10}$`)
// 	if re.MatchString(number) {
// 		return true
// 	} else {
// 		return false
// 	}
// }

func MobileNum_validator(number string) error {
	re := regexp.MustCompile(`^[0-9]{10}$`)
	if re.MatchString(number) {
		return nil
	} else {
		return errors.New("invalid mobile number")
	}
}
