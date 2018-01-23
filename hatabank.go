package hatabank_gae

import (
	"net/http"
	"hatabank_gae/lib"
)

func init() {
	http.HandleFunc("/", lib.Entry)
	http.HandleFunc("/logout", lib.ShowLogout)
	http.HandleFunc("/user/list", lib.ShowList)
	http.HandleFunc("/user/income", lib.ShowIncome)
	http.HandleFunc("/user/expense", lib.ShowExpense)
	http.HandleFunc("/user/summary", lib.ShowSummary)
	http.HandleFunc("/user/doincome", lib.DoIncome)
	http.HandleFunc("/user/doexpense", lib.DoExpense)
	http.HandleFunc("/user/dodelete", lib.DoDelete)
	http.HandleFunc("/dologout", lib.DoLogout)
	http.Handle("/css", http.FileServer(http.Dir("css")))
}
