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

func TestBuildRegisterInvoice(t *testing.T) {
	config.Now = func() time.Time { return time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC) }

	t.Run("General case", func(t *testing.T) {
		af := generateAuditFile(t)

		res, err := realtime.BuildRegisterInvoice(af)
		require.NoError(t, err)
		require.NotNil(t, res)

		// Header fields
		assert.Equal(t, "0.0.1", res.EFaturaMDVersion)
		assert.Equal(t, "770013813", res.TaxRegistrationNumber)
		assert.Equal(t, "Global", res.TaxEntity)
		assert.Equal(t, 123, res.SoftwareCertificateNumber)

		// Invoice data fields
		if assert.NotNil(t, res.InvoiceData) {
			inv := af.SourceDocuments.SalesInvoices.Invoices[0]
			assert.Equal(t, inv.InvoiceNo, res.InvoiceData.InvoiceNo)
			assert.Equal(t, inv.ATCUD, res.InvoiceData.ATCUD)
			assert.Equal(t, inv.InvoiceDate, res.InvoiceData.InvoiceDate)
			assert.Equal(t, inv.InvoiceType, res.InvoiceData.InvoiceType)
			assert.Equal(t, inv.CustomerID, res.InvoiceData.CustomerTaxID)
			assert.Equal(t, inv.SystemEntryDate, res.InvoiceData.SystemEntryDate)

			// Document totals
			if assert.NotNil(t, res.InvoiceData.DocumentTotals) {
				assert.Equal(t, inv.DocumentTotals.TaxPayable, res.InvoiceData.DocumentTotals.TaxPayable)
				assert.Equal(t, inv.DocumentTotals.NetTotal, res.InvoiceData.DocumentTotals.NetTotal)
				assert.Equal(t, inv.DocumentTotals.GrossTotal, res.InvoiceData.DocumentTotals.GrossTotal)
			}

			// Line summaries
			if assert.NotNil(t, res.InvoiceData.LineSummaries) {
				assert.Greater(t, len(res.InvoiceData.LineSummaries), 0)
			}
		}
	})

	t.Run("AuditFile is nil", func(t *testing.T) {
		_, err := realtime.BuildRegisterInvoice(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must not be nil")
	})

	t.Run("AuditFile without header", func(t *testing.T) {
		af := &saft.AuditFile{}
		_, err := realtime.BuildRegisterInvoice(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file header must not be nil")
	})

	t.Run("AuditFile without sales invoices", func(t *testing.T) {
		af := generateAuditFile(t)

		// Remove sales invoices
		af.SourceDocuments.SalesInvoices = nil

		_, err := realtime.BuildRegisterInvoice(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain sales invoices")
	})

	t.Run("AuditFile without invoices", func(t *testing.T) {
		af := generateAuditFile(t)

		// Remove invoices
		af.SourceDocuments.SalesInvoices.Invoices = nil

		_, err := realtime.BuildRegisterInvoice(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one invoice")
	})

	t.Run("AuditFile without customers", func(t *testing.T) {
		af := generateAuditFile(t)

		// Remove customers
		af.MasterFiles.Customers = nil

		_, err := realtime.BuildRegisterInvoice(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one customer")
	})

	t.Run("With withholding taxes", func(t *testing.T) {
		af := generateAuditFile(t)

		af.SourceDocuments.SalesInvoices.Invoices[0].WithholdingTaxes = []*saft.WithholdingTax{
			{WithholdingTaxType: "IRS", WithholdingTaxDescription: "Art. 101.º CIRS", WithholdingTaxAmount: "250.00"},
		}

		res, err := realtime.BuildRegisterInvoice(af)
		require.NoError(t, err)
		require.NotNil(t, res)

		if assert.NotNil(t, res.InvoiceData) && assert.Len(t, res.InvoiceData.WithholdingTaxes, 1) {
			wt := res.InvoiceData.WithholdingTaxes[0]
			assert.Equal(t, "IRS", wt.WithholdingTaxType)
			assert.Equal(t, "250.00", wt.WithholdingTaxAmount)
		}
	})
}

func TestBuildRegisterInvoiceLineSummaries(t *testing.T) {
	config.Now = func() time.Time { return time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC) }

	// baseLine returns a line that matches every other baseLine for compaction purposes
	baseLine := func(amount string) *saft.Line {
		return &saft.Line{
			OrderReferences: []*saft.OrderReference{{OriginatingON: "ORD-1"}},
			TaxPointDate:    "2024-08-01",
			References:      []*saft.References{{Reference: "FT FT/1"}},
			CreditAmount:    amount,
			Tax: &saft.Tax{
				TaxType:          "IVA",
				TaxCountryRegion: "PT",
				TaxCode:          "NOR",
				TaxPercentage:    "23.00",
			},
		}
	}

	t.Run("Matching lines are compacted", func(t *testing.T) {
		af := generateAuditFile(t)

		af.SourceDocuments.SalesInvoices.Invoices[0].Lines = []*saft.Line{
			baseLine("100.00"),
			baseLine("50.50"),
			baseLine("0.25"),
		}

		res, err := realtime.BuildRegisterInvoice(af)
		require.NoError(t, err)

		if assert.Len(t, res.InvoiceData.LineSummaries, 1) {
			ls := res.InvoiceData.LineSummaries[0]
			assert.Equal(t, "150.75", ls.Amount)
			assert.Equal(t, "C", ls.DebitCreditIndicator)
			assert.Equal(t, "2024-08-01", ls.TaxPointDate)
			assert.Equal(t, []string{"FT FT/1"}, ls.References)
			if assert.Len(t, ls.OrderReferences, 1) {
				assert.Equal(t, "ORD-1", ls.OrderReferences[0].OriginatingON)
			}
			if assert.NotNil(t, ls.Tax) {
				assert.Equal(t, "IVA", ls.Tax.TaxType)
				assert.Equal(t, "PT", ls.Tax.TaxCountryRegion)
				assert.Equal(t, "NOR", ls.Tax.TaxCode)
				assert.Equal(t, "23.00", ls.Tax.TaxPercentage)
			}
		}
	})

	t.Run("Lines with different tax codes are not compacted", func(t *testing.T) {
		af := generateAuditFile(t)

		l1 := baseLine("100.00")
		l2 := baseLine("50.00")
		l2.Tax.TaxCode = "RED"
		l2.Tax.TaxPercentage = "6.00"
		af.SourceDocuments.SalesInvoices.Invoices[0].Lines = []*saft.Line{l1, l2}

		res, err := realtime.BuildRegisterInvoice(af)
		require.NoError(t, err)

		if assert.Len(t, res.InvoiceData.LineSummaries, 2) {
			assert.Equal(t, "NOR", res.InvoiceData.LineSummaries[0].Tax.TaxCode)
			assert.Equal(t, "100.00", res.InvoiceData.LineSummaries[0].Amount)
			assert.Equal(t, "RED", res.InvoiceData.LineSummaries[1].Tax.TaxCode)
			assert.Equal(t, "50.00", res.InvoiceData.LineSummaries[1].Amount)
		}
	})

	t.Run("Lines differing in any key field are not compacted", func(t *testing.T) {
		tests := []struct {
			name   string
			mutate func(l *saft.Line)
		}{
			{
				name:   "TaxType",
				mutate: func(l *saft.Line) { l.Tax.TaxType = "IS" },
			},
			{
				name:   "TaxCountryRegion",
				mutate: func(l *saft.Line) { l.Tax.TaxCountryRegion = "PT-MA" },
			},
			{
				name:   "TaxCode",
				mutate: func(l *saft.Line) { l.Tax.TaxCode = "INT" },
			},
			{
				name: "TaxExemptionCode",
				mutate: func(l *saft.Line) {
					l.Tax.TaxCode = "ISE"
					l.Tax.TaxPercentage = "0.00"
					l.TaxExemptionCode = "M01"
				},
			},
			{
				name:   "TaxPointDate",
				mutate: func(l *saft.Line) { l.TaxPointDate = "2024-08-02" },
			},
			{
				name:   "References",
				mutate: func(l *saft.Line) { l.References = []*saft.References{{Reference: "FT FT/2"}} },
			},
			{
				name:   "Number of references",
				mutate: func(l *saft.Line) { l.References = nil },
			},
			{
				name:   "OrderReferences",
				mutate: func(l *saft.Line) { l.OrderReferences = []*saft.OrderReference{{OriginatingON: "ORD-2"}} },
			},
			{
				name:   "Number of order references",
				mutate: func(l *saft.Line) { l.OrderReferences = nil },
			},
			{
				name: "DebitCreditIndicator",
				mutate: func(l *saft.Line) {
					l.DebitAmount = l.CreditAmount
					l.CreditAmount = ""
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				af := generateAuditFile(t)

				l1 := baseLine("100.00")
				l2 := baseLine("50.00")
				tt.mutate(l2)
				af.SourceDocuments.SalesInvoices.Invoices[0].Lines = []*saft.Line{l1, l2}

				res, err := realtime.BuildRegisterInvoice(af)
				require.NoError(t, err)

				if assert.Len(t, res.InvoiceData.LineSummaries, 2) {
					assert.Equal(t, "100.00", res.InvoiceData.LineSummaries[0].Amount)
					assert.Equal(t, "50.00", res.InvoiceData.LineSummaries[1].Amount)
				}
			})
		}
	})
}

func TestBuildChangeInvoiceStatus(t *testing.T) {
	config.Now = func() time.Time { return time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC) }

	t.Run("General case", func(t *testing.T) {
		af := generateAuditFile(t)

		res, err := realtime.BuildChangeInvoiceStatus(af)
		require.NoError(t, err)
		require.NotNil(t, res)

		// Header fields
		assert.Equal(t, "0.0.1", res.EFaturaMDVersion)
		assert.Equal(t, "770013813", res.TaxRegistrationNumber)

		// Invoice header
		inv := af.SourceDocuments.SalesInvoices.Invoices[0]
		if assert.NotNil(t, res.InvoiceHeader) {
			assert.Equal(t, inv.InvoiceNo, res.InvoiceHeader.InvoiceNo)
			assert.Equal(t, inv.ATCUD, res.InvoiceHeader.ATCUD)
			assert.Equal(t, inv.InvoiceDate, res.InvoiceHeader.InvoiceDate)
			assert.Equal(t, inv.InvoiceType, res.InvoiceHeader.InvoiceType)
		}

		// Invoice status
		if assert.NotNil(t, res.InvoiceStatus) {
			assert.Equal(t, inv.DocumentStatus.InvoiceStatus, res.InvoiceStatus.InvoiceStatus)
			assert.Equal(t, inv.DocumentStatus.InvoiceStatusDate, res.InvoiceStatus.InvoiceStatusDate)
		}
	})

	t.Run("AuditFile is nil", func(t *testing.T) {
		_, err := realtime.BuildChangeInvoiceStatus(nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must not be nil")
	})

	t.Run("AuditFile without header", func(t *testing.T) {
		af := &saft.AuditFile{}
		_, err := realtime.BuildChangeInvoiceStatus(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file header must not be nil")
	})

	t.Run("AuditFile without sales invoices", func(t *testing.T) {
		af := generateAuditFile(t)

		// Remove sales invoices
		af.SourceDocuments.SalesInvoices = nil

		_, err := realtime.BuildChangeInvoiceStatus(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain sales invoices")
	})

	t.Run("AuditFile without invoices", func(t *testing.T) {
		af := generateAuditFile(t)

		// Remove invoices
		af.SourceDocuments.SalesInvoices.Invoices = nil

		_, err := realtime.BuildChangeInvoiceStatus(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one invoice")
	})

	t.Run("AuditFile without customers", func(t *testing.T) {
		af := generateAuditFile(t)

		// Remove customers
		af.MasterFiles.Customers = nil

		_, err := realtime.BuildChangeInvoiceStatus(af)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "audit file must contain at least one customer")
	})
}

func generateAuditFile(t *testing.T) *saft.AuditFile {
	t.Helper()

	env, err := test.LoadEnvelope("invoice-b2b.json")
	require.NoError(t, err)

	cfg := &config.Config{
		CertificateNumber: 123,
		SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
	}
	_, err = saft.Seal(env, nil, "AAJFJHNC82", cfg)
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

	return af
}
