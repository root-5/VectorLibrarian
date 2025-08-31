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
	// InputText      string      `json:"input_text"`
	// NormalizedText string      `json:"normalized_text"`
	// IsQuery        bool        `json:"is_query"`
	// ModelName      string      `json:"model_name"`
	// Dimensions     int         `json:"dimensions"`
	MaxTokenLength     int         `json:"max_token_length"`
	OverlapTokenLength int         `json:"overlap_token_length"`
	ModelName          string      `json:"model_name"`
	ModelVectorLength  int         `json:"model_vector_length"`
	Chunks             []string    `json:"chunks"`
	Vectors            [][]float32 `json:"vectors"`
}

/*
nlp サーバーにテキストを送信してベクトルに変換する関数
正規化も nlp サーバー側で行う
  - text)		変換するテキスト
  - isQuery)	クエリかどうかの真偽値（True なら「query: 」、False なら「passage: 」のプレフィックスが文頭に付与される）
  - return)		最大トークン長、オーバーラップトークン長、モデル名、モデル特有のベクトル長、チャンクの配列、ベクトルの2次元配列、エラー
*/
func ConvertToVector(text string, isQuery bool) (
	maxTokenLength int,
	overlapTokenLength int,
	modelName string,
	modelVectorLength int,
	chunks []string,
	vectors [][]float32,
	err error,
) {
	// リクエストボディを作成
	requestBody := ConvertRequest{
		Text:    text,
		IsQuery: isQuery,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Error(err)
		return 0, 0, "", 0, nil, nil, err
	}

	// POSTリクエストを送信
	resp, err := http.Post("http://"+os.Getenv("NLP_HOST")+":8000/convert", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(err)
		return 0, 0, "", 0, nil, nil, err
	}
	defer resp.Body.Close()

	// 構造体に直接デコード
	var result ConvertResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error(err)
		return 0, 0, "", 0, nil, nil, err
	}

	maxTokenLength = result.MaxTokenLength
	overlapTokenLength = result.OverlapTokenLength
	modelName = result.ModelName
	modelVectorLength = result.ModelVectorLength
	chunks = result.Chunks
	vectors = result.Vectors

	return maxTokenLength, overlapTokenLength, modelName, modelVectorLength, chunks, vectors, nil
}
