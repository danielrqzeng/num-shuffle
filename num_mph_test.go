package shuffle

import (
	"testing"
)

func TestNumMPH(t *testing.T) {

	const max = 50000
	mph := NumMPHType{}
	mph.Init(max, "test")
	rangeFrom, rangeTo := uint64(0), mph.GetRange()
	t.Log(mph.String(), ",from=", rangeFrom, ",to=", rangeTo)
	has := make(map[uint64]bool)

	for i := rangeFrom; i < rangeTo; i++ {
		plainNum := uint64(i)
		encodeNum, err := mph.Encode(plainNum)
		if err != nil {
			t.Fatalf("encode err=" + err.Error())
		}
		decodeNum, err := mph.Decode(encodeNum)
		if err != nil {
			t.Fatalf("decode err=" + err.Error())

		}
		if plainNum != decodeNum {
			t.Fatalf("plainNum=%d,plain2=%d", plainNum, decodeNum)

		}
		//t.Logf("plainNum=%d,encodeNum=%d\n", plainNum, encodeNum)
		if _, ok := has[encodeNum]; ok {
			t.Fatalf("plainNum=%d,encodeNum=%d\n", plainNum, encodeNum)
		}
	}
}
