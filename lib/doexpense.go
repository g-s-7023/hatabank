package lib

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 入力実行がクリックされたときに実行する関数
func DoExpense(w http.ResponseWriter, r *http.Request) {
	var err error
	//===
	//=== postされた値の受け取り
	//===
	// r.Formに各値を格納
	r.ParseForm()
	// 各キーごとの値を格納する配列
	inputDates := r.Form["Date"]
	inputDetails := r.Form["Detail"]
	inputPrices := r.Form["Price"]
	//===
	//=== 入力チェックと登録するデータの作成
	//===
	// エラーチェック用の関数
	errorCheck := func(date, detail, price string) (*ParamToInsert, bool) {
		tempParam := new(ParamToInsert)
		// 入力有無のチェック
		if date == "" {
			if detail != "" || price != "" {
				// 日付がないのに他の値が入っていたらエラー
				return nil, false
			}
			// 全て空欄ならとばして次の行へ
			return nil, true
		} else {
			if price == "" {
				// 日付が入っていて、金額が入っていなければエラー
				return nil, false
			}
		}
		// 日付のチェック
		if inputDate := strings.Split(date, "-"); len(inputDate) != 3 {
			// 年・月・日の3要素のみでなかった場合、エラー
			return nil, false
		} else {
			if _, err = time.Parse("2006-01-02", date); err != nil {
				// 日付が現実にあるものでなければエラー
				return nil, false
			}
			tempParam.Year, err = strconv.Atoi(inputDate[0])
			tempParam.Month, err = strconv.Atoi(inputDate[1])
			tempParam.Day, err = strconv.Atoi(inputDate[2])
		}
		// 価格のチェック
		if price, err := strconv.Atoi(price); err != nil {
			// 価格が数値でなければエラー
			return nil, false
		} else {
			if price < 0 {
				// 価格が0以上でなければエラー
				return nil, false
			} else {
				// 価格が0以上なら、DB登録用のパラメータにコピー
				tempParam.Price = price
			}
		}
		// エラーがなければDetailを追加してparamToInsertを返却
		tempParam.Detail = detail
		tempParam.Type = EXPENSE
		return tempParam, true
	}
	// DB登録用の関数に渡すパラメータ
	paramToInsert := make([]ParamToInsert, 0)
	// 各入力行のエラー有無
	hasError := make([]bool, INPUTLINES)
	// 1行でもエラーがあるかどうか
	hasAnyError := false
	// 各行に対してエラーチェック関数を実行
	for i := 0; i < INPUTLINES; i++ {
		p, e := errorCheck(inputDates[i],inputDetails[i], inputPrices[i])
		if !e {
			// エラーがある場合、hasErrorとhasAnyErrorをtrueにする
			hasAnyError = true
			hasError[i] = true
		} else {
			// エラーがない場合
			if p != nil {
				// 登録用のパラメータが返ってきたらDB登録用のパラメータに追加
				paramToInsert = append(paramToInsert, *p)
			}
			// hasErrorをfalseにする
			hasError[i] = false
		}
	}
	//===
	//=== 入力エラー時の画面遷移
	//===
	// hasAnyErrorがtrueの場合、登録せずに元の画面に戻る
	if hasAnyError {
		// 元の画面にセットするパラメータの作成
		paramToReturn := ParamToShowInput{
			Lines:        make([]int, INPUTLINES),
			Date:         inputDates,
			Detail:       inputDetails,
			Price:        inputPrices,
			HasAnyError:  hasAnyError,
			HasError:     hasError,
		}
		// 行番号の追加
		for i := 0; i < INPUTLINES; i++ {
			paramToReturn.Lines[i] = i
		}
		// htmlファイルを読み込み
		html := template.Must(template.ParseFiles(DIR_HTML + "expense.html"))
		if err := html.ExecuteTemplate(w, "expense.html", paramToReturn); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	//===
	//=== データの登録
	//===
	err = InsertData(r, paramToInsert)
	//===
	//=== ページ遷移
	//===
	// データ登録エラーの場合
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/user/expense", http.StatusFound)
}

