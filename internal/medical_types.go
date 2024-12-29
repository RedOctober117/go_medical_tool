package main

import (
	"database/sql"
	"reflect"
)

type Queryable interface {
	Insert(db *sql.DB) (sql.Result, error)
	Update(db *sql.DB) (sql.Result, error)
	Delete(db *sql.DB) (*sql.Rows, error)
}

type ReportProviders struct {
	Id       int
	Provider string
}

func (r ReportProviders) Insert(db *sql.DB) (sql.Result, error) {
	return db.Exec("insert into report_providers (Provider) values (?)", r.Provider)
}

func GetAllReportProviders(db *sql.DB) ([]ReportProviders, error) {
	rows, err := db.Query("select * from report_providers")

	if err != nil {
		return nil, err
	}

	results := make([]ReportProviders, 0)
	for rows.Next() {
		temp := new(ReportProviders)
		Ok(rows.Scan(&temp.Id, &temp.Provider))
		results = append(results, *temp)
	}

	return results, nil
}

type ReportTypes struct {
	Id   int
	Type string
}

type Report struct {
	Id   int
	Date string
}

func GetReportHistoryPretty(db *sql.DB) (TableData, error) {
	rows, err := db.Query(`select report_providers.provider as 'Provider',
    report_types.type as 'Report Type',
    report.ID as "Report ID",
    report.date as "Date"
from report
    JOIN report_provider ON report_provider.report_id = report.ID
    JOIN report_type ON report_type.report_id = report.ID
    JOIN report_types ON report_types.ID = report_type.report_type_id
    JOIN report_providers ON report_providers.ID = report_provider.provider_id
order by
	report.date DESC`)

	if err != nil {
		return TableData{}, err
	}

	headers := []any{
		"Provider",
		"Report Type",
		"Report ID",
		"Date",
	}

	report_provider_slice := make([]any, 0)
	report_type_slice := make([]any, 0)
	report_id_slice := make([]any, 0)
	report_date_slice := make([]any, 0)
	line_item_id_slice := make([]any, 0)
	item_class_abbreviation_slice := make([]any, 0)
	item_measurement_slice := make([]any, 0)
	item_units_shorthand_slice := make([]any, 0)
	for rows.Next() {
		var (
			report_provider         string
			report_type             string
			report_id               string
			report_date             string
			line_item_id            string
			item_class_abbreviation string
			item_measurement        string
			item_units_shorthand    string
		)
		rows.Scan(&report_provider, &report_type, &report_id, &report_date, &line_item_id, &item_class_abbreviation, &item_measurement, &item_units_shorthand)

		report_provider_slice = append(report_provider_slice, report_provider)
		report_type_slice = append(report_type_slice, report_type)
		report_id_slice = append(report_id_slice, report_id)
		report_date_slice = append(report_date_slice, report_date)
		line_item_id_slice = append(line_item_id_slice, line_item_id)
		item_class_abbreviation_slice = append(item_class_abbreviation_slice, item_class_abbreviation)
		item_measurement_slice = append(item_measurement_slice, item_measurement)
		item_units_shorthand_slice = append(item_units_shorthand_slice, item_units_shorthand)
	}
	parsed_rows_slice := [][]any{
		report_provider_slice,
		report_type_slice,
		report_id_slice,
		report_date_slice,
		line_item_id_slice,
		item_class_abbreviation_slice,
		item_measurement_slice,
		item_units_shorthand_slice,
	}

	return TableData{
		Headings: headers,
		Rows:     parsed_rows_slice,
	}, nil
}

type ReportProvider struct {
	ReportId   int
	ProviderId int
}

type ReportType struct {
	ReportId     int
	ReportTypeId int
}

type ItemUnits struct {
	Id          int
	FullName    string
	Description string
	Shorthand   string
}

type ItemClasses struct {
	Id           int
	Name         string
	Abbreviation string
}

type LineItem struct {
	Id       int
	ReportId int
}

type LineItemClass struct {
	LineItemId  int
	ItemClassId int
}

type LineItemUnit struct {
	LineItemId int
	ItemUnitId int
}

type LineItemMeasurement struct {
	LineItemId      int
	ItemMeasurement float32
}

// Generic table struct to wrap values
type TableData struct {
	Headings []any
	Rows     [][]any
}

func TransposeStructs[T any](structs []T) [][]any {

	transposed := make([][]any, 0)
	for i := range len(structs) {
		refValue := reflect.ValueOf(structs[i])

		temp_list := make([]any, 0)

		for j := range refValue.NumField() {
			temp_list = append(temp_list, refValue.Field(j).Interface())
		}

		transposed = append(transposed, temp_list)
	}

	return transposed
}

func TransposeStructsToTableData[T any](structs []T) TableData {
	rows := TransposeStructs(structs)

	refType := reflect.ValueOf(structs[0]).Type()

	headers := make([]any, 0)
	for i := range refType.NumField() {
		headers = append(headers, refType.Field(i).Name)
	}

	return TableData{Headings: headers, Rows: rows}
}
