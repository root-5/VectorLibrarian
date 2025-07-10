package main

import (
	"app/controller/log"
	"app/controller/postgres"
	"crypto/sha1"
	"encoding/hex"
	"regexp"
	"strings"
	"time"

	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/base"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/commonmark"
	"github.com/JohannesKaufmann/html-to-markdown/v2/plugin/table"
	"github.com/gocolly/colly/v2"
	_ "github.com/lib/pq"
)

func main() {
	targetDomain := "www.city.hamura.tokyo.jp" // ターゲットドメインを設定
	allowedPaths := []string{                  // URLパスの制限（特定のパス以外をスキップ）
		"/prsite/",
	}
	collyCacheDir := "./cache" // Colly のキャッシュディレクトリを設定

	// =======================================================================
	// データベース接続とテーブル初期化
	// =======================================================================
	err := postgres.Connect()
	if err != nil {
		return
	}
	err = postgres.InitTable()
	if err != nil {
		return
	}

	// =======================================================================
	// html-to-markdown のコンバーターを作成
	// =======================================================================
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

	// =======================================================================
	// Colly のコレクターを作成
	// =======================================================================
	// デフォルトのコレクターを作成
	c := colly.NewCollector(
		colly.AllowedDomains(targetDomain), // 許可するドメインを設定
		colly.MaxDepth(2),                  // 最大深度を 2 に設定
		colly.CacheDir(collyCacheDir),      // キャッシュディレクトリを指定
		// colly.Async(true),                  // 非同期モードを有効にする
		// colly.IgnoreRobotsTxt(),            // robots.txt を無視
	)

	// リクエスト間で1~2秒の時間を空ける
	c.Limit(&colly.LimitRule{
		DomainGlob: targetDomain, // 対象ドメインを指定
		Delay:      time.Second,  // リクエスト間の最小遅延
	})

	// リクエスト前に "アクセス >> " を表示
	c.OnRequest(func(r *colly.Request) {
		log.Info(">> URL:" + r.URL.String())
	})

	// テキストコンテンツを抽出するためのコールバック
	c.OnHTML("html", func(e *colly.HTMLElement) {
		domain := e.Request.URL.Hostname()
		path := e.Request.URL.Path

		// ページタイトル、ディスクリプション、キーワードを取得（それぞれ存在しない場合は "--" を設定）
		pageTitle := e.ChildText("title")
		if pageTitle == "" {
			pageTitle = "--"
		}
		description := e.ChildAttr("meta[name=description]", "content")
		if description == "" {
			description = "--"
		}
		keywords := e.ChildAttr("meta[name=keywords]", "content")
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
		markdown, err := conv.ConvertString(html)
		if err != nil {
			log.Error(err)
			return
		}

		// markdown のハッシュを計算
		hashBin := sha1.Sum([]byte(markdown))
		hash := hex.EncodeToString(hashBin[:])

		// =======================================================================
		// データベースに保存
		// =======================================================================
		err = postgres.SaveCrawledData(domain, path, pageTitle, description, keywords, markdown, hash)
		if err != nil {
			return
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// .pdf で終わるリンク、mailto:/javascript:/# 始まるリンク、空のリンクはスキップ
		if matched, _ := regexp.MatchString(`(?i)\.pdf$|^mailto:|^javascript:|^$|^#`, link); matched {
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
			return
		}
		// 特定のパス以外をスキップ
		for _, allowedPath := range allowedPaths {
			if !strings.Contains(link, allowedPath) {
				return
			}
		}

		// ページ内で見つかったリンクを訪問
		e.Request.Visit(link)
	})

	// 指定ドメインからスクレイピングを開始
	c.Visit("https://" + targetDomain + "/")
}
