package lib

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 入力実行がクリックされたときに実行する関数
func DoIncome(w http.ResponseWriter, r *http.Request) {
	var err error
	//===
	//=== postされた値の受け取り
	//===
	// r.Formに各値を格納
	r.ParseForm()
	// 各キーごとの値を格納する配列
	inputDates := r.Form["Date"]
	inputPayerIndexes := r.Form["Payer"]
	inputDetails := r.Form["Detail"]
	inputPrices := r.Form["Price"]
	//===
	//=== 入力チェックと登録するデータの作成
	//===
	// エラーチェック用の関数
	errorCheck := func(date, payer, detail, price string) (*ParamToInsert, bool) {
		tempParam := new(ParamToInsert)
		// 入力有無のチェック
		if date == "" {
			if payer != "0" || detail != "" || price != "" {
				// 日付がないのに他の値が入っていたらエラー
				return nil, false
			}
			// 全て空欄ならとばして次の行へ
			return nil, true
		} else {
			if payer == "0" || price == "" {
				// 日付が入っていて、費目または金額が入っていなければエラー
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
		// 費目のチェック
		if p, err := strconv.Atoi(payer); err != nil {
			// 費目が数値でなければエラー
			return nil, false
		} else {
			if _, ok := Payers[p]; ok {
				// 費目が定義されていれば、DB登録用のパラメータにコピー
				tempParam.Payer = p
			} else {
				// 費目が定義されていなければエラー
				return nil, false
			}
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
		tempParam.Type = INCOME
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
		p, e := errorCheck(inputDates[i], inputPayerIndexes[i], inputDetails[i], inputPrices[i])
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
		// payerIndexはintに直す
		pIndex := make([]int, INPUTLINES)
		for i := 0; i < INPUTLINES; i++ {
			if v, err := strconv.Atoi(inputPayerIndexes[i]); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				pIndex[i] = v
			}
		}
		paramToReturn := ParamToShowInput{
			PayerList:   make(map[int]string),
			Lines:       make([]int, INPUTLINES),
			Date:        inputDates,
			Payer:       pIndex,
			Detail:      inputDetails,
			Price:       inputPrices,
			HasAnyError: hasAnyError,
			HasError:    hasError,
		}
		// 選択肢用の費目
		for key, val := range Payers {
			paramToReturn.PayerList[key] = val.Name
		}
		// 行番号の追加
		for i := 0; i < INPUTLINES; i++ {
			paramToReturn.Lines[i] = i
		}
		// htmlファイルを読み込み
		html := template.Must(template.ParseFiles(DIR_HTML + "income.html"))
		if err := html.ExecuteTemplate(w, "income.html", paramToReturn); err != nil {
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
	http.Redirect(w, r, "/user/income", http.StatusFound)
}
