package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}

func NBytes(base64string string)(int, error){
	b, err := base64.URLEncoding.DecodeString(base64string)
	if err != nil{
		return -1, nil
	}
	return len(b), nil
}
