# me2tlk

Go library for reading Mass Effect 2 TLK (talk table) files.

## Installation

```bash
go get github.com/commander-spaceman/me2tlk
```

## Packages

### reader

Parse and decode TLK files.

```go
import "github.com/commander-spaceman/me2tlk/reader"

f, err := reader.ReadFile("BIOGame_INT.tlk")
text, ok := reader.ResolveString(f, stringID, true)
```

### resolver

DLC-aware string resolution across multiple TLK files.

```go
import "github.com/commander-spaceman/me2tlk/resolver"

r, err := resolver.BuildResolver("BIOGame_INT.tlk", "DLC", "INT", false)
text, sourceTLK := r.ResolveWithSource(stringID)
```

## License

MIT
