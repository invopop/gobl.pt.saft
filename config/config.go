// Package config provides the configuration for the other gobl.pt.saft packages
package config

import (
	"os"
	"time"

	"github.com/invopop/gobl/cal"
)

// Now is the function that the all packages use to get the current time. It is
// set to time.Now by default, but can be changed for testing purposes.
var Now = time.Now

// TodayIn returns the current date in the given time zone.
func TodayIn(loc *time.Location) cal.Date {
	return cal.DateOf(Now().In(loc))
}

// Config stores the configuration needed by the gobl.pt.saft packages
type Config struct {
	ProducerTaxID     string `json:"producer_tax_id"`
	ProducerName      string `json:"producer_name"`
	CertificateNumber int    `json:"certificate_number"`
	AppName           string `json:"app_name"`
	AppVersion        string `json:"app_version"`
	SigningKeyPath    string `json:"signing_key_path"`
	SigningKeyVersion string `json:"signing_key_version"`

	// cache
	productID  string `json:"-"`
	signingKey []byte `json:"-"`
}

// ProductID returns the product ID in the expected format
func (i *Config) ProductID() string {
	if i.productID == "" {
		i.productID = i.AppName + "/" + i.ProducerName
	}
	return i.productID
}

// SigningKey reads and returns the configured signing key
func (i *Config) SigningKey() ([]byte, error) {
	if i.signingKey == nil {
		skey, err := os.ReadFile(i.SigningKeyPath)
		if err != nil {
			return nil, err
		}
		i.signingKey = skey
	}
	return i.signingKey, nil
}
