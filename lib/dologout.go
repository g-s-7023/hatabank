package lib

import (
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

// ログアウトクリック時に実行される
func DoLogout(w http.ResponseWriter, r *http.Request) {
	//===
	//=== ログインの確認とページ遷移
	//===
	// contextの作成
	ctx := appengine.NewContext(r)
	// ユーザの取得
	u := user.Current(ctx)
	if u != nil {
		// ログイン中なら、ログアウトしてログアウトページに遷移
		url, _ := user.LogoutURL(ctx, "/logout")
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		// 既にログアウトしている場合
		http.Error(w, "already logouted", http.StatusInternalServerError)
		return
	}
}

