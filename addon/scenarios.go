package addon

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
)

var scenarios = []*tax.ScenarioSet{
	invoiceScenarios,
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyInvoiceType: InvoiceTypeStandard,
			}),
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{tax.TagSimplified},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyInvoiceType: InvoiceTypeSimplified,
			}),
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Tags:  []cbc.Key{TagCashVAT},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyCashVAT: "1",
			}),
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Filter: func(doc any) bool {
				inv, ok := doc.(*bill.Invoice)
				if !ok {
					return false
				}
				if inv.HasTags(pt.TagInvoiceReceipt) {
					return true
				}
				return inv.Totals.Paid() && invoiceDocType(inv) != InvoiceTypeStandard
			},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyInvoiceType: InvoiceTypeInvoiceReceipt,
			}),
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeDebitNote},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyInvoiceType: InvoiceTypeDebitNote,
			}),
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyInvoiceType: InvoiceTypeCreditNote,
			}),
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeProforma},
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyWorkType: WorkTypeProforma,
			}),
		},
	},
}
