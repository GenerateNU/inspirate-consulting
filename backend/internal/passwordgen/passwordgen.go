package passwordgen

import (
	"crypto/rand"
	"math/big"
)

const (
	tempPasswordLength = 16
	lowerChars         = "abcdefghijklmnopqrstuvwxyz"
	upperChars         = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars         = "0123456789"
	specialChars       = "!@#~$%^&*()+|_.,;<>?/{}-"
)

func GenerateTempPassword() (string, error) {
	classes := []string{lowerChars, upperChars, digitChars, specialChars}
	all := lowerChars + upperChars + digitChars + specialChars

	password := make([]byte, 0, tempPasswordLength)
	for _, class := range classes {
		c, err := randomChar(class)
		if err != nil {
			return "", err
		}
		password = append(password, c)
	}
	for len(password) < tempPasswordLength {
		c, err := randomChar(all)
		if err != nil {
			return "", err
		}
		password = append(password, c)
	}

	for num := len(password) - 1; num > 0; num-- {
		numTwo, err := randomInt(num + 1)
		if err != nil {
			return "", err
		}
		password[num], password[numTwo] = password[numTwo], password[num]
	}

	return string(password), nil
}

func randomChar(chars string) (byte, error) {
	num, err := randomInt(len(chars))
	if err != nil {
		return 0, err
	}
	return chars[num], nil
}

func randomInt(n int) (int, error) {
	v, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0, err
	}
	return int(v.Int64()), nil
}
