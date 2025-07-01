package model

// ページコンテンツ情報
type PageContent struct {
	Domain      string `json:"domain" column:"domain" dbType:"varchar(255)" unique:"true"` // ドメイン
	Path        string `json:"path" column:"path" dbType:"varchar(255)" unique:"true"`     // パス
	Title       string `json:"title" column:"title" dbType:"varchar(255)"`                 // ページタイトル
	Description string `json:"description" column:"description" dbType:"varchar(255)"`     // ページ説明
	Keywords    string `json:"keywords" column:"keywords" dbType:"varchar(255)"`           // キーワード
	Markdown    string `json:"markdown" column:"markdown" dbType:"text"`                   // HTML コンテンツのマークダウン
}
