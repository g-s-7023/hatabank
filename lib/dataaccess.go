package lib

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"strconv"
	"time"
)

// 入力されたデータをDBに登録する関数
func InsertData(r *http.Request, p []ParamToInsert) error {
	var err error
	//===
	//=== DB登録用のパラメータの作成
	//===
	// 登録するエントリ数
	putNum := len(p)
	// DBに登録する家計簿のエントリ
	var entryToInsert = make([]Hatabank, putNum)
	// 登録用エントリへの値のコピー
	for i := 0; i < putNum; i++ {
		// 曜日の算出
		dayOfWeek := time.Date(p[i].Year, time.Month(p[i].Month), p[i].Day, 0, 0, 0, 0, jst).Weekday()
		entryToInsert[i] = Hatabank{
			Day:       p[i].Day,
			Month:     p[i].Month,
			Year:      p[i].Year,
			DayOfWeek: int(dayOfWeek),
			Payer:     p[i].Payer,
			Type:      p[i].Type,
			Price:     p[i].Price,
			Detail:    p[i].Detail,
		}
	}
	//===
	//=== DB登録
	//===
	// Contextの作成
	ctx := appengine.NewContext(r)
	// 登録用のKeyをputNum分作成
	keys := make([]*datastore.Key, putNum)
	for i := 0; i < putNum; i++ {
		keys[i] = datastore.NewIncompleteKey(ctx, ENTRY,
			getMonthKey(ctx, hatabank, entryToInsert[i].Year, entryToInsert[i].Month))
	}
	// トランザクションで更新
	err = datastore.RunInTransaction(ctx, func(c context.Context) error {
		_, e := datastore.PutMulti(c, keys, entryToInsert)
		return e
	}, nil)
	if err != nil {
		// トランザクションがエラーならerrを返却
		return err
	}
	return nil
}

// 一覧画面に表示するデータの検索を行う関数
func ReadList(r *http.Request, year, month int) (*ParamToShowList, error) {
	var err error
	//===
	//=== エントリの取得
	//===
	// contextの作成
	ctx := appengine.NewContext(r)
	// queryの作成
	// 指定した家計簿中の年・月で指定したエントリを全て表示
	query := datastore.NewQuery(ENTRY).
		Ancestor(getMonthKey(ctx, hatabank, year, month)).
		Order("-Day")
	// queryの結果格納用のsliceの作成
	var tempArray []Hatabank
	// queryの実行結果のキー格納用のsliceの作成
	var keys []*datastore.Key

	// トランザクションでクエリを実行
	err = datastore.RunInTransaction(ctx, func(c context.Context) error {
		var e error
		keys, e = query.GetAll(c, &tempArray)
		return e
	}, nil)
	if err != nil {
		// エラーならnilを返す
		return nil, err
	}
	//===
	//=== エントリの出力
	//===
	// 出力用の構造体の作成
	p := new(ParamToShowList)
	p.Year = strconv.Itoa(year)
	if month < 10 {
		// 10以下の場合は0詰め
		p.Month = "0" + strconv.Itoa(month)
	} else {
		p.Month = strconv.Itoa(month)
	}
	p.EntryToShow = make([]WebHatabank, len(keys))
	// tempの内容とkeyを出力用の配列に格納
	for i := 0; i < len(keys); i++ {
		// Key
		p.EntryToShow[i].DatastoreId = strconv.FormatInt(keys[i].IntID(), 10)
		// 年
		p.EntryToShow[i].Year = strconv.Itoa(tempArray[i].Year)
		// 月
		if tempArray[i].Month < 10 {
			// 10以下の場合は0詰め
			p.EntryToShow[i].Month = "0" + strconv.Itoa(tempArray[i].Month)
		} else {
			p.EntryToShow[i].Month = strconv.Itoa(tempArray[i].Month)
		}
		// 日
		if tempArray[i].Day < 10 {
			// 10以下の場合は0詰め
			p.EntryToShow[i].Day = "0" + strconv.Itoa(tempArray[i].Day)
		} else {
			p.EntryToShow[i].Day = strconv.Itoa(tempArray[i].Day)
		}
		// 曜日
		p.EntryToShow[i].DayOfWeek = getWeekString(tempArray[i].DayOfWeek)
		// 費目
		if tempArray[i].Type == INCOME {
			p.EntryToShow[i].Type = "収入"
		} else if tempArray[i].Type == EXPENSE {
			p.EntryToShow[i].Type = "支出"
		}
		// 支払者
		if v, ok := Payers[tempArray[i].Payer]; ok {
			p.EntryToShow[i].Payer = v.Name
		} else {
			// 対応する費目が見つからなかった場合、空白
			p.EntryToShow[i].Payer = ""
		}
		// 支払者のインデックス
		p.EntryToShow[i].PayerIndex = tempArray[i].Payer
		p.EntryToShow[i].Detail = tempArray[i].Detail
		p.EntryToShow[i].Price = strconv.Itoa(tempArray[i].Price)
	}
	return p, nil
}

