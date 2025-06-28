package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	// ターゲットドメインを設定
	targetDomain := "hotel-example-site.takeyaqa.dev"

	// ログファイルを作成
	logFile, err := os.OpenFile("log/scraper.md", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("ログファイルの作成に失敗しました:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// ログファイルの文字をすべて削除
	if err := os.Truncate("log/scraper.md", 0); err != nil {
		log.Fatalln("ログファイルの初期化に失敗しました:", err)
	}

	// デフォルトのコレクターを作成
	c := colly.NewCollector(
		colly.AllowedDomains(targetDomain), // 許可するドメインを設定
		colly.MaxDepth(2),                  // 最大深度を 2 に設定
		colly.CacheDir("./cache"),          // キャッシュディレクトリを指定
		// colly.Async(true),                  // 非同期モードを有効にする
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
		fmt.Println(">> URL: ", r.URL.String())
	})

	// テキストコンテンツを抽出するためのコールバック
	c.OnHTML("html", func(e *colly.HTMLElement) {
		// // ページのタイトルを表示
		// pageTitle := e.ChildText("title")
		// if pageTitle == "" {
		// 	pageTitle = "--"
		// }
		// fmt.Println(">> title:", pageTitle)

		// // ディスクリプションを表示
		// description := e.ChildAttr("meta[name=description]", "content")
		// if description == "" {
		// 	description = "--"
		// }
		// fmt.Println(">> description:", description)

		// head, header, footer, script タグを削除（ルートのみ header, footer は残す）
		if e.Request.URL.Path != "/" {
			e.DOM.Find("header").Remove()
			e.DOM.Find("footer").Remove()
		}
		e.DOM.Find("head").Remove()
		e.DOM.Find("script").Remove()

		// HTMLをマークダウン形式に変換して取得
		html, err := e.DOM.Html()
		if err != nil {
			fmt.Println("HTMLの取得に失敗しました:", err)
			return
		}
		markdown := convertHtmlToMarkdown(html)
		// log.Println(">> メイン部のマークダウン:")
		log.Println("\n" + markdown)
	})

	// href属性を持つa要素ごとにコールバックを実行
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// .pdf で終わるリンク、mailto:/javascript:/# 始まるリンク、空のリンクはスキップ
		if matched, _ := regexp.MatchString(`(?i)\.pdf$|^mailto:|^javascript:|^$|^#`, link); matched {
			// fmt.Println(">>     - skip: ", link)
			return
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
			// fmt.Println(">>     - skip: ", link)
			return
		}

		// ページ内で見つかったリンクを訪問
		// fmt.Println(">>     - link: ", link)
		e.Request.Visit(link)
	})

	// 指定ドメインからスクレイピングを開始
	c.Visit("https://" + targetDomain + "/")
}

