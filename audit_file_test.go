package saft_test

import (
	"testing"
	"time"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildAuditFile(t *testing.T) {
	per := cal.Period{
		Start: cal.NewDate(2012, time.October, 1),
		End:   cal.NewDate(2012, time.October, 31),
	}

	t.Run("General case", func(t *testing.T) {
		isr := validIssuer()

		inf := &config.Config{
			ProducerTaxID:     "999999999",
			ProducerName:      "Software Company",
			CertificateNumber: 123,
			AppName:           "App Name",
			AppVersion:        "1.0",
		}

		cus := []*saft.Customer{{}}
		pds := []*saft.Product{{}}
		tes := []*saft.TaxTableEntry{{}}
		sds := []saft.SourceDocument{
			&saft.Invoice{
				Lines: []*saft.Line{
					{
						DebitAmount:  "123.45",
						CreditAmount: "111.11",
					},
					{
						DebitAmount:  "111.11",
						CreditAmount: "11.11",
					},
				},
			},
			&saft.Invoice{
				Lines: []*saft.Line{
					{
						DebitAmount:  "111.11",
						CreditAmount: "100.00",
					},
				},
			},
			&saft.Payment{
				Lines: []*saft.PaymentLine{
					{
						DebitAmount:  "543.21",
						CreditAmount: "111.11",
					},
					{
						DebitAmount:  "111.11",
						CreditAmount: "100.00",
					},
				},
			},
			&saft.Payment{
				Lines: []*saft.PaymentLine{
					{
						DebitAmount:  "12.34",
						CreditAmount: "11.11",
					},
				},
			},
			&saft.WorkDocument{
				Lines: []*saft.Line{
					{
						DebitAmount:  "222.22",
						CreditAmount: "111.11",
					},
					{
						DebitAmount:  "111.11",
						CreditAmount: "111.11",
					},
				},
			},
			&saft.WorkDocument{
				Lines: []*saft.Line{
					{
						DebitAmount:  "223.34",
						CreditAmount: "112.23",
					},
				},
			},
			&saft.StockMovement{
				Lines: []*saft.Line{
					{
						Quantity: "7.00",
					},
					{
						Quantity: "3.00",
					},
				},
			},
			&saft.StockMovement{
				Lines: []*saft.Line{
					{
						Quantity: "5.00",
					},
				},
			},
		}

		af, err := saft.BuildAuditFile(&saft.AuditFileParams{
			Config:          *inf,
			Issuer:          isr,
			Period:          per,
			Products:        pds,
			Customers:       cus,
			TaxTableEntries: tes,
			SourceDocuments: sds,
		})
		require.NoError(t, err)

		require.NotNil(t, af)
		if h := af.Header; assert.NotNil(t, h) {
			assert.Equal(t, "1.04_01", h.AuditFileVersion)
			assert.Equal(t, "123456789", h.CompanyID)
			assert.Equal(t, "123456789", h.TaxRegistrationNumber)
			assert.Equal(t, "F", h.TaxAccountingBasis)
			assert.Equal(t, "Company Name", h.CompanyName)
			assert.Equal(t, "Commercial Name", h.BusinessName)
			if a := h.CompanyAddress; assert.NotNil(t, a) {
				assert.Equal(t, "Avenida da Liberdade", a.StreetName)
				assert.Equal(t, "Lisboa", a.City)
				assert.Equal(t, "1250-146", a.PostalCode)
				assert.Equal(t, "PT", a.Country)
			}
			assert.Equal(t, "2012", h.FiscalYear)
			assert.Equal(t, "2012-10-01", h.StartDate)
			assert.Equal(t, "2012-10-31", h.EndDate)
			assert.Equal(t, "EUR", h.CurrencyCode)
			assert.Equal(t, time.Now().Format("2006-01-02"), h.DateCreated)
			assert.Equal(t, "Global", h.TaxEntity)
			assert.Equal(t, "999999999", h.ProductCompanyTaxID)
			assert.Equal(t, 123, h.SoftwareCertificateNumber)
			assert.Equal(t, "App Name/Software Company", h.ProductID)
			assert.Equal(t, "1.0", h.ProductVersion)
			assert.Equal(t, "+351 912 345 678", h.Telephone)
			assert.Equal(t, "company@example.com", h.Email)
			assert.Equal(t, "https://company.example.com", h.Website)
		}
		assert.Len(t, af.MasterFiles.Customers, 1)
		assert.Len(t, af.MasterFiles.Products, 1)
		assert.Len(t, af.MasterFiles.TaxTable.TaxTableEntries, 1)
		if si := af.SourceDocuments.SalesInvoices; assert.NotNil(t, si) {
			assert.Equal(t, 2, si.NumberOfEntries)
			assert.Equal(t, "345.67", si.TotalDebit)
			assert.Equal(t, "222.22", si.TotalCredit)
			assert.Len(t, si.Invoices, 2)
		}
		if py := af.SourceDocuments.Payments; assert.NotNil(t, py) {
			assert.Equal(t, 2, py.NumberOfEntries)
			assert.Equal(t, "666.66", py.TotalDebit)
			assert.Equal(t, "222.22", py.TotalCredit)
			assert.Len(t, py.Payments, 2)
		}
		if wd := af.SourceDocuments.WorkingDocuments; assert.NotNil(t, wd) {
			assert.Equal(t, 2, wd.NumberOfEntries)
			assert.Equal(t, "556.67", wd.TotalDebit)
			assert.Equal(t, "334.45", wd.TotalCredit)
			assert.Len(t, wd.WorkDocuments, 2)
		}
		if mg := af.SourceDocuments.MovementOfGoods; assert.NotNil(t, mg) {
			assert.Equal(t, 3, mg.NumberOfMovementLines) // 2 movements, 3 lines between them
			assert.Equal(t, "15.00", mg.TotalQuantityIssued)
			assert.Len(t, mg.StockMovements, 2)
		}
	})

	t.Run("Canceled documents", func(t *testing.T) {
		sds := []saft.SourceDocument{
			&saft.Invoice{
				DocumentStatus: &saft.InvoiceStatus{InvoiceStatus: saft.StatusCanceled},
				Lines: []*saft.Line{
					{DebitAmount: "100.00", CreditAmount: "50.00"},
				},
			},
			&saft.Invoice{
				DocumentStatus: &saft.InvoiceStatus{InvoiceStatus: saft.StatusBilled},
				Lines: []*saft.Line{
					{DebitAmount: "200.00", CreditAmount: "100.00"},
				},
			},
			&saft.Payment{
				DocumentStatus: &saft.PaymentStatus{PaymentStatus: saft.StatusCanceled},
				Lines: []*saft.PaymentLine{
					{DebitAmount: "100.00", CreditAmount: "50.00"},
				},
			},
			&saft.Payment{
				DocumentStatus: &saft.PaymentStatus{PaymentStatus: saft.StatusBilled},
				Lines: []*saft.PaymentLine{
					{DebitAmount: "200.00", CreditAmount: "100.00"},
				},
			},
			&saft.WorkDocument{
				DocumentStatus: &saft.WorkStatus{WorkStatus: saft.StatusCanceled},
				Lines: []*saft.Line{
					{DebitAmount: "300.00", CreditAmount: "200.00"},
				},
			},
			&saft.WorkDocument{
				DocumentStatus: &saft.WorkStatus{WorkStatus: saft.StatusBilled},
				Lines: []*saft.Line{
					{DebitAmount: "100.00", CreditAmount: "400.00"},
				},
			},
			&saft.StockMovement{
				DocumentStatus: &saft.MovementStatus{MovementStatus: saft.StatusCanceled},
				Lines: []*saft.Line{
					{Quantity: "10.00"},
				},
			},
		}

		af, err := saft.BuildAuditFile(&saft.AuditFileParams{
			Issuer:          validIssuer(),
			Period:          per,
			SourceDocuments: sds,
		})

		require.NoError(t, err)
		require.NotNil(t, af)

		// Check SalesInvoices
		if si := af.SourceDocuments.SalesInvoices; assert.NotNil(t, si) {
			assert.Equal(t, 2, si.NumberOfEntries) // Count all documents
			assert.Equal(t, "0", si.TotalDebit)    // But sum only non-canceled/non-billed
			assert.Equal(t, "0", si.TotalCredit)   // But sum only non-canceled/non-billed
			assert.Len(t, si.Invoices, 2)          // Documents are still included
		}

		// Check Payments
		if p := af.SourceDocuments.Payments; assert.NotNil(t, p) {
			assert.Equal(t, 2, p.NumberOfEntries) // Count all documents
			assert.Equal(t, "0", p.TotalDebit)    // But sum only non-canceled/non-billed
			assert.Equal(t, "0", p.TotalCredit)   // But sum only non-canceled/non-billed
			assert.Len(t, p.Payments, 2)          // Documents are still included
		}

		// Check WorkDocuments
		if wd := af.SourceDocuments.WorkingDocuments; assert.NotNil(t, wd) {
			assert.Equal(t, 2, wd.NumberOfEntries) // Count all documents
			assert.Equal(t, "0", wd.TotalDebit)    // But sum only non-canceled/non-billed
			assert.Equal(t, "0", wd.TotalCredit)   // But sum only non-canceled/non-billed
			assert.Len(t, wd.WorkDocuments, 2)     // Documents are still included
		}

		// Check MovementOfGoods
		if mg := af.SourceDocuments.MovementOfGoods; assert.NotNil(t, mg) {
			assert.Equal(t, 1, mg.NumberOfMovementLines) // Count all lines, canceled included
			assert.Equal(t, "0", mg.TotalQuantityIssued) // But sum only non-canceled
			assert.Len(t, mg.StockMovements, 1)          // Documents are still included
		}
	})
}

func TestBuildAuditFilePaymentsTotalCreditOverflow(t *testing.T) {
	per := cal.Period{
		Start: cal.NewDate(2026, time.May, 1),
		End:   cal.NewDate(2026, time.May, 31),
	}

	sds := []saft.SourceDocument{
		&saft.Payment{
			Lines: []*saft.PaymentLine{
				{CreditAmount: "3.8499999999999996"},
			},
		},
		&saft.Payment{
			Lines: []*saft.PaymentLine{
				{CreditAmount: "1000.00"},
			},
		},
	}

	af, err := saft.BuildAuditFile(&saft.AuditFileParams{
		Issuer:          validIssuer(),
		Period:          per,
		SourceDocuments: sds,
	})
	require.NoError(t, err)
	require.NotNil(t, af)

	if py := af.SourceDocuments.Payments; assert.NotNil(t, py) {
		assert.Equal(t, 2, py.NumberOfEntries)
		assert.Equal(t, "0", py.TotalDebit)
		assert.Equal(t, "1003.85000000", py.TotalCredit)
	}
}

func validIssuer() *org.Party {
	return &org.Party{
		Name:  "Company Name",
		Alias: "Commercial Name",
		TaxID: &tax.Identity{
			Code:    "123456789",
			Country: "PT",
		},
		Addresses: []*org.Address{
			{
				Street:   "Avenida da Liberdade",
				Locality: "Lisboa",
				Code:     "1250-146",
				Country:  "PT",
			},
		},
		Telephones: []*org.Telephone{{Number: "+351 912 345 678"}},
		Emails:     []*org.Email{{Address: "company@example.com"}},
		Websites:   []*org.Website{{URL: "https://company.example.com"}},
	}
}
