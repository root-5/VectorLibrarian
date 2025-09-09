# Notion

あとで Notion へ転記する内容を Markdown に変換しておく。

## DEFAULT CURRENT_TIMESTAMP

`~ TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP` の場合、新規挿入時のみ自動で現在時刻が入る。
「INSERT時のみ」自動セットされ、UPDATE では自動的に変更されないので、 `created_at` は通常問題ないが、
`updated_at` はアプリ側で明示的に `SET updated_at = CURRENT_TIMESTAMP` を指定する必要がある。

## DB のトリガー

データベースにはトリガーと呼ばれる機能があり、特定のイベント（INSERT、UPDATE、DELETEなど）が発生したときに自動的に実行される処理を定義できます。

## ORM の効果

### メリット

- モデルとしてまとめて DB を定義するので、SQL を書くことによる二度書きがなくなる
- DB の管理をアプリ側に持ってきて一元化できる

### デメリット

- ORM の抽象化により、SQL の細かい制御が難しい
- 見た目上のコード記述は減る一方、冗長な SQL で処理をかけてしまう可能性がある
- 独特の記述が必要

### 使用感メモ

適当に Go コード上で string 型とだけ指定したカラムを使ってテーブルを作成すると、 PostgreSQL では TEXT 型として扱われる。
実際には以下の様な記述が必要。これは Gorm でも同様。

```go
type User struct {
  Id        int64     `bun:"id,pk,autoincrement"`
  Name      string    `bun:"name,type:varchar(50),notnull"`
  Email     string    `bun:"email,type:varchar(100),unique,notnull"`
}
```

## 中間構造体

```go
type PageContentBase struct {
  Domain      sql.NullString
  Path        sql.NullString
  Title       sql.NullString
  Description sql.NullString
  Keywords    sql.NullString
  Markdown    sql.NullString
  Hash        sql.NullString
}

type PageContentInsert struct {
  PageContentBase
  CreatedAt sql.NullTime
  UpdatedAt sql.NullTime
}

type PageContentUpdate struct {
  PageContentBase
  UpdatedAt sql.NullTime
  DeletedAt time.Time
}
```

## distroless イメージの種類

- gcr.io/distroless/static-debian12
  - **最も小さいイメージ**
  - glibc や SSL ライブラリすら含まない
  - 動的リンクのバイナリは動かせない（完全な静的リンクバイナリ専用）
  - シェルもパッケージマネージャもなし
  - 主に **Go や Rust の静的ビルド済みバイナリ**を入れて配布する用途

- gcr.io/distroless/base-debian12
  - glibc と最低限の動的リンクライブラリが入っている
  - OpenSSL も含まれており、HTTPS通信可能
  - 動的リンクの Linux バイナリ（C/C++/Java/Node 等）を動かす用途
  - **SSL 必要なアプリケーションの標準的ベース**

- gcr.io/distroless/base-nossl-debian12
  - base-debian12 から **SSL/TLS 関連のライブラリ（OpenSSL, CA証明書）を除いたもの**
  - HTTPS や TLS が不要な場面でさらにサイズを減らすために使う
  - 社内LAN内でHTTPオンリー通信しかしないバッチ処理などで軽量化に有効

- gcr.io/distroless/cc-debian12
  - base-debian12 に加えて \*\*C/C++ランタイムライブラリ（libstdc++ 等）\*\*が含まれる
  - g++/clang でビルドされた C++ アプリを動かす用途
  - SSLも利用可能
  - 主に **動的リンクのC++アプリ** や **一部の機械学習推論バイナリ（ONNX Runtimeなど）** を実行するために使う

## データ構造の分離

データ構造のうち、ドメイン領域のものを domain/model、そうでないものを usecase/entity に分離する構成をとっている。
この時、本来の意味でドメイン領域を考えるなら以下のようになり、同じ id でも表現したい対象によってドメイン領域になったり、ならなかったりする。こうなると、正直ドメインの原理にはのっとっているもののかなり見づらしく、後に編集するときに基準が読み手にわかりづらくなる。

正直、テーブル定義を分離する必要は薄いかもしれない。

```go
// 取引ドメインモデル
type Trade struct {
  ProductIds   []int64    // 取引と商品の関係はドメイン（商売）上の概念
  Amount       float64
}

// 取引DBモデル
type TradeModel struct {
  Trade
  ID        int64         // 技術的な主キー、ドメイン（商売）上は不要
  CreatedAt time.Time     // 取引発生時刻参照はドメイン（商売）上不要
  UpdatedAt time.Time     // 取引更新時刻参照はドメイン（商売）上不要
  DeletedAt sql.NullTime  // 技術的な論理削除など
}

// 商品ドメインモデル
type Product struct {
  Name        string
  Price       float64
}

// 商品DBモデル
type ProductModel struct {
  Product
  ID        int64         // 技術的な主キー、ドメイン（商売）上は不要
  CreatedAt time.Time     // Trade とは異なり、商品登録時刻参照はドメイン（商売）上不要
  UpdatedAt time.Time     // Trade とは異なり、商品更新時刻参照はドメイン（商売）上不要
  DeletedAt sql.NullTime  // 技術的な論理削除など
}
