package saft_test

import (
	"context"
	"testing"
	"time"

	"github.com/invopop/gobl"
	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl.pt.saft/test"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateXSD(t *testing.T) {
	config.Now = func() time.Time { return time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC) }

	ctx := context.Background()

	seal := func(t *testing.T, name string) *gobl.Envelope {
		t.Helper()

		env, err := test.LoadEnvelope(name)
		require.NoError(t, err)

		_, err = saft.Seal(env, nil, "AAJFJHNC82", &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		return env
	}

	convert := func(t *testing.T, name string) *saft.ConvertResult {
		t.Helper()

		res, err := saft.Convert(seal(t, name), "SOURCE_ID", &config.Config{
			SigningKeyVersion: "3",
		})
		require.NoError(t, err)

		return res
	}

	generate := func(t *testing.T, name string) *saft.AuditFile {
		t.Helper()

		env := seal(t, name)
		sup := env.Extract().(*bill.Invoice).Supplier

		per := cal.Period{
			Start: cal.NewDate(2024, 8, 1),
			End:   cal.NewDate(2024, 8, 3),
		}
		af, err := test.GenerateAuditFile(sup, per, &config.Config{
			ProducerTaxID:     "999999999",
			ProducerName:      "Software Company",
			CertificateNumber: 123,
			AppName:           "App Name",
			AppVersion:        "1.0",
			SigningKeyVersion: "3",
		}, "SOURCE_ID", env)
		require.NoError(t, err)

		return af
	}

	t.Run("valid invoice", func(t *testing.T) {
		res := convert(t, "invoice-b2b.json")
		assert.Empty(t, res.ValidateXSD(ctx))
	})

	t.Run("valid payment", func(t *testing.T) {
		res := convert(t, "receipt.json")
		assert.Empty(t, res.ValidateXSD(ctx))
	})

	t.Run("valid work document", func(t *testing.T) {
		res := convert(t, "invoice-proforma.json")
		assert.Empty(t, res.ValidateXSD(ctx))
	})

	t.Run("valid work document from an order", func(t *testing.T) {
		res := convert(t, "order-other.json")
		assert.Empty(t, res.ValidateXSD(ctx))
	})

	t.Run("valid stock movement", func(t *testing.T) {
		res := convert(t, "delivery-note.json")
		assert.Empty(t, res.ValidateXSD(ctx))
	})

	t.Run("negative line amount", func(t *testing.T) {
		res := convert(t, "invoice-b2b.json")

		inv := res.Document.(*saft.Invoice)
		inv.Lines[0].CreditAmount = "-10.00"

		errs := res.ValidateXSD(ctx)
		require.Len(t, errs, 1)
		assert.Contains(t, errs[0].Error(), "CreditAmount")
		assert.Contains(t, errs[0].Error(), "minInclusive")
	})

	t.Run("sub-cent negative amounts", func(t *testing.T) {
		res := convert(t, "invoice-b2b.json")

		inv := res.Document.(*saft.Invoice)
		inv.Lines[0].UnitPrice = "-0.000050"
		inv.Lines[0].CreditAmount = "-0.0001"

		errs := res.ValidateXSD(ctx)
		require.Len(t, errs, 2)
		assert.Contains(t, errs[0].Error(), "UnitPrice")
		assert.Contains(t, errs[1].Error(), "CreditAmount")
	})

	t.Run("invalid document number", func(t *testing.T) {
		res := convert(t, "invoice-b2b.json")

		inv := res.Document.(*saft.Invoice)
		inv.InvoiceNo = "FT SERIES A/1"

		errs := res.ValidateXSD(ctx)
		require.Len(t, errs, 1)
		assert.Contains(t, errs[0].Error(), "InvoiceNo")
	})

	t.Run("invalid customer", func(t *testing.T) {
		res := convert(t, "invoice-b2b.json")

		res.Customer.CustomerTaxID = ""

		errs := res.ValidateXSD(ctx)
		require.Len(t, errs, 1)
		assert.Contains(t, errs[0].Error(), "CustomerTaxID")
	})

	t.Run("invalid product", func(t *testing.T) {
		res := convert(t, "invoice-b2b.json")

		res.Products[0].ProductType = "X"

		errs := res.ValidateXSD(ctx)
		require.Len(t, errs, 1)
		assert.Contains(t, errs[0].Error(), "ProductType")
	})

	t.Run("invalid tax table entry", func(t *testing.T) {
		res := convert(t, "invoice-b2b.json")

		res.TaxTableEntries[0].Description = ""

		errs := res.ValidateXSD(ctx)
		require.Len(t, errs, 1)
		assert.Contains(t, errs[0].Error(), "Description")
	})

	t.Run("valid audit file", func(t *testing.T) {
		af := generate(t, "invoice-b2b.json")
		assert.Empty(t, af.ValidateXSD(ctx))
	})

	t.Run("invalid audit file", func(t *testing.T) {
		af := generate(t, "invoice-b2b.json")

		af.Header.TaxRegistrationNumber = "ABC"

		errs := af.ValidateXSD(ctx)
		require.Len(t, errs, 1)
		assert.Contains(t, errs[0].Error(), "TaxRegistrationNumber")
	})

	t.Run("multiple errors are accumulated", func(t *testing.T) {
		res := convert(t, "invoice-b2b.json")

		inv := res.Document.(*saft.Invoice)
		inv.Lines[0].CreditAmount = "-10.00"
		res.Customer.CustomerTaxID = ""

		errs := res.ValidateXSD(ctx)
		assert.Len(t, errs, 2)
	})
}
