package sqla

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

// Filter for a page. This Filter type is like a socket to connect database, backend, and even frontend parts.
//
// ClassFilter allows to filter by a list of sevaral integers.
// ClassFilterOR has the same functionality, however is allows to put OR operator in SQL statement between different ClassFilterOR filters (which have the same name but different columns).
// See descriptions of other filter types for details.
type Filter struct {
	ClassFilter       []ClassFilter
	ClassFilterOR     []ClassFilter
	DateFilter        []DateFilter
	SumFilter         []SumFilter
	TextFilterName    string
	TextFilter        string
	TextFilterColumns []string
}

// ClassFilter to filter types, statuses, etc.
// Selector defines some options list name in user interface, e.g. id of a <select> element.
// InJSON allows to search for an integer value inside JSON stored in a text-type or varchar-type column.
type ClassFilter struct {
	Name     string
	Selector string
	InJSON   bool
	Column   string
	List     []int
}

// DateFilter to filter dates and datetime.
// Dates should be stored as timestamps with value type of int64.
// DatesStr contains strings to represent values in user interface.
type DateFilter struct {
	Name     string
	Column   string
	Relation string
	Dates    []int64
	DatesStr []string
}

// SumFilter to filter currency amounts.
// Sums are stored as integers to avoid loss of accuracy (due to the nature of floats). They all are multiplied by 100. E.g. 0.1 in UI will be searched as 10 in DB and 1 in UI will be searched as 100 in DB.
// SumsStr contains strings to represent values in user interface.
type SumFilter struct {
	Name           string
	Column         string
	CurrencyColumn string
	CurrencyCode   int
	Relation       string
	Sums           []int
	SumsStr        []string
}

// GetFilterFromJSON unmarshals JSON to Filter struct and then only converts dates and sums from strings to integer representation.
// dateConvFunc and dateTimeConvFunc - are used to convert string-typed dates to int64-datestamps or int64-timestamps. These may be the same - it ia a developer's choice.
func (f *Filter) GetFilterFromJSON(JSON []byte,
	dateConvFunc func(string) int64,
	dateTimeConvFunc func(string) int64) {

	err := json.Unmarshal(JSON, f)
	if err != nil {
		log.Println(err)
	}

	dtRegExp := regexp.MustCompile("-[0-9]{1,2}[T ][0-9]{1,2}:[0-9]")
	for i, df := range f.DateFilter {
		df.Dates = nil
		for i := range df.DatesStr {
			if dtRegExp.MatchString(df.DatesStr[i]) {
				df.Dates = append(df.Dates, dateConvFunc(df.DatesStr[i]))
			} else {
				df.Dates = append(df.Dates, dateTimeConvFunc(df.DatesStr[i]))
			}
		}
		f.DateFilter[i].Dates = df.Dates
	}

	for i, sf := range f.SumFilter {
		sf.Sums = nil
		for i := range sf.SumsStr {
			sf.Sums = append(sf.Sums, processJSONSum(sf.SumsStr[i]))
		}
		f.SumFilter[i].Sums = sf.Sums
	}
}

