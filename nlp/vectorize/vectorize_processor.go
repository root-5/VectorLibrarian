package vectorize

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/daulet/tokenizers"
	"github.com/yalue/onnxruntime_go"
	"golang.org/x/text/unicode/norm"
)

/*
replaceLinks は Markdown 中のリンクと画像リンクを除去し、リンクテキスト、alt テキストに置換して返す。
  - markdown	マークダウン形式のテキスト
  - return)		リンクテキスト、alt テキストに置換した文字列
*/
func replaceLinks(markdown string) (replacedMarkdown string) {
	// 画像を先に置換: ![alt](url) -> alt
	imgRe := regexp.MustCompile(`!\[([^\]]*)\]\([^)]*\)`)
	replacedMarkdown = imgRe.ReplaceAllString(markdown, "$1")

	// 通常のリンクを置換: [text](url) -> text
	linkRe := regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
	replacedMarkdown = linkRe.ReplaceAllString(replacedMarkdown, "$1")

	return replacedMarkdown
}

// テキストの正規化関数
func normalizeText(text string) (normalizedText string) {
	// Unicode正規化（NFKC）、全角カタカナに統一、半角英数に統一など
	normalizedText = norm.NFKC.String(text)

	// 小文字化
	normalizedText = strings.ToLower(normalizedText)

	// 3 回以上の連続する改行を 2 回の改行に置換
	normalizedText = regexp.MustCompile(`(\n{3,})`).ReplaceAllString(normalizedText, "\n\n")

	// 「., !, ?, 、, 。」が複数連続する場合は1つに置換
	punctRe := regexp.MustCompile(`([.!?、。ー])(?:[.!?、。ー])+`)
	normalizedText = punctRe.ReplaceAllString(normalizedText, "$1")

	// 最初と最後の空白と改行を削除
	normalizedText = strings.TrimSpace(normalizedText)

	return normalizedText
}

// チャンキング関数
func chunkText(text string, maxToken int, overlapMaxToken int) (chunks []string) {
	// 簡易的に文単位で分割、正規化段階で句読点の連続を1つにしているため、ここでは単純に句点と改行で分割
	sentences := regexp.MustCompile(`(?m)([^\n。!?]*[。!?\n]|[^\n。!?]+$)`).FindAllString(text, -1)

	var currentTokenCount = 0 // 現在のトークン数（計算用）
	var currentChunk []string // 現在のチャンク（計算用）
	var tokenCounts = []int{} // 各文のトークン数を保持するためのスライス

	for i, sentence := range sentences {
		// トークン数を簡易的に文字数で代用（正確にはトークナイザーで計測するのが望ましい）
		tokenIds, err := tokenize(sentence)
		if err != nil {
			fmt.Printf("トークナイズエラー: %v\n", err)
			return nil
		}
		sentenceTokenCount := len(tokenIds)

		// 文追加で最大トークン数を上回るときは、チャンクを保存してオーバーラップ分を持ち越す
		if currentTokenCount+sentenceTokenCount > maxToken {
			// 現在のチャンクを保存
			chunks = append(chunks, strings.Join(currentChunk, ""))

			// オーバーラップ分のトークン数を計算し、持ち越し分として次の現在のトークン数とする
			overlapTokenCount := 0
			roopCount := 0
			for j := 0; j < len(tokenCounts); j++ {
				additionalTokenCount := tokenCounts[len(tokenCounts)-1-j]
				if overlapTokenCount+additionalTokenCount >= overlapMaxToken {
					break
				} else {
					overlapTokenCount += additionalTokenCount
				}
				roopCount++
			}
			currentTokenCount = overlapTokenCount

			// オーバーラップ計算時のループ数から持ち越す文を計算、次の現在のチャンクとする
			overlapChunk := []string{}
			for j := 0; j < roopCount; j++ {
				overlapChunk = append([]string{sentences[i-1-j]}, overlapChunk...)
			}
			currentChunk = overlapChunk
		}

		// 文とトークン数を現在の値に追加
		currentTokenCount += sentenceTokenCount
		currentChunk = append(currentChunk, sentence)

		// 各文のトークン数に新しい文のトークン数を追加
		tokenCounts = append(tokenCounts, sentenceTokenCount)
	}

	// 最後のチャンクを追加
	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.Join(currentChunk, ""))
	}

	return chunks
}

