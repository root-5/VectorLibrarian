package crawler

import (
	"app/controller/log"
	"app/controller/nlp"
	"app/controller/postgres"
	"time"

	"github.com/gocolly/colly/v2"
	_ "github.com/lib/pq"
)

// スケジューラーから呼び出すための関数、一旦定数などもここで設定
func Start() (err error) {
	// 初期設定・定数
	targetDomain := "www.city.hamura.tokyo.jp"
	startPath := "/"
	allowedPaths := []string{
		"/",
	}
	maxScrapeDepth := 7
	isTest := false

	// クロールを開始
	err = CrawlDomain(targetDomain, startPath, allowedPaths, maxScrapeDepth, isTest)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

/*
対象ドメインをクロールする関数
  - targetDomain	クロール対象のドメイン
  - startPath		クロールを開始するパス
  - allowedPaths	パスに必ず含まれなければならない文字列のリスト
  - maxScrapeDepth	最大スクレイピング深度
  - isTest			テストモードの真偽値
  - return) err		エラー

※ allowedPaths について
["/docs/", "/articles/"] なら "~/docs/abc", "~/articles/xyz" は許可されるが "~/blog/123" は許可されない
["/"] 指定であれば全て許可される
*/
func CrawlDomain(targetDomain string, startPath string, allowedPaths []string, maxScrapeDepth int, isTest bool) (err error) {

	// デフォルトのコレクターを作成
	c := colly.NewCollector(
		colly.AllowedDomains(targetDomain), // 許可するドメインを設定
		colly.MaxDepth(maxScrapeDepth),     // 最大深度を設定
	)

	// Colly のキャッシュディレクトリを設定（テストモード時はキャッシュしない）
	if !isTest {
		c.CacheDir = "./cache"
	}

	// リクエスト間で 1 秒の時間を空ける
	c.Limit(&colly.LimitRule{
		DomainGlob: targetDomain, // 対象ドメインを指定
		Delay:      time.Second,  // リクエスト間の最小遅延
	})

	// リクエスト前に "アクセス >> " を表示
	c.OnRequest(func(r *colly.Request) {
		if isTest {
			log.Info(">> URL:" + r.URL.String())
		}
	})

	// html タグを見つけたときの処理
	c.OnHTML("html", func(e *colly.HTMLElement) {
		// ページデータを抽出
		domain, path, pageTitle, description, keywords, markdown, hash, err := htmlToPageData(e)
		if err != nil {
			log.Error(err)
			return
		}

		if isTest {
			log.Info(">> path:" + path)
			log.Info(">> pageTitle:" + pageTitle)
			log.Info(">> description:" + description)
			log.Info(">> keywords:" + keywords)
			// log.Info(">> markdown:" + markdown)
			log.Info(">> hash:" + hash)
			log.Info("\n")
		}

		// ハッシュ値を照合
		isHashExists, err := postgres.CheckHashExists(hash)
		if err != nil {
			log.Error(err)
			return
		}

		if isHashExists && !isTest {
			return // 既に保存されているハッシュがあればスキップ（テストモード時はスキップしない）
		}

		// model に記載した通り、見出しをマークダウンから抽出して箇条書きに変換
		itemization := extractHeadings(markdown)
		targetText := "## page title\n\n" + pageTitle +
			// "\n\n## page description\n\n" + description +
			// "\n\n## page keywords\n\n" + keywords +
			"\n\n## page itemization\n\n" + itemization

		// 箇条書きをテキスト正規化、ベクトル化のリクエストを NLP サーバーに送信
		vector, err := nlp.ConvertToVector(targetText, false)
		if err != nil {
			log.Error(err)
			return
		}

		// ページデータをデータベースに保存
		err = postgres.SaveCrawledData(domain, path, pageTitle, description, keywords, markdown, hash, vector)
		if err != nil {
			log.Error(err)
			return
		}
	})

	// a タグを見つけたときの処理
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// URL を取得
		url, isValid := validateAndFormatLinkUrl(e, targetDomain, allowedPaths)
		if !isValid {
			return // 無効なリンクはスキップ
		}

		// ページ内で見つかったリンクを訪問
		e.Request.Visit(url)
	})

	// 指定ドメインからスクレイピングを開始
	c.Visit("https://" + targetDomain + startPath)

	return nil
}
