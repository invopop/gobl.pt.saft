package addon_test

import (
	"testing"

	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestItemValidation(t *testing.T) {
	tests := []struct {
		name string
		item *org.Item
		err  string
	}{
		{
			name: "valid item",
			item: &org.Item{
				Name: "Test Item",
				Unit: "kg",
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyProductType: "P",
				}),
			},
		},
		{
			name: "nil item",
			item: nil,
		},
		{
			name: "missing extensions",
			item: &org.Item{},
			err:  "product type is required",
		},
		{
			name: "empty extensions",
			item: &org.Item{
				Ext: tax.ExtensionsOf(cbc.CodeMap{}),
			},
			err: "product type is required",
		},
		{
			name: "missing extension",
			item: &org.Item{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					"random": "12345678",
				}),
			},
			err: "product type is required",
		},
		{
			name: "missing unit",
			item: &org.Item{},
			err:  "cannot be blank",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			err := rules.Validate(ts.item, withAddonContext())
			if ts.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, ts.err)
			}
		})
	}
}

func TestItemExtProductTypeNormalization(t *testing.T) {
	tests := []struct {
		name string
		item *org.Item
		out  cbc.Code
	}{
		{
			name: "extension present",
			item: &org.Item{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					addon.ExtKeyProductType: "P",
				}),
			},
			out: "P",
		},
		{
			name: "nil item",
			item: nil,
		},
		{
			name: "empty extensions",
			item: &org.Item{
				Ext: tax.ExtensionsOf(cbc.CodeMap{}),
			},
			out: "S",
		},
		{
			name: "missing extension",
			item: &org.Item{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					"random": "12345678",
				}),
			},
			out: "S",
		},
		{
			name: "goods unit set",
			item: &org.Item{
				Unit: "kg",
			},
			out: "P",
		},
		{
			name: "service unit set",
			item: &org.Item{
				Unit: "service",
			},
			out: "S",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			norm.Normalize(tt.item, tax.AddonContext(addon.V1))
			if tt.item != nil {
				assert.Equal(t, tt.out, tt.item.Ext.Get(addon.ExtKeyProductType))
			}
		})
	}
}

func TestItemUnitNormalization(t *testing.T) {
	tests := []struct {
		name string
		item *org.Item
		out  cbc.Key
		ext  cbc.Code
	}{
		{
			name: "unit present",
			item: &org.Item{
				Unit: "kg",
			},
			out: "kg",
		},
		{
			name: "nil item",
			item: nil,
		},
		{
			name: "unit not present",
			item: &org.Item{},
			out:  "one",
		},
		{
			name: "legacy UN/ECE unit code",
			item: &org.Item{
				Unit: "KGM",
			},
			out: "kg",
			ext: "KGM",
		},
		{
			name: "UN/ECE unit code in extension",
			item: &org.Item{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyUnit: "HUR",
				}),
			},
			out: "h",
			ext: "HUR",
		},
		{
			name: "unknown UN/ECE unit code in extension",
			item: &org.Item{
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyUnit: "ZZ",
				}),
			},
			out: "one",
			ext: "ZZ",
		},
	}

	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			norm.Normalize(ts.item, tax.AddonContext(addon.V1))
			if ts.out != "" {
				assert.Equal(t, ts.out, ts.item.Unit)
			}
			if ts.ext != "" {
				assert.Equal(t, ts.ext, ts.item.Ext.Get(untdid.ExtKeyUnit))
			}
		})
	}
}

func TestItemProductTypeFromUNTDIDUnit(t *testing.T) {
	// A UN/ECE code moved into the extension must still drive the product
	// type, or goods would be reported to the AT as services.
	item := &org.Item{Unit: "KGM"}
	norm.Normalize(item, tax.AddonContext(addon.V1))
	assert.Equal(t, cbc.Code("P"), item.Ext.Get(addon.ExtKeyProductType))

	item = &org.Item{Unit: "HUR"}
	norm.Normalize(item, tax.AddonContext(addon.V1))
	assert.Equal(t, cbc.Code("S"), item.Ext.Get(addon.ExtKeyProductType))
}
