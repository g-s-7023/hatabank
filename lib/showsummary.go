package lib

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

func ShowSummary(w http.ResponseWriter, r *http.Request) {
	var err error
	//===
	//=== postで渡された値の格納
	//===
	requestedYear := r.FormValue("selected_year")
	//===
	//=== 入力チェックと表示する年の設定
	//===
	// 表示する年を格納する変数
	var yearToShow int
	if requestedYear == "" {
		// requestedYearがセットされていない場合、現在の年をセット
		yearToShow = time.Now().UTC().In(jst).Year()
	} else {
		// requestMonthがセットされている場合、その値を設定
		yearToShow, err = strconv.Atoi(requestedYear)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	//===
	//=== 表示するデータの取得
	//===
	paramToShowSummary, err := ReadSummary(r, yearToShow)
	//===
	//=== エントリの出力
	//===
	// エラーの場合
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// htmlテンプレートを読み込み
	html := template.Must(template.ParseFiles(DIR_HTML + "summary.html"))
	// htmlの出力
	if err := html.ExecuteTemplate(w, "summary.html", paramToShowSummary); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
