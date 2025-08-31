package vectorize

import (
	"fmt"
)

func MarkdownToVector(markdown string) (vectors [][]float32, err error) {
	// マークダウンのリンクを置換
	replacedMarkdown := replaceLinks(markdown)

	// テキストを正規化
	text := normalizeText(replacedMarkdown)

	// テキストを分割（チャンキング）
	sentences := chunkText(text, 512, 128)

	// Initialize vectors slice with the length of sentences
	vectors = make([][]float32, len(sentences))

	for i, sentence := range sentences {
		// プレフィックスの付与
		sentence = "passage: " + sentence

		// トークン化
		ids, err := tokenize(sentence)
		if err != nil {
			fmt.Printf("トークナイズエラー: %v\n", err)
			return nil, err
		}

		// ベクトル化
		vectors[i], err = vectorize(ids)
		if err != nil {
			fmt.Printf("ONNX推論実行エラー: %v\n", err)
			return nil, err
		}
	}

	return vectors, nil
}

func QueryToVector(text string) (vector []float32, err error) {
	// テキストを正規化
	text = normalizeText(text)

	// プレフィックスの付与
	text = "query: " + text

	// トークン化
	ids, err := tokenize(text)
	if err != nil {
		fmt.Printf("トークナイズエラー: %v\n", err)
		return nil, err
	}

	// ベクトル化
	vector, err = vectorize(ids)
	if err != nil {
		fmt.Printf("ONNX推論実行エラー: %v\n", err)
		return nil, err
	}
	return vector, nil
}
