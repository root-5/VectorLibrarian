package crawler

import (
	"app/controller/log"
	"crypto/sha1"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
	"github.com/gocolly/colly/v2"
)

/*
colly.HTMLElement からページに関する各データを抽出する関数
  - e					HTMLElement
  - return)	domain		ドメイン
  - return)	path		パス
  - return)	pageTitle	ページタイトル
  - return)	description	ディスクリプション
  - return)	keywords	キーワード
  - return)	markdown	マークダウン形式のコンテンツ
  - return)	hash		コンテンツのハッシュ値
  - return) err			エラー
*/
func htmlToPageData(e *colly.HTMLElement) (domain, path, pageTitle, description, keywords, markdown, hash string, err error) {
	// html-to-markdown のコンバーターを作成
	conv := converter.NewConverter(
		converter.WithPlugins(
			base.NewBasePlugin(),             // ベースプラグイン（HTMLの基本的な変換を行う）
			commonmark.NewCommonmarkPlugin(), // マークダウンの変換プラグイン
			table.NewTablePlugin( // テーブルの変換プラグイン
				table.WithHeaderPromotion(true),      // false だとヘッダー行がなかった時にテーブル用のマークダウンが生成されない
				table.WithSpanCellBehavior("mirror"), // 結合されたセルがある場合、内容を複数セルに複製する
			),
		),
	)

	// URL からドメインとパスを取得
	domain = e.Request.URL.Hostname()
	path = e.Request.URL.Path

	// ページタイトル、ディスクリプション、キーワードを取得（それぞれ存在しない場合は "--" を設定）
	pageTitle = e.ChildText("title")
	if pageTitle == "" {
		pageTitle = "--"
	}
	description = e.ChildAttr("meta[name=description]", "content")
	if description == "" {
		description = "--"
	}
	keywords = e.ChildAttr("meta[name=keywords]", "content")
	if keywords == "" {
		keywords = "--"
	}

	// head, header, footer, script タグを削除（ルートのみ header, footer は残す）
	if e.Request.URL.Path != "/" {
		e.DOM.Find("header").Remove()
		e.DOM.Find("footer").Remove()
	}
	e.DOM.Find("head").Remove()
	e.DOM.Find("script").Remove()

	// HTML をマークダウン形式に変換して取得
	html, err := e.DOM.Html()
	if err != nil {
		log.Error(err)
		return
	}
	markdown, err = conv.ConvertString(html)
	if err != nil {
		log.Error(err)
		return
	}

	// markdown のハッシュを計算
	hashBin := sha1.Sum([]byte(markdown))
	hash = hex.EncodeToString(hashBin[:])

	return domain, path, pageTitle, description, keywords, markdown, hash, nil
}

/*
URL パスの検証と絶対パス変換、https 変換を行う関数
  - e						HTMLElement
  - targetDomain			対象ドメイン
  - allowedPaths			許可するパスの配列
  - return)	formattedLink	フォーマット済みのリンクURL
  - return) isValid			URLが有効かどうか
*/
func validateAndFormatLinkUrl(e *colly.HTMLElement, targetDomain string, allowedPaths []string) (formattedLink string, isValid bool) {
	link := e.Attr("href")

	// .pdf で終わるリンク、mailto:/javascript:/# 始まるリンク、空のリンクはスキップ
	if matched, _ := regexp.MatchString(`(?i)\.pdf$|^mailto:|^javascript:|^$|^#`, link); matched {
		return "", false
	}
	// http:// を https:// に変換
	if strings.HasPrefix(link, "http://") {
		link = strings.Replace(link, "http://", "https://", 1)
	}
	// 相対パスを絶対パスに変換
	if !strings.HasPrefix(link, "https://") {
		link = e.Request.AbsoluteURL(link)
	}
	// 外部ドメインはスキップ
	if !strings.HasPrefix(link, "https://"+targetDomain) {
		return "", false
	}
	// 特定のパス以外をスキップ
	matched := false
	for _, allowedPath := range allowedPaths {
		if strings.Contains(link, allowedPath) {
			matched = true
			break
		}
	}
	if !matched {
		return "", false
	}
	formattedLink = link
	isValid = true

	return formattedLink, isValid
}

/*
マークダウンから見出しを抽出し、箇条書き形式に変換する関数
  - markdown	マークダウン形式のテキスト
  - return)		見出しを箇条書き形式にした文字列
*/
func ExtractHeadings(markdown string) (itemizationStr string) {
	// 正規表現で見出しを抽出
	re := regexp.MustCompile(`(?m)^(#{1,6})\s+(.*)$`)
	matches := re.FindAllStringSubmatch(markdown, -1)

	// 見出しを箇条書き形式に変換
	var itemization []string
	for _, match := range matches {
		level := len(match[1])                                  // # の数でレベルを取得
		item := strings.Repeat("  ", level-1) + "- " + match[2] // レベルに応じてインデント
		itemization = append(itemization, item)
	}
	// 箇条書き形式の文字列を結合
	itemizationStr = strings.Join(itemization, "\n")

	return itemizationStr
}
