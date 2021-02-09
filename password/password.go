package password

import (
	"golang.org/x/crypto/bcrypt"
)

type CostFunc func() int

type Cipher struct {
	cost CostFunc
}

func NewCipher(fn CostFunc) *Cipher {
	if fn == nil {
		fn = func() int {
			return bcrypt.MinCost
		}
	}
	return &Cipher{cost: fn}
}

func (c *Cipher) EncryptPassword(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), c.cost())
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (c *Cipher) ComparePassword(encodePassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(encodePassword), []byte(password))
}
