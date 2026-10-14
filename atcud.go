package saft

import (
	"regexp"
)

// BuildATCUD builds the ATCUD code for a given document and validation code.
func BuildATCUD(doc SourceDocument, valCode string) string {
	seq := regexp.MustCompile(`\d+$`).FindString(doc.GetDocumentNo())
	return valCode + "-" + seq
}
