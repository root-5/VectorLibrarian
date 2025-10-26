// OpenAI APIを利用するための関数をまとめたパッケージ
package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// OpenAI APIのリクエスト構造体
type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAI APIのレスポンス構造体
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

/*
OpenAI APIを呼び出してRAG応答を生成する関数
  - query				ユーザーの質問
  - contextMarkdowns	検索結果のMarkdownコンテンツ（上位3件など）
  - return)				生成された回答
  - return) err			エラー
*/
func GenerateRAGResponse(query string, contextMarkdowns []string) (answer string, err error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	modelName := os.Getenv("OPENAI_MODEL_NAME")

	// コンテキストを結合
	context := ""
	for i, md := range contextMarkdowns {
		context += fmt.Sprintf("## 参照情報 %d\n%s\n\n", i+1, md)
	}

	// システムプロンプトとユーザーメッセージを構築
	systemPrompt := `あなたは羽村市の公式サイト情報に基づいて質問に答えるアシスタントです。
以下の参照情報を基に、ユーザーの質問に正確かつ簡潔に日本語で答えてください。
参照情報に含まれていない内容については、「提供された情報には含まれていません」と答えてください。`

	userMessage := fmt.Sprintf("参照情報:\n%s\n\n質問: %s", context, query)

	// リクエストボディを構築
	reqBody := ChatCompletionRequest{
		Model: modelName,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// HTTPクライアントを作成
	client := &http.Client{Timeout: 180 * time.Second}

	// リクエストを作成
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// リクエストを送信
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// レスポンスを読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error: %d - %s", resp.StatusCode, string(body))
	}

	// レスポンスをパース
	var chatResp ChatCompletionResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI API")
	}

	return chatResp.Choices[0].Message.Content, nil
}
