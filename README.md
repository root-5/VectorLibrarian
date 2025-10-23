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
- [実行したコマンドメモ](./documents/実行したコマンドメモ.md)
- [全体構造](./documents/全体構造.md)

## 中期目標とタスク

### 中期目標

1. ベクトルDB × RAG構築（一旦プロダクトとして形にする）
2. MLOps（継続的改善の仕組みを作る）
3. ファインチューニング（モデルの精度向上を図る）

### タスク

- [x] crawler を定期実行としてパッケージ化
- [x] 外部からのクエリを受け付ける api パッケージの実装
- [x] AWS 本番環境構築
- [ ] 精度向上
  - [x] トークン上限の確認
  - [ ] markdown 入力時の精度検証
  - [ ] multilingual の別モデルを試す
  - [x] nlp の API レスポンスにモデル名等も含める
- [ ] CI/CD 強化
  - [x] 頻繁にトライ&エラーが発生する箇所に単体テストとテストモードを追加
  - [ ] ベクトルデータは単独のテーブルとして切り出して、モデル切り替え時にテーブルごと入れ替えられるようにする
  - [ ] main ブランチを使用したデプロイ自動化
- [ ] 検索機能改善
  - [ ] 全文検索など他の検索機能との統合
  - [ ] プロンプトと埋め込み元データの改善
  - [ ] 最初に対象ドメインを選択できるようにする
- [ ] RAG 化
  - [ ] gpt-oss を利用しての LLM 回答機能の実装
- [ ] セキュリティ強化
- [ ] 本番環境強化
  - [ ] gRPC 化
  - [ ] 監視コンテナの導入
  - [ ] テストコード追加
- [ ] インフラ構成変更
  - [ ] GCP, VertexAI 化
  - [ ] Terraform での構築？

## コマンド

- `docker compose up -d`: 開発環境コンテナの起動
- `docker compose down`: 開発環境コンテナの停止
- `docker compose down --rmi all`: 開発環境コンテナの停止とイメージの削除
- `docker-compose -f="compose.prod.yml" up -d`: 本番環境コンテナの起動
- `docker-compose -f="compose.prod.yml" down`: 本番環境コンテナの停止

### バックアップ

コード、DBデータ、クローラーキャッシュ等全データバックアップ

```sh
sudo cp -rp VectorLibrarian VectorLibrarian.backup.`date "+%Y-%m-%d_%H-%M"`
```

### app コンテナ用

- `docker compose exec app sh`: 開発環境コンテナ内でシェルを開く
- `docker compose exec app go run main.go`: 開発環境コンテナ内でアプリケーションを実行
- `docker compose exec app curl -X POST "http://nlp:8000/convert" -H "Content-Type: application/json" -d '{ "text": "これは日本語の文章です。", "is_query": true}'`: ベクトル化 API をテスト
- `docker compose exec app go test ./controller/crawler`: 単体テストを実行
- `docker compose exec app go run main.go -mode=test`: テストモードでアプリケーションを実行（統合的なテスト用）

### db コンテナ用

- `docker compose exec db sh -c 'psql -U $POSTGRES_USER -d $POSTGRES_DB'`: 開発環境コンテナ内で PostgreSQL に接続
  - `SELECT * FROM pages;`: データベースの内容を確認
  - `\dx`: 拡張機能の確認
- `docker compose exec db sh -c 'pg_dump -U $POSTGRES_USER $POSTGRES_DB > /backup/backup_$(date +%Y-%m-%d_%H-%M).sql'`: データベースのバックアップを取得
- `docker compose exec db sh -c 'psql -U $POSTGRES_USER $POSTGRES_DB < /backup/backup.sql'`: データベースのバックアップを復元

ドメインテーブルの初期設定用ドメインのINSERT文例：

```sql
INSERT INTO "public"."domains" ("id", "domain", "created_at", "updated_at", "deleted_at") VALUES
(1, 'www.city.hamura.tokyo.jp', '2025-10-12 22:44:44.038714+09', '2025-10-12 22:44:44.038714+09', '0001-01-01 09:18:59+09:18:59');
```

### nlp コンテナ用

- `docker compose exec nlp sh`: NLP コンテナ内でシェルを開く
  - `curl -X POST "http://localhost:8000/convert" -H "Content-Type: application/json" -d '{ "text": "これは日本語の文章です。", "is_query": true}'`: ベクトル化 API をテスト
- `docker compose exec nlp go test ./vectorize`: 単体テストを実行
