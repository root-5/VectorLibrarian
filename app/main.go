package main

import (
	"app/controller/log"
	"app/controller/postgres"
	"app/usecase/processor"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
	_ "github.com/lib/pq"
)

// NLPサーバーへのリクエスト用の構造体
type ConvertRequest struct {
	Text    string `json:"text"`
	IsQuery bool   `json:"is_query"`
}

// NLPサーバーからのレスポンス用の構造体
type ConvertResponse struct {
	InputText      string    `json:"input_text"`
	NormalizedText string    `json:"normalized_text"`
	IsQuery        bool      `json:"is_query"`
	ModelName      string    `json:"model_name"`
	Dimensions     int       `json:"dimensions"`
	Vector         []float32 `json:"vector"`
}

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
			// return // 既に保存されているハッシュがあればスキップ
		}

		// ページデータをデータベースに保存
		err = postgres.SaveCrawledData(domain, path, pageTitle, description, keywords, markdown, hash)
		if err != nil {
			return
		}

		// テキスト正規化

		// チャンク化

		// ベクトル化
		result, err := requestNlpServer(markdown, false)
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("NLP サーバーからのレスポンスの中身: " + fmt.Sprintf("%+v", result.Vector))
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

func requestNlpServer(text string, isQuery bool) (*ConvertResponse, error) {
	// リクエストボディを作成
	requestBody := ConvertRequest{
		Text:    text,
		IsQuery: isQuery,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("JSONエンコードに失敗: %w", err)
	}

	// POSTリクエストを送信
	resp, err := http.Post("http://"+os.Getenv("NLP_HOST")+":8000/convert", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("NLP サーバーへのリクエストに失敗: %w", err)
	}
	defer resp.Body.Close()

	log.Info("NLP request: " + resp.Status)

	// 構造体に直接デコード
	var result ConvertResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("NLP サーバーからのレスポンスの解析に失敗: %w", err)
	}

	return &result, nil
}
