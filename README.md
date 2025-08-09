# VectorLibrarian

文章をベクトル化して保存し、ベクトル検索とソース表示ができる API のテスト。
サブ目的として、AI技術を活用したシステム構築のテンプレート作成を目指す。

## ドキュメント

- [Notion_AI関係](./documents/Notion_AI関係.md)
- [Notion](./documents/Notion.md)
- [アイディア](./documents/アイディア.md)
- [テーブル設計](./documents/テーブル設計.md)
- [使用ツールとライブラリ](./documents/使用ツールとライブラリ.md)
- [全体構造](./documents/全体構造.md)

## タスク

- [ ] crawler を定期実行としてパッケージ化
- [ ] 外部からのクエリを受け付ける api パッケージの実装
- [ ] 本番環境構築
  - [ ] GCP 初期設定？
  - [ ] Terraform での構築？
- [ ] 検索機能改善
  - [ ] 全文検索など他の検索機能との統合
  - [ ] プロンプトと埋め込み元データの改善

## 実行コマンド

### app の実行コマンド

1. `go mod init github.com/root-5/VectorLibrarian`: モジュールの初期化
2. `go get github.com/gocolly/colly/v2`: クローリングライブラリのインストール
3. `go run main.go`: アプリケーションの実行
4. `go get -u github.com/JohannesKaufmann/html-to-markdown/v2`: HTMLをマークダウンに変換するライブラリのインストール
5. `go mod tidy`: 依存関係の整理
6. `npm install -g @google/gemini-cli`: Gemini CLI のインストール（Node はインストール済み）
7. `go get github.com/lib/pq`: PostgreSQL ドライバのインストール（Gemini CLI が実行）
8. `go get github.com/uptrace/bun`: Bun ORM のインストール
9. `go get github.com/uptrace/bun/dialect/pgdialect`: PostgreSQL 用の Bun ダイアレクトのインストール
10. `go get github.com/uptrace/bun/driver/pgdriver`: PostgreSQL ドライバのインストール
11. `go mod tidy`: 依存関係の整理、便宜上最後のコマンドとして記載しているがライブラリのインストール後に適宜実行した

### nlp の実行コマンド

1. `uv init`: uv の初期化
2. `uv add transformers`: transformers の追加
3. `uv add optimum`: optimum の追加
4. `uv add onnx`: onnx の追加
5. `uv add onnxruntime`: onnxruntime の追加
6. `uv run main.py`: main.py の実行

7. `uv add fastapi`: fastapi の追加
8. `uv add "uvicorn[standard]"`: "uvicorn[standard]" の追加
9. `uv add neologdn`: neologdn の追加

## Docker 関係コマンド

- `docker compose up -d`: 開発環境コンテナの起動
- `docker compose down`: 開発環境コンテナの停止
- `docker compose down --rmi all`: 開発環境コンテナの停止とイメージの削除

### app の Docker コマンド

- `docker compose exec app sh`: 開発環境コンテナ内でシェルを開く
- `docker compose exec app go run main.go`: 開発環境コンテナ内でアプリケーションを実行
- `docker compose exec app curl -X POST "http://nlp:8000/convert" -H "Content-Type: application/json" -d '{ "text": "機械学習とは何ですか？", "is_query": true}'`: ベクトル化 API をテスト

### db の Docker コマンド

- `docker compose exec db sh -c 'psql -U $POSTGRES_USER -d $POSTGRES_DB'`: 開発環境コンテナ内で PostgreSQL に接続
  - `SELECT * FROM pages;`: データベースの内容を確認
- `docker compose exec db sh -c 'pg_dump -U $POSTGRES_USER $POSTGRES_DB > /backup/backup.sql'`: データベースのバックアップを取得
- `docker compose exec db sh -c 'psql -U $POSTGRES_USER $POSTGRES_DB < /backup/backup.sql'`: データベースのバックアップを復元

### nlp の Docker コマンド

- `docker compose exec nlp sh`: NLP コンテナ内でシェルを開く
- `docker compose exec nlp uv run main.py`: NLP コンテナ内で uv を使って main.py を実行
- `docker compose exec nlp uv run uvicorn main:app --reload --host 0.0.0.0`: ホットリロードを有効にして FastAPI アプリケーションを実行
  - `curl -X POST "http://localhost:8000/convert" -H "Content-Type: application/json" -d '{ "text": "機械学習とは何ですか？", "is_query": true}'`: ベクトル化 API をテスト

hugging face version 指定が望ましい
https://zenn.dev/yagiyuki/articles/load_pretrained

