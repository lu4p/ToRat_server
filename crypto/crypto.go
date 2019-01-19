package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
)

var private = privateKey()

func privateKey() *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(privatekey))
	if block == nil {
		fmt.Println("Could not decode rsakey")
		return nil
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("Could not parse rsakey", err)
		return nil
	}
	return priv
}

func DecRsa(encData []byte) ([]byte, error) {
	rng := rand.Reader
	decData, err := rsa.DecryptOAEP(sha256.New(), rng, private, encData, nil)

	if err != nil {
		log.Println("[!] Rsa:", err)
		return nil, err
	}
	return decData, nil
}

func DecAes(encData []byte, aeskey []byte) ([]byte, error) {
	block, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := aesgcm.Open(nil, encData[:12], encData[12:], nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func DecAsym(encData []byte) ([]byte, error) {
	encAeskey := encData[:512]
	encContent := encData[512:]
	aeskey, err := DecRsa(encAeskey)
	if err != nil {
		return nil, err
	}
	return DecAes(encContent, aeskey)
}

func DecFile(path string) error {
	encData, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	decData, err := DecAsym(encData)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, decData, 0666)
	if err != nil {
		return err
	}
	return nil
}
