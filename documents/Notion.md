# Notion

あとで Notion へ転記する内容を Markdown に変換しておく。

## DEFAULT CURRENT_TIMESTAMP

はい、created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP の場合、
新規挿入時のみ自動で現在時刻が入ります。
created_at は「INSERT時のみ」自動セットされ、UPDATE では自動的に変更されません。

一方、updated_at も DEFAULT CURRENT_TIMESTAMP だと「INSERT時のみ」自動セットされますが、
UPDATE時に自動で更新するには、アプリ側で明示的に updated_at = CURRENT_TIMESTAMP を指定する必要があります。

**まとめ**
created_at の DEFAULT CURRENT_TIMESTAMP は通常問題ありません（UPDATE時に自動で変わることはありません）。
updated_at はUPDATE時にアプリ側で SET updated_at = CURRENT_TIMESTAMP するのが一般的です。
DEFAULT がついていても、UPDATE時に自動で値が変わることはありません。
（PostgreSQLのトリガー等を使えば自動更新も可能ですが、SQL定義だけでは自動更新されません。）

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
