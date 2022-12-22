package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	mathRand "math/rand"
	"net/url"
	"strings"
	"time"
	"util_api/util/log"
)

var (
	//key 需要 16 或 32
	key = ""
)

func KeyEncrypt(cryptoText string) (string, error) {
	keyBytes := sha256.Sum256([]byte(key))
	return encrypt(keyBytes[:], cryptoText)
}

// encrypt string to base64 crypto using AES
func encrypt(key []byte, text string) (string, error) {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error(err)
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	padding := 0
	result := len(plaintext) % aes.BlockSize
	if result != 0 {
		padding = aes.BlockSize - result
		for i := 0; i < padding; i++ {
			plaintext = append(plaintext, 0)
		}
	}

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Error(err)
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func KeyDecrypt(cryptoText string) (string, error) {
	keyBytes := sha256.Sum256([]byte(key))
	return decrypt(keyBytes[:], cryptoText)
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		log.Error(err)
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error(err)
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(cipherText) < aes.BlockSize {
		log.Error(err)
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)

	for i := len(cipherText) - 1; i >= 0; i-- {
		if cipherText[i] != 0 {
			cipherText = cipherText[:i+1]
			break
		}
	}

	return string(cipherText), nil
}

//生成随机字符串
func GetRandomString(lenght int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lenght; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func pkcs7Padding(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

func Aes256CBCDecryptTypeI(dataByte string, keyStr string, ivStr string) (string, error) {
	/*
	 * 藍新金流用
	 */

	keyByte := []byte(keyStr)
	ivByte := []byte(ivStr)

	hex2binByte, _ := hex.DecodeString(dataByte)

	block, err := aes.NewCipher(keyByte)
	log.Error(err)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, len(hex2binByte))
	mode := cipher.NewCBCDecrypter(block, ivByte)
	mode.CryptBlocks(cipherText, hex2binByte)

	return string(pkcs7UnPadding(cipherText)), nil
}

func HttpParseQuery(paramsStr string) (paramsMap map[string]interface{}) {
	/*
	 * 解析 HttpQueryString
	 */

	if params, err := url.ParseQuery(paramsStr); err != nil {
		log.Error(err)
		return
	} else {
		paramsMap = make(map[string]interface{}, 0)
		for key, values := range params {
			paramsMap[key] = strings.Join(values, "")
		}
		return paramsMap
	}
}
