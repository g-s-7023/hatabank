package lib

import (
	"html/template"
	"net/http"
)

// 入力画面の表示
func ShowExpense(w http.ResponseWriter, r *http.Request) {
	//===
	//=== htmlに渡すパラメータの作成(日付・費目・詳細・価格)
	//===
	paramToShowInput := ParamToShowInput{
		Date:         make([]string, INPUTLINES),
		Detail:       make([]string, INPUTLINES),
		Price:        make([]string, INPUTLINES),
		HasAnyError:  false,
		HasError:     make([]bool, INPUTLINES),
	}
	// 行番号の追加
	for i := 0; i < INPUTLINES; i++ {
		paramToShowInput.Lines = append(paramToShowInput.Lines, i)
	}
	//===
	//=== ページ遷移
	//===
	// htmlファイルを読み込み
	html := template.Must(template.ParseFiles(DIR_HTML + "expense.html"))
	if err := html.ExecuteTemplate(w, "expense.html", paramToShowInput); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
