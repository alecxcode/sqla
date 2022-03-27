package sqla

import (
	"testing"
)

func TestToDecimalStr(t *testing.T) {
	var res string

	res = toDecimalStr("-1")
	if res != "-0.01" {
		t.Errorf("Expected:%s, received:%s", "-0.01", res)
	}

	res = toDecimalStr("-10")
	if res != "-0.10" {
		t.Errorf("Expected:%s, received:%s", "-0.10", res)
	}
}

func TestProcessFormSum(t *testing.T) {
	var res int

	res, _ = processFormSum("-0.001")
	if res != 0 {
		t.Errorf("Expected:%d, received:%d", 0, res)
	}

	res, _ = processFormSum("-0.01")
	if res != -1 {
		t.Errorf("Expected:%d, received:%d", -1, res)
	}

	res, _ = processFormSum("-0.1")
	if res != -10 {
		t.Errorf("Expected:%d, received:%d", -10, res)
	}

	res, _ = processFormSum("-1.0")
	if res != -100 {
		t.Errorf("Expected:%d, received:%d", -100, res)
	}

	res, _ = processFormSum("sometext")
	if res != 0 {
		t.Errorf("Expected:%d, received:%d", 0, res)
	}

	res, _ = processFormSum("0.9.9")
	if res != 990 {
		t.Errorf("Expected:%d, received:%d", 990, res)
	}
}
