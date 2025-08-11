// 主にスプレッドシートからの利用を想定したAPIを提供する
package api

import (
	"encoding/json"
	"fmt"
	"net/http"

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
type ConvertResponse struct {
	Vector []float32 `json:"vector"`
}

// ====================================================================================
// API関数
// ====================================================================================

// APIサーバーを起動する関数
func StartServer() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
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

	// リクエストパスによって処理を分岐
	switch path {
	case "/convert":
		var req ConvertRequest

		// リクエストボディを構造体に変換
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "無効なリクエストボディ", http.StatusBadRequest)
			return
		}

		// テキストをベクトル化
		vector, err := vectorize.TextToVector(req.Text, req.IsQuery)
		if err != nil {
			http.Error(w, fmt.Sprintf("ベクトル化エラー: %v", err), http.StatusInternalServerError)
			return
		}

		// レスポンスを構造体に変換
		response := ConvertResponse{Vector: vector}

		// レスポンスをJSON形式で返す
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	default:
		fmt.Fprintf(w, "Not found")
	}
}
