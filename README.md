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
6. `uv add tokenizers`: tokenizers の追加
7. `uv add "huggingface_hub[cli]"`: huggingface_hub[cli] の追加
8. `uv run main.py`: main.py の実行
9. `go mod init github.com/root-5/VectorLibrarian/nlp`: モジュールの初期化
10. `go get github.com/yalue/onnxruntime_go`: ONNX Runtime Go ライブラリのインストール
11. `go get github.com/daulet/tokenizers`: トークナイザライブラリのインストール

12. `uv add fastapi`: fastapi の追加
13. `uv add "uvicorn[standard]"`: "uvicorn[standard]" の追加
14. `uv add neologdn`: neologdn の追加

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
バージョンが変わってしまうと、以前生成したベクトルデータと新しいモデルで生成したベクトルデータの生成プロセスが変化し、検索結果が変わってしまう可能性があるため。
また、再ダウンロードが生じる可能性もある。
<https://zenn.dev/yagiyuki/articles/load_pretrained>

uv run hf download sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2 onnx/model.onnx tokenizer.json --revision 86741b4e3f5cb7765a600d3a3d55a0f6a6cb443d --local-dir onnx_model
バージョンの指定ができ、おそらく専用のバリデーションとエラーハンドリングが行われるため wget や curl よりも適している
<https://huggingface.co/docs/huggingface_hub/main/en/guides/download>

<https://huggingface.co/sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2>

model_qint8_arm64.onnx
model_qint8_avx512.onnx
model_qint8_avx512_vnni.onnx
model_quint8_avx2.onnx

# CPU の命令セットを確認する

lscpu | grep Flags  # Linux

Flags:                                fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ss ht syscall nx pdpe1gb rdtscp lm constant_tsc rep_good nopl xtopology tsc_reliable nonstop_tsc cpuid tsc_known_freq pni pclmulqdq vmx ssse3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm abm 3dnowprefetch ssbd ibrs ibpb stibp ibrs_enhanced tpr_shadow ept vpid ept_ad fsgsbase tsc_adjust bmi1 avx2 smep bmi2 erms invpcid rdseed adx smap clflushopt clwb sha_ni xsaveopt xsavec xgetbv1 xsaves avx_vnni vnmi umip waitpkg gfni vaes vpclmulqdq rdpid movdiri movdir64b fsrm md_clear serialize flush_l1d arch_capabilities

avx2 があるので model_quint8_avx2.onnx を使用する

| 項目        | FP32 | INT8量子化      |
| --------- | ---- | ------------ |
| モデルサイズ    | 100% | 約25%         |
| 推論速度(CPU) | 基準   | 1.5〜4倍       |
| 推論速度(GPU) | 基準   | 1.2〜2倍（対応次第） |
| 精度        | 最高   | 0〜数%低下       |
| 対応性       | 高い   | ハード依存あり      |

# libonnxruntime.so

onnxruntime.so は以下をビルドしたものの一部？onnxモデルを動かすために必要。gitで以下の公式を落としてきてビルドしたり、pipを使ってインストールする必要ことも可能だが、コンテナでの利用だと重たくなるのでリリースのほうから対象バイナリを含んだtgzをダウンロードして解凍するのが良さそう。
リリースバージョンによってAssetsに含まれるOSの種類やGPUの有無などが異なるのでほしい物がない場合は過去バージョンを探す必要がある。
今回は linux x64 gpu なし（onnxruntime-linux-x64-1.22.0.tgz）を使用する。
<https://github.com/microsoft/onnxruntime>
<https://github.com/microsoft/onnxruntime/releases>
