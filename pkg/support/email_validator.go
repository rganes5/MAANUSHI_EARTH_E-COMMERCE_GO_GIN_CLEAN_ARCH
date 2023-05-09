package support

import "regexp"

func Email_validator(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if re.MatchString(email) {
		return true
	} else {
		return false
	}
}
