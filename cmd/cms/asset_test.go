package cms

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	hypertemplates "github.com/herdworks/hypertemplates-go"
	hypertexting "github.com/herdworks/hypertexting-go"
	href "github.com/herdworks/hypertexting-go/href"
)

func writeAssetListTestFile(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir asset test file dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("write asset test file: %v", err)
	}
}

func TestAssetListRowsUseWinningAssetSource(t *testing.T) {
	root := t.TempDir()
	website := &hypertexting.Website{
		URI: &href.Href{URL: url.URL{Path: filepath.Join(root, "site.yaml")}},
		Theme: hypertemplates.Theme{
			URI: &url.URL{Path: filepath.Join(root, "theme", "theme.json")},
		},
	}
	writeAssetListTestFile(t, filepath.Join(website.Theme.StaticDirPath(), "shared.css"), "theme")
	writeAssetListTestFile(t, filepath.Join(website.StaticDirPath(), "shared.css"), "website")
	writeAssetListTestFile(t, filepath.Join(website.ContentDirPath(), "shared.css"), "page")
	website.Theme.Assets = []string{"shared.css"}
	website.Assets = []string{"shared.css"}
	page := &hypertexting.Page{Assets: []string{"shared.css"}}

	entries, err := hypertexting.AssetEntries(&website.Theme, website, hypertexting.Pages{page})
	if err != nil {
		t.Fatalf("AssetEntries returned error: %v", err)
	}
	assets, err := assetListRows(website, entries)
	if err != nil {
		t.Fatalf("assetListRows returned error: %v", err)
	}

	if len(assets) != 1 {
		t.Fatalf("asset rows len = %d, want 1", len(assets))
	}
	if assets[0].Path != "/shared.css" {
		t.Fatalf("asset path = %q, want /shared.css", assets[0].Path)
	}
	if want := filepath.Join("content", "shared.css"); assets[0].Source != want {
		t.Fatalf("asset source = %q, want %q", assets[0].Source, want)
	}
	if assets[0].Size != int64(len("page")) {
		t.Fatalf("asset size = %d, want %d", assets[0].Size, len("page"))
	}
}
