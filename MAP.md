# Project Map

**Purpose:** Go library for parsing Mass Effect 2 TLK (talk table) files — a Huffman-compressed binary format storing localized game strings. Provides single-file reading and DLC-aware multi-file resolution.

## Notes for AI Agents

- **Entry points:** `reader.ReadFile` / `reader.Parse`, `resolver.BuildResolver`
- **Main patterns:** Single-package Go library with no external dependencies beyond stdlib. Data flows from binary bytes → parsed `File` struct → decoded strings. The resolver composes reader primitives for multi-file lookups.
- **General rule:** The `reader` package is the low-level engine; the `resolver` package is the higher-level consumer. Never introduce circular imports between them.

---

## 1. Reader

Low-level TLK file parser and Huffman decoder. Reads the binary header, entry tables, Huffman tree, and bitstream, then decodes individual strings on demand.

```text
reader/
  reader.go        — Parse(), DecodeString(), ResolveString(), File methods
  types.go         — Structs (Header, File, Node, Entry), TLKMagic constant
  testutil.go      — BuildTestFile() synthetic TLK constructor for tests
  reader_test.go   — Unit tests for parsing, decoding, bit ops, iteration, search
```

**Main responsibilities:**

- Parse TLK binary format (header, male/female entry tables, Huffman tree nodes, bitstream)
- Decode Huffman-compressed UTF-8 strings from the bitstream
- Expose iteration, search, and ID-listing APIs on parsed `File`

**Key files:**

- `reader/reader.go:19` — `Parse()` is the core deserialization routine; validates magic, sizes, and builds the in-memory representation.
- `reader/types.go:25` — `File` struct holds all parsed state (path, header, entry maps, nodes, bitstream).
- `reader/reader.go:112` — `DecodeString()` walks the Huffman tree bit-by-bit to reconstruct UTF-8 text.
- `reader/testutil.go:5` — `BuildTestFile()` constructs a minimal valid TLK in memory for both reader and resolver tests.

**Relationships:**

- Imported by `resolver` and `tests` packages.
- Depends only on Go stdlib (`encoding/binary`, `os`, `strings`).

---

## 2. Resolver

DLC-aware string resolution across multiple TLK files. Discovers DLC talk tables via `BIOEngine.ini` modules, `Mount.dlc` priority, and file globs; resolves string IDs with override precedence.

```text
resolver/
  resolver.go        — BuildResolver(), FindDlcTlkFiles(), Resolver methods
  resolver_test.go   — Unit tests for resolution, priority, BIOEngine parsing, mount priority
```

**Main responsibilities:**

- Build a priority-ordered list of TLK files (base game + DLCs)
- Resolve string IDs across multiple files with first-match-wins semantics
- Parse `BIOEngine.ini` for DLC module mappings and `Mount.dlc` for load priority
- Deduplicate and search across all loaded TLKs

**Key files:**

- `resolver/resolver.go:261` — `BuildResolver()` is the primary entry point; loads base TLK, discovers DLC TLKs, and returns a ready-to-use `Resolver`.
- `resolver/resolver.go:192` — `FindDlcTlkFiles()` does the heavy lifting of DLC folder traversal, priority sorting, and glob-based fallback.
- `resolver/resolver.go:25` — `Resolver.Resolve()` resolves a string ID against the file list in priority order.

**Relationships:**

- Depends on `reader` for `Parse`/`ReadFile`/`ResolveString`.
- Imports Go stdlib (`os`, `path/filepath`, `sort`, `strings`, `encoding/binary`).

---

## 3. Integration Tests

Cross-package integration tests exercising the full read-and-resolve pipeline, including optional tests against real TLK files on disk.

```text
tests/
  tlk_test.go   — Round-trip test, multi-file resolver, priority override, iteration, optional real-file test
```

**Key files:**

- `tests/tlk_test.go:12` — `TestReaderRoundTrip` validates the full parse→decode chain.
- `tests/tlk_test.go:113` — `TestRealTLKFile` reads an actual TLK from disk if available (skipped otherwise).

**Relationships:**

- Depends on both `reader` and `resolver`.

---

## 4. Root Files

```text
go.mod        — Module: github.com/commander-spaceman/me2tlk, Go 1.25.5
README.md     — Usage examples for reader and resolver packages
LICENSE       — MIT
.gitignore    — Go binaries, test artifacts, IDE files
```
