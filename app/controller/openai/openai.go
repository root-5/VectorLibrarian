// OpenAI APIを利用するための関数をまとめたパッケージ
package openai

import (
	"app/controller/log"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// OpenAI APIのリクエスト構造体
type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAI APIのストリーミングレスポンス構造体
type StreamResponse struct {
	Choices []StreamChoice `json:"choices"`
}

type StreamChoice struct {
	Delta        Delta  `json:"delta"`
	FinishReason string `json:"finish_reason"`
}

type Delta struct {
	Content string `json:"content"`
}

/*
OpenAI APIを呼び出してRAG応答をストリーミングで生成する関数
  - query				ユーザーの質問
  - contextMarkdowns	検索結果のMarkdownコンテンツ（上位3件など）
  - writer				ストリーミング結果を書き込むWriter
  - return) err			エラー
*/
func GenerateRAGResponseStream(query string, contextMarkdowns []string, writer io.Writer) error {
	apiKey := os.Getenv("OPENAI_API_KEY")
	modelName := os.Getenv("OPENAI_MODEL_NAME")
	if modelName == "" {
		modelName = "gpt-4o-mini"
	}

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

	// リクエストボディを構築（ストリーミング有効化）
	reqBody := ChatCompletionRequest{
		Model: modelName,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
		Stream: true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Error(err)
		return err
	}

	// HTTPクライアントを作成
	client := &http.Client{Timeout: 180 * time.Second}

	// リクエストを作成
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// リクエストを送信
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("OpenAI API error: %d - %s", resp.StatusCode, string(body))
	}

	// ストリーミングレスポンスを読み取り
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Error(err)
			return err
		}

		// "data: " プレフィックスを削除
		lineStr := strings.TrimSpace(string(line))
		if !strings.HasPrefix(lineStr, "data: ") {
			continue
		}
		lineStr = strings.TrimPrefix(lineStr, "data: ")

		// [DONE] で終了
		if lineStr == "[DONE]" {
			break
		}

		// JSONをパース
		var streamResp StreamResponse
		if err := json.Unmarshal([]byte(lineStr), &streamResp); err != nil {
			continue
		}

		// コンテンツを書き込み
		if len(streamResp.Choices) > 0 && streamResp.Choices[0].Delta.Content != "" {
			content := streamResp.Choices[0].Delta.Content
			// JSON文字列として正しくエンコード
			jsonContent, err := json.Marshal(content)
			if err != nil {
				continue
			}
			// SSE形式でデータを送信
			fmt.Fprintf(writer, "data: %s\n\n", string(jsonContent))
			if flusher, ok := writer.(http.Flusher); ok {
				flusher.Flush()
			}
		}
	}

	return nil
}
