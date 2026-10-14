package saft_test

import (
	"testing"
	"time"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl.pt.saft/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildSignature(t *testing.T) {
	prev := &saft.Signature{
		Hash: "aX1zuk8FOTLY2UCFH8oufJbOF8Icapxb3OACySJUTbZHqSrmrhbA/L/Zt+uD3M4/oQJy/sv+7Do1K35DQ7gc6BIOARdLh/3ORo67sFFBCLqM9IOztq5WxgX7ezJXow+Gz+J2XpKk1TUptSORgdBoT05BA3spMAt/QWUDCxG78hM=", //nolint:goconst
	}
	config.Now = func() time.Time { return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC) }

	pkey, err := test.LoadFile("test", "certs", "signing_key.pem")
	require.NoError(t, err)

	t.Run("General case (Invoice)", func(t *testing.T) {
		doc := &saft.Invoice{
			InvoiceDate: "2021-01-01",
			InvoiceNo:   "FT TEST-1/33",
			DocumentTotals: &saft.DocumentTotals{
				GrossTotal: "1234.56",
			},
		}

		sig, err := saft.BuildSignature(doc, prev, pkey)

		require.NoError(t, err)
		assert.Equal(t, "2021-01-01", sig.IssueDate)
		assert.Equal(t, "2021-01-01T00:00:00", sig.SystemEntryDate)
		assert.Equal(t, "FT TEST-1/33", sig.DocumentNo)
		assert.Equal(t, "1234.56", sig.GrossTotal)
		assert.Equal(t, prev.Hash, sig.PreviousHash)

		assert.Equal(t, "2021-01-01;2021-01-01T00:00:00;FT TEST-1/33;1234.56;aX1zuk8FOTLY2UCFH8oufJbOF8Icapxb3OACySJUTbZHqSrmrhbA/L/Zt+uD3M4/oQJy/sv+7Do1K35DQ7gc6BIOARdLh/3ORo67sFFBCLqM9IOztq5WxgX7ezJXow+Gz+J2XpKk1TUptSORgdBoT05BA3spMAt/QWUDCxG78hM=", sig.Msg())
		assert.Equal(t, "05WGKniCIC85QpCPDkb80pYdYQuaOFC/MQo+xBgHv0zxHln+Uuo4uH2DXISXkHGJRrCCBsKv80xA6F8SpVBGfsMeh+Ev9dtrgVwDzftUNM2ZnLFnukcX4N3l45dv92ljJYHIW6eIC9dGzf0xlVBGIH5HRkwOh1zCRiDET2vk8NQ=", sig.Hash)
		assert.Equal(t, "080C", sig.HashCharacters())
	})

	t.Run("No previous signature (Invoice)", func(t *testing.T) {
		doc := &saft.Invoice{
			InvoiceDate: "2021-01-01",
			InvoiceNo:   "FT TEST-1/33",
			DocumentTotals: &saft.DocumentTotals{
				GrossTotal: "1234.56",
			},
		}

		sig, err := saft.BuildSignature(doc, nil, pkey)
		require.NoError(t, err)

		assert.Equal(t, "2021-01-01;2021-01-01T00:00:00;FT TEST-1/33;1234.56;", sig.Msg())
	})

	t.Run("General case (Payment)", func(t *testing.T) {
		doc := &saft.Payment{
			TransactionDate: "2021-01-01",
			PaymentRefNo:    "RG TEST-1/33",
			DocumentTotals: &saft.DocumentTotals{
				GrossTotal: "1234.56",
			},
		}

		sig, err := saft.BuildSignature(doc, prev, pkey)
		require.NoError(t, err)

		assert.Equal(t, "2021-01-01", sig.IssueDate)
		assert.Equal(t, "2021-01-01T00:00:00", sig.SystemEntryDate)
		assert.Equal(t, "RG TEST-1/33", sig.DocumentNo)
		assert.Equal(t, "1234.56", sig.GrossTotal)
		assert.Equal(t, prev.Hash, sig.PreviousHash)

		assert.Equal(t, "2021-01-01;2021-01-01T00:00:00;RG TEST-1/33;1234.56;aX1zuk8FOTLY2UCFH8oufJbOF8Icapxb3OACySJUTbZHqSrmrhbA/L/Zt+uD3M4/oQJy/sv+7Do1K35DQ7gc6BIOARdLh/3ORo67sFFBCLqM9IOztq5WxgX7ezJXow+Gz+J2XpKk1TUptSORgdBoT05BA3spMAt/QWUDCxG78hM=", sig.Msg())
		assert.Equal(t, "LJhL8qxuf8g8NrtT1Ve3XqDNcQUWJARLda1DFmP94nvVovlZsFH1gkPA8jH2r3UqhS+0D6fYR4Yg0MNqekdg+JMiIw1Vn0KboaH2UoIc7cibZBe5LxDtgyDCILmj0ql/0tFrt6BX2uWUFybV8awoiJtzTgvwjV0OvYPxgBLONJ0=", sig.Hash)
		assert.Equal(t, "LgXR", sig.HashCharacters())
	})

	t.Run("General case (Stock Movement)", func(t *testing.T) {
		doc := &saft.StockMovement{
			MovementDate:   "2021-01-01",
			DocumentNumber: "GR TEST-1/33",
			DocumentTotals: &saft.DocumentTotals{
				GrossTotal: "1234.56",
			},
		}

		sig, err := saft.BuildSignature(doc, prev, pkey)
		require.NoError(t, err)

		assert.Equal(t, "2021-01-01", sig.IssueDate)
		assert.Equal(t, "2021-01-01T00:00:00", sig.SystemEntryDate)
		assert.Equal(t, "GR TEST-1/33", sig.DocumentNo)
		assert.Equal(t, "1234.56", sig.GrossTotal)
		assert.Equal(t, prev.Hash, sig.PreviousHash)
		assert.Equal(t, "HQS7rXwOD8LOyAb3sfT7FAbtE7k6hQGyk3xRK24UTPrk6CqhM7i4s40UqVSTsKCsA8DwiIA2IWdUyvYiSXOlplko5HR2/3Lyx0+h8SbtPGY/9utXT9jxUA7/kSIWUwV/sy4QK+O97vJ2XKMnsXW/vyLG2LKmUJxlUmaQLKGMpnU=", sig.Hash)
		assert.Equal(t, "HLFG", sig.HashCharacters())
	})

	t.Run("General case (Work Document)", func(t *testing.T) {
		doc := &saft.WorkDocument{
			WorkDate:       "2021-01-01",
			DocumentNumber: "PF TEST-1/25",
			DocumentTotals: &saft.WorkDocumentTotals{
				GrossTotal: "1234.56",
			},
		}

		sig, err := saft.BuildSignature(doc, prev, pkey)
		require.NoError(t, err)

		assert.Equal(t, "2021-01-01", sig.IssueDate)
		assert.Equal(t, "2021-01-01T00:00:00", sig.SystemEntryDate)
		assert.Equal(t, "PF TEST-1/25", sig.DocumentNo)
		assert.Equal(t, "1234.56", sig.GrossTotal)
		assert.Equal(t, prev.Hash, sig.PreviousHash)
		assert.Equal(t, "2021-01-01;2021-01-01T00:00:00;PF TEST-1/25;1234.56;aX1zuk8FOTLY2UCFH8oufJbOF8Icapxb3OACySJUTbZHqSrmrhbA/L/Zt+uD3M4/oQJy/sv+7Do1K35DQ7gc6BIOARdLh/3ORo67sFFBCLqM9IOztq5WxgX7ezJXow+Gz+J2XpKk1TUptSORgdBoT05BA3spMAt/QWUDCxG78hM=", sig.Msg())
		assert.Equal(t, "VS1f", sig.HashCharacters())
	})
}
