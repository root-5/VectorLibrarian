// 主にスプレッドシートからの利用を想定したAPIを提供する
package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"nlp/vectorize"
)

// ====================================================================================
// API関数
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

// APIサーバーを起動する関数
func StartServer() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}

// リクエストを処理する関数
func handler(w http.ResponseWriter, r *http.Request) {
	// リクエストのメソッドによって処理を分岐
	switch r.Method {
	case "GET":
		getHandler(w, r)
	case "POST":
		postHandler(w, r)
	default:
		fmt.Fprintf(w, "Method not allowed")
	}
}

// GETリクエストを処理する関数
func getHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストパスを取得
	path := r.URL.Path

	// リクエストパスによって処理を分岐
	switch path {
	default:
		fmt.Fprintf(w, "Not found")
	}
}

// POSTリクエストを処理する関数
func postHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストパスを取得
	path := r.URL.Path

	// リクエストパスによって処理を分岐
	switch path {
	case "/convert":
		// リクエストボディを構造体に変換
		var request ConvertRequest
		if err := decodeRequestBody(r, &request); err != nil {
			http.Error(w, fmt.Sprintf("リクエストボディの解析に失敗: %v", err), http.StatusBadRequest)
			return
		}
		// ConvertToVector関数を呼び出してベクトルに変換
		vector, err := vectorize.TextToVector(request.Text, request.IsQuery)
		if err != nil {
			http.Error(w, fmt.Sprintf("ベクトル変換に失敗: %v", err), http.StatusInternalServerError)
			return
		}
		// レスポンス用の構造体を作成
		response := ConvertResponse{
			Vector: vector,
		}
		sendJsonResponse(w, response)

	default:
		fmt.Fprintf(w, "Not found")
	}
}

// ====================================================================================
// リクエストボディの処理関数
// ====================================================================================
// リクエストボディを構造体に変換する関数
func decodeRequestBody(r *http.Request, v interface{}) error {
	// リクエストボディを直接JSONとして読み込む
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

// ====================================================================================
// レスポンスの処理関数
// ====================================================================================
// 構造体をjson形式の文字列に変換してレスポンスを返す関数
func sendJsonResponse(w http.ResponseWriter, data interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(w, "JSON変換に失敗: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
