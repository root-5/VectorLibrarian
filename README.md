# VectorLibrarian

WEBページの文章をベクトル化して保存し、ベクトル検索ができる API。
サブ目的として、AI技術を活用したシステム構築のイメージをつかむ。

## 特徴

- `sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2` モデルを使用
- 文章での検索が可能で、文章のほうが精度が良くなる
- アラビア語（لقد تلقيت إشعارًا ضريبيًا）など本来使用されていない語句・言語での検索が可能
- 逆に完全一致していなくても上位に表示される可能性がある
- 現状は 20 件だけ出力
- 意味抽出にはタイトルと見出しを使用しているので、構造化が乏しいページや文章が長いページに弱い
- データはある程度クローリングしていっているが、全てではないし pdf などは未対応
- Google 検索がいかに素晴らしいかを改めて感じられる

## ドキュメント

- [Notion_AI関係](./documents/Notion_AI関係.md)
- [Notion](./documents/Notion.md)
- [アイディア](./documents/アイディア.md)
- [テーブル設計](./documents/テーブル設計.md)
- [使用ツールとライブラリ](./documents/使用ツールとライブラリ.md)
- [全体構造](./documents/全体構造.md)

## タスク

- [x] crawler を定期実行としてパッケージ化
- [x] 外部からのクエリを受け付ける api パッケージの実装
- [x] AWS 本番環境構築
- [ ] 精度向上
  - [x] トークン上限の確認
  - [ ] markdown 入力時の精度検証
  - [ ] multilingual の別モデルを試す
- [ ] CI/CD 強化
  - [ ] main ブランチを使用したデプロイ自動化
- [ ] 検索機能改善
  - [ ] 全文検索など他の検索機能との統合
  - [ ] プロンプトと埋め込み元データの改善
- [ ] RAG 化
  - [ ] gpt-oss を利用しての LLM 回答機能の実装
- [ ] セキュリティ強化
- [ ] 本番環境強化
  - [ ] 監視コンテナの導入
  - [ ] テストコード追加
- [ ] インフラ構成変更
  - [ ] GCP 化
  - [ ] Terraform での構築？

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
6. `uv add tokenizers`: tokenizers の追加
7. `uv add "huggingface_hub[cli]"`: huggingface_hub[cli] の追加
8. `uv run main.py`: main.py の実行
9. `go mod init github.com/root-5/VectorLibrarian/nlp`: モジュールの初期化
10. `go get github.com/yalue/onnxruntime_go`: ONNX Runtime Go ライブラリのインストール
11. `go get github.com/daulet/tokenizers`: トークナイザライブラリのインストール

## Docker 関係コマンド

- `docker compose up -d`: 開発環境コンテナの起動
- `docker compose down`: 開発環境コンテナの停止
- `docker compose down --rmi all`: 開発環境コンテナの停止とイメージの削除
- `docker-compose -f="compose.prod.yml" up -d`: 本番環境コンテナの起動
- `docker-compose -f="compose.prod.yml" down`: 本番環境コンテナの停止

### app の Docker コマンド

- `docker compose exec app sh`: 開発環境コンテナ内でシェルを開く
- `docker compose exec app go run main.go`: 開発環境コンテナ内でアプリケーションを実行
- `docker compose exec app curl -X POST "http://nlp:8000/convert" -H "Content-Type: application/json" -d '{ "text": "これは日本語の文章です。", "is_query": true}'`: ベクトル化 API をテスト

### db の Docker コマンド

- `docker compose exec db sh -c 'psql -U $POSTGRES_USER -d $POSTGRES_DB'`: 開発環境コンテナ内で PostgreSQL に接続
  - `SELECT * FROM pages;`: データベースの内容を確認
  - `\dx`: 拡張機能の確認
- `docker compose exec db sh -c 'pg_dump -U $POSTGRES_USER $POSTGRES_DB > /backup/backup_$(date +%Y-%m-%d_%H-%M).sql'`: データベースのバックアップを取得
- `docker compose exec db sh -c 'psql -U $POSTGRES_USER $POSTGRES_DB < /backup/backup.sql'`: データベースのバックアップを復元

### nlp の Docker コマンド

- `docker compose exec nlp sh`: NLP コンテナ内でシェルを開く
  - `curl -X POST "http://localhost:8000/convert" -H "Content-Type: application/json" -d '{ "text": "これは日本語の文章です。", "is_query": true}'`: ベクトル化 API をテスト
