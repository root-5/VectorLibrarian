// 主にスプレッドシートからの利用を想定したAPIを提供する
package api

import (
	"app/controller/log"
	"app/usecase/usecase"
	"encoding/json"
	"fmt"
	"net/http"
)

// ====================================================================================
// 各エンドポイントのハンドラ関数
// ====================================================================================

// ベクトル検索
func searchHandler(w http.ResponseWriter, r *http.Request) {
	// 検索クエリを取得
	query := r.URL.Query().Get("q")
	if query == "" {
		log.Info("query parameter 'q' is required")
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// ベクトル検索を実行
	resultLimit := 20
	similarPages, err := usecase.VectorSearch(query, resultLimit)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJsonResponse(w, similarPages)
}

// RAG検索（ベクトル検索 + OpenAI API）- ストリーミング対応
func ragSearchHandler(w http.ResponseWriter, r *http.Request) {
	// 検索クエリを取得
	query := r.URL.Query().Get("q")
	if query == "" {
		log.Info("query parameter 'q' is required")
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// ベクトル検索を実行（上位5件）
	resultLimit := 5
	similarPages, err := usecase.VectorSearch(query, resultLimit)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 検索結果のMarkdownを収集
	contextMarkdowns := make([]string, 0, len(similarPages))
	for _, page := range similarPages {
		contextMarkdowns = append(contextMarkdowns, page.Markdown)
	}

	// SSE（Server-Sent Events）のヘッダーを設定
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// まず参照元情報を送信
	sourcesJSON, _ := json.Marshal(map[string]interface{}{
		"sources":      similarPages,
		"source_count": len(similarPages),
	})
	fmt.Fprintf(w, "data: {\"type\":\"sources\",\"data\":%s}\n\n", sourcesJSON)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	// OpenAI APIでRAG応答をストリーミング生成
	fmt.Fprintf(w, "data: {\"type\":\"start\"}\n\n")
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	err = generateRAGResponseStream(query, contextMarkdowns, w)
	if err != nil {
		log.Error(err)
		fmt.Fprintf(w, "data: {\"type\":\"error\",\"message\":\"%s\"}\n\n", err.Error())
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		return
	}

	// 完了メッセージを送信
	fmt.Fprintf(w, "data: {\"type\":\"done\"}\n\n")
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// ====================================================================================
// レスポンスの処理関数
// ====================================================================================
// 構造体をjson形式の文字列に変換してレスポンスを返す関数
func sendJsonResponse(w http.ResponseWriter, data interface{}) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
