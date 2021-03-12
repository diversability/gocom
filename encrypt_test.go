package gocom

import (
	"fmt"
	"testing"
)

// go test -v encrypt_test.go encrypt.go -test.run Test_Encrypt
// go test -v encrypt_test.go encrypt.go -test.run Test_Decrypt
func Test_Encrypt(t *testing.T) {
	aes := NewAesCbcPKCS7("XGYUZj78QvlvyHQ1eKeSeNhCJcJRQOyQ")
	out, err := aes.Encrypt([]byte("1234567890abcdefg"))
	if err != nil {
		t.Fatalf("fail: %s", err.Error())
	}

	fmt.Println(out)

	out, err = aes.Encrypt([]byte("1234567890abcdefg"))
	if err != nil {
		t.Fatalf("fail: %s", err.Error())
	}

	fmt.Println(out)

	out, err = aes.Encrypt([]byte("1234567890abcdefg"))
	if err != nil {
		t.Fatalf("fail: %s", err.Error())
	}

	fmt.Println(out)
}

func Test_Decrypt(t *testing.T) {
	aes := NewAesCbcPKCS7("XGYUZj78QvlvyHQ1eKeSeNhCJcJRQOyQ")
	str := "35vbMVKu5Lnjdh3m2BMuTBGhCocSH7GNCOalTe+oz3bzf3rtAl99KBui4QS9NxEs"

	out, err := aes.Decrypt(str)
	if err != nil {
		t.Fatalf("fail: %s", err.Error())
	}

	fmt.Println(out)

	str = "9yuf+1y4yEo2CXXODPJ3+faPFNsMqB+BUXY11TD6XxjGbVsGYV/oqPfOKJuYejSs"

	out1, err := aes.Decrypt(str)
	if err != nil {
		t.Fatalf("fail: %s", err.Error())
	}

	fmt.Println(out1)

	str = "TmceU11kLEyS4bbZ7ViTg9yrNcCCNpT9QwlDVcfc+CeMYtmCeuvqEp/lS2dvs6Bu"

	out2, err := aes.Decrypt(str)
	if err != nil {
		t.Fatalf("fail: %s", err.Error())
	}

	fmt.Println(out2)
}
