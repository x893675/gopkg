package password

import (
	"github.com/x893675/gopkg/stringplus"
	"golang.org/x/crypto/bcrypt"
)

var cost int

func init() {
	cost = bcrypt.DefaultCost
}

func SetCost(c int) {
	if c < bcrypt.MinCost || c > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	} else {
		cost = c
	}
}

func GeneratePassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(stringplus.String2bytes(pw), cost)
	if err != nil {
		return "", err
	}
	return stringplus.Bytes2string(hash), nil
}

func ComparePassword(encodePw, password string) error {
	return bcrypt.CompareHashAndPassword(stringplus.String2bytes(encodePw), stringplus.String2bytes(password))
}
