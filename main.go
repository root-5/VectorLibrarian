package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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
	targetDomain := "hotel-example-site.takeyaqa.dev" // ターゲットドメインを設定
	logFilePath := "log/scraper.md"                   // ログファイルのパスを設定

	// =======================================================================
	// ログの設定
	// =======================================================================
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("ログファイルの作成に失敗しました:", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	err = os.Truncate(logFilePath, 0)
	if err != nil {
		fmt.Println("ログファイルの初期化に失敗しました:", err)
	}

	// =======================================================================
	// データベース接続
	// =======================================================================
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	psqlInfo := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("データベースへの接続に失敗しました:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("データベースへのPingに失敗しました:", err)
	}
	log.Println("データベースへの接続に成功しました。")

	// =======================================================================
	// テーブル作成
	// =======================================================================
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS pages (
        id SERIAL PRIMARY KEY,
        url TEXT NOT NULL UNIQUE,
        title TEXT,
        description TEXT,
        keywords TEXT,
        markdown_content TEXT,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("テーブルの作成に失敗しました:", err)
	}
	log.Println("テーブル'pages'を正常に作成または確認しました。")

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
		log.Println(">> URL:", r.URL.String())
	})

	// テキストコンテンツを抽出するためのコールバック
	c.OnHTML("html", func(e *colly.HTMLElement) {
		pageTitle := e.ChildText("title") // ページのタイトルを取得
		if pageTitle == "" {
			pageTitle = "--"
		}
		description := e.ChildAttr("meta[name=description]", "content") // ディスクリプションを取得
		if description == "" {
			description = "--"
		}
		keywords := e.ChildAttr("meta[name=keywords]", "content") // キーワードを取得
		if keywords == "" {
			keywords = "--"
		}
		log.Printf(">> Page Info:\n- Title: %s\n- Description: %s\n- Keywords: %s", pageTitle, description, keywords)

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
		markdown, err := conv.ConvertString(html)
		if err != nil {
			fmt.Println("マークダウンへの変換に失敗しました:", err)
			return
		}
		log.Println(">> markdown:\n" + markdown)

		// =======================================================================
		// データベースに保存
		// =======================================================================
		insertSQL := `
        INSERT INTO pages (url, title, description, keywords, markdown_content)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (url) DO UPDATE SET
            title = EXCLUDED.title,
            description = EXCLUDED.description,
            keywords = EXCLUDED.keywords,
            markdown_content = EXCLUDED.markdown_content,
            created_at = CURRENT_TIMESTAMP;`

		_, err = db.Exec(insertSQL, e.Request.URL.String(), pageTitle, description, keywords, markdown)
		if err != nil {
			log.Printf("データベースへのデータ挿入に失敗しました: %v", err)
		} else {
			log.Println("データベースにデータを保存しました:", e.Request.URL.String())
		}
	})

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
