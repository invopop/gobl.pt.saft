package saft

import (
	"github.com/invopop/gobl/num"
)

// DocTotals holds the line count and line-amount sums of a single source document,
// computed with the same precision rules as the SAF-T report totals. Lines is counted
// for every class even though only MovementOfGoods reports it, as NumberOfMovementLines;
// the other three wrappers count documents, with NumberOfEntries.
type DocTotals struct {
	Lines    int
	Debit    num.Amount
	Credit   num.Amount
	Quantity num.Amount
}

// ComputeDocTotals sums the line amounts of a document.
func ComputeDocTotals(doc SourceDocument) DocTotals {
	var t DocTotals
	switch d := doc.(type) {
	case *Invoice:
		t.Lines = len(d.Lines)
		for _, line := range d.Lines {
			t.Debit = addDocAmount(t.Debit, line.DebitAmount)
			t.Credit = addDocAmount(t.Credit, line.CreditAmount)
		}
	case *WorkDocument:
		t.Lines = len(d.Lines)
		for _, line := range d.Lines {
			t.Debit = addDocAmount(t.Debit, line.DebitAmount)
			t.Credit = addDocAmount(t.Credit, line.CreditAmount)
		}
	case *Payment:
		t.Lines = len(d.Lines)
		for _, line := range d.Lines {
			t.Debit = addDocAmount(t.Debit, line.DebitAmount)
			t.Credit = addDocAmount(t.Credit, line.CreditAmount)
		}
	case *StockMovement:
		t.Lines = len(d.Lines)
		for _, line := range d.Lines {
			t.Quantity = addDocAmount(t.Quantity, line.Quantity)
		}
	}
	return t
}

// Merge returns the result of adding the other document totals to these ones.
func (t DocTotals) Merge(o DocTotals) DocTotals {
	return DocTotals{
		Lines:    t.Lines + o.Lines,
		Debit:    t.Debit.MatchPrecision(o.Debit).Add(o.Debit),
		Credit:   t.Credit.MatchPrecision(o.Credit).Add(o.Credit),
		Quantity: t.Quantity.MatchPrecision(o.Quantity).Add(o.Quantity),
	}
}

// addDocAmount parses a SAF-T amount string and adds it to the running total, ignoring
// empty and unparseable values.
func addDocAmount(total num.Amount, s string) num.Amount {
	if s == "" {
		return total
	}
	a, err := num.AmountFromString(s)
	if err != nil {
		return total
	}
	a = a.RescaleDown(8) // Cap to avoid overflows.
	return total.MatchPrecision(a).Add(a)
}
