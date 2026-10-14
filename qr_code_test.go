package saft_test

import (
	"testing"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/addon"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/pt"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestBuildQRCode(t *testing.T) {
	sup := &saft.Supplier{
		SupplierTaxID: "165902612",
	}
	cus := &saft.Customer{
		CustomerTaxID: "151040222",
		BillingAddress: &saft.Address{
			Country: "PT",
		},
	}
	taxTotal := func() *tax.Total {
		return &tax.Total{
			Categories: []*tax.CategoryTotal{
				{
					Code: tax.CategoryVAT,
					Rates: []*tax.RateTotal{
						{
							Ext: tax.ExtensionsOf(cbc.CodeMap{
								pt.ExtKeyRegion:     "PT",
								addon.ExtKeyTaxRate: "RED",
							}),
							Base:   num.MakeAmount(122222, 2),
							Amount: num.MakeAmount(1234, 2),
						},
					},
					Amount: num.MakeAmount(1234, 2),
				},
			},
		}
	}

	invoice := func() *saft.Invoice {
		return &saft.Invoice{
			InvoiceType: "FT",
			InvoiceNo:   "FT TEST-1/123",
			InvoiceDate: "2021-01-01",
			DocumentStatus: &saft.InvoiceStatus{
				InvoiceStatus: "N",
			},
			DocumentTotals: &saft.DocumentTotals{
				GrossTotal: "1234.56",
				TaxPayable: "12.34",
			},
			GOBLTaxTotal: taxTotal(),
		}
	}

	sig := &saft.Signature{
		Hash: "a---------b---------c---------d",
	}
	atcud := "AAJFJ3N5YS-123"
	swCertNum := 456

	t.Run("General case (Invoice)", func(t *testing.T) {
		doc := invoice()

		qr := saft.BuildQRCode(&saft.BuildQRCodeParams{
			Supplier:  sup,
			Customer:  cus,
			Document:  doc,
			ATCUD:     atcud,
			Signature: sig,
			SWCertNum: swCertNum,
		})

		assert.Equal(t, qr, "A:165902612*B:151040222*C:PT*D:FT*E:N*F:20210101*G:FT TEST-1/123*H:AAJFJ3N5YS-123*I1:PT*I3:1222.22*I4:12.34*N:12.34*O:1234.56*Q:abcd*R:456")
	})

	t.Run("Final customer", func(t *testing.T) {
		doc := invoice()
		cus := &saft.Customer{
			CustomerTaxID: "999999990",
			BillingAddress: &saft.Address{
				Country: "PT",
			},
		}

		qr := saft.BuildQRCode(&saft.BuildQRCodeParams{
			Supplier:  sup,
			Customer:  cus,
			Document:  doc,
			ATCUD:     atcud,
			Signature: sig,
			SWCertNum: swCertNum,
		})

		assert.Contains(t, qr, "B:999999990")
		assert.Contains(t, qr, "C:PT")
	})

	t.Run("VAT exempt", func(t *testing.T) {
		doc := invoice()
		rate := doc.GOBLTaxTotal.Categories[0].Rates[0]
		rate.Ext = rate.Ext.Set(addon.ExtKeyTaxRate, "ISE")
		rate.Amount = num.AmountZero

		qr := saft.BuildQRCode(&saft.BuildQRCodeParams{
			Supplier:  sup,
			Customer:  cus,
			Document:  doc,
			ATCUD:     atcud,
			Signature: sig,
			SWCertNum: swCertNum,
		})

		assert.Contains(t, qr, "I1:PT*I2:1222.22*N")
	})

	t.Run("VAT reduced", func(t *testing.T) {
		doc := invoice()
		rate := doc.GOBLTaxTotal.Categories[0].Rates[0]
		rate.Ext = rate.Ext.Set(addon.ExtKeyTaxRate, "RED")

		qr := saft.BuildQRCode(&saft.BuildQRCodeParams{
			Supplier:  sup,
			Customer:  cus,
			Document:  doc,
			ATCUD:     atcud,
			Signature: sig,
			SWCertNum: swCertNum,
		})

		assert.Contains(t, qr, "I1:PT*I3:1222.22*I4:12.34")
	})

	t.Run("VAT intermediate", func(t *testing.T) {
		doc := invoice()
		rate := doc.GOBLTaxTotal.Categories[0].Rates[0]
		rate.Ext = rate.Ext.Set(addon.ExtKeyTaxRate, "INT")

		qr := saft.BuildQRCode(&saft.BuildQRCodeParams{
			Supplier:  sup,
			Customer:  cus,
			Document:  doc,
			ATCUD:     atcud,
			Signature: sig,
			SWCertNum: swCertNum,
		})

		assert.Contains(t, qr, "I1:PT*I5:1222.22*I6:12.34")
	})

	t.Run("VAT normal", func(t *testing.T) {
		doc := invoice()
		rate := doc.GOBLTaxTotal.Categories[0].Rates[0]
		rate.Ext = rate.Ext.Set(addon.ExtKeyTaxRate, "NOR")

		qr := saft.BuildQRCode(&saft.BuildQRCodeParams{
			Supplier:  sup,
			Customer:  cus,
			Document:  doc,
			ATCUD:     atcud,
			Signature: sig,
			SWCertNum: swCertNum,
		})

		assert.Contains(t, qr, "I1:PT*I7:1222.22*I8:12.34")
	})

	t.Run("PT-AT region", func(t *testing.T) {
		doc := invoice()
		rate := doc.GOBLTaxTotal.Categories[0].Rates[0]
		rate.Ext = rate.Ext.Set(pt.ExtKeyRegion, "PT-AT")

		qr := saft.BuildQRCode(&saft.BuildQRCodeParams{
			Supplier:  sup,
			Customer:  cus,
			Document:  doc,
			ATCUD:     atcud,
			Signature: sig,
			SWCertNum: swCertNum,
		})

		assert.Contains(t, qr, "I1:PT-AT*I3:1222.22*I4:12.34*N")
	})

	t.Run("Multiple regions", func(t *testing.T) {
		doc := invoice()
		vat := doc.GOBLTaxTotal.Categories[0]
		vat.Rates = append(vat.Rates, &tax.RateTotal{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				addon.ExtKeyTaxRate: "NOR",
				pt.ExtKeyRegion:     "PT-AT",
			}),
			Base:   num.MakeAmount(10000, 2),
			Amount: num.MakeAmount(1000, 2),
		})

		qr := saft.BuildQRCode(&saft.BuildQRCodeParams{
			Supplier:  sup,
			Customer:  cus,
			Document:  doc,
			ATCUD:     atcud,
			Signature: sig,
			SWCertNum: swCertNum,
		})

		assert.Contains(t, qr, "J1:")
	})

	t.Run("No VAT", func(t *testing.T) {
		doc := invoice()
		doc.GOBLTaxTotal.Categories = nil

		qr := saft.BuildQRCode(&saft.BuildQRCodeParams{
			Supplier:  sup,
			Customer:  cus,
			Document:  doc,
			ATCUD:     atcud,
			Signature: sig,
			SWCertNum: swCertNum,
		})

		assert.Contains(t, qr, "I1:0")
	})

	t.Run("Multiple entries of same rate", func(t *testing.T) {
		doc := invoice()
		vat := doc.GOBLTaxTotal.Categories[0]
		vat.Rates = append(vat.Rates, []*tax.RateTotal{
			{
				Base:   num.MakeAmount(10000, 2),
				Amount: num.MakeAmount(1000, 2),
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					pt.ExtKeyRegion:     "PT",
					addon.ExtKeyTaxRate: "RED",
				}),
			},
			{
				Base: num.MakeAmount(10000, 2),
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					pt.ExtKeyRegion:     "PT",
					addon.ExtKeyTaxRate: "ISE",
				}),
			},
			{
				Base: num.MakeAmount(10000, 2),
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					pt.ExtKeyRegion:     "PT",
					addon.ExtKeyTaxRate: "ISE",
				}),
			},
		}...)

		qr := saft.BuildQRCode(&saft.BuildQRCodeParams{
			Supplier:  sup,
			Customer:  cus,
			Document:  doc,
			ATCUD:     atcud,
			Signature: sig,
			SWCertNum: swCertNum,
		})

		assert.Contains(t, qr, "I2:200.00*I3:1322.22*I4:22.34")
	})

	t.Run("General case (Payment)", func(t *testing.T) {
		doc := &saft.Payment{
			PaymentType:     "RG",
			PaymentRefNo:    "RG TEST-1/33",
			TransactionDate: "2021-01-01",
			DocumentStatus: &saft.PaymentStatus{
				PaymentStatus: "N",
			},
			DocumentTotals: &saft.DocumentTotals{
				GrossTotal: "1234.56",
				TaxPayable: "12.34",
			},
			GOBLTaxTotal: taxTotal(),
		}

		qr := saft.BuildQRCode(&saft.BuildQRCodeParams{
			Supplier:  sup,
			Customer:  cus,
			Document:  doc,
			ATCUD:     atcud,
			Signature: sig,
			SWCertNum: swCertNum,
		})

		assert.Equal(t, "A:165902612*B:151040222*C:PT*D:RG*E:N*F:20210101*G:RG TEST-1/33*H:AAJFJ3N5YS-123*I1:PT*I3:1222.22*I4:12.34*N:12.34*O:1234.56*Q:abcd*R:456", qr)
	})

	t.Run("General case (Stock Movement)", func(t *testing.T) {
		doc := &saft.StockMovement{
			MovementType:   "GR",
			DocumentNumber: "GR TEST-1/33",
			MovementDate:   "2021-01-01",
			DocumentStatus: &saft.MovementStatus{
				MovementStatus: "N",
			},
			DocumentTotals: &saft.DocumentTotals{
				GrossTotal: "1234.56",
				TaxPayable: "12.34",
			},
			GOBLTaxTotal: taxTotal(),
		}
		qr := saft.BuildQRCode(&saft.BuildQRCodeParams{
			Supplier:  sup,
			Customer:  cus,
			Document:  doc,
			ATCUD:     atcud,
			Signature: sig,
			SWCertNum: swCertNum,
		})

		assert.Equal(t, "A:165902612*B:151040222*C:PT*D:GR*E:N*F:20210101*G:GR TEST-1/33*H:AAJFJ3N5YS-123*I1:PT*I3:1222.22*I4:12.34*N:12.34*O:1234.56*Q:abcd*R:456", qr)
	})
	t.Run("General case (Work Document)", func(t *testing.T) {
		doc := &saft.WorkDocument{
			WorkType:       "PF",
			WorkDate:       "2021-01-01",
			DocumentNumber: "PF TEST-1/33",
			DocumentStatus: &saft.WorkStatus{
				WorkStatus: "N",
			},
			DocumentTotals: &saft.WorkDocumentTotals{
				GrossTotal: "1234.56",
				TaxPayable: "12.34",
			},
			GOBLTaxTotal: taxTotal(),
		}

		qr := saft.BuildQRCode(&saft.BuildQRCodeParams{
			Supplier:  sup,
			Customer:  cus,
			Document:  doc,
			ATCUD:     atcud,
			Signature: sig,
			SWCertNum: swCertNum,
		})

		assert.Equal(t, "A:165902612*B:151040222*C:PT*D:PF*E:N*F:20210101*G:PF TEST-1/33*H:AAJFJ3N5YS-123*I1:PT*I3:1222.22*I4:12.34*N:12.34*O:1234.56*Q:abcd*R:456", qr)
	})
}
