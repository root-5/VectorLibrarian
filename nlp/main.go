package main

import (
	"fmt"

	ort "github.com/yalue/onnxruntime_go"
)

func main() {
	ort.SetSharedLibraryPath("/usr/local/lib/libonnxruntime.so")

	err := ort.InitializeEnvironment()
	if err != nil {
		panic(err)
	}
	defer ort.DestroyEnvironment()

	// わずかなパフォーマンス向上と、既存のテンソルを再利用する際の利便性のため、
	// このライブラリでは、セッションを作成する前にすべての入力・出力テンソルを
	// 作成することを想定しています。この方法があなたのユースケースに最適でない場合は、
	// ドキュメントの DynamicAdvancedSession 型を参照してください。この型では、
	// セッション初期化時ではなく、Run() 呼び出し時に入力・出力テンソルを指定できます。
	inputData := []float32{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9}
	inputShape := ort.NewShape(2, 5)
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	defer inputTensor.Destroy()
	// この仮想的なネットワークは 2x5 の入力 -> 2x3x4 の出力にマッピングします。
	outputShape := ort.NewShape(1, 384)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	defer outputTensor.Destroy()

	session, err := ort.NewAdvancedSession("onnx_model/onnx/model.onnx",
		[]string{"Input 1 Name"}, []string{"Output 1 Name"},
		[]ort.Value{inputTensor}, []ort.Value{outputTensor}, nil)
	defer session.Destroy()

	err = session.Run()
	outputData := outputTensor.GetData()

	fmt.Printf("Output data: %v\n", outputData)

	// 異なる入力でネットワークを実行したい場合は、入力テンソルのデータ
	// (inputTensor.GetData() 経由で利用可能) を変更し、再度 Run() を呼び出
	// すればよいだけです。
}
