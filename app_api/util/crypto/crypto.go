package crypto

import (
	"app_api/util/log"
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
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.

	if len(cipherText) < aes.BlockSize {
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

func AesCBC256Encrpty(keyStr string, textByte []byte) (string, error) {
	key := []byte(keyStr)
	plaintext := pkcs7Padding(textByte)

	block, err := aes.NewCipher(key)
	log.Error(err)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, len(plaintext))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Error(err)
		return "", err
	}

	stream := cipher.NewCBCEncrypter(block, iv)

	var newCipherText []uint8

	for i := 0; i < len(iv); i++ {
		newCipherText = append(newCipherText, iv[i])
	}

	stream.CryptBlocks(cipherText, plaintext)

	for i := 0; i < len(cipherText); i++ {
		newCipherText = append(newCipherText, cipherText[i])
	}

	return base64.StdEncoding.EncodeToString(newCipherText), nil
}

func AesCBC256Decrpty(keyStr string, textByte []byte) (string, error) {
	key := []byte(keyStr)
	plaintext := pkcs7Padding(textByte)

	block, err := aes.NewCipher(key)
	log.Error(err)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, len(plaintext))
	iv := cipherText[:aes.BlockSize]

	stream := cipher.NewCBCDecrypter(block, iv)
	stream.CryptBlocks(cipherText, cipherText)

	return string(cipherText), nil
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

func HttpBuildQuery(params map[string]interface{}) (paramsStr string) {
	/*
	 * golang 取代 PHP http_build_query 的做法
	 * 會把 key 按照 a~z 排序
	 */

	var u url.URL
	q := u.Query()

	for k, v := range params {
		q.Add(k, fmt.Sprintf("%v", v))
	}
	paramsStr = q.Encode()

	return paramsStr
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

func Aes256CBCEncryptTypeI(dataByte string, keyStr string, ivStr string) (string, error) {
	/*
	 * 藍新金流用
	 * 可對應到 php openssl_encrypt();
	 * pkcs7Padding 對應 PHP OPENSSL_RAW_DATA
	 * hex.EncodeToString 對應 PHP bin2hex
	 */

	keyByte := []byte(keyStr)
	ivByte := []byte(ivStr)
	plaintextByte := pkcs7Padding([]byte(dataByte))

	block, err := aes.NewCipher(keyByte)
	log.Error(err)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, len(plaintextByte))
	mode := cipher.NewCBCEncrypter(block, ivByte)
	mode.CryptBlocks(cipherText, plaintextByte)

	return hex.EncodeToString(cipherText), nil
}

func CheckcodeChcek(paramsMap map[string]interface{}, checkCode string, hashKey string, hashIV string) bool {
	/*
	 * 藍新驗證 CheckCode
	 */

	httpQueryStr := HttpBuildQuery(paramsMap)

	strSlices := []string{
		"HashIV=" + hashIV,
		httpQueryStr,
		"HashKey=" + hashKey,
	}
	paramsStr := strings.Join(strSlices, "&")

	hash := sha256.Sum256([]byte(paramsStr))
	hashStr := strings.ToUpper(fmt.Sprintf("%x", hash))

	if strCompare := strings.Compare(hashStr, checkCode); strCompare == 0 {
		return true
	} else {
		return false
	}
}
