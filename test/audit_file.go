package test

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/invopop/gobl"
	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/config"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
)

// GenerateAuditFile converts a list of GOBL envelopes to SAF-T structures and builds a
// SAF-T report with them. The service streams reports straight from stored documents, so
// this whole-file build exists only to give tests a report to assert against.
func GenerateAuditFile(sup *org.Party, per cal.Period, cfg *config.Config, sourceID string, envs ...*gobl.Envelope) (*saft.AuditFile, error) {
	var custs []*saft.Customer
	var prods []*saft.Product
	var ttens []*saft.TaxTableEntry
	var sdocs []saft.SourceDocument

	for _, env := range envs {
		res, err := saft.Convert(env, sourceID, cfg)
		if err != nil {
			return nil, err
		}
		if res.Customer != nil {
			custs = append(custs, res.Customer)
		}
		prods = append(prods, res.Products...)
		ttens = append(ttens, res.TaxTableEntries...)
		sdocs = append(sdocs, res.Document)
	}

	return saft.BuildAuditFile(&saft.AuditFileParams{
		Config:          *cfg,
		Issuer:          sup,
		Period:          per,
		Products:        uniqueBy(prods, func(p *saft.Product) any { return p.ProductCode }),
		Customers:       uniqueBy(custs, func(c *saft.Customer) any { return c.CustomerID }),
		TaxTableEntries: uniqueBy(ttens, func(t *saft.TaxTableEntry) any { return hashID(t) }),
		SourceDocuments: sdocs,
	})
}

// hashID generates a unique ID for an entity from the hash of its XML representation.
func hashID(e any) string {
	data, _ := saft.MarshalXML(e, false)
	hash := sha256.Sum256(data)
	id := hex.EncodeToString(hash[:])
	if len(id) > 30 {
		id = id[:30]
	}
	return id
}

func uniqueBy[T any](items []*T, comparer func(*T) any) []*T {
	unique := make([]*T, 0, len(items))
	seen := make(map[any]bool)
	for _, i := range items {
		key := comparer(i)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, i)
		}
	}
	return unique
}
