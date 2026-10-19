package saft_test

import (
	"testing"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestBuildCustomer(t *testing.T) {
	t.Run("General case", func(t *testing.T) {
		cus := &org.Party{
			Name: "Innovatech Lda",
			TaxID: &tax.Identity{
				Code:    "123456789",
				Country: "PT",
			},
			People: []*org.Person{
				{
					Name: &org.Name{
						Given:   "José",
						Surname: "Mourinho",
					},
				},
			},
			Addresses: []*org.Address{
				{
					Street:   "Avenida da Liberdade",
					Number:   "152",
					Floor:    "1º",
					Block:    "B",
					Door:     "D",
					Locality: "Lisboa",
					Code:     "1250-146",
					Region:   "Distrito de Lisboa",
					Country:  "PT",
				},
			},
			Telephones: []*org.Telephone{
				{
					Number: "123456789",
				},
				{
					Label:  "Fax",
					Number: "987654321",
				},
			},
			Emails: []*org.Email{
				{
					Address: "innovatech@example.com",
				},
			},
			Websites: []*org.Website{
				{
					URL: "https://innovatech.example.com",
				},
			},
		}

		c := saft.BuildCustomer(cus)

		assert.Equal(t, "123456789", c.CustomerID)
		assert.Equal(t, "Desconhecido", c.AccountID)
		assert.Equal(t, "123456789", c.CustomerTaxID)
		assert.Equal(t, "Innovatech Lda", c.CompanyName)
		assert.Equal(t, "José Mourinho", c.Contact)
		if assert.NotNil(t, c.BillingAddress) {
			assert.Equal(t, "152", c.BillingAddress.BuildingNumber)
			assert.Equal(t, "Avenida da Liberdade", c.BillingAddress.StreetName)
			assert.Equal(t, "Avenida da Liberdade 152 B 1º D", c.BillingAddress.AddressDetail)
			assert.Equal(t, "Lisboa", c.BillingAddress.City)
			assert.Equal(t, "1250-146", c.BillingAddress.PostalCode)
			assert.Equal(t, "Distrito de Lisboa", c.BillingAddress.Region)
			assert.Equal(t, "PT", c.BillingAddress.Country)
		}
		assert.Equal(t, "123456789", c.Telephone)
		assert.Equal(t, "987654321", c.Fax)
		assert.Equal(t, "innovatech@example.com", c.Email)
		assert.Equal(t, "https://innovatech.example.com", c.Website)
		assert.Equal(t, 0, c.SelfBillingIndicator)
	})

	finalConsumers := []struct {
		name    string
		cus     *org.Party
		expID   string
		expName string
	}{
		{
			name:    "missing",
			cus:     nil,
			expID:   "Consumidor final",
			expName: "Consumidor final",
		},
		{
			name:    "empty",
			cus:     &org.Party{},
			expID:   "Consumidor final",
			expName: "Consumidor final",
		},
		{
			name:    "missing tax ID",
			cus:     &org.Party{Name: "José Mourinho"},
			expID:   "6c0d40b4e4d56ad360f0b8b9670acd",
			expName: "José Mourinho",
		},
		{
			name: "missing tax ID code",
			cus: &org.Party{
				TaxID: &tax.Identity{},
				Name:  "José Mourinho",
			},
			expID:   "6c0d40b4e4d56ad360f0b8b9670acd",
			expName: "José Mourinho",
		},
	}

	for _, tc := range finalConsumers {
		t.Run("Final consumer – "+tc.name, func(t *testing.T) {
			c := saft.BuildCustomer(tc.cus)

			assert.Equal(t, tc.expID, c.CustomerID)
			assert.Equal(t, tc.expName, c.CompanyName)
			assert.Equal(t, "Desconhecido", c.AccountID)
			assert.Equal(t, "999999990", c.CustomerTaxID)
			if assert.NotNil(t, c.BillingAddress) {
				assert.Equal(t, "Desconhecido", c.BillingAddress.AddressDetail)
				assert.Equal(t, "Desconhecido", c.BillingAddress.City)
				assert.Equal(t, "Desconhecido", c.BillingAddress.PostalCode)
				assert.Equal(t, "Desconhecido", c.BillingAddress.Country)
			}
		})
	}

	t.Run("Final consumer – generic PT tax ID", func(t *testing.T) {
		genericTaxID := func() *tax.Identity {
			return &tax.Identity{Country: "PT", Code: "999999990"}
		}

		guest1 := saft.BuildCustomer(&org.Party{
			Name:  "José Mourinho",
			TaxID: genericTaxID(),
		})
		guest2 := saft.BuildCustomer(&org.Party{
			Name:  "Cristiano Ronaldo",
			TaxID: genericTaxID(),
			Addresses: []*org.Address{
				{Locality: "London", Country: "GB"},
			},
		})
		anonymous := saft.BuildCustomer(&org.Party{TaxID: genericTaxID()})

		for _, c := range []*saft.Customer{guest1, guest2, anonymous} {
			assert.NotEqual(t, "999999990", c.CustomerID)
			assert.Equal(t, "999999990", c.CustomerTaxID)
		}

		assert.Equal(t, "José Mourinho", guest1.CompanyName)
		assert.Equal(t, "PT", guest1.BillingAddress.Country)
		assert.Equal(t, "Cristiano Ronaldo", guest2.CompanyName)
		assert.Equal(t, "GB", guest2.BillingAddress.Country)
		assert.Equal(t, "Consumidor final", anonymous.CompanyName)
		assert.Equal(t, "PT", anonymous.BillingAddress.Country)

		assert.NotEqual(t, guest1.CustomerID, guest2.CustomerID)
		assert.NotEqual(t, guest1.CustomerID, anonymous.CustomerID)
		assert.Equal(t, anonymous.CustomerID, saft.BuildCustomer(&org.Party{TaxID: genericTaxID()}).CustomerID)
	})

	t.Run("Missing customer address", func(t *testing.T) {
		cus := &org.Party{
			Name: "Innovatech Lda",
		}

		c := saft.BuildCustomer(cus)

		if assert.NotNil(t, c.BillingAddress) {
			assert.Equal(t, "Desconhecido", c.BillingAddress.AddressDetail)
			assert.Equal(t, "Desconhecido", c.BillingAddress.City)
			assert.Equal(t, "Desconhecido", c.BillingAddress.PostalCode)
			assert.Equal(t, "Desconhecido", c.BillingAddress.Country)
		}

		cus.Addresses = []*org.Address{{}}

		c = saft.BuildCustomer(cus)

		if assert.NotNil(t, c.BillingAddress) {
			assert.Equal(t, "Desconhecido", c.BillingAddress.AddressDetail)
			assert.Equal(t, "Desconhecido", c.BillingAddress.City)
			assert.Equal(t, "Desconhecido", c.BillingAddress.PostalCode)
			assert.Equal(t, "Desconhecido", c.BillingAddress.Country)
		}
	})
}
