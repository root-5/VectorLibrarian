# VectorLibrarian
文章をベクトル化して保存し、ベクトル検索とソース表示ができる API のテスト

## タスク
- [ ] ベクトル化と検索のテスト
  - [ ] ベクトル化のためのライブラリ？Postgres拡張？を選定
  - [ ] 精度・計算資源の検証
- [ ] フォーマルな形に再構成
  - [ ] DB の型やカラムを適切なものに修正
  - [ ] データの型情報を domain/model に切り出し
- [ ] 本番環境構築
  - [ ] GCP 初期設定？
  - [ ] Terraform での構築？

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
