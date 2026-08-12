package cms

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"

	httheme "github.com/herdworks/httt/theme"
	"github.com/herdworks/htxt"
	href "github.com/herdworks/htxt/href"
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
	website := &htxt.Website{
		URI: &href.Href{URL: url.URL{Path: filepath.Join(root, "site.yaml")}},
		Theme: httheme.Theme{
			URI: &url.URL{Path: filepath.Join(root, "theme", "theme.json")},
		},
	}
	writeAssetListTestFile(t, filepath.Join(website.Theme.StaticDirPath(), "shared.css"), "theme")
	writeAssetListTestFile(t, filepath.Join(website.StaticDirPath(), "shared.css"), "website")
	writeAssetListTestFile(t, filepath.Join(website.ContentDirPath(), "shared.css"), "page")
	website.Theme.Assets = []string{"shared.css"}
	website.Assets = []string{"shared.css"}
	page := &htxt.Page{Assets: []string{"shared.css"}}

	entries, err := htxt.AssetEntries(&website.Theme, website, htxt.Pages{page})
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

func TestAssetLookupMatchesOrdersThemeWebsitePageWinnerLast(t *testing.T) {
	root := t.TempDir()
	website := &htxt.Website{
		URI: &href.Href{URL: url.URL{Path: filepath.Join(root, "site.yaml")}},
		Theme: httheme.Theme{
			URI: &url.URL{Path: filepath.Join(root, "theme", "theme.json")},
		},
	}
	website.Theme.Assets = []string{"shared.css", "theme-only.css"}
	website.Assets = []string{"shared.css"}
	page := &htxt.Page{Assets: []string{"shared.css"}}

	candidates, err := htxt.AssetCandidates(&website.Theme, website, htxt.Pages{page})
	if err != nil {
		t.Fatalf("AssetCandidates returned error: %v", err)
	}

	matches := assetLookupMatches(candidates, "shared.css")
	if len(matches) != 3 {
		t.Fatalf("matches len = %d, want 3", len(matches))
	}
	wantSources := []htxt.AssetSource{
		htxt.AssetSourceTheme,
		htxt.AssetSourceWebsite,
		htxt.AssetSourcePage,
	}
	for i, want := range wantSources {
		if matches[i].Source != want {
			t.Fatalf("matches[%d].Source = %q, want %q", i, matches[i].Source, want)
		}
	}

	if got := assetLookupMatches(candidates, "/shared.css"); len(got) != 3 {
		t.Fatalf("leading-slash query matches len = %d, want 3", len(got))
	}

	themeOnly := assetLookupMatches(candidates, "theme-only.css")
	if len(themeOnly) != 1 || themeOnly[0].Source != htxt.AssetSourceTheme {
		t.Fatalf("theme-only matches = %#v, want one theme entry", themeOnly)
	}

	if got := assetLookupMatches(candidates, "missing.css"); len(got) != 0 {
		t.Fatalf("missing.css matches len = %d, want 0", len(got))
	}
}
