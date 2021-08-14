package sdcrypto

import (
	"encoding/json"

	"github.com/gaorx/stardust3/sdencoding"
	"github.com/gaorx/stardust3/sderr"
)

// bytes

func EncryptBytes(e Encrypter, key, data []byte) (sdencoding.Bytes, error) {
	r, err := e.Encrypt(key, data)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

func DecryptBytes(e Encrypter, key, crypted []byte) (sdencoding.Bytes, error) {
	r, err := e.Decrypt(key, crypted)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

func MustEncryptBytes(e Encrypter, key, data []byte) sdencoding.Bytes {
	crypted, err := EncryptBytes(e, key, data)
	if err != nil {
		panic(sderr.WithStack(err))
	}
	return crypted
}

func MustDecryptBytess(e Encrypter, key, crypted []byte) sdencoding.Bytes {
	data, err := DecryptBytes(e, key, crypted)
	if err != nil {
		panic(sderr.WithStack(err))
	}
	return data
}

// string

func EncryptStr(e Encrypter, encoding sdencoding.Encoding, key, data []byte) (string, error) {
	data, err := EncryptBytes(e, key, data)
	if err != nil {
		return "", sderr.WithStack(err)
	}
	return encoding.EncodeStr(data), nil
}

func DecryptStr(e Encrypter, encoding sdencoding.Encoding, key []byte, crypted string) ([]byte, error) {
	crypted1, err := encoding.DecodeStr(crypted)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	r, err := e.Decrypt(key, crypted1)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

func MustEncryptStr(e Encrypter, encoding sdencoding.Encoding, key, data []byte) string {
	crypted, err := EncryptStr(e, encoding, key, data)
	if err != nil {
		panic(sderr.WithStack(err))
	}
	return crypted
}

func MustDecryptStr(e Encrypter, encoding sdencoding.Encoding, key []byte, crypted string) []byte {
	data, err := DecryptStr(e, encoding, key, crypted)
	if err != nil {
		panic(sderr.WithStack(err))
	}
	return data
}

// Json string

func EncryptJsonStr(e Encrypter, encoding sdencoding.Encoding, key []byte, data interface{}) (string, error) {
	j, err := json.Marshal(data)
	if err != nil {
		return "", sderr.WithStack(err)
	}
	s, err := EncryptStr(e, encoding, key, j)
	if err != nil {
		return "", sderr.WithStack(err)
	}
	return s, nil
}

func DecryptJsonStr(e Encrypter, encoding sdencoding.Encoding, key []byte, crypted string, to interface{}) error {
	j, err := DecryptStr(e, encoding, key, crypted)
	if err != nil {
		return sderr.WithStack(err)
	}
	err = json.Unmarshal(j, to)
	if err != nil {
		return sderr.WithStack(err)
	}
	return nil
}

func MustEncryptJsonStr(e Encrypter, encoding sdencoding.Encoding, key []byte, data interface{}) string {
	cryped, err := EncryptJsonStr(e, encoding, key, data)
	if err != nil {
		panic(sderr.WithStack(err))
	}
	return cryped
}

func MustDecryptJsonStr(e Encrypter, encoding sdencoding.Encoding, key []byte, crypted string, to interface{}) {
	err := DecryptJsonStr(e, encoding, key, crypted, to)
	if err != nil {
		panic(sderr.WithStack(err))
	}
}
