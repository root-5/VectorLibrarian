package crawler

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// 単体テスト（外部依存がない関数のテスト）を定義
// `docker compose exec app go test ./controller/crawler`

func TestExtractHeadings(t *testing.T) {
	testCases := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name: "通常ケース",
			markdown: `
# Title 1
Some text here.
## Title 2
More text.
### Title 3
`,
			expected: "- Title 1\n  - Title 2\n    - Title 3",
		},
		{
			name:     "見出しなし",
			markdown: "This is a text without any headings.",
			expected: "",
		},
		{
			name: "同レベルの複数見出し",
			markdown: `
# Title 1
## Title 2.1
## Title 2.2
`,
			expected: "- Title 1\n  - Title 2.1\n  - Title 2.2",
		},
		{
			name:     "空のマークダウン",
			markdown: "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := extractHeadings(tc.markdown)
			if actual != tc.expected {
				t.Errorf("期待値:\n%s\n実際:\n%s", tc.expected, actual)
			}
		})
	}
}

func TestValidateAndFormatLinkUrl(t *testing.T) {
	targetDomain := "example.com"
	allowedPaths := []string{"/allowed"}

	// モックリクエスト
	u, _ := url.Parse("https://example.com/")
	req := &colly.Request{
		URL:    u,
		Method: "GET",
		Ctx:    colly.NewContext(),
	}

	testCases := []struct {
		name          string
		href          string
		expectedLink  string
		expectedValid bool
	}{
		{"有効な相対パス", "/allowed/page1", "https://example.com/allowed/page1", true},
		{"有効な絶対パス", "https://example.com/allowed/page2", "https://example.com/allowed/page2", true},
		{"HTTP を HTTPS に変換", "http://example.com/allowed/page3", "https://example.com/allowed/page3", true},
		{"外部ドメイン", "https://another.com/page", "", false},
		{"PDF リンク", "/allowed/document.pdf", "", false},
		{"Mailto リンク", "mailto:test@example.com", "", false},
		{"JavaScript リンク", "javascript:void(0)", "", false},
		{"空のリンク", "", "", false},
		{"アンカーリンク", "#section1", "", false},
		{"許可されていないパス", "/disallowed/page", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モック Response を作成
			resp := &colly.Response{
				Request: req,
			}

			// モック HTMLElement を作成
			doc, _ := goquery.NewDocumentFromReader(strings.NewReader(`<a href="` + tc.href + `">Link</a>`))
			aNode := doc.Find("a").Get(0)
			e := colly.NewHTMLElementFromSelectionNode(resp, doc.Find("a"), aNode, 0)

			formattedLink, isValid := validateAndFormatLinkUrl(e, targetDomain, allowedPaths)

			if formattedLink != tc.expectedLink {
				t.Errorf("期待されるリンク '%s' ですが、実際は '%s' でした", tc.expectedLink, formattedLink)
			}
			if isValid != tc.expectedValid {
				t.Errorf("期待される有効性 %v ですが、実際は %v でした", tc.expectedValid, isValid)
			}
		})
	}
}

func TestHtmlToPageData(t *testing.T) {
	htmlContent := `
<html>
<head>
    <title>Test Title</title>
    <meta name="description" content="Test Description">
    <meta name="keywords" content="Test, Keywords">
</head>
<body>
    <header>Header content</header>
    <h1>Main Heading</h1>
    <p>This is a paragraph.</p>
    <footer>Footer content</footer>
</body>
</html>`

	// 手動で colly.Request を作成
	u, _ := url.Parse("https://example.com/test-path")
	req := &colly.Request{
		URL:    u,
		Method: "GET",
		Ctx:    colly.NewContext(),
	}

	// モックの http.Response を作成
	httpResp := &http.Response{
		StatusCode: http.StatusOK,
		// Body フィールドは現状未使用のため省略
		// Request フィールドは現状未使用のため省略
	}
	// colly.Response を作成してリクエストに紐付ける
	resp := &colly.Response{
		Request:    req,
		StatusCode: httpResp.StatusCode,
		Body:       []byte(htmlContent),
		Ctx:        req.Ctx,
		Headers:    &httpResp.Header,
	}

	// goquery.Document を作成
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("goquery ドキュメントの作成に失敗しました: %v", err)
	}

	// ルートの html ノードを使って HTMLElement を作成
	htmlNode := doc.Find("html").Get(0)
	if htmlNode == nil {
		t.Fatal("html ノードが見つかりません")
	}
	e := colly.NewHTMLElementFromSelectionNode(resp, doc.Find("html"), htmlNode, 0)

	domain, path, pageTitle, description, keywords, markdown, hash, err := htmlToPageData(e)

	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}

	expectedDomain := "example.com"
	if domain != expectedDomain {
		t.Errorf("期待されるドメイン '%s' ですが、実際は '%s' でした", expectedDomain, domain)
	}

	expectedPath := "/test-path"
	if path != expectedPath {
		t.Errorf("期待されるパス '%s' ですが、実際は '%s' でした", expectedPath, path)
	}

	expectedTitle := "Test Title"
	if pageTitle != expectedTitle {
		t.Errorf("期待されるタイトル '%s' ですが、実際は '%s' でした", expectedTitle, pageTitle)
	}

	expectedDescription := "Test Description"
	if description != expectedDescription {
		t.Errorf("期待されるディスクリプション '%s' ですが、実際は '%s' でした", expectedDescription, description)
	}

	expectedKeywords := "Test, Keywords"
	if keywords != expectedKeywords {
		t.Errorf("期待されるキーワード '%s' ですが、実際は '%s' でした", expectedKeywords, keywords)
	}

	// path が "/" でないため header と footer は除去されるはず
	if strings.Contains(markdown, "Header content") {
		t.Errorf("Markdown に header の内容が含まれてはいけません")
	}
	if strings.Contains(markdown, "Footer content") {
		t.Errorf("Markdown に footer の内容が含まれてはいけません")
	}

	// 実際に生成された markdown から期待されるハッシュを計算
	hashBin := sha1.Sum([]byte(markdown))
	expectedHash := hex.EncodeToString(hashBin[:])
	if hash != expectedHash {
		t.Errorf("期待されるハッシュ '%s' ですが、実際は '%s' でした", expectedHash, hash)
	}
}
