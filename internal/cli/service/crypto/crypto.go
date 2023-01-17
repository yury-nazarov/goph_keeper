package crypto

import (
	b64 "encoding/base64"
)

type crypto struct {

}

func New() *crypto{
	c := &crypto{}
	return c
}


type Crypto interface {
	Encrypt(data []byte) string
	Decrypt(encData string) (string, error)
}

func (c *crypto) Encrypt(data []byte) string {
	encData := b64.StdEncoding.EncodeToString(data)
	return encData
}
func (c *crypto) Decrypt(encData string) (string, error) {
	data, err := b64.StdEncoding.DecodeString(encData)
	if err != nil {
		return "", err
	}
	return string(data), nil
}