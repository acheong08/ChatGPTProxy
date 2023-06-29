package arkose

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"hash"
)

type encryptionData struct {
	Ct string `json:"ct"`
	Iv string `json:"iv"`
	S  string `json:"s"`
}

func encrypt(data string, key string) string {
	encData, _ := aesEncrypt(data, key)

	encDataJson, err := json.Marshal(encData)
	if err != nil {
		panic(err)
	}

	return string(encDataJson)
}

func aesEncrypt(content string, password string) (*encryptionData, error) {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	key, iv, err := defaultEvpKDF([]byte(password), salt)

	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	cipherBytes := pKCS5Padding([]byte(content), aes.BlockSize)
	mode.CryptBlocks(cipherBytes, cipherBytes)

	//TODO: remove redundant code
	md5Hash := md5.New()
	salted := ""
	var dx []byte

	for i := 0; i < 3; i++ {
		md5Hash.Write(dx)
		md5Hash.Write([]byte(password))
		md5Hash.Write(salt)

		dx = md5Hash.Sum(nil)
		md5Hash.Reset()

		salted += hex.EncodeToString(dx)
	}

	cipherText := base64.StdEncoding.EncodeToString(cipherBytes)
	encData := &encryptionData{
		Ct: cipherText,
		Iv: salted[64 : 64+32],
		S:  hex.EncodeToString(salt),
	}
	return encData, nil
}

// https://stackoverflow.com/questions/27677236/encryption-in-javascript-and-decryption-with-php/27678978#27678978
// https://github.com/brix/crypto-js/blob/8e6d15bf2e26d6ff0af5277df2604ca12b60a718/src/evpkdf.js#L55
func evpKDF(password []byte, salt []byte, keySize int, iterations int, hashAlgorithm string) ([]byte, error) {
	var block []byte
	var hasher hash.Hash
	derivedKeyBytes := make([]byte, 0)
	switch hashAlgorithm {
	case "md5":
		hasher = md5.New()
	default:
		return []byte{}, errors.New("not implement hasher algorithm")
	}
	for len(derivedKeyBytes) < keySize*4 {
		if len(block) > 0 {
			hasher.Write(block)
		}
		hasher.Write(password)
		hasher.Write(salt)
		block = hasher.Sum([]byte{})
		hasher.Reset()

		for i := 1; i < iterations; i++ {
			hasher.Write(block)
			block = hasher.Sum([]byte{})
			hasher.Reset()
		}
		derivedKeyBytes = append(derivedKeyBytes, block...)
	}
	return derivedKeyBytes[:keySize*4], nil
}

func defaultEvpKDF(password []byte, salt []byte) (key []byte, iv []byte, err error) {
	// https://github.com/brix/crypto-js/blob/8e6d15bf2e26d6ff0af5277df2604ca12b60a718/src/cipher-core.js#L775
	keySize := 256 / 32
	ivSize := 128 / 32
	derivedKeyBytes, err := evpKDF(password, salt, keySize+ivSize, 1, "md5")
	if err != nil {
		return []byte{}, []byte{}, err
	}
	return derivedKeyBytes[:keySize*4], derivedKeyBytes[keySize*4:], nil
}
func pKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}
