package proxy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"time"
)

func generate(certPath string, keyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	tmpl := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"Vite Proxy"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(100, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	keyFile, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	err = pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	if err != nil {
		return err
	}

	certFile, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certFile.Close()

	return pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
}

func GetAutoCert() (*tls.Certificate, error) {
	dir, err := Store.Dir()
	if err != nil {
		return nil, err
	}

	s1, err := os.Stat(dir + "/cert.pem")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	s2, err := os.Stat(dir + "/key.pem")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if s1 == nil && s2 == nil {
		err = generate(dir+"/cert.pem", dir+"/key.pem")
		if err != nil {
			return nil, err
		}
	}

	cert, err := tls.LoadX509KeyPair(dir+"/cert.pem", dir+"/key.pem")
	if err != nil {
		return nil, err
	}

	return &cert, nil
}