// HTML をマークダウン形式に変換する関数、対応するのは h1~h6, ul, ol, li, table, tr, td, th, p, div タグ
func convertHtmlToMarkdown(html string) (markdown string) {
	targetDomain := "hotel-example-site.takeyaqa.dev"

	// HTML タグとマークダウンの変換作業前に文字数を削減
	html = regexp.MustCompile(`(?s)>\s+<`).ReplaceAllString(html, "><")    // タグ間のスペースと改行を削除
	html = regexp.MustCompile(`\s{2,}`).ReplaceAllString(html, " ")        // 2 つ以上の連続する空白を 1 つにまとめる
	html = regexp.MustCompile(`(?s)<!--.*?-->`).ReplaceAllString(html, "") // コメントを削除

	// タグの属性を削除（href と alt 属性は残す）
	html = regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)[^>]*?(?:\s+(?:href|alt)=["'][^"']*["'])?[^>]*>`).ReplaceAllStringFunc(html, func(match string) string {
		re := regexp.MustCompile(`<([a-zA-Z][a-zA-Z0-9]*)[^>]*?((?:\s+(?:href|alt)=["'][^"']*["'])*)`)
		matches := re.FindStringSubmatch(match)
		if len(matches) >= 3 {
			return "<" + matches[1] + matches[2] + ">"
		}
		return "<" + matches[1] + ">"
	})

	// タグとマークダウン記法の変換
	replacer := []struct {
		reOpen  *regexp.Regexp
		reClose *regexp.Regexp
		mdOpen  string
		mdClose string
	}{
		{regexp.MustCompile(`(?i)<h1\b[^>]*>`), regexp.MustCompile(`(?i)</h1>`), "# ", "\n\n"},
		{regexp.MustCompile(`(?i)<h2\b[^>]*>`), regexp.MustCompile(`(?i)</h2>`), "## ", "\n\n"},
		{regexp.MustCompile(`(?i)<h3\b[^>]*>`), regexp.MustCompile(`(?i)</h3>`), "### ", "\n\n"},
		{regexp.MustCompile(`(?i)<h4\b[^>]*>`), regexp.MustCompile(`(?i)</h4>`), "#### ", "\n\n"},
		{regexp.MustCompile(`(?i)<h5\b[^>]*>`), regexp.MustCompile(`(?i)</h5>`), "##### ", "\n\n"},
		{regexp.MustCompile(`(?i)<h6\b[^>]*>`), regexp.MustCompile(`(?i)</h6>`), "###### ", "\n\n"},
		{regexp.MustCompile(`(?i)<ul\b[^>]*>`), regexp.MustCompile(`(?i)</ul>`), "", "\n"},
		{regexp.MustCompile(`(?i)<ol\b[^>]*>`), regexp.MustCompile(`(?i)</ol>`), "", "\n"},
		{regexp.MustCompile(`(?i)<li\b[^>]*>`), regexp.MustCompile(`(?i)</li>`), "- ", "\n"},
		{regexp.MustCompile(`(?i)<table\b[^>]*>`), regexp.MustCompile(`(?i)</table>`), "\n", "\n\n"},
		{regexp.MustCompile(`(?i)<caption\b[^>]*>`), regexp.MustCompile(`(?i)</caption>`), "", "\n"},
		{regexp.MustCompile(`(?i)<tr\b[^>]*>`), regexp.MustCompile(`(?i)</tr>`), "", "|\n"},
		{regexp.MustCompile(`(?i)<th\b[^>]*>`), regexp.MustCompile(`(?i)</th>`), "| ", " "},
		{regexp.MustCompile(`(?i)<td\b[^>]*>`), regexp.MustCompile(`(?i)</td>`), "| ", " "},
		{regexp.MustCompile(`(?i)<p\b[^>]*>`), regexp.MustCompile(`(?i)</p>`), "", "\n\n"},
		{regexp.MustCompile(`(?i)<div\b[^>]*>`), regexp.MustCompile(`(?i)</div>`), "", "\n\n"},
	}
	markdown = html
	for _, r := range replacer {
		markdown = r.reOpen.ReplaceAllString(markdown, r.mdOpen)
		markdown = r.reClose.ReplaceAllString(markdown, r.mdClose)
	}

	// 特殊なタグの処理
	// a タグはリンクをマークダウン形式に変換
	markdown = regexp.MustCompile(`(?i)<a\b[^>]*>(.*?)</a>`).ReplaceAllStringFunc(markdown, func(match string) string {
		re := regexp.MustCompile(`(?i)<a\s[^>]*?href=["']([^"']+)["'][^>]*?>(.*?)<\/a>`)
		matches := re.FindStringSubmatch(match)
		if len(matches) < 3 {
			return match // マッチしない場合はそのまま返す
		}
		url := matches[1]
		text := matches[2]

		// リンクが相対パスの場合は絶対パスに変換
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + targetDomain + url
		}
		return fmt.Sprintf(" [%s](%s) ", text, url) // マークダウン形式に変換
	})

	// img タグは画像のリンクをマークダウン形式に変換
	markdown = regexp.MustCompile(`(?i)<img\b[^>]*?src=["']([^"']+)["'][^>]*?>`).ReplaceAllStringFunc(markdown, func(match string) string {
		reSrc := regexp.MustCompile(`(?i)src=["']([^"']+)["']`)
		reAlt := regexp.MustCompile(`(?i)alt=["']([^"']+)["']`)
		srcMatch := reSrc.FindStringSubmatch(match)
		altMatch := reAlt.FindStringSubmatch(match)
		if len(srcMatch) < 2 {
			return match // src 属性がない場合はそのまま返す
		}
		src := srcMatch[1]
		alt := "image"
		if len(altMatch) >= 2 {
			alt = altMatch[1] // alt 属性があればそれを使用
		}
		// リンクが相対パスの場合は絶対パスに変換
		if !strings.HasPrefix(src, "http://") && !strings.HasPrefix(src, "https://") {
			src = "https://" + targetDomain + src
		}
		return fmt.Sprintf(" ![%s](%s) ", alt, src) // マークダウン形式に変換
	})

	// 残ったタグをすべて除去
	reTag := regexp.MustCompile(`(?s)<[^>]+>`)
	markdown = reTag.ReplaceAllString(markdown, "")

	// html 中でよく使われる特殊文字を置換
	markdown = regexp.MustCompile(`(&nbsp;|&#160;)`).ReplaceAllString(markdown, " ")
	markdown = regexp.MustCompile(`(&amp;|&#38;)`).ReplaceAllString(markdown, "&")
	markdown = regexp.MustCompile(`(&lt;|&#60;)`).ReplaceAllString(markdown, "<")
	markdown = regexp.MustCompile(`(&gt;|&#62;)`).ReplaceAllString(markdown, ">")
	markdown = regexp.MustCompile(`(&quot;|&#34;)`).ReplaceAllString(markdown, "\"")
	markdown = regexp.MustCompile(`(&apos;|&#39;)`).ReplaceAllString(markdown, "'")

	// マークダウンを整形
	markdown = strings.TrimSpace(markdown)                                         // 先頭と末尾の空白を削除
	markdown = regexp.MustCompile(`[^\S\n]{2,}`).ReplaceAllString(markdown, " ")   // 連続する空白を1つにまとめる
	markdown = regexp.MustCompile(`\n\s{1,}\n`).ReplaceAllString(markdown, "\n\n") // 連続する改行の前後に空白がある場合は削除
	markdown = markdown + "\n"                                                     // 最後に改行を追加

	return markdown
}
