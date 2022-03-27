package sqla

import (
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

func currentFunction() string {
	counter, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(counter).Name()
}

func intSlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func isStringASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func firstLetterIndex(r []rune) int {
	ri := 0
	for i := 0; i < len(r); i++ {
		if unicode.IsLetter(r[i]) {
			ri = i
			break
		}
	}
	return ri
}

func processJSONSum(s string) (sum int) {
	if !strings.Contains(s, ".") {
		s += ".00"
	}
	if s[len(s)-2] == '.' {
		s += "0"
	}
	if strings.Contains(s, ".") && s[len(s)-3] != '.' {
		tsa := strings.Split(s, ".")
		if len(tsa[1]) == 1 {
			tsa[1] = tsa[1] + "0"
		}
		s = tsa[0] + tsa[1][:2]
	}
	sum, _ = strconv.Atoi(strings.Replace(s, ".", "", -1))
	return sum
}

func processFormSum(s string) (sum int, decimal string) {
	sum = processFormSumInt(s)
	decimal = toDecimalStr(strconv.Itoa(sum))
	return sum, decimal
}

func processFormSumInt(s string) (sum int) {
	if !strings.Contains(s, ".") {
		s += ".00"
	} else if s[len(s)-2] == '.' {
		s += "0"
	}
	if strings.Contains(s, ".") && s[len(s)-3] != '.' {
		tsa := strings.Split(s, ".")
		if len(tsa[1]) == 1 {
			tsa[1] = tsa[1] + "0"
		}
		s = tsa[0] + tsa[1][:2]
	}
	sum, _ = strconv.Atoi(strings.Replace(s, ".", "", -1))
	return sum
}

func toDecimalStr(s string) string {
	var addminus = false
	if strings.HasPrefix(s, "-") {
		s = strings.TrimPrefix(s, "-")
		addminus = true
	}
	if len(s) == 1 {
		s = "0.0" + s
	} else if len(s) == 2 {
		s = "0." + s
	} else {
		i := len(s) - 2
		s = s[:i] + "." + s[i:]
	}
	if addminus {
		s = "-" + s
	}
	return s
}
