# GOBL ➡️ Portugal SAF-T

Portuguese SAF-T (PT) tax addon for [GOBL](https://github.com/invopop/gobl).

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
service. This addon (`pt-saft-v1`) ensures that GOBL documents carry all the
fields the AT requires.

Unlike the format converters in the GOBL ecosystem, this is a true **addon**: it
registers extensions, normalizers, scenarios, and validation rules into GOBL's
global registry. It lives in its own module so that only projects handling
Portuguese SAF-T documents take on its weight.

## Layout

- `addon/` — the GOBL addon: extensions, normalizers, scenarios, and validation
  rules that register into GOBL on import. This package is kept dependency-light
  so importing it never pulls in conversion tooling.
- the module root (and future subpackages) is reserved for converters and other
  SAF-T logic that build on the addon.

## Usage

Add a blank import of the **addon** so it registers itself, then use GOBL as
normal:

```go
import (
	"github.com/invopop/gobl"
	_ "github.com/invopop/gobl.pt.saft/addon"
)
```

Declare the addon on a document (or let the regime/scenario add it) and
`Calculate` + `Validate` will run the full SAF-T normalization and rules.

> **Note**: the `pt-saft-v1` key is listed in GOBL core's approved
> external-addon registry, so it is recognised as a valid `$addons` value in the
> JSON Schema. The runtime check stays strict, however: a document declaring the
> `pt-saft-v1` addon will fail validation with `add-on must be registered`
> unless this module is imported. Any service that processes Portuguese SAF-T
> documents must import it.

## Development

The addon builds on core GOBL features (the approved external-addon registry)
that are not yet in a tagged release. The `go.mod` therefore pins
`github.com/invopop/gobl` to the core checkout via a `replace` directive; bump it
to the release tag and drop the replace once core is published.

```sh
go test ./...
```

### Examples

`examples/` holds sample documents, with their expected JSON envelopes under
`examples/out/`. They are verified via GOBL's shared `pkg/examples` helpers.
Regenerate the golden output after intentional changes with:

```sh
go test . -run TestExamples -update
```

## License

Apache 2.0 — see [LICENSE](./LICENSE).
