package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

// func main() {
// 	// デフォルトのコレクターを作成
// 	c := colly.NewCollector(
// 		// 訪問するドメインを hackerspaces.org, wiki.hackerspaces.org のみに制限
// 		colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
// 	)

// 	// href属性を持つa要素ごとにコールバックを実行
// 	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
// 		link := e.Attr("href")
// 		// リンクを表示
// 		fmt.Printf("リンク発見: %q -> %s\n", e.Text, link)
// 		// ページ内で見つかったリンクを訪問
// 		// AllowedDomains に含まれるリンクのみ訪問される
// 		c.Visit(e.Request.AbsoluteURL(link))
// 	})

// 	// リクエスト前に "Visiting ..." を表示
// 	c.OnRequest(func(r *colly.Request) {
// 		fmt.Println("訪問中", r.URL.String())
// 	})

// 	// https://hackerspaces.org からスクレイピングを開始
// 	c.Visit("https://nextjs-root-5.vercel.app/")
// }

func main() {
	// デフォルトのコレクターを作成
	c := colly.NewCollector(
		// 最大深度を1に設定（スクレイピング対象ページ内のリンクのみ訪問し、それ以上は辿らない）
		colly.AllowedDomains("nextjs-root-5.vercel.app"),
		colly.MaxDepth(1),
	)

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

		// ルートパス以外では header, footer, script タグを削除
		if e.Request.URL.Path != "/" {
			e.DOM.Find("header").Remove()
			e.DOM.Find("footer").Remove()
		}

		// head タグと script タグは常時削除
		e.DOM.Find("head").Remove()
		e.DOM.Find("script").Remove()

		// テキストコンテンツをマークダウン形式に変換しながら取得
		textContent := e.DOM.Text()
		// markdown = markdownify(textContent)

		fmt.Println("ページタイトル:", pageTitle)
		fmt.Println("")
		fmt.Println("ページディスクリプション:", description)
		fmt.Println("")
		fmt.Println("ページテキスト:", textContent)
		fmt.Println("")
		// fmt.Println("ページマークダウン:", markdown)
		// fmt.Println("")
	})

	// href属性を持つa要素ごとにコールバックを実行
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// リンクを表示
		fmt.Println("- リンク発見：", link)

		// ページ内で見つかったリンクを訪問
		// e.Request.Visit(link)

		// AllowedDomains に含まれるリンクのみ訪問される
		// c.Visit(e.Request.AbsoluteURL(link))
	})

	// リクエスト前に "アクセス >> " を表示
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("アクセス >> ", r.URL.String())
		fmt.Println("")
		fmt.Println("")
	})

	// https://nextjs-root-5.vercel.app/ からスクレイピングを開始
	c.Visit("https://nextjs-root-5.vercel.app/")
}
