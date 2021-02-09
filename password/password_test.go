package password

import (
	"testing"
)

func TestCipher(t *testing.T) {
	cipher := NewCipher(nil)
	password := "passw0rd"
	encodePassword, _ := cipher.EncryptPassword(password)
	err := cipher.ComparePassword(encodePassword, password)
	if err != nil {
		t.Errorf("expected no error")
	}
}
