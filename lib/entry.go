package lib

import (
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

// 最初に実行される
func Entry(w http.ResponseWriter, r *http.Request) {
	//===
	//=== ログインの確認とページ遷移
	//===
	// contextの作成
	ctx := appengine.NewContext(r)
	// ユーザの取得
	u := user.Current(ctx)
	if u == nil {
		// ログインしていなければ、ログインページに遷移
		url, _ := user.LoginURL(ctx, "/user/list")
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		http.Redirect(w, r, "/user/list", http.StatusFound)
	}
}


