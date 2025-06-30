package model

import (
	"database/sql"
)

// ページコンテンツ情報
type PageContent struct {
	ID          sql.NullInt64  `json:"id" db:"id"`                   // ID（自動増分）
	Domain      sql.NullString `json:"domain" db:"domain"`           // ドメイン
	Path        sql.NullString `json:"path" db:"path"`               // パス
	Title       sql.NullString `json:"title" db:"title"`             // ページタイトル
	Description sql.NullString `json:"description" db:"description"` // ページ説明
	Keywords    sql.NullString `json:"keywords" db:"keywords"`       // キーワード
	Markdown    sql.NullString `json:"markdown" db:"markdown"`       // コンテンツのマークダウン
	CreatedAt   sql.NullTime   `json:"created_at" db:"created_at"`   // 作成日時
	UpdatedAt   sql.NullTime   `json:"updated_at" db:"updated_at"`   // 更新日時
	Hash        sql.NullString `json:"hash" db:"hash"`               // コンテンツのハッシュ値
}