// GetFilterFromForm analyses http.Request and fills the Filter by reqired values (lists, dates, sums, textfilter) from HTML form.
// Any filters with empty lists will be removed from Filter.
// Before executing this method some initial values should be set: filter names and table's columns.
//
// Developer is required to provide dateConvFunc and dateTimeConvFunc. They used to convert string-typed dates from a form to int64-datestamps or int64-timestamps. These may be the same - it ia a developer's choice.
// keywords allow to replace some string from related HTML form with integer value for any ClassFilter.
func (f *Filter) GetFilterFromForm(r *http.Request,
	dateConvFunc func(string) int64,
	dateTimeConvFunc func(string) int64,
	keywords map[string]int) {

	r.ParseForm()

	cfListToReplace := []ClassFilter{}

	for i := range f.ClassFilter {
		classes := r.Form[f.ClassFilter[i].Name]
		curlength := len(classes)
		if curlength > 0 {
			classList := make([]int, curlength, curlength)
			for j := 0; j < len(classes); j++ {
				intval, _ := strconv.Atoi(classes[j])
				if _, ok := keywords[classes[j]]; ok {
					intval = keywords[classes[j]]
				}
				classList[j] = intval
			}
			f.ClassFilter[i].List = classList
			cfListToReplace = append(cfListToReplace, f.ClassFilter[i])
		}
	}
	f.ClassFilter = cfListToReplace

	cfListToReplace = nil
	for i := range f.ClassFilterOR {
		classes := r.Form[f.ClassFilterOR[i].Name]
		curlength := len(classes)
		if curlength > 0 {
			classList := make([]int, curlength, curlength)
			for j := 0; j < len(classes); j++ {
				intval, _ := strconv.Atoi(classes[j])
				if _, ok := keywords[classes[j]]; ok {
					intval = keywords[classes[j]]
				}
				classList[j] = intval
			}
			f.ClassFilterOR[i].List = classList
			cfListToReplace = append(cfListToReplace, f.ClassFilterOR[i])
		}
	}
	f.ClassFilterOR = cfListToReplace

	dfListToReplace := []DateFilter{}

	dtRegExp := regexp.MustCompile("-[0-9]{1,2}[T ][0-9]{1,2}:[0-9]")
	for i := range f.DateFilter {
		datesStr := r.Form[f.DateFilter[i].Name]
		dates := []int64{}
		relation := r.FormValue(f.DateFilter[i].Name + "Relation")
		for i := range datesStr {
			if dtRegExp.MatchString(datesStr[i]) {
				dates = append(dates, dateTimeConvFunc(datesStr[i]))
			} else {
				dates = append(dates, dateConvFunc(datesStr[i]))
			}
		}
		if len(datesStr) == 1 && datesStr[0] != "" {
			f.DateFilter[i] = DateFilter{f.DateFilter[i].Name, f.DateFilter[i].Column, relation, dates, datesStr}
			dfListToReplace = append(dfListToReplace, f.DateFilter[i])
		} else if len(datesStr) == 2 && datesStr[0] != "" && datesStr[1] != "" {
			f.DateFilter[i] = DateFilter{f.DateFilter[i].Name, f.DateFilter[i].Column, relation, dates, datesStr}
			dfListToReplace = append(dfListToReplace, f.DateFilter[i])
		}
	}
	f.DateFilter = dfListToReplace

	sfListToReplace := []SumFilter{}
	for i := range f.SumFilter {
		sumsStr := r.Form[f.SumFilter[i].Name]
		relation := r.FormValue(f.SumFilter[i].Name + "Relation")
		curcode, _ := strconv.Atoi(r.FormValue(f.SumFilter[i].Name + "CurrencyCode"))
		if len(sumsStr) == 1 && sumsStr[0] != "" {
			i0, s0 := processFormSum(sumsStr[0])
			f.SumFilter[i] = SumFilter{
				f.SumFilter[i].Name,
				f.SumFilter[i].Column,
				f.SumFilter[i].CurrencyColumn,
				curcode, relation,
				[]int{i0}, []string{s0}}
			sfListToReplace = append(sfListToReplace, f.SumFilter[i])
		} else if len(sumsStr) == 2 && sumsStr[0] != "" && sumsStr[1] != "" {
			i0, s0 := processFormSum(sumsStr[0])
			i1, s1 := processFormSum(sumsStr[1])
			f.SumFilter[i] = SumFilter{
				f.SumFilter[i].Name,
				f.SumFilter[i].Column,
				f.SumFilter[i].CurrencyColumn,
				curcode, relation,
				[]int{i0, i1}, []string{s0, s1}}
			sfListToReplace = append(sfListToReplace, f.SumFilter[i])
		}
	}
	f.SumFilter = sfListToReplace

	f.TextFilter = r.FormValue(f.TextFilterName)

}

// ClearColumnsValues removes names of (an sql table) columns form a filter if you are going to pass the filter entirely into a response and do not wish to show these names.
func (f *Filter) ClearColumnsValues() {
	for i := range f.ClassFilter {
		f.ClassFilter[i].Column = ""
	}
	for i := range f.ClassFilterOR {
		f.ClassFilterOR[i].Column = ""
	}
	for i := range f.DateFilter {
		f.DateFilter[i].Column = ""
	}
	for i := range f.SumFilter {
		f.SumFilter[i].Column = ""
		f.SumFilter[i].CurrencyColumn = ""
	}
	f.TextFilterColumns = []string{}
}
