package saft_test

import (
	"testing"

	saft "github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNoteValidation(t *testing.T) {
	t.Run("nil note", func(t *testing.T) {
		var note *org.Note
		err := rules.Validate(note, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("invalid exemption note - too long", func(t *testing.T) {
		note := &org.Note{
			Key:  org.NoteKeyLegal,
			Src:  saft.ExtKeyExemption,
			Text: "1234567890123456789012345678901234567890123456789012345678901", // 61 chars
		}
		err := rules.Validate(note, withAddonContext())
		assert.ErrorContains(t, err, "the length must be between 6 and 60")
	})

	t.Run("invalid exemption note - too short", func(t *testing.T) {
		note := &org.Note{
			Key:  org.NoteKeyLegal,
			Src:  saft.ExtKeyExemption,
			Text: "12345",
		}
		err := rules.Validate(note, withAddonContext())
		assert.ErrorContains(t, err, "the length must be between 6 and 60")
	})

	t.Run("valid exemption note - min length", func(t *testing.T) {
		note := &org.Note{
			Key:  org.NoteKeyLegal,
			Src:  saft.ExtKeyExemption,
			Text: "123456",
		}

		err := rules.Validate(note, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("valid exemption note - max length", func(t *testing.T) {
		note := &org.Note{
			Key:  org.NoteKeyLegal,
			Src:  saft.ExtKeyExemption,
			Text: "123456789012345678901234567890123456789012345678901234567890", // 60 chars
		}

		err := rules.Validate(note, withAddonContext())
		assert.NoError(t, err)
	})

	t.Run("other note", func(t *testing.T) {
		note := &org.Note{
			Text: "1234567890123456789012345678901234567890123456789012345678901", // 61 chars
		}
		err := rules.Validate(note, withAddonContext())
		assert.NoError(t, err)
	})
}

func TestWithholdingNoteValidation(t *testing.T) {
	const text61 = "1234567890123456789012345678901234567890123456789012345678901" // 61 chars
	const text60 = "123456789012345678901234567890123456789012345678901234567890"  // 60 chars

	t.Run("nil note", func(t *testing.T) {
		var note *tax.Note
		err := rules.Validate(note, withAddonContext())
		assert.NoError(t, err)
	})

	tests := []struct {
		name    string
		cat     cbc.Code
		text    string
		wantErr bool
	}{
		{name: "IRS valid min length", cat: pt.TaxCategoryIRS, text: "1"},
		{name: "IRS valid max length", cat: pt.TaxCategoryIRS, text: text60},
		{name: "IRS invalid too long", cat: pt.TaxCategoryIRS, text: text61, wantErr: true},
		{name: "IRC valid min length", cat: pt.TaxCategoryIRC, text: "1"},
		{name: "IRC valid max length", cat: pt.TaxCategoryIRC, text: text60},
		{name: "IRC invalid too long", cat: pt.TaxCategoryIRC, text: text61, wantErr: true},
		{name: "VAT not affected", cat: tax.CategoryVAT, text: text61},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note := &tax.Note{
				Category: tt.cat,
				Text:     tt.text,
			}
			err := rules.Validate(note, withAddonContext())
			if tt.wantErr {
				assert.ErrorContains(t, err, "the length must be at most 60")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