// トークナイズ関数
func tokenize(text string) (ids []uint32, err error) {
	// 環境変数から設定を読み込む
	tokenizerPath := os.Getenv("DOWNLOAD_DIR") + "/" + os.Getenv("SAVED_TOKENIZER_PATH")

	tokenizerData, err := os.ReadFile(tokenizerPath)
	if err != nil {
		// エラーハンドリング
		return nil, err
	}

	// 最大トークン数を環境変数から取得
	maxTokenLengthStr := os.Getenv("MAX_TOKEN_LENGTH")
	maxTokenLength, err := strconv.ParseUint(maxTokenLengthStr, 10, 32)
	if err != nil {
		fmt.Printf("MAX_TOKEN_LENGTH の変換エラー: %v\n", err)
		return nil, err
	}

	// トランケーション（長すぎるトークンの切り捨て）方向は右側
	tk, err := tokenizers.FromBytesWithTruncation(tokenizerData, uint32(maxTokenLength), tokenizers.TruncationDirectionRight)
	if err != nil {
		fmt.Printf("tokenizer.json ロードエラー: %v\n", err)
		return nil, err
	}
	defer tk.Close()

	// トークン化（デバッグ時は戻り値を ids, tokens に保存）
	ids, _ = tk.Encode(text, true) // 第二引数 addSpecialTokens を true にしないと python の結果と異なってしまう

	// トークン数を表示
	// fmt.Printf("%s", text)
	// fmt.Printf(">> トークン数: %d\n", len(ids))

	return ids, nil
}

// ONNX推論用のヘルパー関数
func vectorize(tokenIds []uint32) (sentenceVector []float32, err error) {
	tensorLengthStr := os.Getenv("TENSOR_LENGTH")
	tensorLength, err := strconv.ParseInt(tensorLengthStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("TENSOR_LENGTH 環境変数が設定されていないか無効です: %v", err)
	}
	onnxruntimePath := os.Getenv("LIBRARY_PATH") + "/libonnxruntime.so"
	modelPath := os.Getenv("DOWNLOAD_DIR") + "/" + os.Getenv("SAVED_MODEL_PATH")

	// ONNX Runtime 環境の初期化
	onnxruntime_go.SetSharedLibraryPath(onnxruntimePath)
	err = onnxruntime_go.InitializeEnvironment()
	if err != nil {
		return nil, fmt.Errorf("ONNX Runtime の初期化に失敗しました: %v", err)
	}
	defer onnxruntime_go.DestroyEnvironment()

	// トークンIDを int64 に変換（BERT系モデルで一般的）
	inputData := make([]int64, len(tokenIds))
	for i, id := range tokenIds {
		inputData[i] = int64(id)
	}

	// token_type_ids を作成（単一文なので全て0で初期化）、attention_mask （全トークン有効のため、1で初期化）
	tokenTypeData := make([]int64, len(tokenIds))
	attentionMaskData := make([]int64, len(tokenIds))
	for i := range attentionMaskData {
		attentionMaskData[i] = 1
	}

	// 入力テンソルの形状を実際のトークン数に合わせて調整
	inputShape := onnxruntime_go.NewShape(1, int64(len(tokenIds)))

	// input_ids, token_type_ids, attention_mask テンソルを作成
	inputTensor, _ := onnxruntime_go.NewTensor(inputShape, inputData)
	defer inputTensor.Destroy()
	tokenTypeTensor, _ := onnxruntime_go.NewTensor(inputShape, tokenTypeData)
	defer tokenTypeTensor.Destroy()
	attentionMaskTensor, _ := onnxruntime_go.NewTensor(inputShape, attentionMaskData)
	defer attentionMaskTensor.Destroy()

	// 出力テンソルの形状: [batch_size, sequence_length, hidden_size]
	outputShape := onnxruntime_go.NewShape(1, int64(len(tokenIds)), tensorLength)
	outputTensor, err := onnxruntime_go.NewEmptyTensor[float32](outputShape)
	if err != nil {
		return nil, fmt.Errorf("出力テンソルの作成に失敗しました: %v", err)
	}
	defer outputTensor.Destroy()

	// セッション作成（実際のモデルファイルパスと入出力名を使用）
	inputNames := []string{"input_ids", "attention_mask", "token_type_ids"} // 3つの入力を指定
	outputNames := []string{"last_hidden_state"}                            // BERT系モデルの一般的な出力名
	session, err := onnxruntime_go.NewAdvancedSession(
		modelPath, inputNames, outputNames,
		[]onnxruntime_go.Value{inputTensor, attentionMaskTensor, tokenTypeTensor},
		[]onnxruntime_go.Value{outputTensor}, nil)
	if err != nil {
		return nil, fmt.Errorf("セッションの作成に失敗しました: %v", err)
	}
	defer session.Destroy()

	// モデル推論実行
	err = session.Run()
	if err != nil {
		return nil, fmt.Errorf("推論の実行に失敗しました: %v", err)
	}

	// 出力データを取得、ライブラリの使用か何らかの理由で1次元の float32 スライス（384*トークン数）として返される
	outputData := outputTensor.GetData()

	// 各トークンの埋め込みを平均化して文全体の埋め込みを計算、マスクは考慮しない（すべてのトークンが有効と仮定）
	hiddenSize := int(tensorLength) // モデルの隠れ層の次元数
	sentenceVector = make([]float32, hiddenSize)
	for i := 0; i < hiddenSize; i++ {
		for j := 0; j < len(tokenIds); j++ {
			sentenceVector[i] += outputData[j*hiddenSize+i]
		}
		sentenceVector[i] /= float32(len(tokenIds)) // 平均化
	}

	return sentenceVector, nil
}
