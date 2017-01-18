package token

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/molibei/cibepay-backend/data"
	"github.com/molibei/cibepay-backend/module"
)

type TokenService struct{}

func init() {
	module.Token = &TokenService{}
}

// Decode 解密令牌
func (s *TokenService) Decode(token, secret string) (tk data.Token, ok bool) {
	token = strings.Replace(token, ".000000.", "=", -1)
	token = strings.Replace(token, ".00000.", "/", -1)
	token = strings.Replace(token, ".0000.", "+", -1)
	tokenbyte, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return
	}
	sc := []byte(secret)
	key, err := aes.NewCipher(sc)
	if nil != err {
		return
	}
	decrypter := cipher.NewCBCDecrypter(key, sc)
	in := make([]byte, len(tokenbyte))
	decrypter.CryptBlocks(in, tokenbyte)
	in = UnPKCS7Padding(in)
	tk = data.Token{}
	err = json.Unmarshal(in, &tk)
	if err != nil {
		return
	}
	return tk, true
}

// Encode 加密令牌 .
func (s *TokenService) Encode(tk data.Token, secret string) (token string, err error) {
	sc := []byte(secret)
	tkjson, _ := json.Marshal(tk)
	tkbyte := []byte(tkjson)
	key, err := aes.NewCipher(sc)
	if nil != err {
		return token, err
	}
	encrypter := cipher.NewCBCEncrypter(key, sc)
	tkbyte = PKCS7Padding(tkbyte)
	out := make([]byte, len(tkbyte))
	encrypter.CryptBlocks(out, tkbyte)
	token = base64.StdEncoding.EncodeToString(out)
	token = strings.Replace(token, "+", ".0000.", -1)
	token = strings.Replace(token, "/", ".00000.", -1)
	token = strings.Replace(token, "=", ".000000.", -1)
	return
}

func PKCS7Padding(data []byte) []byte {
	dataLen := len(data)
	var bit16 int
	if dataLen%16 == 0 {
		bit16 = dataLen
	} else {
		bit16 = int(dataLen/16+1) * 16
	}
	paddingNum := bit16 - dataLen
	bitCode := byte(paddingNum)
	padding := make([]byte, paddingNum)
	for i := 0; i < paddingNum; i++ {
		padding[i] = bitCode
	}
	return append(data, padding...)
}

//	去除补位
func UnPKCS7Padding(data []byte) []byte {
	dataLen := len(data)
	endIndex := int(data[dataLen-1])
	if 16 > endIndex {
		return data[:dataLen-endIndex]
	}
	return data
}
