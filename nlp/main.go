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
	fmt.Printf("Embedding: %v\n", embedding)
}

// ONNX推論用のヘルパー関数
func runONNXInference(tokenIds []uint32) ([]float32, error) {
	ort.SetSharedLibraryPath("/usr/local/lib/libonnxruntime.so")

	err := ort.InitializeEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ONNX Runtime: %v", err)
	}
	defer ort.DestroyEnvironment()

    inputData := []float32{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9}
    inputShape := ort.NewShape(2, 5)
    inputTensor, err := ort.NewTensor(inputShape, inputData)
    defer inputTensor.Destroy()
	// This hypothetical network maps a 2x5 input -> 2x3x4 output.
	outputShape := ort.NewShape(2, 3, 4)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	defer outputTensor.Destroy()

	session, err := ort.NewAdvancedSession("onnx_model/model.onnx",
		[]string{"Input 1 Name"}, []string{"Output 1 Name"},
		[]ort.Value{inputTensor}, []ort.Value{outputTensor}, nil)
	defer session.Destroy()

	// Calling Run() will run the network, reading the current contents of the
	// input tensors and modifying the contents of the output tensors.
	err = session.Run()

	// Get a slice view of the output tensor's data.
	outputData := outputTensor.GetData()

	return outputData, nil
}
