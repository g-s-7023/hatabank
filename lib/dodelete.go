package lib

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
	"errors"
)

// 削除処理を行う関数
func DoDelete(w http.ResponseWriter, r *http.Request) {
	var err error
	//===
	//=== postされた値の受け取りと入力チェック
	//===
	// 各キーごとの値を格納する変数
	idToDelete, _ := strconv.Atoi(r.FormValue("ID"))
	dateBefore := r.FormValue("DateBefore")
	// 日付の入力値チェック
	var yearBefore, monthBefore int
	// 入力値チェック用の関数
	errorCheck := func() error {
		if date := strings.Split(dateBefore, "-"); len(date) != 3 {
			// 年月日の3要素以外の場合、エラー
			return errors.New("invalid format date")
		} else {
			if _, err = time.Parse("2006-01-02", dateBefore); err != nil {
				// 日付が現実にあるものでなければエラー
				return err
			}
			yearBefore, _ = strconv.Atoi(date[0])
			monthBefore, _ = strconv.Atoi(date[1])
		}
		return nil
	}
	// yearかmonthが数値でない場合、エラー
	if err := errorCheck(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//===
	//=== 削除処理
	//===
	err = Delete(r, idToDelete, yearBefore, monthBefore)
	//===
	//=== ページ遷移
	//===
	// データ削除エラーの場合
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// htmlに渡すパラメータを作成
	src := "delete"
	// htmlファイルを読み込み
	html := template.Must(template.ParseFiles(DIR_HTML + "succeed.html"))
	if err = html.ExecuteTemplate(w, "succeed.html", src); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
