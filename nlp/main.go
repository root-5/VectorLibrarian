package main

import (
	"fmt"

	"github.com/daulet/tokenizers"
	ort "github.com/yalue/onnxruntime_go"
)

func main() {

	tk, err := tokenizers.FromFile("onnx_model/tokenizer.json")
	if err != nil {
		fmt.Printf("Error loading tokenizer: %v\n", err)
		return
	}
	// release native resources
	defer tk.Close()

	text := "これは日本語の文章です。"

	// トークン化
	ids, tokens := tk.Encode(text, true) // 第二引数 addSpecialTokens を true にしないと python の結果と異なってしまう
	fmt.Println("")
	fmt.Printf("Token IDs: %v\n", ids)
	fmt.Printf("Tokens: %v\n", tokens)
	fmt.Println("")

	// ベクトル化
	embedding, err := runONNXInference(ids)
	if err != nil {
		fmt.Printf("Error running ONNX inference: %v\n", err)
		return
	}
	// 最初の5次元を表示
	fmt.Println("")
	fmt.Printf("Embedding shape: %d\n", len(embedding))
	fmt.Println("")
	fmt.Printf("First 5 dimensions: %v\n", embedding[:5]) // [-0.050542377 0.16145682 0.010935326 -0.02508498 0.15268864]
}

// ONNX推論用のヘルパー関数
func runONNXInference(tokenIds []uint32) ([]float32, error) {
	ort.SetSharedLibraryPath("/usr/local/lib/libonnxruntime.so")

	err := ort.InitializeEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ONNX Runtime: %v", err)
	}
	defer ort.DestroyEnvironment()

	// トークンIDを int64 に変換（BERT系モデルで一般的）
	inputData := make([]int64, len(tokenIds))
	for i, id := range tokenIds {
		inputData[i] = int64(id)
	}

	// token_type_ids を作成（全て0で初期化、単一文の場合）
	tokenTypeData := make([]int64, len(tokenIds))
	// 全て0のまま（単一文なので）

	// attention_mask を作成（全て1で初期化、すべてのトークンに注意を向ける）
	attentionMaskData := make([]int64, len(tokenIds))
	for i := range attentionMaskData {
		attentionMaskData[i] = 1 // すべてのトークンが有効
	}

	// 入力テンソルの形状を実際のトークン数に合わせて調整
	// [batch_size, sequence_length] = [1, len(tokenIds)]
	inputShape := ort.NewShape(1, int64(len(tokenIds)))

	// input_ids, token_type_ids, attention_mask テンソルを作成
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	if err != nil {
		return nil, fmt.Errorf("failed to create input tensor: %v", err)
	}
	defer inputTensor.Destroy()

	tokenTypeTensor, err := ort.NewTensor(inputShape, tokenTypeData)
	if err != nil {
		return nil, fmt.Errorf("failed to create token_type_ids tensor: %v", err)
	}
	defer tokenTypeTensor.Destroy()

	attentionMaskTensor, err := ort.NewTensor(inputShape, attentionMaskData)
	if err != nil {
		return nil, fmt.Errorf("failed to create attention_mask tensor: %v", err)
	}
	defer attentionMaskTensor.Destroy()

	// 出力テンソルの形状: [batch_size, sequence_length, hidden_size]
	// sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2 の hidden_size は 384
	outputShape := ort.NewShape(1, int64(len(tokenIds)), 384)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		return nil, fmt.Errorf("failed to create output tensor: %v", err)
	}
	defer outputTensor.Destroy()

	// セッション作成（実際のモデルファイルパスと入出力名を使用）
	modelPath := "onnx_model/onnx/model.onnx"
	inputNames := []string{"input_ids", "attention_mask", "token_type_ids"} // 3つの入力を指定
	outputNames := []string{"last_hidden_state"}                            // BERT系モデルの一般的な出力名

	session, err := ort.NewAdvancedSession(
		modelPath,
		inputNames,
		outputNames,
		[]ort.Value{inputTensor, attentionMaskTensor, tokenTypeTensor},
		[]ort.Value{outputTensor},
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Destroy()

	// モデル推論実行
	err = session.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run inference: %v", err)
	}

	// 出力データを取得
	outputData := outputTensor.GetData() // なぜか単一の float32 スライス（384*len(tokenIds)の長さの1次元配列）として返される

	hiddenSize := 384
	sentenceEmbedding := make([]float32, hiddenSize)

	for i := 0; i < hiddenSize; i++ {
		for j := 0; j < len(tokenIds); j++ {
			sentenceEmbedding[i] += outputData[j*hiddenSize+i]
		}
		sentenceEmbedding[i] /= float32(len(tokenIds)) // 平均化
	}

	return sentenceEmbedding, nil
}
