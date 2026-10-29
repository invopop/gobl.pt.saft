package realtime_test

import (
	"testing"
	"time"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl.pt.saft/realtime"
	"github.com/invopop/gobl.pt.saft/test"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildRegisterPayment(t *testing.T) {
	config.Now = func() time.Time { return time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC) }

	t.Run("General case", func(t *testing.T) {
		env, err := test.LoadEnvelope("receipt.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gpmt := env.Extract().(*bill.Payment)
		gpmt.Ext = gpmt.Ext.Set(addon.ExtKeyPaymentType, "RC") // Only RC is allowed for realtime reporting
		per := cal.Period{
			Start: &gpmt.IssueDate,
			End:   &gpmt.IssueDate,
		}
		af, err := test.GenerateAuditFile(gpmt.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		res, err := realtime.BuildRegisterPayment(af)
		require.NoError(t, err)
		require.NotNil(t, res)

		// Header fields
		assert.Equal(t, "0.0.1", res.EFaturaMDVersion)
		assert.Equal(t, "770013813", res.TaxRegistrationNumber)
		assert.Equal(t, "Global", res.TaxEntity)
		assert.Equal(t, 123, res.SoftwareCertificateNumber)

		// Payment data fields
		if assert.NotNil(t, res.PaymentData) {
			pmt := af.SourceDocuments.Payments.Payments[0]
			assert.Equal(t, pmt.PaymentRefNo, res.PaymentData.PaymentRefNo)
			assert.Equal(t, pmt.ATCUD, res.PaymentData.ATCUD)
			assert.Equal(t, pmt.TransactionDate, res.PaymentData.TransactionDate)
			assert.Equal(t, pmt.PaymentType, res.PaymentData.PaymentType)
			assert.Equal(t, pmt.SystemEntryDate, res.PaymentData.SystemEntryDate)

			// Document totals
			if assert.NotNil(t, res.PaymentData.DocumentTotals) {
				assert.Equal(t, pmt.DocumentTotals.TaxPayable, res.PaymentData.DocumentTotals.TaxPayable)
				assert.Equal(t, pmt.DocumentTotals.NetTotal, res.PaymentData.DocumentTotals.NetTotal)
				assert.Equal(t, pmt.DocumentTotals.GrossTotal, res.PaymentData.DocumentTotals.GrossTotal)
			}

			// Line summaries
			if assert.NotNil(t, res.PaymentData.LineSummaries) {
				assert.GreaterOrEqual(t, len(res.PaymentData.LineSummaries), 0)
			}
		}
	})

	t.Run("AuditFile is nil", func(t *testing.T) {
		_, err := realtime.BuildRegisterPayment(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must not be nil")
	})

	t.Run("AuditFile without header", func(t *testing.T) {
		af := &saft.AuditFile{}
		_, err := realtime.BuildRegisterPayment(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file header must not be nil")
	})

	t.Run("AuditFile without payments", func(t *testing.T) {
		env, err := test.LoadEnvelope("receipt.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gpmt := env.Extract().(*bill.Payment)
		per := cal.Period{
			Start: &gpmt.IssueDate,
			End:   &gpmt.IssueDate,
		}
		af, err := test.GenerateAuditFile(gpmt.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		// Remove payments
		af.SourceDocuments.Payments = nil

		_, err = realtime.BuildRegisterPayment(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain payments")
	})

	t.Run("AuditFile without payment entries", func(t *testing.T) {
		env, err := test.LoadEnvelope("receipt.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gpmt := env.Extract().(*bill.Payment)
		per := cal.Period{
			Start: &gpmt.IssueDate,
			End:   &gpmt.IssueDate,
		}
		af, err := test.GenerateAuditFile(gpmt.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		// Remove payment entries
		af.SourceDocuments.Payments.Payments = nil

		_, err = realtime.BuildRegisterPayment(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one payment")
	})

	t.Run("AuditFile without customers", func(t *testing.T) {
		env, err := test.LoadEnvelope("receipt.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gpmt := env.Extract().(*bill.Payment)
		per := cal.Period{
			Start: &gpmt.IssueDate,
			End:   &gpmt.IssueDate,
		}
		af, err := test.GenerateAuditFile(gpmt.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		// Remove customers
		af.MasterFiles.Customers = nil

		_, err = realtime.BuildRegisterPayment(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one customer")
	})

	t.Run("Payment type RG is rejected", func(t *testing.T) {
		env, err := test.LoadEnvelope("receipt.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gpmt := env.Extract().(*bill.Payment)
		gpmt.Ext = gpmt.Ext.Set(addon.ExtKeyPaymentType, "RG") // Only RC is allowed for realtime reporting
		per := cal.Period{
			Start: &gpmt.IssueDate,
			End:   &gpmt.IssueDate,
		}
		af, err := test.GenerateAuditFile(gpmt.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		// Payment type should be RG by default in the test data
		assert.Equal(t, "RG", af.SourceDocuments.Payments.Payments[0].PaymentType)

		_, err = realtime.BuildRegisterPayment(af)
		require.Error(t, err)
		assert.ErrorContains(t, err, "payment type RG is not supported for real-time reporting")
	})

	t.Run("With withholding taxes", func(t *testing.T) {
		env, err := test.LoadEnvelope("receipt.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gpmt := env.Extract().(*bill.Payment)
		gpmt.Ext = gpmt.Ext.Set(addon.ExtKeyPaymentType, "RC") // Only RC is allowed for realtime reporting
		per := cal.Period{
			Start: &gpmt.IssueDate,
			End:   &gpmt.IssueDate,
		}
		af, err := test.GenerateAuditFile(gpmt.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		af.SourceDocuments.Payments.Payments[0].WithholdingTaxes = []*saft.WithholdingTax{
			{WithholdingTaxType: "IRS", WithholdingTaxDescription: "Art. 101.º CIRS", WithholdingTaxAmount: "250.00"},
		}

		res, err := realtime.BuildRegisterPayment(af)
		require.NoError(t, err)
		require.NotNil(t, res)

		if assert.NotNil(t, res.PaymentData) && assert.Len(t, res.PaymentData.WithholdingTaxes, 1) {
			wt := res.PaymentData.WithholdingTaxes[0]
			assert.Equal(t, "IRS", wt.WithholdingTaxType)
			assert.Equal(t, "250.00", wt.WithholdingTaxAmount)
		}
	})
}

func TestBuildChangePaymentStatus(t *testing.T) {
	config.Now = func() time.Time { return time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC) }

	t.Run("General case", func(t *testing.T) {
		env, err := test.LoadEnvelope("receipt.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gpmt := env.Extract().(*bill.Payment)
		gpmt.Ext = gpmt.Ext.Set(addon.ExtKeyPaymentType, "RC") // Only RC is allowed for realtime reporting
		per := cal.Period{
			Start: &gpmt.IssueDate,
			End:   &gpmt.IssueDate,
		}
		af, err := test.GenerateAuditFile(gpmt.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		res, err := realtime.BuildChangePaymentStatus(af)
		require.NoError(t, err)
		require.NotNil(t, res)

		// Header fields
		assert.Equal(t, "0.0.1", res.EFaturaMDVersion)
		assert.Equal(t, "770013813", res.TaxRegistrationNumber)

		// Payment header
		pmt := af.SourceDocuments.Payments.Payments[0]
		if assert.NotNil(t, res.PaymentHeader) {
			assert.Equal(t, pmt.PaymentRefNo, res.PaymentHeader.PaymentRefNo)
			assert.Equal(t, pmt.ATCUD, res.PaymentHeader.ATCUD)
			assert.Equal(t, pmt.TransactionDate, res.PaymentHeader.TransactionDate)
			assert.Equal(t, pmt.PaymentType, res.PaymentHeader.PaymentType)
		}

		// Payment status
		if assert.NotNil(t, res.PaymentStatus) {
			assert.Equal(t, pmt.DocumentStatus.PaymentStatus, res.PaymentStatus.PaymentStatus)
			assert.Equal(t, pmt.DocumentStatus.PaymentStatusDate, res.PaymentStatus.PaymentStatusDate)
		}
	})

	t.Run("AuditFile is nil", func(t *testing.T) {
		_, err := realtime.BuildChangePaymentStatus(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must not be nil")
	})

	t.Run("AuditFile without header", func(t *testing.T) {
		af := &saft.AuditFile{}
		_, err := realtime.BuildChangePaymentStatus(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file header must not be nil")
	})

	t.Run("AuditFile without payments", func(t *testing.T) {
		env, err := test.LoadEnvelope("receipt.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gpmt := env.Extract().(*bill.Payment)
		per := cal.Period{
			Start: &gpmt.IssueDate,
			End:   &gpmt.IssueDate,
		}
		af, err := test.GenerateAuditFile(gpmt.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		// Remove payments
		af.SourceDocuments.Payments = nil

		_, err = realtime.BuildChangePaymentStatus(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain payments")
	})

	t.Run("AuditFile without payment entries", func(t *testing.T) {
		env, err := test.LoadEnvelope("receipt.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gpmt := env.Extract().(*bill.Payment)
		per := cal.Period{
			Start: &gpmt.IssueDate,
			End:   &gpmt.IssueDate,
		}
		af, err := test.GenerateAuditFile(gpmt.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		// Remove payment entries
		af.SourceDocuments.Payments.Payments = nil

		_, err = realtime.BuildChangePaymentStatus(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one payment")
	})

	t.Run("AuditFile without customers", func(t *testing.T) {
		env, err := test.LoadEnvelope("receipt.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gpmt := env.Extract().(*bill.Payment)
		per := cal.Period{
			Start: &gpmt.IssueDate,
			End:   &gpmt.IssueDate,
		}
		af, err := test.GenerateAuditFile(gpmt.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		// Remove customers
		af.MasterFiles.Customers = nil

		_, err = realtime.BuildChangePaymentStatus(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one customer")
	})

	t.Run("Payment type RG is rejected", func(t *testing.T) {
		env, err := test.LoadEnvelope("receipt.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gpmt := env.Extract().(*bill.Payment)
		per := cal.Period{
			Start: &gpmt.IssueDate,
			End:   &gpmt.IssueDate,
		}
		af, err := test.GenerateAuditFile(gpmt.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		// Payment type should be RG by default in the test data
		assert.Equal(t, "RG", af.SourceDocuments.Payments.Payments[0].PaymentType)

		_, err = realtime.BuildChangePaymentStatus(af)
		require.Error(t, err)
		assert.ErrorContains(t, err, "payment type RG is not supported for real-time reporting")
	})
}
