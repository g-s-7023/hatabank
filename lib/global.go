package lib

import (
	"time"
)

const (
	//===
	//=== htmlファイルのあるディレクトリ
	//===
	DIR_HTML = "html/"
	//===
	//=== 一覧(Input.html)に表示する行数
	//===
	INPUTLINES = 15
	//===
	//=== 使いはじめる年
	//===
	STARTYEAR = 2017
	//===
	//=== 収入・支出
	//===
	INCOME  = 0
	EXPENSE = 1
	//===
	//=== カインド
	//===
	// ルートエンティティのカインド(ユーザごとに異なる家計簿のID)
	ID = "Hatabank"
	// 年エンティティのカインド(親キー：KakeiboIDカインドのキー)
	YEAR = "HatabankYear"
	// 月エンティティのカインド(親キー：KakeiboYearカインドのキー)
	MONTH = "HatabankMonth"
	// 家計簿の各エントリが入るカインド(親キー：KakeiboMonthカインドのキー)
	ENTRY = "HatabankEntry"
	//===
	//=== 一括削除する際の月の値
	//===
	ALL = 13
	//===
	//=== 各メンバーに該当する番号
	//===
	NONE     = 0
	EZOE     = 1
	OODATE   = 2
	KITAHARA = 3
	KUBOTA   = 4
	SHIMIZU  = 5
	SUGAYA   = 6
	NAKAMURA = 7
	HATA     = 8
	HATANO   = 9
	YAMAZOE  = 10
	OTHER    = 99
)

// DBに登録・参照するための構造体
type Hatabank struct {
	Day       int
	Month     int
	Year      int
	DayOfWeek int
	Payer     int
	Type      int
	Price     int
	Detail    string
}

// 一覧ページで出力するために各値を格納する構造体
type WebHatabank struct {
	DatastoreId string
	Day         string
	Month       string
	Year        string
	DayOfWeek   string
	Payer       string
	PayerIndex  int
	Type        string
	Price       string
	Detail      string
}

// 入力ページで出力するために各値を格納する構造体
type ParamToShowInput struct {
	PayerList   map[int]string
	Lines       []int
	HasAnyError bool
	HasError    []bool
	Date        []string
	Payer       []int
	Detail      []string
	Price       []string
}

// DB登録を行う関数に渡す構造体
type ParamToInsert struct {
	Day    int
	Month  int
	Year   int
	Payer  int
	Price  int
	Type   int
	Detail string
}

// 更新ページを出力するために書く値を格納する構造体
type ParamToShowUpdate struct {
	PayerList   map[int]string
	HasAnyError bool
	DatastoreId string
	Date        string
	Payer       string
	PayerIndex  int
	Price       string
	Detail      string
}

// DB更新を行う関数に渡す構造体
type ParamToUpdate struct {
	ID          int
	Day         int
	Month       int
	MonthBefore int
	Year        int
	YearBefore  int
	Payer       int
	Detail      string
	Price       int
}

// 一覧を表示するためにhtmlTemplateに渡す構造体
type ParamToShowList struct {
	Year        string
	Month       string
	EntryToShow []WebHatabank
}

// 支出者とまとめの計算有無
type HatabankPayer struct {
	Name          string
	IsCalcSummary bool
}

// 各人からの名称と各月の収入を格納する構造体
type ResultOfMonth struct {
	// 費目に対応する文字列
	Name string
	// 各月の金額
	Summary []int
}

// まとめを表示するためにhtmlTemplateに渡す構造体
type ParamToShowSummary struct {
	// 表示する年
	Year int
	// 集計対象の年のリスト
	YearList []int
	// 収支
	TotalBudget int
	// 各月の収入合計金額を格納
	SumOfMonth []int
	// 各月の支出金額を格納
	ExpenseOfMonth []int
	// 表示対象の費目に対応した文字列と各月の金額
	Results map[int]*ResultOfMonth
}

var (
	//===
	//=== 家計簿IDカインドのキー
	//===
	hatabank = "hatabank"
	//===
	//=== タイムゾーン
	//===
	jst = time.FixedZone("Asia/Tokyo", 9*60*60)
	//===
	//=== 支払者で選択可能な文字列とまとめの計算有無
	//===
	Payers = map[int]HatabankPayer{
		NONE:     {"", false},
		EZOE:     {"江添", true},
		OODATE:   {"大舘", true},
		KITAHARA: {"北原", true},
		KUBOTA:   {"久保田", true},
		SHIMIZU:  {"清水", true},
		SUGAYA:   {"菅谷", true},
		NAKAMURA: {"中村", true},
		HATA:     {"畑", true},
		HATANO:   {"羽田野", true},
		YAMAZOE:  {"山添", true},
		OTHER:    {"その他", false},
	}
)

// 曜日に対応する文字列を返す関数
func getWeekString(w int) string {
	wday := []string{"(日)", "(月)", "(火)", "(水)", "(木)", "(金)", "(土)"}
	if w >= 0 && w < 7 {
		return wday[w]
	} else {
		return ""
	}
}
