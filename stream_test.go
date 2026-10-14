package saft_test

import (
	"bytes"
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

func TestStreamEncoder(t *testing.T) {
	per := cal.Period{
		Start: cal.NewDate(2012, time.October, 1),
		End:   cal.NewDate(2012, time.October, 31),
	}
	cfg := config.Config{
		ProducerTaxID:     "999999999",
		ProducerName:      "Software Company",
		CertificateNumber: 123,
		AppName:           "App Name",
		AppVersion:        "1.0",
	}

	t.Run("matches the built audit file byte for byte", func(t *testing.T) {
		canceled := &saft.Invoice{
			Lines: []*saft.Line{{DebitAmount: "999.99", CreditAmount: "999.99"}},
		}
		canceled.SetStatus(saft.StatusCanceled, "Emitido por engano", "source")

		af, err := saft.BuildAuditFile(&saft.AuditFileParams{
			Config: cfg,
			Issuer: accentedIssuer(),
			Period: per,
			Customers: []*saft.Customer{
				{CustomerID: "C1", CompanyName: "Gonçalves & Simões Lda"},
				{CustomerID: "C2", CompanyName: "Consumidor final"},
			},
			Products: []*saft.Product{
				{ProductCode: "P1", ProductDesc: "Serviço de manutenção"},
			},
			TaxTableEntries: []*saft.TaxTableEntry{
				{TaxType: "IVA", TaxCountryRegion: "PT", TaxCode: "NOR", Description: "Taxa normal"},
			},
			SourceDocuments: []saft.SourceDocument{
				&saft.Invoice{
					InvoiceNo: "FT A/1",
					Lines: []*saft.Line{
						{DebitAmount: "10.50", CreditAmount: "20.25"},
						{CreditAmount: "5.5"},
					},
				},
				canceled,
				&saft.StockMovement{
					DocumentNumber: "GR A/1",
					Lines:          []*saft.Line{{Quantity: "2.000"}, {Quantity: "3"}},
				},
				&saft.WorkDocument{
					DocumentNumber: "PF A/1",
					Lines:          []*saft.Line{{CreditAmount: "135.00"}},
				},
				&saft.Payment{
					PaymentRefNo: "RG A/1",
					Lines:        []*saft.PaymentLine{{CreditAmount: "153.85"}},
				},
			},
		})
		require.NoError(t, err)

		expected, err := saft.MarshalXML(af, true)
		require.NoError(t, err)

		got := streamFromAuditFile(t, af, cfg, accentedIssuer(), per)
		assert.Equal(t, string(expected), string(got))
	})

	t.Run("no master files and no documents", func(t *testing.T) {
		af, err := saft.BuildAuditFile(&saft.AuditFileParams{
			Config: cfg,
			Issuer: accentedIssuer(),
			Period: per,
		})
		require.NoError(t, err)

		expected, err := saft.MarshalXML(af, true)
		require.NoError(t, err)

		got := streamFromAuditFile(t, af, cfg, accentedIssuer(), per)
		assert.Equal(t, string(expected), string(got))
		assert.Contains(t, string(got), "<MasterFiles></MasterFiles>")
		assert.Contains(t, string(got), "<SourceDocuments></SourceDocuments>")
	})

	t.Run("documents without master files", func(t *testing.T) {
		af, err := saft.BuildAuditFile(&saft.AuditFileParams{
			Config: cfg,
			Issuer: accentedIssuer(),
			Period: per,
			SourceDocuments: []saft.SourceDocument{
				&saft.Invoice{
					InvoiceNo: "FT A/1",
					Lines:     []*saft.Line{{CreditAmount: "1.00"}},
				},
			},
		})
		require.NoError(t, err)

		expected, err := saft.MarshalXML(af, true)
		require.NoError(t, err)

		got := streamFromAuditFile(t, af, cfg, accentedIssuer(), per)
		assert.Equal(t, string(expected), string(got))
	})

	t.Run("master files without documents", func(t *testing.T) {
		af, err := saft.BuildAuditFile(&saft.AuditFileParams{
			Config:          cfg,
			Issuer:          accentedIssuer(),
			Period:          per,
			Customers:       []*saft.Customer{{CustomerID: "C1"}},
			TaxTableEntries: []*saft.TaxTableEntry{{TaxType: "IVA"}},
		})
		require.NoError(t, err)

		expected, err := saft.MarshalXML(af, true)
		require.NoError(t, err)

		got := streamFromAuditFile(t, af, cfg, accentedIssuer(), per)
		assert.Equal(t, string(expected), string(got))
	})

	t.Run("class with only canceled documents keeps zero totals", func(t *testing.T) {
		canceled := &saft.Invoice{
			InvoiceNo: "FT A/1",
			Lines:     []*saft.Line{{DebitAmount: "10.00", CreditAmount: "10.00"}},
		}
		canceled.SetStatus(saft.StatusCanceled, "Emitido por engano", "source")

		af, err := saft.BuildAuditFile(&saft.AuditFileParams{
			Config:          cfg,
			Issuer:          accentedIssuer(),
			Period:          per,
			SourceDocuments: []saft.SourceDocument{canceled},
		})
		require.NoError(t, err)

		expected, err := saft.MarshalXML(af, true)
		require.NoError(t, err)

		got := streamFromAuditFile(t, af, cfg, accentedIssuer(), per)
		assert.Equal(t, string(expected), string(got))
		assert.Contains(t, string(got), "<NumberOfEntries>1</NumberOfEntries>")
		assert.Contains(t, string(got), "<TotalDebit>0</TotalDebit>")
		assert.Contains(t, string(got), "<TotalCredit>0</TotalCredit>")
	})

	t.Run("newlines in text are escaped, so indenting a raw fragment is safe", func(t *testing.T) {
		af, err := saft.BuildAuditFile(&saft.AuditFileParams{
			Config: cfg,
			Issuer: accentedIssuer(),
			Period: per,
			Customers: []*saft.Customer{
				{CustomerID: "C1", CompanyName: "Line one\nLine two"},
			},
			SourceDocuments: []saft.SourceDocument{
				&saft.Invoice{
					InvoiceNo: "FT A/1",
					Lines: []*saft.Line{
						{Description: "first\nsecond", CreditAmount: "1.00"},
					},
				},
			},
		})
		require.NoError(t, err)

		expected, err := saft.MarshalXML(af, true)
		require.NoError(t, err)

		got := streamFromAuditFile(t, af, cfg, accentedIssuer(), per)
		assert.Equal(t, string(expected), string(got))

		// The encoder escapes newlines in character data, so the only newline bytes a
		// stored fragment carries are structural and may be safely indented.
		assert.Contains(t, string(got), "Line one&#xA;Line two")
		assert.Contains(t, string(got), "first&#xA;second")
	})

	t.Run("callers are expected to reject an invalid issuer", func(t *testing.T) {
		// The encoder trusts its parameters, so the checks it used to repeat live with
		// the caller that builds them.
		require.Error(t, saft.ValidateIssuer(&org.Party{Name: "No Tax ID Lda"}))
	})
}

// streamFromAuditFile streams a built audit file through the StreamEncoder using the
// marshalled fragments of its members, mimicking how stored fragments are streamed.
func streamFromAuditFile(t *testing.T, af *saft.AuditFile, cfg config.Config, issuer *org.Party, per cal.Period) []byte {
	t.Helper()

	var buf bytes.Buffer
	enc := saft.NewStreamEncoder(&buf, cfg, issuer, per)
	require.NoError(t, enc.Begin())

	if mf := af.MasterFiles; mf != nil {
		for _, c := range mf.Customers {
			require.NoError(t, enc.WriteMasterFile(marshalFrag(t, c)))
		}
		for _, pr := range mf.Products {
			require.NoError(t, enc.WriteMasterFile(marshalFrag(t, pr)))
		}
		if mf.TaxTable != nil {
			for _, te := range mf.TaxTable.TaxTableEntries {
				require.NoError(t, enc.WriteTaxTableEntry(marshalFrag(t, te)))
			}
		}
	}

	if sd := af.SourceDocuments; sd != nil {
		if si := sd.SalesInvoices; si != nil {
			require.NoError(t, enc.BeginClass(saft.DocumentClassInvoice, saft.ClassTotals{
				Entries: si.NumberOfEntries,
				Debit:   si.TotalDebit,
				Credit:  si.TotalCredit,
			}))
			for _, d := range si.Invoices {
				require.NoError(t, enc.WriteDocument(marshalFrag(t, d)))
			}
		}
		if mg := sd.MovementOfGoods; mg != nil {
			require.NoError(t, enc.BeginClass(saft.DocumentClassTransport, saft.ClassTotals{
				Entries:  len(mg.StockMovements),
				Lines:    mg.NumberOfMovementLines,
				Quantity: mg.TotalQuantityIssued,
			}))
			for _, d := range mg.StockMovements {
				require.NoError(t, enc.WriteDocument(marshalFrag(t, d)))
			}
		}
		if wd := sd.WorkingDocuments; wd != nil {
			require.NoError(t, enc.BeginClass(saft.DocumentClassWork, saft.ClassTotals{
				Entries: wd.NumberOfEntries,
				Debit:   wd.TotalDebit,
				Credit:  wd.TotalCredit,
			}))
			for _, d := range wd.WorkDocuments {
				require.NoError(t, enc.WriteDocument(marshalFrag(t, d)))
			}
		}
		if py := sd.Payments; py != nil {
			require.NoError(t, enc.BeginClass(saft.DocumentClassPayment, saft.ClassTotals{
				Entries: py.NumberOfEntries,
				Debit:   py.TotalDebit,
				Credit:  py.TotalCredit,
			}))
			for _, d := range py.Payments {
				require.NoError(t, enc.WriteDocument(marshalFrag(t, d)))
			}
		}
	}

	require.NoError(t, enc.End())
	return buf.Bytes()
}

func marshalFrag(t *testing.T, v any) []byte {
	t.Helper()
	data, err := saft.MarshalXML(v, false)
	require.NoError(t, err)
	return data
}

func accentedIssuer() *org.Party {
	return &org.Party{
		Name:  "Sabores de São João Lda",
		Alias: "Cantinho das Promoções",
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
	}
}
