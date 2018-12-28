package mini

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type User struct {
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarURL string `json:"avatarUrl"`
	UnionID   string `json:"unionId"`
	Watermark struct {
		Timestamp int64  `json:"timestamp"`
		Appid     string `json:"appid"`
	} `json:"watermark"`
}

func DecryptUserInfo(encryptedData, sessionKey string) (u *User, err error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		err = fmt.Errorf("the encryptedData not base64")
		return
	}
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		err = fmt.Errorf("the sessionKey not base64")
		return
	}
	const (
		BlockSize = 32 // PKCS#7
	)

	if len(ciphertext) < BlockSize {
		err = fmt.Errorf("the length of ciphertext too short: %d", len(ciphertext))
		return
	}

	plaintext := make([]byte, len(ciphertext)) // len(plaintext) >= BLOCK_SIZE

	// 解密
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, aesKey[:16])
	mode.CryptBlocks(plaintext, ciphertext)

	buf := bytes.NewBufferString("{")
	i := bytes.Index(plaintext, []byte("nickName"))
	buf.Write(bytes.Trim(plaintext[i-1:], "\x05"))

	u = &User{}
	err = json.NewDecoder(buf).Decode(u)
	if err != nil {
		return
	}
	return
}
