package main

import (
	"app/controller/log"
	"app/controller/nlp"
	"app/controller/postgres"
	"app/usecase/processor"
	"time"

	"github.com/gocolly/colly/v2"
	_ "github.com/lib/pq"
)

func main() {
	// =======================================================================
	// 初期設定・定数
	// =======================================================================
	targetDomain := "www.city.hamura.tokyo.jp" // ターゲットドメインを設定
	allowedPaths := []string{                  // URLパスの制限（特定のパス以外をスキップ）
		"/prsite/",
	}
	maxScrapeDepth := 1        // 最大スクレイピング深度を設定
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
	// Colly のコレクターを作成
	// =======================================================================
	// デフォルトのコレクターを作成
	c := colly.NewCollector(
		colly.AllowedDomains(targetDomain), // 許可するドメインを設定
		colly.MaxDepth(maxScrapeDepth),     // 最大深度を設定
		colly.CacheDir(collyCacheDir),      // キャッシュディレクトリを指定
	)

	// リクエスト間で 1 秒の時間を空ける
	c.Limit(&colly.LimitRule{
		DomainGlob: targetDomain, // 対象ドメインを指定
		Delay:      time.Second,  // リクエスト間の最小遅延
	})

	// リクエスト前に "アクセス >> " を表示
	c.OnRequest(func(r *colly.Request) {
		log.Info(">> URL:" + r.URL.String())
	})

	// html タグを見つけたときの処理
	c.OnHTML("html", func(e *colly.HTMLElement) {
		// ページデータを抽出
		domain, path, pageTitle, description, keywords, markdown, hash, err := processor.HtmlToPageData(e)
		if err != nil {
			return
		}

		// ハッシュ値を照合
		exists, err := postgres.CheckHashExists(hash)
		if err != nil {
			log.Error(err)
			return
		}
		if exists {
			return // 既に保存されているハッシュがあればスキップ
		}

		// model に記載した通り、見出しをマークダウンから抽出して箇条書きに変換
		itemization := processor.ExtractHeadings(markdown)

		// 箇条書きをテキスト正規化、ベクトル化のリクエストを NLP サーバーに送信
		vector, err := nlp.ConvertToVector(itemization, false)
		if err != nil {
			log.Error(err)
			return
		}

		// ページデータをデータベースに保存
		err = postgres.SaveCrawledData(domain, path, pageTitle, description, keywords, markdown, hash, vector)
		if err != nil {
			return
		}
	})

	// a タグを見つけたときの処理
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// URL を取得
		url, isValid := processor.GetLinkUrl(e, targetDomain, allowedPaths)
		if !isValid {
			return // 無効なリンクはスキップ
		}

		// ページ内で見つかったリンクを訪問
		e.Request.Visit(url)
	})

	// 指定ドメインからスクレイピングを開始
	c.Visit("https://" + targetDomain + "/")
}