// まとめ画面に表示するデータの検索を行う関数
func ReadSummary(r *http.Request, year int) (*ParamToShowSummary, error) {
	var err error
	//===
	//=== 表示用のパラメータの初期化
	//===
	p := new(ParamToShowSummary)
	p.SumOfMonth = make([]int, 12)
	p.ExpenseOfMonth = make([]int, 12)
	p.Results = make(map[int]*ResultOfMonth)
	for k, v := range Payers {
		if v.IsCalcSummary {
			// 計算対象となっている費目だけResults[]を作成
			p.Results[k] = &ResultOfMonth{
				Name:    v.Name,
				Summary: make([]int, 12),
			}
		}
	}
	p.Year = year
	p.YearList = make([]int, 0)
	// いったんUTCで現在時刻をフォーマットして、JSTに変換
	thisYear := time.Now().UTC().In(jst).Year()
	// 年のリストに表示する年を格納
	for y := thisYear; y >= STARTYEAR; y-- {
		p.YearList = append(p.YearList, y)
	}
	//===
	//=== エントリの取得
	//===
	// contextの作成
	ctx := appengine.NewContext(r)
	// queryの作成
	// 指定した家計簿中の年・月で指定したエントリを全て表示
	query := datastore.NewQuery(ENTRY).
		Ancestor(getKakeiboKey(ctx, hatabank)).
		Filter("Year =", p.Year)
	// queryの結果格納用のsliceの作成
	var tempArray []Hatabank
	// トランザクションでクエリを実行
	err = datastore.RunInTransaction(ctx, func(c context.Context) error {
		_, e := query.GetAll(c, &tempArray)
		return e
	}, nil)
	if err != nil {
		// エラーの場合errを返却
		return nil, err
	}
	//===
	//=== 集計
	//===
	// queryの結果を月と費目ごとに集計
	for _, v := range tempArray {
		if v.Type == INCOME {
			// 収支に加算
			p.TotalBudget += v.Price
			// 金額を合計金額に加算
			p.SumOfMonth[v.Month-1] += v.Price
			// 費目別に金額を加算
			if _, ok := Payers[v.Payer]; ok {
				if Payers[v.Payer].IsCalcSummary {
					// 集計対象としている人であれば、それに加算
					p.Results[v.Payer].Summary[v.Month-1] += v.Price
				}
			}
		} else if v.Type == EXPENSE {
			// 収支から減算
			p.TotalBudget -= v.Price
			// 金額を合計金額に加算
			p.ExpenseOfMonth[v.Month-1] += v.Price
		}
	}
	return p, nil
}

// DBからデータ削除を行う関数
func Delete(r *http.Request, id, year, month int) error {
	var err error
	//===
	//=== エントリの削除
	//===
	// Contextの作成
	ctx := appengine.NewContext(r)
	// 削除対象のエンティティのキーの取得
	keyForDelete := datastore.NewKey(ctx, ENTRY, "", int64(id), getMonthKey(ctx, hatabank, year, month))
	err = datastore.RunInTransaction(ctx, func(c context.Context) error {
		e := datastore.Delete(c, keyForDelete)
		return e
	}, nil)
	if err != nil {
		return err
	}
	// エラーがなければnilを返す
	return nil
}

// 家計簿の名前に対応するキーを返す関数
func getKakeiboKey(c context.Context, kakeiboName string) *datastore.Key {
	return datastore.NewKey(c, ID, kakeiboName, 0, nil)
}

// 年に対応するキーを返す関数
func getYearKey(c context.Context, kakeiboName string, year int) *datastore.Key {
	return datastore.NewKey(c, YEAR, "", int64(year), getKakeiboKey(c, kakeiboName))
}

// 月に対応するキーを返す関数
func getMonthKey(c context.Context, kakeiboName string, year, month int) *datastore.Key {
	return datastore.NewKey(c, MONTH, "", int64(month), getYearKey(c, kakeiboName, year))
}
