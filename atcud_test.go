package saft_test

import (
	"testing"

	saft "github.com/invopop/gobl.pt.saft"
	"github.com/stretchr/testify/assert"
)

func TestBuildATCUD(t *testing.T) {
	t.Run("General case", func(t *testing.T) {
		doc := &saft.Invoice{
			InvoiceNo: "FT TEST-1/123",
		}
		valCode := "AAJFJ3N5YS"

		atcud := saft.BuildATCUD(doc, valCode)

		assert.Equal(t, "AAJFJ3N5YS-123", atcud)
	})
}
