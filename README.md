# GOBL ➡️ Portugal SAF-T

Portuguese SAF-T (PT) support for [GOBL](https://github.com/invopop/gobl): the GOBL **addon** (extensions, normalizers, scenarios, and validation rules), the **converter** that turns GOBL documents into SAF-T (PT) structures and seals them (AT signature, ATCUD and QR code), and the **client** for the AT web services used to report documents in real time.

Released under the Apache 2.0 [LICENSE](https://github.com/invopop/gobl.pt.saft/blob/main/LICENSE), Copyright 2026 [Invopop S.L.](https://invopop.com).

[![Lint](https://github.com/invopop/gobl.pt.saft/actions/workflows/lint.yaml/badge.svg)](https://github.com/invopop/gobl.pt.saft/actions/workflows/lint.yaml)
[![Test Go](https://github.com/invopop/gobl.pt.saft/actions/workflows/test.yaml/badge.svg)](https://github.com/invopop/gobl.pt.saft/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/invopop/gobl.pt.saft)](https://goreportcard.com/report/github.com/invopop/gobl.pt.saft)
[![codecov](https://codecov.io/gh/invopop/gobl.pt.saft/graph/badge.svg)](https://codecov.io/gh/invopop/gobl.pt.saft)
[![GoDoc](https://godoc.org/github.com/invopop/gobl.pt.saft?status.svg)](https://godoc.org/github.com/invopop/gobl.pt.saft)
![Latest Tag](https://img.shields.io/github/v/tag/invopop/gobl.pt.saft)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/invopop/gobl.pt.saft)

Portugal doesn't have an e-invoicing format per se. Tax information is reported
electronically to the AT (Autoridade Tributária e Aduaneira) either periodically
in batches via a SAF-T (PT) report or individually in real time via a web
service. Documents must also be sealed by certified software: numbered in
sequence, signed with an RSA key (whose hash is chained to the previous
document), identified with an ATCUD code and printed with a QR code.

## Layout

- [`addon/`](addon) — the GOBL addon (`pt-saft-v1`): extensions, normalizers,
  scenarios, and validation rules that register into GOBL on import. This
  package is kept dependency-light so importing it never pulls in conversion
  tooling.
- the module root — the converter that turns GOBL invoices, payments, orders
  and delivery documents into SAF-T (PT) structures (`Convert`), builds and
  streams complete SAF-T audit files (`BuildAuditFile`, `StreamEncoder`), and
  seals documents with the AT signature, ATCUD and QR code (`Seal`).
- [`config/`](config) — the configuration shared by the converter and the
  real-time reporting: certified software details and signing key.
- [`goblwrap/`](goblwrap) — wrappers that give the supported GOBL document types
  (invoices, payments, orders and deliveries) a common set of accessors.
- [`atws/`](atws) — a generic client for the AT web services (series
  registration and real-time reporting of invoices, payments, work documents
  and stock movements), usable independently of the rest of the module.
- [`realtime/`](realtime) — builders that map SAF-T structures to the requests
  of the real-time reporting web services.
- [`test/`](test) — test helpers and the sample documents under `test/data`.
- [`docs/`](docs) — reference material: SAF-T data structure, QR code
  specification and the AT's FAQ.

## Usage

### Addon

Unlike the converter, this is a true GOBL **addon**: it registers extensions,
normalizers, and validation rules into GOBL's global registry. Add a blank
import so it registers itself, then use GOBL as normal:

```go
import (
	"github.com/invopop/gobl"
	_ "github.com/invopop/gobl.pt.saft/addon"
)
```

Declare the addon on a document (or let the regime/scenario add it) and
`Calculate` + `Validate` will run the full SAF-T normalization and rules. The
converter imports the addon automatically, so consumers that only convert do not
need the blank import.

> **Note**: the `pt-saft-v1` key is listed in GOBL core's approved
> external-addon registry, so it is recognised as a valid `$addons` value in the
> JSON Schema. The runtime check stays strict, however: a document declaring the
> `pt-saft-v1` addon will fail validation with `add-on must be registered`
> unless this module is imported. Any service that processes Portuguese SAF-T
> documents must import it.

### Sealing

Before a document can be reported it must be sealed. `Seal` assigns the next
sequence number of the series, signs the document chaining it to the previous
signature, generates the ATCUD and QR code, and adds all of them as stamps to
the GOBL envelope. The returned signature must be persisted and passed in when
sealing the next document of the same series.

```go
package main

import (
	"github.com/invopop/gobl"
	saft "github.com/invopop/gobl.pt.saft"
	"github.com/invopop/gobl.pt.saft/config"
)

func main() {
	env := new(gobl.Envelope) // a calculated, unsigned GOBL envelope
	// ...

	cfg := &config.Config{
		CertificateNumber: 1234,                 // software certificate number
		SigningKeyPath:    "./signing_key.pem",  // RSA private key
		SigningKeyVersion: "1",
	}

	var prev *saft.Signature // last signature of the series, nil for the first document

	sig, err := saft.Seal(env, prev, "AAJFJHNC82", cfg) // series validation code
	if err != nil {
		panic(err)
	}

	// TODO: persist sig and sign the envelope
}
```

### Converter

`Convert` maps a sealed GOBL envelope to the SAF-T structures it contributes to
an audit file: the source document plus the customer, products and tax table
entries it references. `BuildAuditFile` assembles a complete audit file from
those, and `StreamEncoder` writes one incrementally from stored XML fragments.

```go
res, err := saft.Convert(env, "SOURCE_ID", cfg)
if err != nil {
	panic(err)
}

af, err := saft.BuildAuditFile(&saft.AuditFileParams{
	Config:          *cfg,
	Issuer:          supplier,
	Period:          period,
	Customers:       []*saft.Customer{res.Customer},
	Products:        res.Products,
	TaxTableEntries: res.TaxTableEntries,
	SourceDocuments: []saft.SourceDocument{res.Document},
})
if err != nil {
	panic(err)
}

out, err := saft.MarshalXML(af, true) // Windows-1252 encoded, as required by the AT
if err != nil {
	panic(err)
}

// TODO: do something with the output
```

### AT web services

The `atws` package is a generic client for the AT web services. It connects
with the client certificate issued by the AT and the credentials of each
taxpayer. The `realtime` package builds the requests for the real-time
reporting services from a SAF-T audit file containing the document to report.

```go
conn, err := atws.New(&atws.Config{
	Environment:         atws.EnvTesting,
	CertificatePath:     "./at-cert.pfx",
	CertificatePassword: "secret",
	PublicKeyPath:       "./at-public-key.pem",
	SWCertNum:           1234,
})
if err != nil {
	panic(err)
}

client := conn.NewClient("123456789/1", "password")

req, err := realtime.BuildRegisterInvoice(af)
if err != nil {
	panic(err)
}

res, err := client.Invoices().Register(req)
if err != nil {
	panic(err)
}

// TODO: handle res.Response
```

## Limitations

Summary of unsupported features:

- SAF-T for accounting
- Invoice-level charges and discounts
- Discounts at the time of payment (Settlements)
- Tax types:
  - IS – Imposto de Selo (Stamp duty)
  - Not subject to VAT/IS
- Customs details
- Self-billing and 3rd-party billing

### Unsupported fields

- QR Codes fields:
  - `L`: Não sujeito / não tributável em IVA / outras situações
  - `M`: Imposto do Selo
- SAF-T:
  - Tables:
    - Accounts table (`GeneralLedgerAccounts`)
    - Account movements table `GeneralLedgerEntries`
    - Suppliers table (`Supplier`)
  - Header fields:
    - `CompanyID` other than the NIF
    - `TaxAccountBasis` other than `F`
    - `TaxEntity` other than `Global`
    - `HeaderComment`
    - `Fax`
  - Products table fields:
    - `ProductGroup`
    - `CustomsDetails`
  - Customers table fields:
    - `ShipToAddress`
    - `SelfBillingIndicator` other than `0`
  - Invoices table fields:
    - DocumentStatus' `InvoiceStatus` other than `N` or `A`
    - DocumentStatus' `SourceBilling` other than `P`
    - SpecialRegimes' `SelfBillingIndicator` other than `0`
    - SpecialRegimes' `ThirdPartiesBillingIndicator` other than `0`
    - `EACCode` (business activity code)
    - `TransactionID` (for accounting)
    - `ShipFrom`
    - `ShipTo`
    - `MovementEndTime`
    - `MovementStartTime`
    - Line's `TaxBase`
    - Line's `ProductSerialNumber`
    - Line's `CustomsInformation`
    - Line Tax's `TaxType` other than `IVA`
    - Line Tax's `TaxAmount` (used for IS only)
    - DocumentTotals' `Settlement`
  - Stock movements table fields:
    - DocumentStatus' `MovementStatus` other than `N` or `A`
    - DocumentStatus' `SourceBilling` other than `P`
    - `TransactionID` (for accounting)
    - `EACCode` (business activity code)
    - ShippingPoint's `DeliveryID`
    - ShippingPoint's `WarehouseID`
    - ShippingPoint's `LocationID`
    - Line's `ProductSerialNumber`
    - Line's `CustomsInformation`
    - DocumentTotals' `Currency` (exchange rates)
  - Working documents table fields:
    - DocumentStatus' `WorkStatus` other than `N` or `A`
    - DocumentStatus' `SourceBilling` other than `P`
    - `TransactionID` (for accounting)
    - Line's `TaxBase`
    - Line's `CustomsInformation`
    - Line Tax's `TaxType` other than `IVA`
    - Line Tax's `TaxAmount` (used for IS only)
  - Payments table fields:
    - DocumentStatus' `PaymentStatus` other than `N` or `A`
    - DocumentStatus' `SourcePayment` other than `P`
    - `TransactionID` (for accounting)
    - Line's `SettlementAmount` (discount at the time of payment)
    - DocumentTotals' `Settlement`
    - DocumentTotals' `Currency` (exchange rates)

## Development

The module tracks tagged releases of `github.com/invopop/gobl`. A commented-out
`replace` directive sits at the bottom of `go.mod` for working against a local
core checkout; uncomment it while developing against unreleased core changes,
and re-comment it before merging.

```sh
go test ./...
```

### Examples

`examples/` holds sample documents for the addon, with their expected JSON
envelopes under `examples/out/`. They are verified via GOBL's shared
`pkg/examples` helpers.

`test/data/` holds sample GOBL envelopes for the converter. `TestConvertExamples`
seals all of them, builds a SAF-T audit file and compares it against
`test/data/out/saft.xml`.

Regenerate the golden outputs after intentional changes with:

```sh
go test . -update
```

With `-update`, the audit file is also validated against the official XSD
schema before being written.

## Resources

- [Portal das Finanças – Apoio ao contribuinte](https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/Pages/default.aspx)
- [Portaria n.º 302/2016 – SAF-T data structure](docs/portaria-no-302-2016-saft-data-structure-taxonomies.en.md)
- [QR code specification](docs/qr_code_specification.md)
- [SAF-T (PT) FAQ](docs/saft-pt-faq.md)

## License

Apache 2.0 — see [LICENSE](./LICENSE).
