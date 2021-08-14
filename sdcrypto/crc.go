package sdcrypto

import (
	"encoding/binary"
	"hash/crc32"

	"github.com/gaorx/stardust3/sdbytes"
	"github.com/gaorx/stardust3/sderr"
)

// CRC32

type Crc32Encrypter struct {
	Encrypter Encrypter
}

func (e *Crc32Encrypter) Encrypt(key, data []byte) ([]byte, error) {
	data = sdbytes.Copy(data)
	sumBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(sumBytes, crc32.ChecksumIEEE(data))
	r, err := e.Encrypter.Encrypt(key, append(data, sumBytes...))
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	return r, nil
}

func (e *Crc32Encrypter) Decrypt(key, crypted []byte) ([]byte, error) {
	decrypted, err := e.Encrypter.Decrypt(key, crypted)
	if err != nil {
		return nil, sderr.WithStack(err)
	}
	n := len(decrypted)
	if n < 4 {
		return nil, sderr.New("decrypted is too short")
	}
	data, sumBytes := decrypted[0:n-4], decrypted[n-4:]
	expectant := binary.LittleEndian.Uint32(sumBytes)
	if crc32.ChecksumIEEE(data) != expectant {
		return nil, sderr.New("crc32 error")
	}
	return data, nil
}
