package saft_test

import (
	"testing"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestBuildAddress(t *testing.T) {
	t.Run("nil party", func(t *testing.T) {
		a := saft.BuildAddress(nil)
		assert.Nil(t, a)
	})

	t.Run("nil address", func(t *testing.T) {
		p := &org.Party{
			Addresses: []*org.Address{nil},
		}
		a := saft.BuildAddress(p)
		assert.Equal(t, saft.ValueUnknown, a.Country)
	})

	t.Run("missing address and tax ID country", func(t *testing.T) {
		p := &org.Party{}
		a := saft.BuildAddress(p)
		assert.Equal(t, saft.ValueUnknown, a.Country)
	})

	t.Run("missing address but with tax ID country", func(t *testing.T) {
		p := &org.Party{
			TaxID: &tax.Identity{
				Country: "EL",
			},
		}
		a := saft.BuildAddress(p)
		assert.Equal(t, "GR", a.Country)
	})

	t.Run("both address and tax ID country", func(t *testing.T) {
		p := &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
			},
			Addresses: []*org.Address{{
				Country: "PT",
			}},
		}
		a := saft.BuildAddress(p)
		assert.Equal(t, "PT", a.Country)
	})
}
