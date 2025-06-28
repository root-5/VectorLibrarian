# VectorLibrarian
文章をベクトル化して保存し、ベクトル検索とソース表示ができる API のテスト

## アイディア
- 主な対象は市役所等行政のホームページ
- HP をクローリングして、ページ内の main 部の文章を効率的にベクトル化して保存
- 高度な検索と LLM による回答を実現する
  - ページをベクトル化して、LLMでそこから解答を生成する
  - ページを簡易的にベクトル化、セマンティック検索で数ページに絞り込んだ上で、LLMに改めて全文を読ませ、解答を生成させる

市役所や市民に使ってもらうのではなく、使いたくなるような

# 作業メモ
## 実行コマンド
1. `go mod init github.com/root-5/VectorLibrarian`: モジュールの初期化
2. `go get github.com/gocolly/colly/v2`: クローリングライブラリのインストール
3. `go run main.go`: アプリケーションの実行
4. `go get -u github.com/JohannesKaufmann/html-to-markdown/v2`: HTMLをマークダウンに変換するライブラリのインストール
5. `go mod tidy`: 依存関係の整理
6. `npm install -g @google/gemini-cli`: Gemini CLI のインストール（Node はインストール済み）

## ライブラリドキュメント
https://pkg.go.dev/github.com/gocolly/colly#section-documentation
https://pkg.go.dev/github.com/JohannesKaufmann/html-to-markdown/v2#section-documentation

## colly
### 使い方
- `OnHTML`: 指定した要素が見つかった時に処理を実行したい
- `OnError`: リクエストでエラーが発生した時に処理を実行したい
- `OnRequest`: 全てのリクエストで処理を実行したい
- `OnResponse`: 全てのレスポンスで処理を実行したい

OnHTML 内の記述例
```go
textContent := e.DOM.Text() // テキストコンテンツを取得
html, err := e.DOM.Html() // HTMLを取得
```

### 挙動の理解
- デフォルトで一読ロールしたページは飛ばしてくれる
- `c.OnHTML("a[href]" ... c.Visit(e.Attr("href"))` のように、リンクをたどった場合、そのページ内で見つかった最初のリンクを訪問するためすべてのリンクを最初に取得するわけではない？
- 現在は全件取得している

## Gemini CLI
1. `npm install -g @google/gemini-cli`

### 設定やコマンドライン引数
https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/configuration.md
https://zenn.dev/schroneko/articles/gemini-cli-tutorial
