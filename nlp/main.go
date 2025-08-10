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

	fmt.Println("Vocab size:", tk.VocabSize())
	fmt.Println(tk.Encode(text, false))
	fmt.Println(tk.Encode(text, true))
	fmt.Println(tk.Decode([]uint32{2829, 4419, 14523, 2058, 1996, 13971, 3899}, true))

	// ONNX推論部分（後で実装）
	fmt.Println("Tokenization completed successfully!")
}

// ONNX推論用のヘルパー関数
func runONNXInference(tokenIds []int32) ([]float32, error) {
	ort.SetSharedLibraryPath("/usr/local/lib/libonnxruntime.so")

	err := ort.InitializeEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ONNX Runtime: %v", err)
	}
	defer ort.DestroyEnvironment()

	// TODO: 実際のモデル推論を実装
	// 現在はダミーベクトルを返す
	embedding := make([]float32, 384)
	for i := range embedding {
		embedding[i] = float32(i) * 0.001
	}

	return embedding, nil
}
