package saft_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl"
	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl.pt.saft/test"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeal(t *testing.T) {
	config.Now = func() time.Time { return time.Date(2024, 8, 3, 0, 0, 0, 0, pt.New().TimeLocation()) }
	certPath := test.Path("test", "certs", "signing_key.pem")
	t.Run("General case (invoice)", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		sig, err := saft.Seal(env, nil, cal.Date{}, "AAJFJ3N5YS", &config.Config{
			CertificateNumber: 456,
			SigningKeyPath:    certPath,
		})
		require.NoError(t, err)

		assert.Equal(t, "tcdL", sig.HashCharacters())
		assertStampContains(t, env, pt.StampProviderATATCUD, "AAJFJ3N5YS-1")
		assertStampContains(t, env, pt.StampProviderATHash, "tcdL")
		assertStampContains(t, env, pt.StampProviderATHashFull, "tY9G1OXXmKcSZN6as3KOdG")
		assertStampContains(t, env, pt.StampProviderATQR, "A:770013813*B:514329874")
		assertStampContains(t, env, pt.StampProviderATAppID, "456")
		assertStampContains(t, env, pt.StampProviderATTimestamp, "2024-08-03T00:00:00")

		inv, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		assert.Equal(t, "1", inv.Code.String())
	})

	t.Run("General case (work document from invoice)", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-proforma.json")
		require.NoError(t, err)

		sig, err := saft.Seal(env, nil, cal.Date{}, "AAJFJ3N5YS", &config.Config{
			CertificateNumber: 456,
			SigningKeyPath:    certPath,
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		assert.Equal(t, "h9KL", sig.HashCharacters())
		assertStampContains(t, env, pt.StampProviderATATCUD, "AAJFJ3N5YS-1")
		assertStampContains(t, env, pt.StampProviderATHash, "h9KL")
		assertStampContains(t, env, pt.StampProviderATHashFull, "hSvbIijMUx9IMR9x9NpZKVl1uKCCEvLddGhH1t")
		assertStampContains(t, env, pt.StampProviderATQR, "A:770013813*B:514329874")
		assertStampContains(t, env, pt.StampProviderATAppID, "456")
		assertStampContains(t, env, pt.StampProviderATTimestamp, "2024-08-03T00:00:00")

		inv, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		assert.Equal(t, "1", inv.Code.String())
	})

	t.Run("General case (work document from order)", func(t *testing.T) {
		env, err := test.LoadEnvelope("purchase-order.json")
		require.NoError(t, err)

		sig, err := saft.Seal(env, nil, cal.Date{}, "AAJFJ3N5YS", &config.Config{
			CertificateNumber: 456,
			SigningKeyPath:    certPath,
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		assert.Equal(t, "AoLo", sig.HashCharacters())
		assertStampContains(t, env, pt.StampProviderATATCUD, "AAJFJ3N5YS-1")
		assertStampContains(t, env, pt.StampProviderATHash, "AoLo")
		assertStampContains(t, env, pt.StampProviderATHashFull, "AFJVoiMNHeoJKIbI7bifLgm+vP/oZw")
		assertStampContains(t, env, pt.StampProviderATQR, "A:770013813*B:514329874")
		assertStampContains(t, env, pt.StampProviderATAppID, "456")
		assertStampContains(t, env, pt.StampProviderATTimestamp, "2024-08-03T00:00:00")

		ord, ok := env.Extract().(*bill.Order)
		require.True(t, ok)
		assert.Equal(t, "1", ord.Code.String())
	})

	t.Run("General case (payment)", func(t *testing.T) {
		env, err := test.LoadEnvelope("receipt.json")
		require.NoError(t, err)

		sig, err := saft.Seal(env, nil, cal.Date{}, "AAJFJ3N5YS", &config.Config{
			CertificateNumber: 456,
			SigningKeyPath:    certPath,
		})
		require.NoError(t, err)

		err = test.SignEnvelope(env)
		require.NoError(t, err)

		assert.Equal(t, "cVoQ", sig.HashCharacters())
		assertStampContains(t, env, pt.StampProviderATATCUD, "AAJFJ3N5YS-1")
		assertStampContains(t, env, pt.StampProviderATHash, "cVoQ")
		assertStampContains(t, env, pt.StampProviderATHashFull, "cm0STGocAZV4UiYBweH3oI0ghaBI4KQQ5d9J7")
		assertStampContains(t, env, pt.StampProviderATQR, "A:770013813*B:514329874")
		assertStampContains(t, env, pt.StampProviderATAppID, "456")
		assertStampContains(t, env, pt.StampProviderATTimestamp, "2024-08-03T00:00:00")

		rec, ok := env.Extract().(*bill.Payment)
		require.True(t, ok)
		assert.Equal(t, "1", rec.Code.String())
	})

	t.Run("General case (delivery)", func(t *testing.T) {
		env, err := test.LoadEnvelope("delivery-note.json")
		require.NoError(t, err)

		sig, err := saft.Seal(env, nil, cal.Date{}, "AAJFJ3N5YS", &config.Config{
			CertificateNumber: 456,
			SigningKeyPath:    certPath,
		})

		require.NoError(t, err)

		assert.Equal(t, "P+rn", sig.HashCharacters())
		assertStampContains(t, env, pt.StampProviderATATCUD, "AAJFJ3N5YS-1")
		assertStampContains(t, env, pt.StampProviderATHash, "P+rn")
		assertStampContains(t, env, pt.StampProviderATHashFull, "PsCFWLbKSs+NXRCU542kreM+wV")
		assertStampContains(t, env, pt.StampProviderATQR, "A:770013813*B:514329874")
		assertStampContains(t, env, pt.StampProviderATAppID, "456")
		assertStampContains(t, env, pt.StampProviderATTimestamp, "2024-08-03T00:00:00")

		dlv, ok := env.Extract().(*bill.Delivery)
		require.True(t, ok)
		assert.Equal(t, "1", dlv.Code.String())
	})
}

func TestCode(t *testing.T) {
	config.Now = func() time.Time { return time.Date(2024, 8, 3, 0, 0, 0, 0, time.UTC) }
	certPath := test.Path("test", "certs", "signing_key.pem")

	t.Run("No previous signature", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJ3N5YS", &config.Config{
			CertificateNumber: 456,
			SigningKeyPath:    certPath,
		})
		require.NoError(t, err)

		inv, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		assert.Equal(t, "1", inv.Code.String())
	})

	t.Run("Previous signature", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		sig := &saft.Signature{
			IssueDate:  "2024-08-01",
			DocumentNo: "FT SERIES-A/123",
		}

		_, err = saft.Seal(env, sig, cal.Date{}, "AAJFJ3N5YS", &config.Config{
			CertificateNumber: 456,
			SigningKeyPath:    certPath,
		})
		require.NoError(t, err)

		inv, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)
		assert.Equal(t, "124", inv.Code.String())
	})

	t.Run("Code given", func(t *testing.T) {
		env, err := test.LoadEnvelope("invoice-b2b.json")
		require.NoError(t, err)

		inv, ok := env.Extract().(*bill.Invoice)
		require.True(t, ok)

		inv.Code = "123"

		_, err = saft.Seal(env, nil, cal.Date{}, "AAJFJ3N5YS", &config.Config{
			CertificateNumber: 456,
			SigningKeyPath:    certPath,
		})
		assert.ErrorContains(t, err, "code: must not be set")
	})
}

