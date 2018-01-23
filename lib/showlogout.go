package lib

import (
	"net/http"
	"html/template"
)

// ログアウトページの表示
func ShowLogout(w http.ResponseWriter, r *http.Request) {
	//===
	//=== ページ遷移
	//===
	// htmlファイルを読み込み
	var html = template.Must(template.ParseFiles(DIR_HTML + "logout.html"))
	if err := html.ExecuteTemplate(w, "logout.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

