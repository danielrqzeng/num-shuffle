package shuffle

import "testing"

func TestRevBit(t *testing.T) {
	for _, data := range [][]uint64{
		{0x3e, 6, 31},
		{32, 7, 2},
		{0x6e, 7, 59},
		{0xef6e, 7, 0xef3b},
	} {
		num := data[0]
		bitNum := data[1]
		expect := data[2]
		revNum := RevBit(num, bitNum)
		t.Log("num=", num, "revNum=", revNum)
		if revNum != expect {
			t.Fatalf("RevBit(%v) expected '%v', but got '%v'", num, expect, revNum)
		}
	}
}

func TestFirstNoneZeroBit(t *testing.T) {
	for _, data := range [][]uint64{
		{0x0, 0},
		{0x1, 1},
		{0x3e, 6},
		{32, 6},
		{0x6e, 7},
		{0xef6e, 16},
		{0xFFFFFFFFFFFFFFFF, 64},
	} {
		num := data[0]
		expect := data[1]
		got := FirstNoneZeroBit(num)
		t.Log("num=", num, "got=", got)
		if got != expect {
			t.Fatalf("FirstNoneZeroBit(%v) expected '%v' for %b, but got '%v'", num, expect, num, got)
		}
	}
}
