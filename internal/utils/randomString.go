package utils

import (
        "crypto/rand"
        "fmt"
        "math"
        "math/big"
)

type CustomError struct {
        Message string
}

func (e *CustomError) Error() string {
        return e.Message
}

// RandomString generates a string with the given length using crypt/rand
func RandomString(length int) (string, error) {
        if length < 0 {
                err := &CustomError{Message: "Length has to be bigger than 0"}
                return "", err
        }
        if length > math.MaxInt32 {
                err := &CustomError{Message: fmt.Sprintf("Length has to be bigger than %d", math.MaxInt32)}
                return "", err
        }
        const characters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
        byteSlice := make([]byte, length)
        for i := 0; i < length; i++ {
                randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
                if err != nil {
                        return "", err
                }
                byteSlice[i] = characters[randomInt.Int64()]
        }
        return string(byteSlice), nil
}
