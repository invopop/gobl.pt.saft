package saft_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/invopop/gobl"
	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl.pt.saft/test"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// Register the Portuguese SAF-T addon so example documents declaring the
	// pt-saft-v1 addon normalize and validate.
	_ "github.com/invopop/gobl.pt.saft/addon"

	"github.com/invopop/gobl/pkg/examples"
)

// TestExamples converts every document under examples/ to a calculated,
// validated JSON envelope and compares it against its golden output, using the
// shared GOBL example helpers. Run with -update to (re)generate the goldens.
func TestExamples(t *testing.T) {
	examples.Run(t, "examples", *test.UpdateOut)
}

// TestConvertExamples seals every envelope under test/data, builds a SAF-T audit
// file with all of them and compares it against test/data/out/saft.xml. Run with
// -update to validate the output against the XSD and (re)generate the golden file.
func TestConvertExamples(t *testing.T) {
	config.Now = func() time.Time {
		return time.Date(2024, 8, 3, 0, 0, 0, 0, pt.New().TimeLocation())
	}

	// Load all test envelopes
	envs, err := test.LoadEnvelopes()
	require.NoError(t, err)

	// Seal the envelopes linking them together
	sigs := make(map[cbc.Code]*saft.Signature) // to store the last signature for each series
	for _, env := range envs {
		series := extractSeries(env)
		sigs[series], err = saft.Seal(env, sigs[series], "AAJFJHNC82", &config.Config{
			CertificateNumber: 123,
			SigningKeyPath:    test.Path("test", "certs", "signing_key.pem"),
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)
	}

	// Extract the supplier from the first invoice
	inv, ok := envs[0].Extract().(*bill.Invoice)
	require.True(t, ok)
	sup := inv.Supplier

	// Generate the SAF-T report
	per := cal.Period{
		Start: cal.NewDate(2024, 8, 1),
		End:   cal.NewDate(2024, 8, 3),
	}
	cfg := &config.Config{
		ProducerTaxID:     "999999999",
		ProducerName:      "Software Company",
		CertificateNumber: 123,
		AppName:           "App Name",
		AppVersion:        "1.0",
		SigningKeyVersion: "3",
	}
	srcID := "SOURCE_ID"

	r, err := test.GenerateAuditFile(sup, per, cfg, srcID, envs...)
	require.NoError(t, err)

	// Marshal the report
	data, err := saft.MarshalXML(r, true)
	require.NoError(t, err)

	// Stream the same report from marshalled fragments and check the output is identical
	streamed := streamReport(t, r, sup, per, cfg)
	require.Equal(t, string(data), string(streamed))

	outPath := test.Path("test", "data", "out", "saft.xml")

	if *test.UpdateOut {
		// Validate against the schema
		errs := r.ValidateXSD(context.Background())
		for _, e := range errs {
			assert.NoError(t, e)
		}
		if len(errs) > 0 {
			assert.Fail(t, "Invalid XML:\n"+string(data))
			return
		}

		// Write the XML file
		err = os.WriteFile(outPath, data, 0644)
		require.NoError(t, err)

		return
	}

	// Compare the output
	expected, err := os.ReadFile(outPath)
	require.False(t, os.IsNotExist(err), "output file %s missing, run tests with `--update` flag to create", filepath.Base(outPath))
	require.NoError(t, err)
	require.Equal(t, string(expected), string(data), "output file %s does not match, run tests with `--update` flag to update", filepath.Base(outPath))
}

// streamReport streams an already built audit file through the StreamEncoder using the
// marshalled fragments of its members, mimicking how stored fragments are streamed.
func streamReport(t *testing.T, af *saft.AuditFile, sup *org.Party, per cal.Period, cfg *config.Config) []byte {
	t.Helper()

	frag := func(v any) []byte {
		data, err := saft.MarshalXML(v, false)
		require.NoError(t, err)
		return data
	}

	var buf bytes.Buffer
	enc := saft.NewStreamEncoder(&buf, *cfg, sup, per)
	require.NoError(t, enc.Begin())

	if mf := af.MasterFiles; mf != nil {
		for _, c := range mf.Customers {
			require.NoError(t, enc.WriteMasterFile(frag(c)))
		}
		for _, p := range mf.Products {
			require.NoError(t, enc.WriteMasterFile(frag(p)))
		}
		if mf.TaxTable != nil {
			for _, te := range mf.TaxTable.TaxTableEntries {
				require.NoError(t, enc.WriteTaxTableEntry(frag(te)))
			}
		}
	}

	if sd := af.SourceDocuments; sd != nil {
		if si := sd.SalesInvoices; si != nil {
			require.NoError(t, enc.BeginClass(saft.DocumentClassInvoice, saft.ClassTotals{
				Entries: si.NumberOfEntries, Debit: si.TotalDebit, Credit: si.TotalCredit,
			}))
			for _, d := range si.Invoices {
				require.NoError(t, enc.WriteDocument(frag(d)))
			}
		}
		if mg := sd.MovementOfGoods; mg != nil {
			require.NoError(t, enc.BeginClass(saft.DocumentClassTransport, saft.ClassTotals{
				Entries: len(mg.StockMovements), Lines: mg.NumberOfMovementLines, Quantity: mg.TotalQuantityIssued,
			}))
			for _, d := range mg.StockMovements {
				require.NoError(t, enc.WriteDocument(frag(d)))
			}
		}
		if wd := sd.WorkingDocuments; wd != nil {
			require.NoError(t, enc.BeginClass(saft.DocumentClassWork, saft.ClassTotals{
				Entries: wd.NumberOfEntries, Debit: wd.TotalDebit, Credit: wd.TotalCredit,
			}))
			for _, d := range wd.WorkDocuments {
				require.NoError(t, enc.WriteDocument(frag(d)))
			}
		}
		if py := sd.Payments; py != nil {
			require.NoError(t, enc.BeginClass(saft.DocumentClassPayment, saft.ClassTotals{
				Entries: py.NumberOfEntries, Debit: py.TotalDebit, Credit: py.TotalCredit,
			}))
			for _, d := range py.Payments {
				require.NoError(t, enc.WriteDocument(frag(d)))
			}
		}
	}

	require.NoError(t, enc.End())
	return buf.Bytes()
}

func extractSeries(env *gobl.Envelope) cbc.Code {
	switch doc := env.Extract().(type) {
	case *bill.Invoice:
		return doc.Series
	case *bill.Payment:
		return doc.Series
	case *bill.Delivery:
		return doc.Series
	case *bill.Order:
		return doc.Series
	default:
		return cbc.CodeEmpty
	}
}
