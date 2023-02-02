package encryption

import (
	"encoding/base64"
	"fmt"
)

type PrivateKey struct {
	n, d uint64
}

func (p *PublicKey) Marshal() []byte {
	return []byte(fmt.Sprintf("%d,%d", p.n, p.e))
}

func (p *PublicKey) Unmarshal(a []byte) error {
	if _, err := fmt.Sscanf(string(a), "%d,%d", &p.n, &p.e); err != nil {
		return err
	}
	return nil
}

type PublicKey struct {
	n, e uint64
}

func (p *PrivateKey) Marshal() []byte {
	return []byte(fmt.Sprintf("%d,%d", p.n, p.d))
}

func (p *PrivateKey) Unmarshal(a []byte) error {
	if _, err := fmt.Sscanf(string(a), "%d,%d", &p.n, &p.d); err != nil {
		return err
	}
	return nil
}

func (p *PublicKey) encrypt(m uint64) uint64 {
	return pow(m, p.e, p.n)
}

func (p *PublicKey) EncryptString(a []byte) string {
	encryptedString := make([]byte, 0)
	for i := 0; i < len(a); i++ {
		currentPart := p.encrypt(uint64(a[i]))
		for j := 0; j < 8; j++ {
			encryptedString = append(encryptedString, uint8(currentPart&0xFF))
			currentPart >>= 8
		}
	}
	return base64.StdEncoding.EncodeToString(encryptedString)
}

func (p *PublicKey) String() string {
	return fmt.Sprintf("<%d, %d>", p.n, p.e)
}

func (p *PrivateKey) decrypt(c uint64) uint64 {
	return pow(c, p.d, p.n)
}

func (p *PrivateKey) DecryptString(a string) []byte {
	encryptedArray, err := base64.StdEncoding.DecodeString(a)
	if err != nil {
		panic(err)
	}
	decryptedString := make([]byte, 0)
	for i := 0; i < len(encryptedArray); i += 8 {
		currentPart := uint64(0)
		for j := 7; j >= 0; j-- {
			currentPart <<= 8
			currentPart += uint64(encryptedArray[i+j])
		}
		decryptedString = append(decryptedString, byte(p.decrypt(currentPart)))
	}
	return decryptedString
}

func (p *PrivateKey) String() string {
	return fmt.Sprintf("<%d, %d>", p.n, p.d)
}
