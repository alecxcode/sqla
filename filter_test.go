package sqla

import (
	"strconv"
	"strings"
	"testing"
)

func dateTimeToInt64(dtstring string) int64 {
	BC := false
	if dtstring == "" {
		return 0
	}
	if dtstring[0] == '-' {
		dtstring = strings.TrimLeft(dtstring, "-")
		BC = true
	}
	dtstring = strings.Replace(dtstring, "T", " ", 1)
	dtarr := strings.SplitN(dtstring, " ", 2)
	if len(dtarr) < 2 {
		return 0
	}
	darr := strings.Split(dtarr[0], "-")
	if len(darr) < 3 {
		return 0
	}
	y, _ := strconv.Atoi(darr[0])
	m, _ := strconv.Atoi(darr[1])
	d, _ := strconv.Atoi(darr[2])
	if BC {
		y = -y
	}
	tarr := strings.Split(dtarr[1], ":")
	if len(tarr) < 2 {
		return 0
	}
	hour, _ := strconv.Atoi(tarr[0])
	minute, _ := strconv.Atoi(tarr[1])
	var AD int64 = 1
	if y < 0 {
		y = -y
		AD = -1
	}
	return (int64(minute) + int64(hour)*60 + int64(d-1)*60*24 + int64(m-1)*60*24*31 + int64(y)*60*24*31*12) * AD
}

func TestGetFilterFromJSON(t *testing.T) {
	JSON := `{
"ClassFilter":[{"Name":"doctypes","Selector":"","InListed":false,"Column":"DocTypes","List":[0,1,2]}],
"ClassFilterOR":[{"Name":"creatorsORassignees","Selector":"userSelector","InListed":false,"Column":"Creators","List":[7,8,9]},{"Name":"creatorsORassignees","Selector":"userSelector","InListed":false,"Column":"Assignees","List":[7,8,9]}],
"DateFilter":[{"Name":"createdDates","Column":"Created","Relation":"=","Dates":[],"DatesStr":["2022-02-01T08:47"]}],
"SumFilter":[{"Name":"sums","Column":"Sum","CurrencyColumn":"Currency","CurrencyCode":840,"Relation":"","Sums":[],"SumsStr":["0.00","1000.00"]}],
"TextFilter":"testphrase","TextFilterColumns":["About","Note"]
}
`
	f := Filter{}
	f.GetFilterFromJSON([]byte(JSON), dateTimeToInt64, dateTimeToInt64)

	t.Logf("Filter: %#v\n", f)
	if f.DateFilter[0].Dates[0] != 1083190127 {
		t.Errorf("Expected:%d, received:%d", 1083190127, f.DateFilter[0].Dates[0])
	}
	if f.SumFilter[0].Sums[1] != 100000 {
		t.Errorf("Expected:%d, received:%d", 100000, f.SumFilter[0].Sums[1])
	}
}
