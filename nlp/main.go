package main

import (
	"fmt"
	"os"

	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/simple"
)

func main() {
	// バックエンドレシーバーを作成
	backend := simple.NewSimpleGraph()
	// モデルを作成し、実行バックエンドを設定
	model := onnx.NewModel(backend)

	// ONNXモデルを読み込み
	b, err := os.ReadFile("onnx_model/model.onnx")
	if err != nil {
		fmt.Println("ONNXファイルの読み込みエラー:", err)
		return
	}

	// モデルにデコード
	err = model.UnmarshalBinary(b)
	if err != nil {
		fmt.Println("ONNXモデルのデコードエラー:", err)
		return
	}

	fmt.Println("ONNXモデルが正常に読み込まれました!")
}
