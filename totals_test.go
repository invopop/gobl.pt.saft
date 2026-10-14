package saft_test

import (
	"testing"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/stretchr/testify/assert"
)

func TestComputeDocTotals(t *testing.T) {
	t.Run("invoice", func(t *testing.T) {
		tot := saft.ComputeDocTotals(&saft.Invoice{
			Lines: []*saft.Line{
				{DebitAmount: "123.45", CreditAmount: "111.11"},
				{DebitAmount: "111.11"},
				{CreditAmount: "11.11"},
			},
		})
		assert.Equal(t, "234.56", tot.Debit.String())
		assert.Equal(t, "122.22", tot.Credit.String())
		assert.Equal(t, "0", tot.Quantity.String())
	})

	t.Run("work document", func(t *testing.T) {
		tot := saft.ComputeDocTotals(&saft.WorkDocument{
			Lines: []*saft.Line{
				{DebitAmount: "10.5", CreditAmount: "20.25"},
			},
		})
		assert.Equal(t, "10.5", tot.Debit.String())
		assert.Equal(t, "20.25", tot.Credit.String())
	})

	t.Run("payment", func(t *testing.T) {
		tot := saft.ComputeDocTotals(&saft.Payment{
			Lines: []*saft.PaymentLine{
				{DebitAmount: "543.21", CreditAmount: "500.00"},
				{CreditAmount: "3.85"},
			},
		})
		assert.Equal(t, "543.21", tot.Debit.String())
		assert.Equal(t, "503.85", tot.Credit.String())
	})

	t.Run("stock movement", func(t *testing.T) {
		tot := saft.ComputeDocTotals(&saft.StockMovement{
			Lines: []*saft.Line{
				{Quantity: "2.000"},
				{Quantity: "3"},
			},
		})
		assert.Equal(t, "5.000", tot.Quantity.String())
		assert.Equal(t, "0", tot.Debit.String())
		assert.Equal(t, "0", tot.Credit.String())
	})

	t.Run("empty and unparseable amounts are ignored", func(t *testing.T) {
		tot := saft.ComputeDocTotals(&saft.Invoice{
			Lines: []*saft.Line{
				{DebitAmount: "", CreditAmount: "not-a-number"},
				{DebitAmount: "1.00"},
			},
		})
		assert.Equal(t, "1.00", tot.Debit.String())
		assert.Equal(t, "0", tot.Credit.String())
	})

	t.Run("amounts beyond 8 decimals are rounded down to 8", func(t *testing.T) {
		tot := saft.ComputeDocTotals(&saft.Invoice{
			Lines: []*saft.Line{
				{DebitAmount: "1.0000000049"},
				{DebitAmount: "1.0000000051"},
			},
		})
		assert.Equal(t, "2.00000001", tot.Debit.String())
	})

	t.Run("status does not affect the totals", func(t *testing.T) {
		inv := &saft.Invoice{
			Lines: []*saft.Line{{DebitAmount: "1.00"}},
		}
		inv.SetStatus(saft.StatusCanceled, "reason", "source")
		tot := saft.ComputeDocTotals(inv)
		assert.Equal(t, "1.00", tot.Debit.String())
	})
}

func TestDocTotalsMerge(t *testing.T) {
	a := saft.ComputeDocTotals(&saft.Invoice{
		Lines: []*saft.Line{{DebitAmount: "10.00", CreditAmount: "5.5"}},
	})
	b := saft.ComputeDocTotals(&saft.Invoice{
		Lines: []*saft.Line{{DebitAmount: "5.5", CreditAmount: "10.00"}},
	})

	tot := a.Merge(b)
	assert.Equal(t, "15.50", tot.Debit.String())
	assert.Equal(t, "15.50", tot.Credit.String())

	tot = saft.DocTotals{}.Merge(a)
	assert.Equal(t, "10.00", tot.Debit.String())
}
