package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	file := string(Must(os.ReadFile("sql/medical_db_build.sqlite")))

	db := Must(sql.Open("sqlite", "test.db"))

	if _, err := os.ReadFile("test.db"); err != nil {
		_ = Must(db.Exec(file))
	}

	index := Must(template.ParseFiles("./html/templates/index.html"))
	table := Must(template.ParseFiles("./html/templates/table.html"))
	_ = http.FileServer(http.Dir("./html/static"))

	// data := TableData{
	// 	Headings: []any{"test heading", "test_heading_2"},
	// 	Rows: [][]any{
	// 		{"test body", "test body 2"},
	// 	},
	// }

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		index.Execute(w, nil)
	})

	http.HandleFunc("GET /report_providers", func(w http.ResponseWriter, r *http.Request) {
		table_data := TransposeStructsToTableData(Must(GetAllReportProviders(db)))

		index.Execute(w, table_data)
	})

	http.HandleFunc("GET /history", func(w http.ResponseWriter, r *http.Request) {
		table_data := Must(GetReportHistoryPretty(db))

		table.Execute(w, table_data)
	})
	// test_result := Must(db.Query("select * from report_providers"))

	// for test_result.Next() {
	// 	var prov = ReportProviders{}
	// 	Ok(test_result.Scan(&prov.id, &prov.provider))
	// 	fmt.Printf("Provider ID: %d, Provider Name: %s\n", prov.id, prov.provider)
	// }

	// _ = Must(ReportProviders.Insert(ReportProviders{id: 0, provider: "Creekside"}, db))

	// test_result = Must(db.Query("select * from report_providers"))

	// for test_result.Next() {
	// 	var prov = ReportProviders{}
	// 	Ok(test_result.Scan(&prov.id, &prov.provider))
	// 	fmt.Printf("Provider ID: %d, Provider Name: %s\n", prov.id, prov.provider)
	// }
	// table := TransposeStructs([]ReportProviders{{Id: 1, Provider: "test"}, {Id: 16, Provider: "bruh"}})
	// for i := range table {
	// 	for j := range table[i] {
	// 		fmt.Println(table[i][j])
	// 	}
	// }
	// fmt.Println(table)

	// fmt.Println(Must(GetReportHistoryPretty(db)))

	log.Println("Starting server at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Must[T any](t T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func Ok(err error) {
	if err != nil {
		log.Fatalf("Return is not Ok: %s", err)
	}
}
