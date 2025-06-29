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
## コマンド
### 実行コマンドメモ
1. `go mod init github.com/root-5/VectorLibrarian`: モジュールの初期化
2. `go get github.com/gocolly/colly/v2`: クローリングライブラリのインストール
3. `go run main.go`: アプリケーションの実行
4. `go get -u github.com/JohannesKaufmann/html-to-markdown/v2`: HTMLをマークダウンに変換するライブラリのインストール
5. `go mod tidy`: 依存関係の整理
6. `npm install -g @google/gemini-cli`: Gemini CLI のインストール（Node はインストール済み）
7. `go get github.com/lib/pq`: PostgreSQL ドライバのインストール（Gemini CLI が実行）
8. `go mod tidy`

### Docker 関係
- `docker compose up -d`: 開発環境コンテナの起動
- `docker compose down`: 開発環境コンテナの停止
- `docker compose down --rmi all --volumes`: 開発環境コンテナの停止とイメージ、ボリュームの削除
- `docker compose exec app sh`: 開発環境コンテナ内でシェルを開く
- `docker compose exec app go run main.go`: 開発環境コンテナ内でアプリケーションを実行
- `docker compose exec db sh`: 開発環境コンテナ内でシェルを開く
  - `psql -U user -d db`: PostgreSQL に接続
  - `SELECT * FROM pages;`: データベースの内容を確認

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

## Gemini 修正内容メモ
あまりにタスクを一気にこなすので、気になった点をメモしておかないと忘れてしまう。
- テーブル作成の SQL, データ挿入の SQL が適切か確認
- `docker compose` で作成される環境が本番用なので、開発用も用意させる
- アプリケーションレイヤーのファイルは `app` ディレクトリに配置する
- air の設定を確認
