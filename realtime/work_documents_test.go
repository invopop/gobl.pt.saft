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

func TestBuildRegisterWork(t *testing.T) {
	config.Now = func() time.Time { return time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC) }

	t.Run("General case", func(t *testing.T) {
		env, err := test.LoadEnvelope("purchase-order.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gord := env.Extract().(*bill.Order)
		per := cal.Period{
			Start: &gord.IssueDate,
			End:   &gord.IssueDate,
		}
		af, err := test.GenerateAuditFile(gord.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		res, err := realtime.BuildRegisterWork(af)
		require.NoError(t, err)
		require.NotNil(t, res)

		// Header fields
		assert.Equal(t, "0.0.1", res.EFaturaMDVersion)
		assert.Equal(t, "770013813", res.TaxRegistrationNumber)
		assert.Equal(t, "Global", res.TaxEntity)
		assert.Equal(t, 123, res.SoftwareCertificateNumber)

		// Work data fields
		if assert.NotNil(t, res.WorkData) {
			wd := af.SourceDocuments.WorkingDocuments.WorkDocuments[0]
			assert.Equal(t, wd.DocumentNumber, res.WorkData.DocumentNumber)
			assert.Equal(t, wd.ATCUD, res.WorkData.ATCUD)
			assert.Equal(t, wd.WorkDate, res.WorkData.WorkDate)
			assert.Equal(t, wd.WorkType, res.WorkData.WorkType)
			assert.Equal(t, wd.SystemEntryDate, res.WorkData.SystemEntryDate)

			// Document totals
			if assert.NotNil(t, res.WorkData.DocumentTotals) {
				assert.Equal(t, wd.DocumentTotals.TaxPayable, res.WorkData.DocumentTotals.TaxPayable)
				assert.Equal(t, wd.DocumentTotals.NetTotal, res.WorkData.DocumentTotals.NetTotal)
				assert.Equal(t, wd.DocumentTotals.GrossTotal, res.WorkData.DocumentTotals.GrossTotal)
			}

			// Line summaries (may be nil or empty)
			if res.WorkData.LineSummaries != nil {
				assert.GreaterOrEqual(t, len(res.WorkData.LineSummaries), 0)
			}
		}
	})

	t.Run("AuditFile is nil", func(t *testing.T) {
		_, err := realtime.BuildRegisterWork(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must not be nil")
	})

	t.Run("AuditFile without header", func(t *testing.T) {
		af := &saft.AuditFile{}
		_, err := realtime.BuildRegisterWork(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file header must not be nil")
	})

	t.Run("AuditFile without working documents", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		ginv := env.Extract().(*bill.Invoice)
		per := cal.Period{
			Start: &ginv.IssueDate,
			End:   &ginv.IssueDate,
		}
		af, err := test.GenerateAuditFile(ginv.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		// Remove working documents
		af.SourceDocuments.WorkingDocuments = nil

		_, err = realtime.BuildRegisterWork(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain working documents")
	})

	t.Run("AuditFile without work document entries", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		ginv := env.Extract().(*bill.Invoice)
		per := cal.Period{
			Start: &ginv.IssueDate,
			End:   &ginv.IssueDate,
		}

		// Convert invoice to work document
		wd, err := saft.BuildWorkDocumentFromInvoice(ginv, &saft.BuildParams{
			ATCUD:             "AAJFJHNC82-1",
			SourceID:          "SOURCE123",
			Hash:              "HASH123",
			SigningKeyVersion: "1",
			SystemEntryDate:   "2024-08-03T10:20:30",
		})
		require.NoError(t, err)

		// Build audit file with work document
		af, err := saft.BuildAuditFile(&saft.AuditFileParams{
			Config:          *cfg,
			Issuer:          ginv.Supplier,
			Period:          per,
			SourceDocuments: []saft.SourceDocument{wd},
		})
		require.NoError(t, err)

		// Remove work document entries
		af.SourceDocuments.WorkingDocuments.WorkDocuments = nil

		_, err = realtime.BuildRegisterWork(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one work document")
	})

	t.Run("AuditFile without customers", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		ginv := env.Extract().(*bill.Invoice)
		per := cal.Period{
			Start: &ginv.IssueDate,
			End:   &ginv.IssueDate,
		}

		// Convert invoice to work document
		wd, err := saft.BuildWorkDocumentFromInvoice(ginv, &saft.BuildParams{
			ATCUD:             "AAJFJHNC82-1",
			SourceID:          "SOURCE123",
			Hash:              "HASH123",
			SigningKeyVersion: "1",
			SystemEntryDate:   "2024-08-03T10:20:30",
		})
		require.NoError(t, err)

		// Build audit file with work document
		af, err := saft.BuildAuditFile(&saft.AuditFileParams{
			Config:          *cfg,
			Issuer:          ginv.Supplier,
			Period:          per,
			SourceDocuments: []saft.SourceDocument{wd},
		})
		require.NoError(t, err)

		// Remove customers
		af.MasterFiles.Customers = nil

		_, err = realtime.BuildRegisterWork(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one customer")
	})
}

func TestBuildChangeWorkStatus(t *testing.T) {
	config.Now = func() time.Time { return time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC) }

	t.Run("General case", func(t *testing.T) {
		env, err := test.LoadEnvelope("purchase-order.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		gord := env.Extract().(*bill.Order)
		per := cal.Period{
			Start: &gord.IssueDate,
			End:   &gord.IssueDate,
		}
		af, err := test.GenerateAuditFile(gord.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		res, err := realtime.BuildChangeWorkStatus(af)
		require.NoError(t, err)
		require.NotNil(t, res)

		// Header fields
		assert.Equal(t, "0.0.1", res.EFaturaMDVersion)
		assert.Equal(t, "770013813", res.TaxRegistrationNumber)

		// Work header
		wd := af.SourceDocuments.WorkingDocuments.WorkDocuments[0]
		if assert.NotNil(t, res.WorkHeader) {
			assert.Equal(t, wd.DocumentNumber, res.WorkHeader.DocumentNumber)
			assert.Equal(t, wd.ATCUD, res.WorkHeader.ATCUD)
			assert.Equal(t, wd.WorkDate, res.WorkHeader.WorkDate)
			assert.Equal(t, wd.WorkType, res.WorkHeader.WorkType)
		}

		// Work status
		if assert.NotNil(t, res.WorkStatus) {
			assert.Equal(t, wd.DocumentStatus.WorkStatus, res.WorkStatus.WorkStatus)
			assert.Equal(t, wd.DocumentStatus.WorkStatusDate, res.WorkStatus.WorkStatusDate)
		}
	})

	t.Run("AuditFile is nil", func(t *testing.T) {
		_, err := realtime.BuildChangeWorkStatus(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must not be nil")
	})

	t.Run("AuditFile without header", func(t *testing.T) {
		af := &saft.AuditFile{}
		_, err := realtime.BuildChangeWorkStatus(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file header must not be nil")
	})

	t.Run("AuditFile without working documents", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		ginv := env.Extract().(*bill.Invoice)
		per := cal.Period{
			Start: &ginv.IssueDate,
			End:   &ginv.IssueDate,
		}
		af, err := test.GenerateAuditFile(ginv.Supplier, per, cfg, "", env)
		require.NoError(t, err)

		// Remove working documents
		af.SourceDocuments.WorkingDocuments = nil

		_, err = realtime.BuildChangeWorkStatus(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain working documents")
	})

	t.Run("AuditFile without work document entries", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		ginv := env.Extract().(*bill.Invoice)
		per := cal.Period{
			Start: &ginv.IssueDate,
			End:   &ginv.IssueDate,
		}

		// Convert invoice to work document
		wd, err := saft.BuildWorkDocumentFromInvoice(ginv, &saft.BuildParams{
			ATCUD:             "AAJFJHNC82-1",
			SourceID:          "SOURCE123",
			Hash:              "HASH123",
			SigningKeyVersion: "1",
			SystemEntryDate:   "2024-08-03T10:20:30",
		})
		require.NoError(t, err)

		// Build audit file with work document
		af, err := saft.BuildAuditFile(&saft.AuditFileParams{
			Config:          *cfg,
			Issuer:          ginv.Supplier,
			Period:          per,
			SourceDocuments: []saft.SourceDocument{wd},
		})
		require.NoError(t, err)

		// Remove work document entries
		af.SourceDocuments.WorkingDocuments.WorkDocuments = nil

		_, err = realtime.BuildChangeWorkStatus(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one work document")
	})

	t.Run("AuditFile without customers", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		cfg := &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		}
		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJHNC82", cfg)
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		ginv := env.Extract().(*bill.Invoice)
		per := cal.Period{
			Start: &ginv.IssueDate,
			End:   &ginv.IssueDate,
		}

		// Convert invoice to work document
		wd, err := saft.BuildWorkDocumentFromInvoice(ginv, &saft.BuildParams{
			ATCUD:             "AAJFJHNC82-1",
			SourceID:          "SOURCE123",
			Hash:              "HASH123",
			SigningKeyVersion: "1",
			SystemEntryDate:   "2024-08-03T10:20:30",
		})
		require.NoError(t, err)

		// Build audit file with work document
		af, err := saft.BuildAuditFile(&saft.AuditFileParams{
			Config:          *cfg,
			Issuer:          ginv.Supplier,
			Period:          per,
			SourceDocuments: []saft.SourceDocument{wd},
		})
		require.NoError(t, err)

		// Remove customers
		af.MasterFiles.Customers = nil

		_, err = realtime.BuildChangeWorkStatus(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one customer")
	})
}
