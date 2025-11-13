package nlp

import (
	"app/controller/log"
	"app/domain/model"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

// NLPサーバーへのリクエスト用の構造体
type ConvertRequest struct {
	Text    string `json:"text"`
	IsQuery bool   `json:"is_query"`
}

// NLPサーバーからのレスポンス用の構造体
// nlp/api/api.go と同じ構造体
type ConvertResponse struct {
	model.NlpConfigInfo
	Chunks  []string    `json:"chunks"`
	Vectors [][]float32 `json:"vectors"`
}

/*
nlp サーバーにテキストを送信してベクトルに変換する関数
正規化も nlp サーバー側で行う
  - text)		変換するテキスト
  - isQuery)	クエリかどうかの真偽値（True なら「query: 」、False なら「passage: 」のプレフィックスが文頭に付与される）
  - return)		最大トークン長、オーバーラップトークン長、モデル名、モデル特有のベクトル長、チャンクの配列、ベクトルの2次元配列、エラー
*/
func ConvertToVector(text string, isQuery bool) (resp ConvertResponse, err error) {
	// リクエストボディを作成
	requestBody := ConvertRequest{
		Text:    text,
		IsQuery: isQuery,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error(err)
		return ConvertResponse{}, err
	}

	// リクエスト URL とタイムアウト設定
	requestUrl := "http://" + os.Getenv("NLP_HOST") + ":" + os.Getenv("NLP_PORT") + "/convert"
	client := &http.Client{Timeout: 1200 * time.Second}

	// POSTリクエストを送信
	httpResp, err := client.Post(requestUrl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(err)
		return ConvertResponse{}, err
	}
	defer httpResp.Body.Close()

	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		log.Error(err)
		return ConvertResponse{}, err
	}

	// 構造体にデコード
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		log.Error(err)
		return ConvertResponse{}, err
	}

	return resp, nil
}
