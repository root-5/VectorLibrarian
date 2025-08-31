package vectorize

import (
	"fmt"
)

func ConvertToVector(text string, isQuery bool) (chunks []string, vectors [][]float32, err error) {
	if !isQuery {
		// マークダウンのリンクを置換
		text = replaceLinks(text)
	}

	// テキストを正規化
	normalizedText := normalizeText(text)

	// テキストを分割（チャンキング）
	chunks = chunkText(normalizedText, 512-3, 103) // -3 はプレフィックス分、-103 はオーバーラップ

	// チャンク数分のスライスを確保
	vectors = make([][]float32, len(chunks))

	for i, chunk := range chunks {
		// プレフィックスの付与
		if isQuery {
			chunk = "query: " + chunk
		} else {
			chunk = "passage: " + chunk
		}

		// トークン化
		ids, err := tokenize(chunk)
		if err != nil {
			fmt.Printf("トークナイズエラー: %v\n", err)
			return nil, nil, err
		}

		// ベクトル化
		vectors[i], err = vectorize(ids)
		if err != nil {
			fmt.Printf("ONNX推論実行エラー: %v\n", err)
			return nil, nil, err
		}
	}

	return chunks, vectors, nil
}
