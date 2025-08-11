package vectorize

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/daulet/tokenizers"
	"github.com/yalue/onnxruntime_go"
	"golang.org/x/text/unicode/norm"
)

func TextToVector(text string, isQuery bool) (vector []float32, err error) {
	// テキストを正規化
	text = normalizeText(text)

	// プレフィックスの付与
	if isQuery {
		text = "query: " + text
	} else {
		text = "passage: " + text
	}

	// トークン化
	ids, err := tokenize(text)

	// ベクトル化
	vector, err = vectorize(ids)
	if err != nil {
		fmt.Printf("ONNX推論実行エラー: %v\n", err)
		return
	}
	return vector, nil
}

// テキストの正規化関数
func normalizeText(text string) (normalizedText string) {
	// Unicode正規化（NFKC）、全角カタカナに統一、半角英数に統一など
	normalizedText = norm.NFKC.String(text)

	// 小文字化
	normalizedText = strings.ToLower(normalizedText)

	// スペースや改行の扱いなども必要に応じて追加

	return normalizedText
}

// トークナイズ関数
func tokenize(text string) (ids []uint32, err error) {
	// 環境変数から設定を読み込む
	tokenizerPath := os.Getenv("DOWNLOAD_DIR") + "/" + os.Getenv("SAVED_TOKENIZER_PATH")

	tk, err := tokenizers.FromFile(tokenizerPath)
	if err != nil {
		fmt.Printf("tokenizer.json ロードエラー: %v\n", err)
		return nil, err
	}
	defer tk.Close()

	// トークン化（デバッグ時は戻り値を ids, tokens に保存）
	ids, _ = tk.Encode(text, true) // 第二引数 addSpecialTokens を true にしないと python の結果と異なってしまう

	return ids, nil
}

// ONNX推論用のヘルパー関数
func vectorize(tokenIds []uint32) (sentenceVector []float32, err error) {
	tensorLength, _ := strconv.ParseInt(os.Getenv("TENSOR_LENGTH"), 10, 64)
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