func TestIssueDateValidation(t *testing.T) {
	certPath := test.Path("test", "certs", "signing_key.pem")

	// The test envelope is issued on 2024-08-03
	tests := []struct {
		name  string
		now   time.Time
		start cal.Date
		prev  *saft.Signature
		src   cbc.Code
		err   string
	}{
		{
			name: "Issue date after today",
			now:  time.Date(2024, 8, 2, 22, 59, 59, 999999999, time.UTC),
			err:  "issue_date: must not be after today in Europe/Lisbon (2024-08-02)",
		},
		{
			name: "Issue date is today (beginning of day)",
			now:  time.Date(2024, 8, 2, 23, 0, 0, 0, time.UTC),
			err:  "",
		},
		{
			name: "Issue date is today (end of day)",
			now:  time.Date(2024, 8, 3, 22, 59, 59, 999999999, time.UTC),
			err:  "",
		},
		{
			name: "Issue date before today",
			now:  time.Date(2024, 8, 3, 23, 0, 0, 0, time.UTC),
			err:  "",
		},
		{
			name:  "Issue date after series start date",
			now:   time.Date(2024, 8, 10, 12, 0, 0, 0, time.UTC),
			start: cal.MakeDate(2024, 8, 1),
			err:   "",
		},
		{
			name:  "Issue date same as series start date",
			now:   time.Date(2024, 8, 10, 12, 0, 0, 0, time.UTC),
			start: cal.MakeDate(2024, 8, 3),
			err:   "",
		},
		{
			name:  "Issue date before series start date",
			now:   time.Date(2024, 8, 10, 12, 0, 0, 0, time.UTC),
			start: cal.MakeDate(2024, 8, 4),
			err:   "issue_date: must not be before the series start date (2024-08-04)",
		},
		{
			name:  "Issue date after previous document",
			now:   time.Date(2024, 8, 10, 12, 0, 0, 0, time.UTC),
			start: cal.MakeDate(2024, 8, 1),
			prev:  &saft.Signature{IssueDate: "2024-08-02", DocumentNo: "FT SERIES-A/1"},
			err:   "",
		},
		{
			name:  "Issue date same as previous document",
			now:   time.Date(2024, 8, 10, 12, 0, 0, 0, time.UTC),
			start: cal.MakeDate(2024, 8, 1),
			prev:  &saft.Signature{IssueDate: "2024-08-03", DocumentNo: "FT SERIES-A/1"},
			err:   "",
		},
		{
			name:  "Issue date before previous document",
			now:   time.Date(2024, 8, 10, 12, 0, 0, 0, time.UTC),
			start: cal.MakeDate(2024, 8, 1),
			prev:  &saft.Signature{IssueDate: "2024-08-04", DocumentNo: "FT SERIES-A/1"},
			err:   "issue_date: must not be before the previous document's issue date (2024-08-04)",
		},
		{
			name: "Invalid previous document issue date",
			now:  time.Date(2024, 8, 10, 12, 0, 0, 0, time.UTC),
			prev: &saft.Signature{IssueDate: "03/08/2024", DocumentNo: "FT SERIES-A/1"},
			err:  "issue_date: invalid previous document's issue date (03/08/2024)",
		},
		{
			name: "Issue date after today - manual",
			now:  time.Date(2024, 8, 2, 22, 59, 59, 999999999, time.UTC),
			src:  "FTM abc/123",
			err:  "",
		},
		{
			name: "Issue date before today - manual",
			now:  time.Date(2024, 8, 3, 23, 0, 0, 0, time.UTC),
			src:  "FTM abc/123",
			err:  "",
		},
		{
			name:  "Issue date before series start date - manual",
			now:   time.Date(2024, 8, 10, 12, 0, 0, 0, time.UTC),
			start: cal.MakeDate(2024, 8, 4),
			src:   "FTM abc/123",
			err:   "",
		},
		{
			name: "Issue date before previous document - manual",
			now:  time.Date(2024, 8, 10, 12, 0, 0, 0, time.UTC),
			prev: &saft.Signature{IssueDate: "2024-08-04", DocumentNo: "FT SERIES-A/1"},
			src:  "FTM abc/123",
			err:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env, err := test.LoadEnvelope("invoice-b2b.json")
			require.NoError(t, err)

			inv, ok := env.Extract().(*bill.Invoice)
			require.True(t, ok)
			if tt.src != "" {
				inv.Tax.Ext = inv.Tax.Ext.Set(addon.ExtKeySourceRef, tt.src)
			}

			config.Now = func() time.Time { return tt.now }

			_, err = saft.Seal(env, tt.prev, tt.start, "AAJFJ3N5YS", &config.Config{
				CertificateNumber: 456,
				SigningKeyPath:    certPath,
			})

			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.err)
			}
		})
	}
}

func assertStampContains(t *testing.T, env *gobl.Envelope, prv cbc.Key, val string) {
	stamp := env.Head.GetStamp(prv)
	if assert.NotNil(t, stamp) {
		assert.Contains(t, stamp.Value, val)
	}
}
