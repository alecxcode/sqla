package sqla

import (
	"encoding/json"
	"log"
)

// anyT struct is not exported, its realization is hidden.
type anyT struct {
	c string // column name
	t int    // defines kind of type: possible values i, b, f, s, nil
	i int64
	b bool
	f float64
	s string
}

// AnyTslice can contain different values (strings, including JSON, integers, nils), although the value type should be specified when adding.
// This slice later used with insert or update to database.
// Appending to this slice is made with defined below functions. You should provide column name of SQL table to most of them.
type AnyTslice []anyT

// AppendInt appends int to AnyTslice
func (a AnyTslice) AppendInt(column string, i int) AnyTslice {
	const I = 0
	a = append(a, anyT{c: column, t: I, i: int64(i)})
	return a
}

// AppendInt64 appends int64 to AnyTslice
func (a AnyTslice) AppendInt64(column string, i int64) AnyTslice {
	const I = 0
	a = append(a, anyT{c: column, t: I, i: i})
	return a
}

// AppendNil appends nil to AnyTslice
func (a AnyTslice) AppendNil(column string) AnyTslice {
	const N = 4
	a = append(a, anyT{c: column, t: N})
	return a
}

// AppendNonEmptyString appends string if it is not empty. If string is empty nothing will be appended (slice unchanged).
func (a AnyTslice) AppendNonEmptyString(column string, s string) AnyTslice {
	const S = 3
	if s != "" {
		a = append(a, anyT{c: column, t: S, s: s})
	}
	return a
}

// AppendStringOrNil appends string if it is not empty. If string is empty nil will be appended.
func (a AnyTslice) AppendStringOrNil(column string, s string) AnyTslice {
	const S = 3
	const N = 4
	if s != "" {
		a = append(a, anyT{c: column, t: S, s: s})
	} else {
		a = append(a, anyT{c: column, t: N})
	}
	return a
}

// AppendJSONList appends JSON made from []string slice if it is not empty. If slice is empty nil will be appended.
func (a AnyTslice) AppendJSONList(column string, sList []string) AnyTslice {
	const S = 3
	const N = 4
	if len(sList) > 0 {
		jsonList, err := json.Marshal(sList)
		if err != nil {
			log.Println(currentFunction()+":", err)
		}
		a = append(a, anyT{c: column, t: S, s: string(jsonList)})
	} else {
		a = append(a, anyT{c: column, t: N})
	}
	return a
}

// AppendJSONListInt appends JSON made from []int slice if it is not empty. If slice is empty nil will be appended.
func (a AnyTslice) AppendJSONListInt(column string, iList []int) AnyTslice {
	const S = 3
	const N = 4
	if len(iList) > 0 {
		jsonList, err := json.Marshal(iList)
		if err != nil {
			log.Println(currentFunction()+":", err)
		}
		a = append(a, anyT{c: column, t: S, s: string(jsonList)})
	} else {
		a = append(a, anyT{c: column, t: N})
	}
	return a
}

// AppendJSONStruct appends JSON made from any struct. It may not be nil pointer.
func (a AnyTslice) AppendJSONStruct(column string, sStruct interface{}) AnyTslice {
	const S = 3
	jsonList, err := json.Marshal(sStruct)
	if err != nil {
		log.Println(currentFunction()+":", err)
	}
	a = append(a, anyT{c: column, t: S, s: string(jsonList)})
	return a
}

// UnmarshalNonEmptyJSONList returns []string slice made before from this kind of slice.
func UnmarshalNonEmptyJSONList(s string) (jsonList []string) {
	if s != "" {
		err := json.Unmarshal([]byte(s), &jsonList)
		if err != nil {
			log.Println(currentFunction()+":", err)
		}
	}
	return jsonList
}

// UnmarshalNonEmptyJSONListInt returns []int slice made before from this kind of slice.
func UnmarshalNonEmptyJSONListInt(s string) (jsonList []int) {
	if s != "" {
		err := json.Unmarshal([]byte(s), &jsonList)
		if err != nil {
			log.Println(currentFunction()+":", err)
		}
	}
	return jsonList
}
