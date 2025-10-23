package vectorize

import (
	"reflect"
	"testing"
)

// 単体テスト（外部依存がない関数のテスト）を定義
// `docker compose exec nlp go test ./vectorize`

func TestReplaceLinks(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{"リンクを置換", "[トップ](https://example.com)", "トップ"},
		{"画像リンクを置換", "![説明](https://example.com/image.png)", "説明"},
		{"ネストされたリンクを置換", "[![alt](https://example.com/a.png)](https://example.com/b)", "alt"},
		{"リンクがない", "テキストのみ", "テキストのみ"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := replaceLinks(tc.input)
			if output != tc.expectedOutput {
				t.Errorf("期待される出力 '%s' ですが、実際は '%s' でした", tc.expectedOutput, output)
			}
		})
	}
}

func TestChunkText(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput []string
	}{
		// {"短いテキスト", "これはテストです。", []string{"これはテストです。"}},
		// {"長いテキスト1", "すもももももももものうち。すもももももももものうち。すもももももももものうち。すもももももももものうち。すもももももももものうち。すもももももももものうち。すもももももももものうち。すもももももももものうち。すもももももももものうち。すもももももももものうち。", []string{"すもももももももものうち。すもももももももものうち。すもももももももものうち。すもももももももものうち。", "すもももももももものうち。すもももももももものうち。すもももももももものうち。すもももももももものうち。", "すもももももももものうち。すもももももももものうち。すもももももももものうち。すもももももももものうち。"}},
		// {"長いテキスト2", `みなさんwikiという言葉を一度は聞いたことがあるかと思いますが、wikiはwikipediaの略語だと認識していませんか？確かに慣用表現としてwikipediaのことをwikiと呼ぶことはありますが、厳密にはwikiとwikipediaは意味が異なっています。この記事ではwikiとは何かといった基本から作り方、ページ作成時に使うMarkdownについて解説します。`, []string{"みなさんwikiという言葉を一度は聞いたことがあるかと思いますが、wikiはwikipediaの略語だと認識していませんか？確かに慣用表現としてwikipediaのことをwikiと呼ぶことはありますが、厳密にはwikiとwikipediaは意味が異なっています。", "この記事ではwikiとは何かといった基本から作り方、ページ作成時に使うMarkdownについて解説します。"}},
		{"長いテキスト3", `すもももももももものうち
すもももももももものうち
すもももももももものうち
すもももももももものうち
すもももももももものうち
すもももももももものうち
すもももももももものうち
すもももももももものうち
すもももももももものうち
すもももももももものうち`, []string{`すもももももももものうち
すもももももももものうち
すもももももももものうち
すもももももももものうち
`, `すもももももももものうち
すもももももももものうち
すもももももももものうち
すもももももももものうち
`, `すもももももももものうち
すもももももももものうち
すもももももももものうち
すもももももももものうち`}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output := chunkText(tc.input, 60, 15)
			if !reflect.DeepEqual(output, tc.expectedOutput) {
				t.Errorf("期待される出力 '%v' ですが、実際は '%v' でした", tc.expectedOutput, output)
			}
		})
	}
}
