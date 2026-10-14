// Package goblwrap provides wrappers to handle GOBL types generically.
package goblwrap

import (
	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

type invoiceWrap bill.Invoice
type paymentWrap bill.Payment
type orderWrap bill.Order
type deliveryWrap bill.Delivery

// For wraps a GOBL type in a type with accessors to handle it generically.
func For(doc any) any {
	switch doc := doc.(type) {
	case *bill.Invoice:
		return (*invoiceWrap)(doc)
	case *bill.Payment:
		return (*paymentWrap)(doc)
	case *bill.Order:
		return (*orderWrap)(doc)
	case *bill.Delivery:
		return (*deliveryWrap)(doc)
	default:
		return nil
	}
}

// Unwrap returns the actual GOBL type from a wrapped type.
func Unwrap(doc any) any {
	switch doc := doc.(type) {
	case *invoiceWrap:
		return (*bill.Invoice)(doc)
	case *paymentWrap:
		return (*bill.Payment)(doc)
	case *orderWrap:
		return (*bill.Order)(doc)
	case *deliveryWrap:
		return (*bill.Delivery)(doc)
	default:
		return doc
	}
}

// Extract returns the content of a GOBL envelope wrapped.
func Extract(env *gobl.Envelope) any {
	return For(env.Extract())
}

// Accessors for invoice
func (i *invoiceWrap) GetRegimeDef() *tax.RegimeDef     { return i.RegimeDef() }
func (i *invoiceWrap) GetSupplier() *org.Party          { return i.Supplier }
func (i *invoiceWrap) GetCustomer() *org.Party          { return i.Customer }
func (i *invoiceWrap) GetCode() cbc.Code                { return i.Code }
func (i *invoiceWrap) SetCode(code cbc.Code)            { i.Code = code }
func (i *invoiceWrap) GetSeries() cbc.Code              { return i.Series }
func (i *invoiceWrap) GetIssueDate() cal.Date           { return i.IssueDate }
func (i *invoiceWrap) GetPreceding() []*org.DocumentRef { return i.Preceding }
func (i *invoiceWrap) GetBillLines() []*bill.Line       { return i.Lines }
func (i *invoiceWrap) GetBillTotals() *bill.Totals      { return i.Totals }
func (i *invoiceWrap) GetCurrency() currency.Code       { return i.Currency }
func (i *invoiceWrap) GetAddons() []cbc.Key             { return i.Addons.GetAddons() }
func (i *invoiceWrap) GetExt() tax.Extensions {
	if i.Tax == nil {
		return tax.Extensions{}
	}
	return i.Tax.Ext
}

// Accessors for payment
func (p *paymentWrap) GetRegimeDef() *tax.RegimeDef     { return p.RegimeDef() }
func (p *paymentWrap) GetSupplier() *org.Party          { return p.Supplier }
func (p *paymentWrap) GetCustomer() *org.Party          { return p.Customer }
func (p *paymentWrap) GetCode() cbc.Code                { return p.Code }
func (p *paymentWrap) SetCode(code cbc.Code)            { p.Code = code }
func (p *paymentWrap) GetSeries() cbc.Code              { return p.Series }
func (p *paymentWrap) GetIssueDate() cal.Date           { return p.IssueDate }
func (p *paymentWrap) GetPreceding() []*org.DocumentRef { return nil }
func (p *paymentWrap) GetBillLines() []*bill.Line       { return nil }
func (p *paymentWrap) GetBillTotals() *bill.Totals      { return nil }
func (p *paymentWrap) GetCurrency() currency.Code       { return p.Currency }
func (p *paymentWrap) GetAddons() []cbc.Key             { return p.Addons.GetAddons() }
func (p *paymentWrap) GetExt() tax.Extensions           { return p.Ext }

// Accessors for order
func (o *orderWrap) GetRegimeDef() *tax.RegimeDef     { return o.RegimeDef() }
func (o *orderWrap) GetSupplier() *org.Party          { return o.Supplier }
func (o *orderWrap) GetCustomer() *org.Party          { return o.Customer }
func (o *orderWrap) GetCode() cbc.Code                { return o.Code }
func (o *orderWrap) SetCode(code cbc.Code)            { o.Code = code }
func (o *orderWrap) GetSeries() cbc.Code              { return o.Series }
func (o *orderWrap) GetIssueDate() cal.Date           { return o.IssueDate }
func (o *orderWrap) GetPreceding() []*org.DocumentRef { return nil }
func (o *orderWrap) GetBillLines() []*bill.Line       { return o.Lines }
func (o *orderWrap) GetBillTotals() *bill.Totals      { return o.Totals }
func (o *orderWrap) GetCurrency() currency.Code       { return o.Currency }
func (o *orderWrap) GetAddons() []cbc.Key             { return o.Addons.GetAddons() }
func (o *orderWrap) GetExt() tax.Extensions {
	if o.Tax == nil {
		return tax.Extensions{}
	}
	return o.Tax.Ext
}

// Accessors for delivery
func (d *deliveryWrap) GetRegimeDef() *tax.RegimeDef     { return d.RegimeDef() }
func (d *deliveryWrap) GetSupplier() *org.Party          { return d.Supplier }
func (d *deliveryWrap) GetCustomer() *org.Party          { return d.Customer }
func (d *deliveryWrap) GetCode() cbc.Code                { return d.Code }
func (d *deliveryWrap) SetCode(code cbc.Code)            { d.Code = code }
func (d *deliveryWrap) GetSeries() cbc.Code              { return d.Series }
func (d *deliveryWrap) GetIssueDate() cal.Date           { return d.IssueDate }
func (d *deliveryWrap) GetPreceding() []*org.DocumentRef { return d.Preceding }
func (d *deliveryWrap) GetBillLines() []*bill.Line       { return d.Lines }
func (d *deliveryWrap) GetBillTotals() *bill.Totals      { return d.Totals }
func (d *deliveryWrap) GetCurrency() currency.Code       { return d.Currency }
func (d *deliveryWrap) GetAddons() []cbc.Key             { return d.Addons.GetAddons() }
func (d *deliveryWrap) GetExt() tax.Extensions {
	if d.Tax == nil {
		return tax.Extensions{}
	}
	return d.Tax.Ext
}
