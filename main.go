package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	// デフォルトのコレクターを作成
	c := colly.NewCollector(
		// 最大深度を1に設定（スクレイピング対象ページ内のリンクのみ訪問し、それ以上は辿らない）
		colly.AllowedDomains("nextjs-root-5.vercel.app"),
		colly.MaxDepth(1),
	)

	// href属性を持つa要素ごとにコールバックを実行
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// リンクを表示
		fmt.Println("    - ", link)

		// ページ内で見つかったリンクを訪問
		// e.Request.Visit(link)

		// AllowedDomains に含まれるリンクのみ訪問される
		// c.Visit(e.Request.AbsoluteURL(link))
	})

	// テキストコンテンツを抽出するためのコールバック
	c.OnHTML("html", func(e *colly.HTMLElement) {
		// ページのタイトルを表示
		pageTitle := e.ChildText("title")
		if pageTitle == "" {
			pageTitle = "--"
		}

		// ディスクリプションを表示
		description := e.ChildAttr("meta[name=description]", "content")
		if description == "" {
			description = "--"
		}

		// head, header, footer, script タグを削除
		e.DOM.Find("head").Remove()
		e.DOM.Find("header").Remove()
		e.DOM.Find("footer").Remove()
		e.DOM.Find("script").Remove()

		// HTMLをマークダウン形式に変換して取得
		html, err := e.DOM.Html()
		if err != nil {
			fmt.Println("HTMLの取得に失敗しました:", err)
			return
		}
		markdown := convertHtmlToMarkdown(html)

		fmt.Println(">> ページタイトル:", pageTitle)
		fmt.Println(">> ページディスクリプション:", description)
		fmt.Println(">> メイン部のマークダウン:\n", markdown)
		fmt.Println("\n\n---\n\n")
	})

	// リクエスト前に "アクセス >> " を表示
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(">> アクセス中... ", r.URL.String())
		fmt.Println(">> リンク取得:")
	})

	// https://nextjs-root-5.vercel.app/ からスクレイピングを開始
	c.Visit("https://nextjs-root-5.vercel.app/")
}

// HTML をマークダウン形式に変換する関数、対応するのは h1~h6, ul, ol, li, table, tr, td, th, p タグ
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
