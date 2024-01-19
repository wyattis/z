package ztls

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/wyattis/z/znet"
	"github.com/wyattis/z/zos"
)

type LoadConfig struct {
	CAPath     string
	CAKeyPath  string
	CertPath   string
	KeyPath    string
	CertConfig *CertConfig
}

func (c *LoadConfig) SetDefaults() (config *LoadConfig, err error) {
	if c == nil {
		c = &LoadConfig{}
	}
	if c.CAPath == "" {
		c.CAPath = "ca.pem"
	}
	if c.CAKeyPath == "" {
		c.CAKeyPath = "ca.key"
	}
	if c.CertPath == "" {
		c.CertPath = "cert.pem"
	}
	if c.KeyPath == "" {
		c.KeyPath = "cert.key"
	}
	config = c
	c.CertConfig, err = c.CertConfig.setDefaults()
	return
}

type CertConfig struct {
	SerialNumber *big.Int
	Ips          []net.IP
	Duration     time.Duration
	Subject      pkix.Name
	SubjectKeyId []byte
}

func (c *CertConfig) setDefaults() (config *CertConfig, err error) {
	if c == nil {
		c = &CertConfig{}
	}
	if c.Duration == 0 {
		c.Duration = time.Hour * 24 * 365 // 1 year
	}
	if len(c.Ips) == 0 {
		outIp, err := znet.GetOutboundIP()
		if err != nil {
			return config, err
		}
		c.Ips = []net.IP{outIp, net.IPv6loopback, net.IPv4(127, 0, 0, 1)}
	}
	if c.SubjectKeyId == nil {
		c.SubjectKeyId = []byte{1, 2, 3, 4, 6}
	}
	if c.SerialNumber == nil {
		c.SerialNumber = big.NewInt(2023)
	}
	config = c
	return
}

// Loads the given TLS config or creates a self-signed one if it doesn't exist.
func LoadOrCreateTLS(config *LoadConfig, persistCerts bool) (tlsConfig *tls.Config, err error) {
	var cert tls.Certificate
	if cert, _, err = LoadCertOrCreateCert(config, persistCerts); err != nil {
		return
	}
	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return
}

// Loads the given certs or creates self-signed ones if they don't exist.
func LoadCertOrCreateCert(config *LoadConfig, persistCerts bool) (cert, caCert tls.Certificate, err error) {
	loadedCa := false
	if config.CAPath != "" && config.CAKeyPath != "" && zos.Exists(config.CAPath) && zos.Exists(config.CAKeyPath) {
		caCert, err = tls.LoadX509KeyPair(config.CAPath, config.CAKeyPath)
		if err != nil {
			return
		}
		loadedCa = true
	}

	// check if we can just load the server certs
	if config.CertPath != "" && config.KeyPath != "" && zos.Exists(config.CertPath) && zos.Exists(config.KeyPath) {
		cert, err = tls.LoadX509KeyPair(config.CertPath, config.KeyPath)
		if err != nil {
			return
		}
		return
	}

	// otherwise, we need to create the certs and write them to disk
	if !loadedCa {
		// We need to create a CA, first, I guess
		caCert, err = CreateCA(config.CertConfig)
		if err != nil {
			return
		}
		if persistCerts {
			if err = writeCert(caCert, config.CAPath, config.CAKeyPath); err != nil {
				return
			}
		}
	}

	cert, err = CertFromCa(config.CertConfig, caCert)
	if err != nil {
		return
	}
	if persistCerts {
		if err = writeCert(cert, config.CertPath, config.KeyPath); err != nil {
			return
		}
	}
	return
}

func writeCert(cert tls.Certificate, certPath string, keyPath string) (err error) {
	var f *os.File
	f, err = os.Create(certPath)
	if err != nil {
		return
	}
	defer f.Close()

	if err = pem.Encode(f, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Certificate[0]}); err != nil {
		return
	}

	f, err = os.Create(keyPath)
	if err != nil {
		return
	}
	defer f.Close()
	err = pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(cert.PrivateKey.(*rsa.PrivateKey)),
	})
	return
}

// Creates a certificate signed by the given CA certificate.
func CertFromCa(config *CertConfig, caCert tls.Certificate) (cert tls.Certificate, err error) {
	if config, err = config.setDefaults(); err != nil {
		return
	}
	// set up our server certificate
	c := &x509.Certificate{
		SerialNumber: config.SerialNumber,
		Subject:      config.Subject,
		IPAddresses:  config.Ips,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(config.Duration),
		SubjectKeyId: config.SubjectKeyId,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return
	}

	if _, ok := caCert.PrivateKey.(crypto.Signer); !ok {
		return cert, fmt.Errorf("CA private key does not implement crypto.Signer")
	}
	// if _, ok := any(&certPrivKey.PublicKey).(crypto.Signer); !ok {
	// 	return cert, fmt.Errorf("certPrivKey public key does not implement crypto.Signer")
	// }

	x509Cert, err := x509.ParseCertificate(caCert.Certificate[0])
	if err != nil {
		return
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, c, x509Cert, &certPrivKey.PublicKey, caCert.PrivateKey)
	if err != nil {
		return
	}

	certPEM := new(bytes.Buffer)
	err = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return
	}

	certPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	return tls.X509KeyPair(certPEM.Bytes(), certPrivKeyPEM.Bytes())
}

func CreateCA(config *CertConfig) (caCert tls.Certificate, err error) {
	if config, err = config.setDefaults(); err != nil {
		return
	}

	ca := &x509.Certificate{
		SerialNumber:          config.SerialNumber,
		Subject:               config.Subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(config.Duration),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	err = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return
	}

	caPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	if err != nil {
		return
	}

	return tls.X509KeyPair(caPEM.Bytes(), caPrivKeyPEM.Bytes())
}
