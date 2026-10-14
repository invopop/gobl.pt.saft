package saft

import (
	"strings"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

// Final consumer constants
const (
	GenericCustomerTaxID = "999999990"
	GenericCustomerID    = "Consumidor final"
	GenericCustomerHash  = "931747646aa67138f346d6ed06fc03"
	ValueUnknown         = "Desconhecido"
)

// Other constants
const (
	LabelFax = "fax"
)

// Customer represents a customer in a SAF-T report
type Customer struct {
	CustomerID           string     `xml:"CustomerID"`
	AccountID            string     `xml:"AccountID"`
	CustomerTaxID        string     `xml:"CustomerTaxID"`
	CompanyName          string     `xml:"CompanyName"`
	Contact              string     `xml:"Contact,omitempty"`
	BillingAddress       *Address   `xml:"BillingAddress"`
	ShipToAddresses      []*Address `xml:"ShipToAddress,omitempty"`
	Telephone            string     `xml:"Telephone,omitempty"`
	Fax                  string     `xml:"Fax,omitempty"`
	Email                string     `xml:"Email,omitempty"`
	Website              string     `xml:"Website,omitempty"`
	SelfBillingIndicator int        `xml:"SelfBillingIndicator"`
}

// BuildCustomer converts a GOBL party to a SAF-T customer
func BuildCustomer(gcus *org.Party) *Customer {
	if gcus == nil {
		gcus = new(org.Party)
	}

	cus := &Customer{
		CustomerTaxID:  customerTaxID(gcus),
		AccountID:      ValueUnknown,
		CompanyName:    companyName(gcus),
		Contact:        contact(gcus),
		BillingAddress: billingAddress(gcus),
		Telephone:      telephone(gcus),
		Fax:            fax(gcus),
		Email:          email(gcus),
		Website:        website(gcus),
	}

	cus.CustomerID = determineCustomerID(gcus, cus)

	return cus
}

func customerTaxID(gcus *org.Party) string {
	if gcus == nil || gcus.TaxID == nil || gcus.TaxID.Code == "" {
		return GenericCustomerTaxID
	}

	return gcus.TaxID.Code.String()
}

func companyName(gcus *org.Party) string {
	if gcus == nil || gcus.Name == "" {
		return GenericCustomerID
	}
	return fit(gcus.Name, 100)
}

func billingAddress(gcus *org.Party) *Address {
	return BuildAddress(gcus)
}

func contact(gcus *org.Party) string {
	if gcus == nil || len(gcus.People) == 0 {
		return ""
	}

	p := gcus.People[0]

	parts := []string{
		p.Name.Given, p.Name.Middle, p.Name.Surname, p.Name.Surname2,
	}

	con := ""
	for _, part := range parts {
		if part == "" {
			continue
		}
		if con != "" {
			con += " "
		}
		con += part
	}

	return fit(con, 50)
}

func fax(gcus *org.Party) string {
	if gcus == nil {
		return ""
	}

	for _, tel := range gcus.Telephones {
		if strings.EqualFold(tel.Label, LabelFax) {
			return fit(tel.Number, 20)
		}
	}

	return ""
}

func determineCustomerID(gcus *org.Party, cus *Customer) string {
	// If the customer is a PT taxpayer, then use the tax ID
	if gcus.TaxID != nil && gcus.TaxID.Code != "" && gcus.TaxID.Country == l10n.PT.Tax() {
		return gcus.TaxID.Code.String()
	}

	// Otherwise, use the hash of the customer data. Any difference in the
	// SAF-T customer data will result in a different customer ID.
	hashID := hashID(cus)
	if hashID == GenericCustomerHash {
		// If all the data is generic, then use the generic customer ID according
		// to 2.2.1 of Portaria 302/2016
		return GenericCustomerID
	}
	return hashID
}
