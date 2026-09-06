# lazyjson

[![CI](https://github.com/headlesslab/lazyjson/actions/workflows/ci.yml/badge.svg)](https://github.com/headlesslab/lazyjson/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/headlesslab/lazyjson.svg)](https://pkg.go.dev/github.com/headlesslab/lazyjson)

A tiny JSON library to read and alter a JSON value. The data structure is lazy: it parses on read, so the parser can be swapped for a faster one where performance matters (see `JSON.Raw`).


## Provenance

`lazyjson` is a Satellite module of [wand](https://github.com/headlesslab/wand). It is a snapshot of [gson](https://github.com/ysmood/gson) `v0.7.3` (commit `c60042b`), imported without git history and maintained here by headlesslab; see [NOTICE](NOTICE). The module path changed, the exported identifiers did not.

Replace `github.com/ysmood/gson` with `github.com/headlesslab/lazyjson` and `gson.` with `lazyjson.`; `JSON`, `New`, `NewFrom`, `Int`, ... keep their names. Error messages carry the `lazyjson:` prefix.

## Usage

```go
import "github.com/headlesslab/lazyjson"
```

The tests are the documentation: see [lazyjson_test.go](lazyjson_test.go).

## Versioning

`v0.x`: no breaking change without a major version bump. Releases are cut by hand when the module changes; wand pins each release exactly.

## License

MIT, see [LICENSE](LICENSE).
