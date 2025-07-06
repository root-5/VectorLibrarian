package model

import (
	"time"

	"github.com/uptrace/bun"
)

// ページコンテンツ情報
type Page struct {
	bun.BaseModel `bun:"table:pages"`

	Id          string    `bun:"id,pk,autoincrement"`                         // ID
	Domain      string    `bun:"domain,notnull,unique:url"`                   // ドメイン
	Path        string    `bun:"path,notnull,unique:url"`                     // パス
	Title       string    `bun:"title,notnull"`                               // ページタイトル
	Description string    `bun:"description,notnull"`                         // ページ説明
	Keywords    string    `bun:"keywords,notnull"`                            // キーワード
	Markdown    string    `bun:"markdown,notnull"`                            // Markdown コンテンツ
	Hash        string    `bun:"hash,notnull"`                                // コンテンツのハッシュ値
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"` // 作成日時
	UpdatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp"` // 更新日時
	DeletedAt   time.Time `bun:",soft_delete"`                                // 削除日時
}
