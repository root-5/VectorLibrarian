# AI 関連

## ベクトル検索の基本

ベクトルのコサイン類似度の計算はベクトル DB の機能を利用するのが高率的。わざわざアプリ側で計算しない。
保存したベクトルをさらに加工して、次元を落とすことでRDBにおけるインデックスと似た効果を持たせる技術も登場している。

### 基本的な流れ

1. 正規化
   - 例: neologdn を利用して、全角英数字を半角に変換、全角スペースを半角スペースに変換、改行コードを統一など
2. トークン化
   - 例: MeCab や Sudachi などの形態素解析器を利用して、文章を単語やフレーズに分割
3. ベクトル化（埋め込み）
   - 例: SentenceTransformer や FastText などの事前学習済みモデルを利用して、トークン化された単語やフレーズをベクトルに変換
4. ベクトルの保存
   - 例: PostgreSQL の pgvector 拡張を利用して、ベクトルをデータベースに保存
5. 検索クエリをベクトル化
   - 例: ユーザーの検索クエリを同様にベクトル化
6. 類似度計算
   - 例: ベクトル DB の機能を利用して、保存されたベクトルと検索クエリのベクトルとのコサイン類似度を計算
7. 結果の取得
   - 例: 類似度が高い順に結果を取得し、ユーザーに返す

## ベクトル DB の選定

