package shuffle

import (
	"testing"
)

func TestShuffle0(t *testing.T) {
	numShuffle := ShuffleType{}
	err := numShuffle.Init(0, 0, "test")
	if err == nil {
		t.Error("unvalid param")
	}
}

func TestShuffle1(t *testing.T) {
	numShuffle := ShuffleType{}
	err := numShuffle.Init(0, 1, "test")
	if err != nil {
		t.Error("unvalid param", err)
	}
	for _, data := range [][]uint64{
		{0, 0},
	} {
		plainNum := data[0]
		expect := data[1]
		encodeNum, err := numShuffle.Encode(plainNum)
		if err != nil {
			t.Fatalf("encode err=" + err.Error())
		}
		if encodeNum != expect {
			t.Fatalf("plainNum=%d,expect=%d,but got=%d", plainNum, expect, encodeNum)
		}
		decodeNum, err := numShuffle.Decode(encodeNum)
		if err != nil {
			t.Fatalf("decode err=" + err.Error())

		}
		if plainNum != decodeNum {
			t.Fatalf("plainNum=%d,expect=%d,but got=%d", plainNum, decodeNum, decodeNum)
		}
		t.Logf("encode %d to %d", plainNum, encodeNum)
	}
}

func TestShuffle2(t *testing.T) {
	numShuffle := ShuffleType{}
	err := numShuffle.Init(0, 0xFFFFFFFFFFFFFFFF, "test")
	if err != nil {
		t.Error("unvalid param", err)
	}
	for _, data := range [][]uint64{
		{0, 0}, //编码-最小值
		{0xFFFFFFFFFFFFFFFE, 12585698967662052142}, //编码-最大值
		{0, 0}, //解码-最大值
		{18041841561391811257, 0xFFFFFFFFFFFFFFFE}, //解码-最小值
	} {
		plainNum := data[0]
		expect := data[1]
		encodeNum, err := numShuffle.Encode(plainNum)
		if err != nil {
			t.Fatalf("encode err=" + err.Error())
		}
		if encodeNum != expect {
			t.Fatalf("plainNum=%d,expect=%d,but got=%d", plainNum, expect, encodeNum)
		}
		decodeNum, err := numShuffle.Decode(encodeNum)
		if err != nil {
			t.Fatalf("decode err=" + err.Error())

		}
		if plainNum != decodeNum {
			t.Fatalf("plainNum=%d,expect=%d,but got=%d", plainNum, decodeNum, decodeNum)
		}
		t.Logf("encode %d to %d", plainNum, encodeNum)

	}
}

func TestShuffle3(t *testing.T) {
	numShuffle := ShuffleType{}
	err := numShuffle.Init(100, 10000, "test")
	if err != nil {
		t.Error(err)
	}
	for _, data := range [][]uint64{
		{100, 103},
		{132, 140},
		{500, 705},
		{10000 - 1, 5120},
	} {
		plainNum := data[0]
		expect := data[1]
		encodeNum, err := numShuffle.Encode(plainNum)
		if err != nil {
			t.Fatalf("encode err=" + err.Error())
		}
		if encodeNum != expect {
			t.Fatalf("plainNum=%d,expect=%d,but got=%d", plainNum, expect, encodeNum)
		}
		decodeNum, err := numShuffle.Decode(encodeNum)
		if err != nil {
			t.Fatalf("decode err=" + err.Error())

		}
		if plainNum != decodeNum {
			t.Fatalf("plainNum=%d,expect=%d,but got=%d", plainNum, decodeNum, decodeNum)
		}
		t.Logf("encode %d to %d", plainNum, encodeNum)
	}
}

func TestShuffle4(t *testing.T) {
	numShuffle := ShuffleType{}
	err := numShuffle.Init(100, 10000, "test")
	if err != nil {
		t.Error(err)
	}
	for _, data := range [][]uint64{
		{0, 0},
		{10000, 0},
		{10001, 0},
	} {
		plainNum := data[0]
		expect := data[1]
		encodeNum, err := numShuffle.Encode(plainNum)
		if err == nil {
			t.Fatalf("encode err=" + err.Error())
		}
		_, _ = expect, encodeNum
		t.Logf("encode %d err=%s", plainNum, err.Error())
	}
}

func TestShuffle5(t *testing.T) {
	numShuffle := ShuffleType{}
	err := numShuffle.Init(0, 35, "test")
	if err != nil {
		t.Error(err)
	}
	for _, data := range [][]uint64{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 28},
		{34, 9},
	} {
		plainNum := data[0]
		expect := data[1]
		encodeNum, err := numShuffle.Encode(plainNum)
		if err != nil {
			t.Fatalf("encode err=" + err.Error())
		}
		if encodeNum != expect {
			t.Fatalf("plainNum=%d,expect=%d,but got=%d", plainNum, expect, encodeNum)
		}
		decodeNum, err := numShuffle.Decode(encodeNum)
		if err != nil {
			t.Fatalf("decode err=" + err.Error())

		}
		if plainNum != decodeNum {
			t.Fatalf("plainNum=%d,expect=%d,but got=%d", plainNum, decodeNum, decodeNum)
		}
		t.Logf("encode %d to %d", plainNum, encodeNum)
	}
}
