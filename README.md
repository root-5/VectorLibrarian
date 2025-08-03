# VectorLibrarian
文章をベクトル化して保存し、ベクトル検索とソース表示ができる API のテスト。
サブ目的として、AI技術を活用したシステム構築のテンプレート作成を目指す。

## タスク
- [ ] ベクトル化と検索のテスト
  - [ ] ベクトル化のためのライブラリ？Postgres拡張？を選定
  - [ ] 精度・計算資源の検証
- [ ] nlp に APIを実装
  - [ ] ベクトル化の API
- [ ] 本番環境構築
  - [ ] GCP 初期設定？
  - [ ] Terraform での構築？

https://chatgpt.com/share/6870edc4-82e0-8003-a163-ac64da6d19e5

## アイディア
### サービス面
- 特定ドメインをクローリングしてデータを蓄積
- 蓄積したデータからのベクトル検索が可能
- 蓄積されたデータをもとに AI による回答が可能
- 当面の主な対象は市役所等行政のホームページ

### 技術面
- コンテナ構造
  - app: フロントエンドへ提供する API や DB とのやり取りを行う Go で実装
  - nlp: テキストデータの正規化、ベクトル化、検索を行う Python で実装
  - db: PostgreSQL、ベクトルも含めて保存
- ORM の是非については諸説あるが、一旦 Bun を使ってみる
  - せっかくの個人プロジェクトなので一般的な Gorm ではなく、新しい ORM を採用してみた
  - 通常の Prisma 等は使ったことがあるので "SQL First" と謳っている Bun を使ってみて使用感を確かめたい
- いろいろな自然言語周辺技術等調べていて思ったが、とてもじゃないがライブラリ使用を避けるのは難しい
- 改善案
  - 「手続き」「申請」「暮らし」「補助金」等を含まれるページを優先してクローリング
  - 「〇〇県〇〇市の～についての～」のようなテンプレートに当て込んでのベクトル化

### 前処理
数字、英字の半角、全角の統一
漢数字をアラビア数字に変換
日付の表記を統一
全角スペースを半角スペースに変換
半角カタカナを全角カタカナに変換
一か月の表記「1か月」「1カ月」「1ヶ月」「1ヵ月」「1箇月」「1ケ月」の統一
記号を全角から半角に変換（例：！→!、？→?、＠→@、￥→¥、”→"）
句読点の統一（例：，→、や．→。）

# 作業メモ
## 実行コマンドメモ
### app
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

### nlp
1. `uv init`: uv の初期化
2. `uv add --dev ruff`: ruff の追加（開発用）
3. `uv add  uv add sentence-transformers`: sentence-transformers の追加
4. `uv run main.py`: main.py の実行
5. `uv run ruff check .`: ruff でコードチェック
6. `uv run ruff format .`: ruff でコードフォーマット

## Docker 関係コマンド
- `docker compose up -d`: 開発環境コンテナの起動
- `docker compose down`: 開発環境コンテナの停止
- `docker compose down --rmi all`: 開発環境コンテナの停止とイメージの削除

### app
- `docker compose exec app sh`: 開発環境コンテナ内でシェルを開く
- `docker compose exec app go run main.go`: 開発環境コンテナ内でアプリケーションを実行

### db
- `docker compose exec db sh -c 'psql -U $POSTGRES_USER -d $POSTGRES_DB'`: 開発環境コンテナ内で PostgreSQL に接続
  - `SELECT * FROM pages;`: データベースの内容を確認
- `docker compose exec db sh -c 'pg_dump -U $POSTGRES_USER $POSTGRES_DB > /backup/backup.sql'`: データベースのバックアップを取得
- `docker compose exec db sh -c 'psql -U $POSTGRES_USER $POSTGRES_DB < /backup/backup.sql'`: データベースのバックアップを復元

### nlp
- `docker compose exec nlp sh`: NLP コンテナ内でシェルを開く
- `docker compose exec nlp python main.py`: NLP コンテナ内で Python スクリプトを実行
