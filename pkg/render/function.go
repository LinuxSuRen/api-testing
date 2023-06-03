package render

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

// BytesToPublicKey bytes to public key
func BytesToPublicKey(data []byte) (pub *rsa.PublicKey, err error) {
	block, _ := pem.Decode(data)
	if block != nil {
		enc := x509.IsEncryptedPEMBlock(block)
		b := block.Bytes
		if enc {
			b, err = x509.DecryptPEMBlock(block, nil)
			if err != nil {
				return
			}
		}
		pub, err = PKIXPublicKey(b)
	}
	return
}

func PKIXPublicKey(data []byte) (pub *rsa.PublicKey, err error) {
	var ifc interface{}
	ifc, err = x509.ParsePKIXPublicKey(data)
	if err == nil {
		pub = ifc.(*rsa.PublicKey)
	}
	return
}

func Base64PKIXPublicKey(data []byte) (pub *rsa.PublicKey, err error) {
	data, err = base64.StdEncoding.DecodeString(string(data))
	if err == nil {
		pub, err = PKIXPublicKey(data)
	}
	return
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) (data []byte, err error) {
	hash := sha512.New()
	data, err = rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	return
}
