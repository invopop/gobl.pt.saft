package atws

import (
	"errors"
)

// Environment represents the AT Web Services environment to use
type Environment string

// Available environments
const (
	EnvProduction Environment = "production"
	EnvTesting    Environment = "testing"
)

// Endpoints holds the URLs of a AT Web Service in each environment.
type Endpoints map[Environment]string

// Config holds the configuration for the AT Web Services client.
type Config struct {
	Environment         Environment `json:"environment"`
	CertificatePath     string      `json:"cert_path"`
	CertificatePassword string      `json:"cert_pass"`
	PublicKeyPath       string      `json:"pub_key_path"`
	SWCertNum           int         `json:"sw_cert_num"`

	// Control
	Debug bool `json:"debug"`
}

// Validate ensures that the configuration is valid.
func (c *Config) Validate() error {
	if c.Environment == "" {
		return errors.New("environment is required")
	}
	if c.CertificatePath == "" {
		return errors.New("certificate path is required")
	}
	if c.CertificatePassword == "" {
		return errors.New("certificate password is required")
	}
	if c.PublicKeyPath == "" {
		return errors.New("public key path is required")
	}
	return nil
}
