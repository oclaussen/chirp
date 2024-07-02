package chirp

import (
	"crypto/tls"
        "strings"
        "os"
	"crypto/x509"
	"fmt"
	"net/url"
	"path/filepath"

	"google.golang.org/grpc/credentials"
)

const (
	ErrInvalidAddressScheme ConfigError = "invalid protocol"
	ErrInvalidCA            ConfigError = "invalid ca"
)

type ConfigError string

func (e ConfigError) Error() string {
	return string(e)
}

type Config struct {
	Address         string
	CertificateFile string
	KeyFile         string
	CAFile          string
}

func (c *Config) DialOptions() (string, string, error) {
  values := strings.SplitN(c.Address, "://", 2)
	if len(values) != 2 {
		return "", "", ErrInvalidAddressScheme
	}

	switch values[0] {
	case "tcp":
		return values[0], values[1], nil

	case "unix":
		addr, err := filepath.Abs(values[1])
		if err != nil {
			return "", "", fmt.Errorf("could not get socket path: %w", err)
		}

		return values[0], addr, nil
	}

	return "", "", ErrInvalidAddressScheme
}

func (c *Config) TLSServerOptions() (credentials.TransportCredentials, error) {
	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{},
		ClientCAs:    x509.NewCertPool(),
	}

	certificate, err := tls.LoadX509KeyPair(c.CertificateFile, c.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not load certificate file: %w", err)
	}

	tlsConfig.Certificates = append(tlsConfig.Certificates, certificate)

	bs, err := os.ReadFile(c.CAFile)
	if err != nil {
		return nil, fmt.Errorf("could not read ca file: %w", err)
	}

	if ok := tlsConfig.ClientCAs.AppendCertsFromPEM(bs); !ok {
		return nil, ErrInvalidCA
	}

	return credentials.NewTLS(tlsConfig), nil
}

func (c *Config) TLSClientOptions() (credentials.TransportCredentials, error) {
	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{},
		RootCAs:      x509.NewCertPool(),
	}

	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	if u.Scheme == "tcp" {
		tlsConfig.ServerName = u.Hostname()
	}

	certificate, err := tls.LoadX509KeyPair(c.CertificateFile, c.KeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not load certificate file: %w", err)
	}

	tlsConfig.Certificates = append(tlsConfig.Certificates, certificate)

	bs, err := os.ReadFile(c.CAFile)
	if err != nil {
		return nil, fmt.Errorf("could not read ca file: %w", err)
	}

	if ok := tlsConfig.RootCAs.AppendCertsFromPEM(bs); !ok {
		return nil, ErrInvalidCA
	}

	return credentials.NewTLS(tlsConfig), nil
}
