package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/commander-spaceman/me2tlk/reader"
	"github.com/commander-spaceman/me2tlk/resolver"
)

func TestReaderRoundTrip(t *testing.T) {
	f := reader.BuildTestFile()

	if f.TotalEntries != 1 {
		t.Errorf("TotalEntries = %d, want 1", f.TotalEntries)
	}

	text, ok := reader.ResolveString(f, 1, true)
	if !ok {
		t.Fatal("failed to resolve test string")
	}
	if text != "AB" {
		t.Errorf("text = %q, want AB", text)
	}

	ids := f.StringIDs()
	if len(ids) != 1 || ids[0] != 1 {
		t.Errorf("StringIDs = %v, want [1]", ids)
	}

	entries := f.Search("B")
	if len(entries) != 1 || entries[0].Text != "AB" {
		t.Errorf("Search('B') = %v, want 1 entry with text 'AB'", entries)
	}

	entries = f.Search("Z")
	if len(entries) != 0 {
		t.Errorf("Search('Z') = %d results, want 0", len(entries))
	}
}

func TestResolverWithMultipleFiles(t *testing.T) {
	primary := reader.BuildTestFile()
	secondary := reader.BuildTestFile()
	secondary.Path = "secondary.tlk"

	r := &resolver.Resolver{Files: []*reader.File{primary, secondary}}
	text, ok := r.Resolve(1)
	if !ok {
		t.Fatal("failed to resolve from multi-file resolver")
	}
	if text != "AB" {
		t.Errorf("text = %q, want AB", text)
	}

	result := r.ResolveWithSource(1)
	if result == nil || !result.Found {
		t.Fatal("ResolveWithSource failed")
	}
	if result.SourceTLK != "test.tlk" {
		t.Errorf("SourceTLK = %q, want test.tlk", result.SourceTLK)
	}

	_, ok = r.Resolve(999)
	if ok {
		t.Fatal("expected resolution failure for id 999")
	}

	n := r.TotalUniqueEntries()
	if n != 1 {
		t.Errorf("TotalUniqueEntries = %d, want 1", n)
	}
}

func TestResolverPriorityOverride(t *testing.T) {
	base := reader.BuildTestFile()
	base.Path = "base.tlk"

	override := reader.BuildTestFile()
	override.Path = "override.tlk"

	r := &resolver.Resolver{Files: []*reader.File{override, base}}
	result := r.ResolveWithSource(1)
	if !result.Found {
		t.Fatal("expected resolution")
	}
	if result.SourceTLK != "override.tlk" {
		t.Errorf("SourceTLK = %q, want override.tlk", result.SourceTLK)
	}
}

func TestResolverIterAllEntries(t *testing.T) {
	primary := reader.BuildTestFile()
	secondary := reader.BuildTestFile()
	secondary.Path = "secondary.tlk"

	r := &resolver.Resolver{Files: []*reader.File{primary, secondary}}

	count := 0
	r.IterAllEntries()(func(id int32, text string, source string) bool {
		count++
		if id != 1 {
			t.Errorf("unexpected id %d", id)
		}
		return true
	})
	if count != 1 {
		t.Errorf("IterAllEntries yielded %d entries, want 1 (deduplication)", count)
	}
}

func TestRealTLKFile(t *testing.T) {
	paths := []string{
		filepath.Join("..", "..", "..", "pcc-toolkit", "output", "BIOGame_INT.tlk"),
		filepath.Join("testdata", "BIOGame_INT.tlk"),
	}
	var tlkPath string
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			tlkPath = p
			break
		}
	}
	if tlkPath == "" {
		t.Skip("real TLK file not available")
	}

	f, err := reader.ReadFile(tlkPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if f.TotalEntries == 0 {
		t.Error("expected non-zero total entries in real TLK")
	}

	if f.Path != tlkPath {
		t.Errorf("Path = %q, want %q", f.Path, tlkPath)
	}
}
