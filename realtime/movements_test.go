package realtime_test

import (
	"testing"
	"time"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl.pt.saft/realtime"
	"github.com/invopop/gobl.pt.saft/test"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildSendStockMovement(t *testing.T) {
	config.Now = func() time.Time { return time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC) }

	t.Run("General case", func(t *testing.T) {
		env, err := test.LoadEnvelope("delivery-note.json")
		require.NoError(t, err)

		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gdlv := env.Extract().(*bill.Delivery)
		per := cal.Period{
			Start: &gdlv.IssueDate,
			End:   &gdlv.IssueDate,
		}
		af, err := test.GenerateAuditFile(gdlv.Supplier, per, &config.Config{}, "", env)
		require.NoError(t, err)

		res, err := realtime.BuildSendStockMovement(af)
		require.NoError(t, err)
		require.NotNil(t, res)

		// Header fields
		assert.Equal(t, "770013813", res.TaxRegistrationNumber)
		assert.Equal(t, "Innovatech Lda", res.CompanyName)
		if assert.NotNil(t, res.CompanyAddress) {
			assert.Equal(t, "Rua das Flores 200 3º Esquerdo", res.CompanyAddress.Addressdetail)
			assert.Equal(t, "Porto", res.CompanyAddress.City)
			assert.Equal(t, "4050-265", res.CompanyAddress.PostalCode)
			assert.Equal(t, "PT", res.CompanyAddress.Country)
		}

		// Document fields
		assert.Equal(t, "GR SERIES-A/1", res.DocumentNumber)
		assert.Equal(t, "AAJFJHNC82-1", res.ATCUD)
		assert.Equal(t, "GR", res.MovementType)
		assert.Equal(t, "2024-08-03", res.MovementDate)
		assert.Equal(t, "N", res.MovementStatus)

		// Customer fields
		assert.Equal(t, "514329874", res.CustomerTaxID)
		assert.Equal(t, "Gusto Unipessoal Lda", res.CustomerName)
		if assert.NotNil(t, res.CustomerAddress) {
			assert.Equal(t, "Avenida da Liberdade 152 2º Direito", res.CustomerAddress.Addressdetail)
			assert.Equal(t, "Lisboa", res.CustomerAddress.City)
			assert.Equal(t, "1250-146", res.CustomerAddress.PostalCode)
			assert.Equal(t, "PT", res.CustomerAddress.Country)
		}

		// Shipping addresses
		if assert.NotNil(t, res.AddressFrom) {
			assert.Equal(t, "Rua da Indústria 45", res.AddressFrom.Addressdetail)
			assert.Equal(t, "Porto", res.AddressFrom.City)
			assert.Equal(t, "4050-300", res.AddressFrom.PostalCode)
			assert.Equal(t, "PT", res.AddressFrom.Country)
		}
		if assert.NotNil(t, res.AddressTo) {
			assert.Equal(t, "Zona Industrial, Lote 23", res.AddressTo.Addressdetail)
			assert.Equal(t, "Lisboa", res.AddressTo.City)
			assert.Equal(t, "1300-300", res.AddressTo.PostalCode)
			assert.Equal(t, "PT", res.AddressTo.Country)
		}

		// Movement times
		assert.Equal(t, "2024-08-03T23:59:59", res.MovementStartTime)
		assert.Empty(t, res.MovementEndTime)

		// Optional fields
		assert.Empty(t, res.ATDocCodeID)

		// Lines
		if assert.Len(t, res.Lines, 2) {
			// First line
			line1 := res.Lines[0]
			assert.Equal(t, "Development services", line1.ProductDescription)
			assert.Equal(t, "1", line1.Quantity)
			assert.Equal(t, "one", line1.UnitOfMeasure)
			assert.Equal(t, "42.5000", line1.UnitPrice)
			if assert.Len(t, line1.OrderReferences, 1) {
				assert.Equal(t, "PO-2024-001", line1.OrderReferences[0].OriginatingON)
				assert.Empty(t, line1.OrderReferences[0].OrderDate)
			}

			// Second line
			line2 := res.Lines[1]
			assert.Equal(t, "Exempt item", line2.ProductDescription)
			assert.Equal(t, "1", line2.Quantity)
			assert.Equal(t, "one", line2.UnitOfMeasure)
			assert.Equal(t, "50.00", line2.UnitPrice)
			assert.Len(t, line2.OrderReferences, 0)
		}

		// Optional fields that may be empty
		assert.Empty(t, res.SupplierTaxID)
		assert.Empty(t, res.VehicleID)
	})

	t.Run("AuditFile is nil", func(t *testing.T) {
		_, err := realtime.BuildSendStockMovement(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must not be nil")
	})

	t.Run("AuditFile without header", func(t *testing.T) {
		af := &saft.AuditFile{}
		_, err := realtime.BuildSendStockMovement(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file header must not be nil")
	})

	t.Run("AuditFile without stock movements", func(t *testing.T) {
		env, err := test.LoadEnvelope("delivery-note.json")
		require.NoError(t, err)

		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gdlv := env.Extract().(*bill.Delivery)
		per := cal.Period{
			Start: &gdlv.IssueDate,
			End:   &gdlv.IssueDate,
		}
		af, err := test.GenerateAuditFile(gdlv.Supplier, per, &config.Config{}, "", env)
		require.NoError(t, err)

		// Remove stock movements
		af.SourceDocuments.MovementOfGoods.StockMovements = nil

		_, err = realtime.BuildSendStockMovement(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one stock movement")
	})

	t.Run("Global transport guide without customer", func(t *testing.T) {
		env, err := test.LoadEnvelope("delivery-waybill-global.json")
		require.NoError(t, err)

		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gdlv := env.Extract().(*bill.Delivery)
		require.Nil(t, gdlv.Customer, "GT global fixture must not declare a customer")

		per := cal.Period{
			Start: &gdlv.IssueDate,
			End:   &gdlv.IssueDate,
		}
		af, err := test.GenerateAuditFile(gdlv.Supplier, per, &config.Config{}, "", env)
		require.NoError(t, err)

		res, err := realtime.BuildSendStockMovement(af)
		require.NoError(t, err)
		require.NotNil(t, res)

		// Customer fields fall back to the AT-mandated generic entry.
		assert.Equal(t, saft.GenericCustomerTaxID, res.CustomerTaxID)
		assert.Equal(t, saft.GenericCustomerID, res.CustomerName)
		// Except the address because the saft.ValueUnknown placeholders aren't allowed in the SOAP schema.
		assert.Nil(t, res.CustomerAddress)

		assert.Equal(t, "GT SERIES-A/1", res.DocumentNumber)
		assert.Equal(t, "GT", res.MovementType)
	})

	t.Run("Preceding GT reference surfaces in line OrderReferences", func(t *testing.T) {
		env, err := test.LoadEnvelope("delivery-with-preceding.json")
		require.NoError(t, err)

		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gdlv := env.Extract().(*bill.Delivery)
		per := cal.Period{
			Start: &gdlv.IssueDate,
			End:   &gdlv.IssueDate,
		}
		af, err := test.GenerateAuditFile(gdlv.Supplier, per, &config.Config{}, "", env)
		require.NoError(t, err)

		res, err := realtime.BuildSendStockMovement(af)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Len(t, res.Lines, 2)

		// First line keeps its per-line gline.Order PO-2024-001 and gains the
		// document-level GT preceding reference. OrderDate is intentionally
		// emptied at the realtime layer (the WSDL has no OrderDate field for
		// stock movements; see realtime/movements.go).
		line1 := res.Lines[0]
		require.Len(t, line1.OrderReferences, 2)
		assert.Equal(t, "PO-2024-001", line1.OrderReferences[0].OriginatingON)
		assert.Empty(t, line1.OrderReferences[0].OrderDate)
		assert.Equal(t, "GT SERIES-A/1", line1.OrderReferences[1].OriginatingON)
		assert.Empty(t, line1.OrderReferences[1].OrderDate)

		// Second line only gains the GT cross-reference.
		line2 := res.Lines[1]
		require.Len(t, line2.OrderReferences, 1)
		assert.Equal(t, "GT SERIES-A/1", line2.OrderReferences[0].OriginatingON)
		assert.Empty(t, line2.OrderReferences[0].OrderDate)
	})

	t.Run("AuditFile without customers", func(t *testing.T) {
		env, err := test.LoadEnvelope("delivery-note.json")
		require.NoError(t, err)

		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gdlv := env.Extract().(*bill.Delivery)
		per := cal.Period{
			Start: &gdlv.IssueDate,
			End:   &gdlv.IssueDate,
		}
		af, err := test.GenerateAuditFile(gdlv.Supplier, per, &config.Config{}, "", env)
		require.NoError(t, err)

		// Remove customers
		af.MasterFiles.Customers = nil

		_, err = realtime.BuildSendStockMovement(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one customer")
	})
}
