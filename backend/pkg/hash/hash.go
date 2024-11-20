package hash

import (
	"crypto/sha512"
	"fmt"
	"io"
	"log"

	"efaturas-xtreme/pkg/errors"
)

func New(input string) string {
	h := sha512.New()
	_, err := io.WriteString(h, input)
	if err != nil {
		log.Println(errors.New(err))
		panic("failed to write hash")
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func NewUserID(uname string, pword string) string {
	return New("hash:" + uname + "-" + pword)
}
