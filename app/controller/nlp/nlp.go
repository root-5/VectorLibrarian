package nlp

import (
	"app/controller/log"
	"bytes"
	"encoding/json"
	"net/http"
	"os"
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

/*
nlp サーバーにテキストを送信してベクトルに変換する関数
正規化も nlp サーバー側で行う
  - text)		変換するテキスト
  - isQuery)	クエリかどうかの真偽値（True なら「query: 」、False なら「passage: 」のプレフィックスが文頭に付与される）
  - return)		変換結果の構造体とエラー
*/
func ConvertToVector(text string, isQuery bool) ([]float32, error) {
	// リクエストボディを作成
	requestBody := ConvertRequest{
		Text:    text,
		IsQuery: isQuery,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// POSTリクエストを送信
	resp, err := http.Post("http://"+os.Getenv("NLP_HOST")+":8000/convert", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Info("NLP request: " + resp.Status)

	// 構造体に直接デコード
	var result ConvertResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error(err)
		return nil, err
	}

	return result.Vector, nil
}