ベクトルに特化したDBも存在はするものの、発展途上であり、小規模においてはRDBを利用するのが現実的である。
[参考](https://zenn.dev/rwcolinpeng/articles/45632994cf8bc1)

## pgvector

PostgreSQL の拡張で、ベクトルデータを扱えるようにするもの。
テーブル設定時に `vector(1024)` のように次元数を指定することで、ベクトルを保存できる。
ただし、このベクトルの次元数は固定であり、異なる次元数のベクトルを同じカラムに保存することはできない。

### 導入方法

基本的には以下のサイトを参考にした。alpine 指定によると思われるえらーは AI に聞いて解決した。
[参考](https://qiita.com/naozo-se/items/0730c8ea650eaa0d51c8)

## ベクトルの距離、類似度計算の種類

- **コサイン距離（<=>）**:意味的な類似性、テキスト埋め込み
- **L2距離（<->）**：位置・距離重視、長さも考慮
- **内積（<#>）**：重み付きベクトルなどでスコアが重要
- **L1距離**：ノイズに強い

### 距離・類似度一覧

| 名称                    | 演算子 | 特徴                                 | 正規化 | 向き/長さ |
|-------------------------|--------|--------------------------------------|--------|-----------|
| L2距離（ユークリッド）  | `<->`  | 直線距離、長さを含めた差             | 不要   | 長さ含む  |
| 内積（ドット積）        | `<#>`  | 向き＋長さの一致度（スコア）         | 不要   | 向き＋長さ|
| コサイン距離            | `<=>`  | 向き（角度）だけで比較、長さ無関係   | 要     | 向きのみ  |
| L1距離（マンハッタン） | ❌未対応 | 各成分差の合計、ノイズに強い         | 不要   | 長さ含む  |

### pgvector 使用例

```sql
-- ユークリッド距離
SELECT * FROM items ORDER BY embedding <-> '[0.1, 0.2, 0.3]' LIMIT 5;

-- 内積
SELECT * FROM items ORDER BY embedding <#> '[0.1, 0.2, 0.3]' LIMIT 5;

-- コサイン距離
SELECT * FROM items ORDER BY embedding <=> '[0.1, 0.2, 0.3]' LIMIT 5;
```

## ベクトル計算の負荷

現状の `paraphrase-multilingual-MiniLM-L12-v2` モデルを利用した場合、ベクトル化の API のレスポンスはほぼなしといっていいレベル。

```log
2025-08-07 05:32:04     NLP request: 200 OK
2025-08-07 05:32:04     >> URL:https://www.city.hamura.tokyo.jp/prsite/0000018548.html
2025-08-07 05:32:04     NLP request: 200 OK
2025-08-07 05:32:04     >> URL:https://www.city.hamura.tokyo.jp/prsite/category/14-4-3-0-0-0-0-0-0-0.html
2025-08-07 05:32:04     NLP request: 200 OK
2025-08-07 05:32:04     >> URL:https://www.city.hamura.tokyo.jp/prsite/0000017650.html
2025-08-07 05:32:04     NLP request: 200 OK
2025-08-07 05:32:04     >> URL:https://www.city.hamura.tokyo.jp/prsite/0000019185.html
2025-08-07 05:32:04     NLP request: 200 OK
2025-08-07 05:32:04     >> URL:https://www.city.hamura.tokyo.jp/prsite/0000019186.html
2025-08-07 05:32:04     NLP request: 200 OK
2025-08-07 05:32:04     >> URL:https://www.city.hamura.tokyo.jp/prsite/0000016876.html
2025-08-07 05:32:04     NLP request: 200 OK
2025-08-07 05:32:04     >> URL:https://www.city.hamura.tokyo.jp/prsite/0000017943.html
2025-08-07 05:32:04     NLP request: 200 OK
2025-08-07 05:32:04     >> URL:https://www.city.hamura.tokyo.jp/prsite/0000005991.html
```

## 計算インフラを考慮したアーキテクチャ選定 3 パターン

1. 通常の EC2, ECS 上で自然言語処理を行う（某インフラ構成と同様）
2. Lambda などのサーバーレスで自然言語処理を行う
   1. パターン1: そのまま自然言語処理を行う（Python, SentenceTransformer, モデルにより重いコールドスタートが発生）
   2. パターン2: ONNX などのコンテナを利用して、モデルを事前にロードしておくことでコールドスタートを回避する
3. SageMaker などのマネージドサービスを利用して自然言語処理を行う

## ONNX

ニューラルネットワークなどの学習済みモデルのデファクトスタンダード。
ONNX（Open Neural Network Exchange）は、異なるフレームワーク間でモデルを共有するためのオープンなフォーマット。これにより、PyTorch や TensorFlow などの異なるフレームワークでトレーニングされたモデルを、他のフレームワークで推論に使用できる。ただし、基本的に静的な計算グラフであるため、動的なモデルには向かないため、主な用途はベクトル化や分類、自然言語生成といった推論に限られる。モデルの学習・訓練、データの前処理、ベクトル間の類似度計算などは通常のフレームワークで行うことが求められる。

### ONNX のメリット

- 異なるフレームワーク間での互換性が高い
- モデルの推論速度が向上することがある
- モデルのサイズが小さくなることがある
- Python 以外の言語（C++, Java など）でも利用可能

### ONNX のデメリット

- 動的な計算グラフには対応していないため、特定のモデルには適用できない
- 一部のフレームワーク固有の機能やオペレーションがサポートされていないことがある
- モデルの変換プロセスが必要

## SageMaker

SageMaker が EC2 などと比較して優れている点は

運用・管理の簡素化

### 自動スケーリング

トラフィックに応じた自動的なインスタンス追加・削除
EC2では手動設定が必要なオートスケーリンググループが自動構築

### マネージドインフラ

OS、Python環境、機械学習ライブラリの管理が不要
セキュリティパッチやアップデートが自動適用

### SageMakerでのモデルデプロイ例

predictor = model.deploy(
    initial_instance_count=1,
    instance_type='ml.t3.medium',
    auto_scaling_enabled=True
)

## トークン化

Golang にはトークン化の公式や強力なライブラリが存在しないため、 github.com/daulet/tokenizers を利用した。以下の方法でのトークン化も試したが、未実装だったり python の tokenizers と同様の結果にならなかったので断念。

トークン数のテスト: <https://huggingface.co/spaces/Xenova/the-tokenizer-playground>

- github.com/sugarme/tokenizer
- AI サポートによる自作

## huggingface モデルのダウンロードについて

hugging face version 指定が望ましい
バージョンが変わってしまうと、以前生成したベクトルデータと新しいモデルで生成したベクトルデータの生成プロセスが変化し、検索結果が変わってしまう可能性があるため。
また、再ダウンロードが生じる可能性もある。
<https://zenn.dev/yagiyuki/articles/load_pretrained>

uv run hf download sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2 onnx/model.onnx tokenizer.json --revision 86741b4e3f5cb7765a600d3a3d55a0f6a6cb443d --local-dir onnx_model
バージョンの指定ができ、おそらく専用のバリデーションとエラーハンドリングが行われるため wget や curl よりも適している
<https://huggingface.co/docs/huggingface_hub/main/en/guides/download>

## CPU の命令セットとモデル選択

<https://huggingface.co/sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2>

model_qint8_arm64.onnx
model_qint8_avx512.onnx
model_qint8_avx512_vnni.onnx
model_quint8_avx2.onnx

`lscpu | grep Flags`

```text
Flags:                                fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ss ht syscall nx pdpe1gb rdtscp lm constant_tsc rep_good nopl xtopology tsc_reliable nonstop_tsc cpuid tsc_known_freq pni pclmulqdq vmx ssse3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm abm 3dnowprefetch ssbd ibrs ibpb stibp ibrs_enhanced tpr_shadow ept vpid ept_ad fsgsbase tsc_adjust bmi1 avx2 smep bmi2 erms invpcid rdseed adx smap clflushopt clwb sha_ni xsaveopt xsavec xgetbv1 xsaves avx_vnni vnmi umip waitpkg gfni vaes vpclmulqdq rdpid movdiri movdir64b fsrm md_clear serialize flush_l1d arch_capabilities
```

avx2 があるので model_quint8_avx2.onnx を使用する
本番は arm64 なので model_qint8_arm64.onnx を使用する

| 項目        | FP32 | INT8量子化      |
| --------- | ---- | ------------ |
| モデルサイズ    | 100% | 約25%         |
| 推論速度(CPU) | 基準   | 1.5〜4倍       |
| 推論速度(GPU) | 基準   | 1.2〜2倍（対応次第） |
| 精度        | 最高   | 0〜数%低下       |
| 対応性       | 高い   | ハード依存あり      |

## モデルの推論プロセス

BERT系Sentence Transformerモデルの推論プロセスを詳しく解説します。

### 1.トークン化フェーズ

処理内容:

- サブワード分割: テキストを語彙に基づいてトークンに分割
- 特殊トークン追加:
  - 文頭に [CLS] トークン (ID: 0)
  - 文末に [SEP] トークン (ID: 2)

結果: [0, 85873, 98449, 154, 20403, 1453, 30, 2]

### 2.入力準備フェーズ

3つの入力の役割:

- input_ids: 実際のトークンID（語彙インデックス）
- attention_mask: 有効トークンを示すマスク（1=有効、0=パディング）
- token_type_ids: 文章セグメント識別（単一文の場合は全て0）

### 3.BERT Encoder処理

Embedding Layer:

Token Embedding + Position Embedding + Segment Embedding
= 各トークンの初期表現ベクトル (768次元)

Multi-Head Self-Attention:

- 各トークンが他の全トークンとの関係性を学習
- アテンションマスクにより、パディング部分は無視
- 複数のアテンションヘッドで異なる観点から関係性を捉える

Feed-Forward Network:

- 各位置で独立に非線形変換を適用
- 表現力を向上させる

Layer Normalization & Residual Connection:

- 訓練の安定化とグラディエント流れの改善

### 4.出力生成

```python
outputs = session.run( ~中略~ )
last_hidden_state = outputs[0]  # [batch_size, seq_length, hidden_size]
print(last_hidden_state.shape)  # 形状: [1, 8, 384]（8トークン、384次元）
```

今回出力が 8 トークンなのは以下が理由

```python
# 入力トークン
Input: "これは日本語の文章です。"
Token IDs: [0, 85873, 98449, 154, 20403, 1453, 30, 2]
Tokens: [<s>, ▁これは, 日本語, の, 文章, です, 。, </s>]
# ↑ 8トークン

# 出力
last_hidden_state.shape: [1, 8, 384]
# ↑ [batch_size, sequence_length, hidden_size]
#   sequence_length = 8（入力トークン数と同じ）
```

- batch_size 1 バッチサイズ（1文）
- sequence_length 8 入力トークン数
- hidden_size 384 各トークンの埋め込み次元数

各トークンが384次元のコンテキスト化された表現ベクトルを持つ。

### 5.センテンス埋め込み生成（プーリング）

プーリングは、BERT系モデルからトークンレベルの埋め込みを文レベルの埋め込みに変換する重要な処理。

Mean Pooling方式:
各トークンの出力ベクトルを単純平均しているように見え文中での重要性が消失しているように見えるが、BERT では各ベクトルに文中の意味の重要度が含まれているため問題ない。

```python
# アテンションマスク（データの有効無効を設定するフィルタ）を考慮した平均プーリング
# アテンションマスクとは、入力シーケンス内の各トークンが有効か無効かを示すバイナリマスク。
# 1は有効トークン、0はパディングトークンを示す。
attention_mask_expanded = np.expand_dims(attention_mask, axis=-1) # アテンションマスクを拡張
sum_embeddings = np.sum(last_hidden_state * attention_mask_expanded, axis=1) # 有効トークンの埋め込みのみを抽出し、合計する
sum_mask = np.clip(attention_mask_expanded.sum(axis=1), a_min=1e-9, a_max=None) #  有効トークン数をカウント（除算用）
sentence_embedding = sum_embeddings / sum_mask # 有効トークン数で除算（平均化）
```

処理詳細:

1. マスク拡張: [1, 8] → [1, 8, 1]
2. 要素積: 有効トークンのみの埋め込みを抽出
3. 合計計算: 有効トークンの埋め込みを合計
4. 正規化: 有効トークン数で除算（平均化）

他のプーリング手法:

- CLS Pooling: [CLS]トークンのみ使用
- Max Pooling: 各次元の最大値を使用

### 6.最終出力

処理フロー図
重要なポイント
コンテキスト考慮: 各トークンは文全体のコンテキストを考慮した表現
位置情報: Position Embeddingにより単語の順序情報を保持
アテンション機構: 重要な単語により多くの注意を向ける
プーリング戦略: タスクに応じて最適な文表現方法を選択
この一連の処理により、意味的に類似した文は近い位置のベクトル空間にマッピングされ、セマンティック検索などの応用が可能になります。

## nlp コンテナのサイズ

最初期の python, sentence-transformers, fastapi, uvicorn, paraphrase-multilingual-MiniLM-L12-v2 通常モデルを使っていたコンテナのサイズは 6GB 超だった。あまりに重たすぎるので以下を行って容量を削減した。
開発用コンテナのサイズは 1GB 程度、本番コンテナサイズは 250MB 程度。

- ONNX モデルを利用することでめちゃくちゃ重い sentence-transformers の依存を削除（4GB?削減）
- sentence-transformers が不要になったことによる python の依存削除（数百MB削減）
- ダウンロードしたあとのモデルキャッシュを削除（数百MB削減）
- 量子化 ONNX モデルの利用（数百MB削減）
- 本番は distroless イメージを利用（数百MB削減）
