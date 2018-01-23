package lib

import (
	"html/template"
	"net/http"
)

// 入力画面の表示
func ShowIncome(w http.ResponseWriter, r *http.Request) {
	//===
	//=== htmlに渡すパラメータの作成(日付・費目・詳細・価格)
	//===
	paramToShowInput := ParamToShowInput{
		PayerList:   make(map[int]string),
		Date:        make([]string, INPUTLINES),
		Payer:       make([]int, INPUTLINES),
		Detail:      make([]string, INPUTLINES),
		Price:       make([]string, INPUTLINES),
		HasAnyError: false,
		HasError:    make([]bool, INPUTLINES),
	}
	// 選択肢用の費目
	for key, val := range Payers {
		paramToShowInput.PayerList[key] = val.Name
	}
	// 行番号の追加
	for i := 0; i < INPUTLINES; i++ {
		paramToShowInput.Lines = append(paramToShowInput.Lines, i)
	}
	//===
	//=== ページ遷移
	//===
	// htmlファイルを読み込み
	html := template.Must(template.ParseFiles(DIR_HTML + "income.html"))
	if err := html.ExecuteTemplate(w, "income.html", paramToShowInput); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
