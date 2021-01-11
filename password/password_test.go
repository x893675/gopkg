package password

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

var pwd = "ThinkBig1!"

func test() {
	_, _ = GeneratePassword1(pwd)
}

func test2() {
	_, _ = GeneratePassword(pwd)
}

func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test()
	}
}

func BenchmarkTestBlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test2()
	}
}

func GeneratePassword1(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
