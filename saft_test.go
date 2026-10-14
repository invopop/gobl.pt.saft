package saft_test

import (
	"testing"
	"time"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl.pt.saft/test"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvert(t *testing.T) {

	config.Now = func() time.Time { return time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC) }

	t.Run("Envelope signed but not sealed", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		inv := env.Extract().(*bill.Invoice)
		inv.Code = "1"

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		env.Head.Stamps = nil

		_, err = saft.Convert(env, "SOURCE_ID", &config.Config{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "ATCUD stamp must be present")
	})

	t.Run("Invalid source ID", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		_, err = saft.Seal(env, nil, "AAJFJHNC82", &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		_, err = saft.Convert(env, uuid.V7().String(), &config.Config{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "sourceID must be less than 30 characters")
	})

	t.Run("General case (invoice)", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		_, err = saft.Seal(env, nil, "AAJFJHNC82", &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		res, err := saft.Convert(env, "SOURCE_ID", &config.Config{
			SigningKeyVersion: "3",
		})
		require.NoError(t, err)

		if assert.NotNil(t, res.Supplier) {
			assert.Equal(t, "770013813", res.Supplier.SupplierTaxID)
		}
		if assert.NotNil(t, res.Customer) {
			assert.Equal(t, "514329874", res.Customer.CustomerTaxID)
		}
		if inv := res.Document.(*saft.Invoice); assert.NotNil(t, inv) {
			assert.Equal(t, "FT SERIES-A/1", inv.InvoiceNo)
			assert.Equal(t, "SOURCE_ID", inv.SourceID)
			assert.Equal(t, "3", inv.HashControl)
		}
		if assert.Len(t, res.Products, 2) {
			assert.Equal(t, "Development services", res.Products[0].ProductDesc)
		}
		if assert.Len(t, res.TaxTableEntries, 2) {
			assert.Equal(t, "NOR-PT-23.0", res.TaxTableEntries[0].Description)
		}
	})

	t.Run("General case (payment)", func(t *testing.T) {
		env, err := test.LoadEnvelope("receipt.json")
		require.NoError(t, err)

		_, err = saft.Seal(env, nil, "AAJFJHNC82", &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		res, err := saft.Convert(env, "SOURCE_ID", &config.Config{
			SigningKeyVersion: "3",
		})
		require.NoError(t, err)

		if assert.NotNil(t, res.Supplier) {
			assert.Equal(t, "770013813", res.Supplier.SupplierTaxID)
		}
		if assert.NotNil(t, res.Customer) {
			assert.Equal(t, "514329874", res.Customer.CustomerTaxID)
		}
		if pay := res.Document.(*saft.Payment); assert.NotNil(t, pay) {
			assert.Equal(t, "RG SERIES-A/1", pay.PaymentRefNo)
			assert.Equal(t, "SOURCE_ID", pay.SourceID)
		}
		assert.Len(t, res.Products, 0)
		assert.Len(t, res.TaxTableEntries, 0)
	})

	t.Run("General case (proforma invoice)", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-proforma.json")
		require.NoError(t, err)

		_, err = saft.Seal(env, nil, "AAJFJHNC82", &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		res, err := saft.Convert(env, "SOURCE_ID", &config.Config{
			SigningKeyVersion: "3",
		})
		require.NoError(t, err)

		if assert.NotNil(t, res.Supplier) {
			assert.Equal(t, "770013813", res.Supplier.SupplierTaxID)
		}
		if assert.NotNil(t, res.Customer) {
			assert.Equal(t, "514329874", res.Customer.CustomerTaxID)
		}
		if wdoc := res.Document.(*saft.WorkDocument); assert.NotNil(t, wdoc) {
			assert.Equal(t, "PF SERIES-A/1", wdoc.DocumentNumber)
			assert.Equal(t, "SOURCE_ID", wdoc.SourceID)
			assert.Equal(t, "3", wdoc.HashControl)
		}
		if assert.Len(t, res.Products, 1) {
			assert.Equal(t, "Development services", res.Products[0].ProductDesc)
		}
		if assert.Len(t, res.TaxTableEntries, 1) {
			assert.Equal(t, "NOR-PT-23.0", res.TaxTableEntries[0].Description)
		}
	})

	t.Run("General case (delivery document)", func(t *testing.T) {
		env, err := test.LoadEnvelope("delivery-note.json")
		require.NoError(t, err)

		_, err = saft.Seal(env, nil, "AAJFJHNC82", &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		res, err := saft.Convert(env, "SOURCE_ID", &config.Config{
			SigningKeyVersion: "3",
		})
		require.NoError(t, err)

		if assert.NotNil(t, res.Supplier) {
			assert.Equal(t, "770013813", res.Supplier.SupplierTaxID)
		}
		if assert.NotNil(t, res.Customer) {
			assert.Equal(t, "514329874", res.Customer.CustomerTaxID)
		}
		if smov := res.Document.(*saft.StockMovement); assert.NotNil(t, smov) {
			assert.Equal(t, "SOURCE_ID", smov.SourceID)
			assert.Equal(t, "3", smov.HashControl)
		}
		if assert.Len(t, res.Products, 2) {
			assert.Equal(t, "Development services", res.Products[0].ProductDesc)
		}
		if assert.Len(t, res.TaxTableEntries, 2) {
			assert.Equal(t, "NOR-PT-23.0", res.TaxTableEntries[0].Description)
		}
	})
}

func TestTimestamp(t *testing.T) {
	t.Run("uses Lisbon time around midnight in summer", func(t *testing.T) {
		config.Now = func() time.Time { return time.Date(2026, 9, 25, 23, 25, 45, 0, time.UTC) }
		assert.Equal(t, "2026-09-26T00:25:45", saft.Timestamp())
	})
}
