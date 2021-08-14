package sdcrypto

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/gaorx/stardust3/sderr"
)

var (
	Aes Encrypter = &EncrypterFunc{
		Encrypter: AesEncrypt,
		Decrypter: AesDecrypt,
	}
	AesCrc32 Encrypter = &Crc32Encrypter{Aes}
)

func AesEncrypt(key, data []byte) ([]byte, error) {
	r, err := AesEncryptPadding(key, data, Pkcs5)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

func AesDecrypt(key, crypted []byte) ([]byte, error) {
	r, err := AesDecryptPadding(key, crypted, UnPkcs5)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

func AesEncryptPadding(key, data []byte, p Padding) ([]byte, error) {
	if p == nil {
		return nil, sderr.New("nil padding")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	data, err = p(data, block.BlockSize())
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	encrypter := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	crypted := make([]byte, len(data))
	encrypter.CryptBlocks(crypted, data)
	return crypted, nil
}

func AesDecryptPadding(key, crypted []byte, p Unpadding) ([]byte, error) {
	if p == nil {
		return nil, sderr.New("nil unpadding")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	decrypter := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	data := make([]byte, len(crypted))
	decrypter.CryptBlocks(data, crypted)
	r, err := p(data, block.BlockSize())
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}
