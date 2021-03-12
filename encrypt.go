package gocom

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type AesCbcPKCS7 struct {
	AesObj cipher.Block
}

func NewAesCbcPKCS7(key string) *AesCbcPKCS7 {
	aesObj, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Printf("gen aes err: %s\n", err.Error())
		return nil
	}

	return &AesCbcPKCS7{AesObj: aesObj}
}

func (obj *AesCbcPKCS7) Encrypt(rawData []byte) (string, error) {
	data, err := obj.AesCBCEncrypt(rawData)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (obj *AesCbcPKCS7) Decrypt(rawData string) (string, error) {
	if rawData == "" {
		return "", fmt.Errorf("no data")
	}

	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return "", err
	}

	dnData, err := obj.AesCBCDecrypt(data)
	if err != nil {
		return "", err
	}

	return string(dnData), nil
}

func (obj *AesCbcPKCS7) AesCBCEncrypt(rawData []byte) ([]byte, error) {
	blockSize := obj.AesObj.BlockSize()

	//填充原文
	rawData = PKCS7Padding(rawData, blockSize)

	//初始向量IV必须是唯一，但不需要保密
	cipherText := make([]byte, blockSize+len(rawData))

	//block大小
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(obj.AesObj, iv)
	mode.CryptBlocks(cipherText[blockSize:], rawData)

	return cipherText, nil
}

func (obj *AesCbcPKCS7) AesCBCDecrypt(encryptData []byte) ([]byte, error) {
	blockSize := obj.AesObj.BlockSize()

	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}

	iv := encryptData[:blockSize]
	encryptData = encryptData[blockSize:]

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(obj.AesObj, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)

	//解填充
	encryptData = PKCS7UnPadding(encryptData)
	return encryptData, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
