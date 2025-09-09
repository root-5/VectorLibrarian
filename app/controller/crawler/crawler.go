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
	// targetDomain := "www.city.hamura.tokyo.jp"
	startPath := "/"
	allowedPaths := []string{
		"/",
	}
	maxScrapeDepth := 7
	isTest := false

	// ドメイン情報を取得
	domains, err := postgres.GetDomains()
	if err != nil {
		log.Error(err)
		return err
	}
	if len(domains) == 0 {
		log.Info("クロール対象のドメインが存在しません。")
		return nil
	}

	// とりあえず最初のドメインだけクロール
	targetDomainId := domains[0].Id
	targetDomain := domains[0].Domain

	log.Info("クロール対象ドメイン: " + targetDomain)

	// クロールを開始
	err = CrawlDomain(targetDomainId, targetDomain, startPath, allowedPaths, maxScrapeDepth, isTest)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

/*
対象ドメインをクロールする関数
  - targetDomainId	クロール対象のドメインのID
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
func CrawlDomain(targetDomainId int64, targetDomain string, startPath string, allowedPaths []string, maxScrapeDepth int, isTest bool) (err error) {

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
		pageInfo, err := htmlToPageData(e)
		if err != nil {
			log.Error(err)
			return
		}

		// ドメインIDを設定
		pageInfo.DomainId = targetDomainId

		if isTest {
			log.Info(">> path:" + pageInfo.Path)
			log.Info(">> pageTitle:" + pageInfo.Title)
			log.Info(">> description:" + pageInfo.Description)
			log.Info(">> keywords:" + pageInfo.Keywords)
			// log.Info(">> markdown:" + pageInfo.Markdown)
			log.Info(">> hash:" + pageInfo.Hash)
			log.Info("\n")
		}

		// ハッシュ値を照合
		isHashExists, err := postgres.CheckHashExists(pageInfo.Hash)
		if err != nil {
			log.Error(err)
			return
		}

		if isHashExists && !isTest {
			return // 既に保存されているハッシュがあればスキップ（テストモード時はスキップしない）
		}

		// 箇条書きをテキスト正規化、ベクトル化のリクエストを NLP サーバーに送信
		convertResult, err := nlp.ConvertToVector(pageInfo.Markdown, false)
		if err != nil {
			log.Error(err)
			return
		}

		// ページデータをデータベースに保存
		err = postgres.SaveCrawledData(pageInfo, convertResult)
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
