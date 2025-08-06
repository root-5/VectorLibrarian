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
