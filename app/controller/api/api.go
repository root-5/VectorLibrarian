// 主にスプレッドシートからの利用を想定したAPIを提供する
package api

import (
	"app/controller/log"
	"fmt"
	"net/http"
)

// ====================================================================================
// API関数
// ====================================================================================

var port = "8080"

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

	// ベクトル検索
	case "/search":
		searchHandler(w, r)

	// RAG検索（ベクトル検索 + OpenAI API）- ストリーミング対応
	case "/rag_search":
		ragSearchHandler(w, r)

	// 静的ファイル
	case "/favicon.ico":
		http.ServeFile(w, r, "controller/api/public/smile.ico")
	case "/style.css":
		http.ServeFile(w, r, "controller/api/public/style.css")
	case "/script.js":
		http.ServeFile(w, r, "controller/api/public/script.js")

	// HTML
	case "/chat":
		http.ServeFile(w, r, "controller/api/public/chat.html")
	case "/":
		http.ServeFile(w, r, "controller/api/public/index.html")

	default:
		// アクセス元のIPアドレスを取得
		// 将来的には複数回アクセスがあった場合に、そのIPアドレスをブロックするようにする
		ip := r.RemoteAddr
		log.Info(path)
		log.Info("Not found: " + ip)
		fmt.Fprintf(w, "Not found")
	}
}

// POSTリクエストを処理する関数
func postHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストパスを取得
	path := r.URL.Path

	// リクエストパスによって処理を分岐
	switch path {
	case "/":
		fmt.Fprintf(w, "Hello, world")
	default:
		fmt.Fprintf(w, "Not found")
	}
}

// ====================================================================================
// リクエストボディの処理関数
// ====================================================================================
// リクエストボディを構造体に変換する関数
// func decodeRequestBody(r *http.Request, v interface{}) error {
// 	// リクエストボディを読み込む
// 	err := r.ParseForm()
// 	if err != nil {
// 		return err
// 	}

// 	// リクエストボディをJSONに変換
// 	body := r.Form.Get("body")
// 	err = json.Unmarshal([]byte(body), v)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
