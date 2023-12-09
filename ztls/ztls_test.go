package ztls

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func testServerClient(t *testing.T, cert, caCert tls.Certificate) {
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprintln(w, "tls works!")
	}))
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	server.StartTLS()
	defer server.Close()
	t.Log(server.URL)
	x509Cert, err := x509.ParseCertificate(caCert.Certificate[0])
	if err != nil {
		t.Error(err)
	}
	certPool := x509.NewCertPool()
	certPool.AddCert(x509Cert)
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
	}
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	b := bytes.Buffer{}
	if _, err = b.ReadFrom(resp.Body); err != nil {
		t.Fatal(err)
	}
	if b.String() != "tls works!\n" {
		t.Fatalf("expected body 'tls works!', got %q", b.String())
	}
}

func TestTLSServer(t *testing.T) {
	cert, caCert, err := LoadCertOrCreateCert(&LoadConfig{
		CAPath:    "test.ca.pem",
		CAKeyPath: "test.ca.key",
		CertPath:  "test.cert.pem",
		KeyPath:   "test.cert.key",
	}, false)
	if err != nil {
		t.Fatal(err)
	}
	testServerClient(t, cert, caCert)
}

func TestTLSPersistence(t *testing.T) {
	config := LoadConfig{
		CAPath:    "test.ca.pem",
		CAKeyPath: "test.ca.key",
		CertPath:  "test.cert.pem",
		KeyPath:   "test.cert.key",
	}
	cert, caCert, err := LoadCertOrCreateCert(&config, true)
	if err != nil {
		t.Fatal(err)
	}
	testServerClient(t, cert, caCert)
	cert, caCert, err = LoadCertOrCreateCert(&config, true)
	if err != nil {
		t.Fatal(err)
	}
	testServerClient(t, cert, caCert)
	for _, path := range []string{config.CAPath, config.CAKeyPath, config.CertPath, config.KeyPath} {
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}
}
