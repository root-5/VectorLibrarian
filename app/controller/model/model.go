// package model

// import (
// 	"app/domain/model"
// 	"database/sql"
// 	"time"
// )

// type PageContentAll struct {
// 	model.PageContent              // 埋め込みで PageContent のフィールドを使用
// 	ID                int64        `json:"id" column:"id"`                           // ID（自動増分）
// 	Hash              string       `json:"hash" column:"hash"`                       // コンテンツのハッシュ値
// 	CreatedAt         time.Time    `json:"created_at" column:"created_at"`           // 作成日時
// 	UpdatedAt         time.Time    `json:"updated_at" column:"updated_at"`           // 更新日時
// 	DeletedAt         sql.NullTime `json:"deleted_at,omitempty" column:"deleted_at"` // 削除日時（論理削除用）
// }

// type PageContentAtInsert struct {
// 	model.PageContent           // 埋め込みで PageContent のフィールドを使用
// 	Hash              string    `json:"hash" column:"hash"`             // コンテンツのハッシュ値
// 	CreatedAt         time.Time `json:"created_at" column:"created_at"` // 作成日時
// 	UpdatedAt         time.Time `json:"updated_at" column:"updated_at"` // 更新日時
// }

// type PageContentAtUpdate struct {
// 	model.PageContent              // 埋め込みで PageContent のフィールドを使用
// 	Hash              string       `json:"hash" column:"hash"`                       // コンテンツのハッシュ値
// 	UpdatedAt         time.Time    `json:"updated_at" column:"updated_at"`           // 更新日時
// 	DeletedAt         sql.NullTime `json:"deleted_at,omitempty" column:"deleted_at"` // 削除日時（論理削除用）
// }
