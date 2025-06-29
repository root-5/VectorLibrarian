# ツールとライブラリ
ツール、ライブラリの使い方や設定方法、ドキュメントへのリンクをまとめる

## ライブラリドキュメント
- [スクレイピング "colly"](https://pkg.go.dev/github.com/gocolly/colly#section-documentation)
- [HTMLをMarkdownに変換 "html-to-markdown"](https://pkg.go.dev/github.com/JohannesKaufmann/html-to-markdown/v2#section-documentation)

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

## TablePlus
Postgres に INSERT したデータを確認する際に使用した GUI ツール。
.env, compose.yml にある DB 接続情報を使って接続する。「SSL mode」は「DISABLE」で設定。
