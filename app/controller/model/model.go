package model

import (
	"app/domain/model"
	"database/sql"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type PageContentAll struct {
	model.PageContent              // 埋め込みで PageContent のフィールドを使用
	ID                int64        `json:"id" column:"id" dbType:"bigint"`                                             // ID（自動増分）
	Hash              string       `json:"hash" column:"hash" dbType:"char(40)"`                                       // コンテンツのハッシュ値
	CreatedAt         time.Time    `json:"created_at" column:"created_at" dbType:"timestamp with time zone"`           // 作成日時
	UpdatedAt         time.Time    `json:"updated_at" column:"updated_at" dbType:"timestamp with time zone"`           // 更新日時
	DeletedAt         sql.NullTime `json:"deleted_at,omitempty" column:"deleted_at" dbType:"timestamp with time zone"` // 削除日時（論理削除用）
}

type PageContentAtInsert struct {
	model.PageContent           // 埋め込みで PageContent のフィールドを使用
	Hash              string    `json:"hash" column:"hash" dbType:"char(40)"`                             // コンテンツのハッシュ値
	CreatedAt         time.Time `json:"created_at" column:"created_at" dbType:"timestamp with time zone"` // 作成日時
	UpdatedAt         time.Time `json:"updated_at" column:"updated_at" dbType:"timestamp with time zone"` // 更新日時
}

type PageContentAtUpdate struct {
	model.PageContent              // 埋め込みで PageContent のフィールドを使用
	Hash              string       `json:"hash" column:"hash" dbType:"char(40)"`                                       // コンテンツのハッシュ値
	UpdatedAt         time.Time    `json:"updated_at" column:"updated_at" dbType:"timestamp with time zone"`           // 更新日時
	DeletedAt         sql.NullTime `json:"deleted_at,omitempty" column:"deleted_at" dbType:"timestamp with time zone"` // 削除日時（論理削除用）
}

// 任意の構造体からカラム名リストを生成できるように修正
func GetColumnNames(v interface{}) (columnNames []string) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("column")
		if tag != "" {
			columnNames = append(columnNames, tag)
		}
	}
	return columnNames
}

// カラム名リストからプレーススホルダを生成
func GetPlaceholders(columnNames []string) (placeholders []string) {
	placeholders = make([]string, len(columnNames))
	placeholderNumber := 1
	for i := range columnNames {
		// カラム名が "_at" の場合は "CURRENT_TIMESTAMP"、それ以外は "$i" の形式でプレースホルダを生成
		if strings.HasSuffix(columnNames[i], "_at") {
			placeholders[i] = "CURRENT_TIMESTAMP"
		} else {
			placeholders[i] = "$" + strconv.Itoa(placeholderNumber)
			placeholderNumber++
		}
	}
	return placeholders
}

// WHERE 句のターゲット文字列を生成
func GetWhereTargetStr(columnNames []string, placeholders []string, whereTargets []interface{}) string {
	whereTargetsStr := ""
	for i, whereTarget := range whereTargets {
		if i > 0 {
			whereTargetsStr += " AND "
		}
		whereTargetsStr += columnName + " = " + placeholders[i]
	}
	return whereTargetsStr
}