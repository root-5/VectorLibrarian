// 主にスプレッドシートからの利用を想定したAPIを提供する
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"nlp/vectorize"
)

// ====================================================================================
// 定数と構造体定義
// ====================================================================================

var port = "8000"

// リクエスト用の構造体
type ConvertRequest struct {
	Text    string `json:"text"`
	IsQuery bool   `json:"is_query"`
}

// NLPサーバーからのレスポンス用の構造体
// app/controller/nlp/nlp.go と同じ構造体
type ConvertResponse struct {
	MaxTokenLength     int         `json:"max_token_length"`
	OverlapTokenLength int         `json:"overlap_token_length"`
	ModelName          string      `json:"model_name"`
	ModelVectorLength  int         `json:"model_vector_length"`
	Chunks             []string    `json:"chunks"`
	Vectors            [][]float32 `json:"vectors"`
}

// ====================================================================================
// API関数
// ====================================================================================

// APIサーバーを起動する関数
func StartServer() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("サーバーの起動に失敗しました: %v\n", err)
	}
}

// リクエストを処理する関数
func handler(w http.ResponseWriter, r *http.Request) {
	// リクエストのメソッドによって処理を分岐
	switch r.Method {
	case "POST":
		postHandler(w, r)
	default:
		fmt.Fprintf(w, "Method not allowed")
	}
}

// POSTリクエストを処理する関数
func postHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストパスを取得
	path := r.URL.Path
	fmt.Println("Access:", path)

	// リクエストパスによって処理を分岐
	switch path {
	case "/convert":
		var req ConvertRequest

		// リクエストボディを構造体に変換
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			fmt.Printf("リクエストボディのデコードエラー: %v\n", err)
			http.Error(w, "無効なリクエストボディ", http.StatusBadRequest)
			return
		}

		chunks, vectors, err := vectorize.ConvertToVector(req.Text, req.IsQuery)
		if err != nil {
			fmt.Printf("ベクトル化エラー: %v\n", err)
			http.Error(w, fmt.Sprintf("ベクトル化エラー: %v", err), http.StatusInternalServerError)
			return
		}

		// 環境変数からの値を取得し、整数に変換
		maxTokenLengthStr := os.Getenv("MAX_TOKEN_LENGTH")
		overlapTokenLengthStr := os.Getenv("OVERLAP_TOKEN_LENGTH")
		modelVectorLengthStr := os.Getenv("MODEL_VECTOR_LENGTH")
		maxTokenLength, err := strconv.Atoi(maxTokenLengthStr)
		if err != nil {
			fmt.Printf("MAX_TOKEN_LENGTH の変換エラー: %v\n", err)
			http.Error(w, fmt.Sprintf("MAX_TOKEN_LENGTH の変換エラー: %v", err), http.StatusInternalServerError)
			return
		}
		overlapTokenLength, err := strconv.Atoi(overlapTokenLengthStr)
		if err != nil {
			fmt.Printf("OVERLAP_TOKEN_LENGTH の変換エラー: %v\n", err)
			http.Error(w, fmt.Sprintf("OVERLAP_TOKEN_LENGTH の変換エラー: %v", err), http.StatusInternalServerError)
			return
		}
		modelVectorLength, err := strconv.Atoi(modelVectorLengthStr)
		if err != nil {
			fmt.Printf("MODEL_VECTOR_LENGTH の変換エラー: %v\n", err)
			http.Error(w, fmt.Sprintf("MODEL_VECTOR_LENGTH の変換エラー: %v", err), http.StatusInternalServerError)
			return
		}

		// レスポンスを構造体に変換
		response := ConvertResponse{
			MaxTokenLength:     maxTokenLength,
			OverlapTokenLength: overlapTokenLength,
			ModelVectorLength:  modelVectorLength,
			ModelName:          os.Getenv("MODEL_NAME"),
			Chunks:             chunks,
			Vectors:            vectors,
		}

		// レスポンスをJSON形式で返す
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	default:
		fmt.Fprintf(w, "Not found")
	}
}
