package valueobject

import (
	"errors"
	"regexp"
)

var ErrInvalidEmail = errors.New("invalid email format")

type Email struct {
	address string
}

func NewEmail(address string) (Email, error) {
	if !isValidEmail(address) {
		return Email{}, ErrInvalidEmail
	}
	return Email{address: address}, nil
}

func (e Email) String() string {
	return e.address
}

func (e Email) Equals(other Email) bool {
	return e.address == other.address
}

func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
