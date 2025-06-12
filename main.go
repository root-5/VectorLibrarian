package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	targetDomain := "www.city.hamura.tokyo.jp"
	// targetDomain := "nextjs-root-5.vercel.app"

	// デフォルトのコレクターを作成
	c := colly.NewCollector(
		colly.AllowedDomains(targetDomain), // 許可するドメインを設定
		colly.MaxDepth(2),                  // 最大深度を 2 に設定
		colly.CacheDir("./cache"),          // キャッシュディレクトリを指定
		colly.Async(true),                  // 非同期モードを有効にする
		// colly.IgnoreRobotsTxt(),            // robots.txt を無視
	)

	// リクエスト間で1~2秒の時間を空ける
	c.Limit(&colly.LimitRule{
		DomainGlob: targetDomain, // 対象ドメインを指定
		Delay:      time.Second,  // リクエスト間の最小遅延
		// RandomDelay: time.Second, // リクエスト間のランダム遅延
	})

	// リクエスト前に "アクセス >> " を表示
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("")
		fmt.Println("")
		fmt.Println(">> ページURL: ", r.URL.String())
	})

	// テキストコンテンツを抽出するためのコールバック
	c.OnHTML("html", func(e *colly.HTMLElement) {
		// ページのタイトルを表示
		pageTitle := e.ChildText("title")
		if pageTitle == "" {
			pageTitle = "--"
		}
		fmt.Println(">> タイトル:", pageTitle)

		// ディスクリプションを表示
		description := e.ChildAttr("meta[name=description]", "content")
		if description == "" {
			description = "--"
		}
		fmt.Println(">> ディスクリプション:", description)

		// head, header, footer, script タグを削除（ルートのみ header, footer は残す）
		if e.Request.URL.Path != "/" {
			e.DOM.Find("header").Remove()
			e.DOM.Find("footer").Remove()
		}
		e.DOM.Find("head").Remove()
		e.DOM.Find("script").Remove()

		// head タグと script タグは常時削除
		e.DOM.Find("head").Remove()
		e.DOM.Find("script").Remove()

		// HTMLをマークダウン形式に変換して取得
		// html, err := e.DOM.Html()
		// if err != nil {
		// 	fmt.Println("HTMLの取得に失敗しました:", err)
		// 	return
		// }
		// markdown := convertHtmlToMarkdown(html)
		// fmt.Println(">> メイン部のマークダウン:")
		// fmt.Println(markdown)
	})

	// href属性を持つa要素ごとにコールバックを実行
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")

		// .pdf で終わるリンク、mailto:/javascript:/# 始まるリンク、空のリンクはスキップ
		if matched, _ := regexp.MatchString(`(?i)\.pdf$|^mailto:|^javascript:|^$|^#`, link); matched {
			fmt.Println(">>     - スキップ: ", link)
			return
		}
		// 外部ドメインはスキップ
		if !strings.Contains(link, targetDomain) {
			fmt.Println(">>     - スキップ: ", link)
			return
		}

		// http:// で始まるリンクは https:// に変換
		if strings.HasPrefix(link, "http://") {
			link = strings.Replace(link, "http://", "https://", 1)
		}
		// リンクが相対パスの場合は絶対URLに変換
		if !strings.HasPrefix(link, "https://") {
			link = e.Request.AbsoluteURL(link)
		}

		// ページ内で見つかったリンクを訪問
		fmt.Println(">>     - リンク: ", link)
		e.Request.Visit(link)
	})

	// 指定ドメインからスクレイピングを開始
	c.Visit("https://" + targetDomain + "/")
}

// HTML をマークダウン形式に変換する関数、対応するのは h1~h6, ul, ol, li, table, tr, td, th, p, div タグ
func convertHtmlToMarkdown(html string) string {
	var markdown strings.Builder

	// タグとマークダウンの対応表
	replacer := []struct {
		reOpen  *regexp.Regexp
		reClose *regexp.Regexp
		mdOpen  string
		mdClose string
	}{
		{regexp.MustCompile(`(?i)<h1[^>]*>`), regexp.MustCompile(`(?i)</h1>`), "# ", "\n\n"},
		{regexp.MustCompile(`(?i)<h2[^>]*>`), regexp.MustCompile(`(?i)</h2>`), "## ", "\n\n"},
		{regexp.MustCompile(`(?i)<h3[^>]*>`), regexp.MustCompile(`(?i)</h3>`), "### ", "\n\n"},
		{regexp.MustCompile(`(?i)<h4[^>]*>`), regexp.MustCompile(`(?i)</h4>`), "#### ", "\n\n"},
		{regexp.MustCompile(`(?i)<h5[^>]*>`), regexp.MustCompile(`(?i)</h5>`), "##### ", "\n\n"},
		{regexp.MustCompile(`(?i)<h6[^>]*>`), regexp.MustCompile(`(?i)</h6>`), "###### ", "\n\n"},
		{regexp.MustCompile(`(?i)<ul[^>]*>`), regexp.MustCompile(`(?i)</ul>`), "", "\n"},
		{regexp.MustCompile(`(?i)<ol[^>]*>`), regexp.MustCompile(`(?i)</ol>`), "", "\n"},
		{regexp.MustCompile(`(?i)<li[^>]*>`), regexp.MustCompile(`(?i)</li>`), "- ", "\n"},
		{regexp.MustCompile(`(?i)<table[^>]*>`), regexp.MustCompile(`(?i)</table>`), "", "\n\n"},
		{regexp.MustCompile(`(?i)<tr[^>]*>`), regexp.MustCompile(`(?i)</tr>`), "", "\n"},
		{regexp.MustCompile(`(?i)<th[^>]*>`), regexp.MustCompile(`(?i)</th>`), "| ", " "},
		{regexp.MustCompile(`(?i)<td[^>]*>`), regexp.MustCompile(`(?i)</td>`), "| ", " "},
		{regexp.MustCompile(`(?i)<p[^>]*>`), regexp.MustCompile(`(?i)</p>`), "", "\n\n"},
		{regexp.MustCompile(`(?i)<div[^>]*>`), regexp.MustCompile(`(?i)</div>`), "", "\n\n"},
	}

	md := html
	for _, r := range replacer {
		md = r.reOpen.ReplaceAllString(md, r.mdOpen)
		md = r.reClose.ReplaceAllString(md, r.mdClose)
	}

	// 残ったタグをすべて除去
	reTag := regexp.MustCompile(`(?s)<[^>]+>`)
	md = reTag.ReplaceAllString(md, "")

	// 改行を適切に処理
	md = strings.TrimSpace(md)
	// 連続する改行を1つにまとめる
	md = regexp.MustCompile(`\n{2,}`).ReplaceAllString(md, "\n\n")

	markdown.WriteString(md)
	return markdown.String()
}
