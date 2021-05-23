package shuffle

import "fmt"

type ShuffleType struct {
	min uint64
	max uint64

	mphs      []*NumMPHType //idx=mph
	mphsRange [][]uint64    //idx=[min,max)
	secretKey string
}

//Init 初始化，范围为[min,max)
func (obj *ShuffleType) Init(min, max uint64, secretKey string) (err error) {
	if max <= min {
		err = fmt.Errorf("unvalid param cuz max<min")
		return
	}
	rangeMax := max - min
	obj.min = min
	obj.max = max
	obj.secretKey = secretKey

	baseMin := obj.min
	pos := uint64(0)
	//rangeLeft := rangeMax
	//for i := uint64(0); i < 64; i++ {
	for i := int64(63); i >= 0; i-- {
		fmt.Println(i)
		maxVal := uint64(1) << i
		if rangeMax&maxVal == 0 {
			continue
		}
		mph := &NumMPHType{}
		mph.Init(maxVal, fmt.Sprintf("%s-%d", obj.secretKey, maxVal))
		obj.mphs = append(obj.mphs, mph)

		//计算mph的范围
		base := baseMin + pos
		incr := mph.GetRange() + 1
		obj.mphsRange = append(obj.mphsRange, []uint64{base, base + incr})
		pos += incr
	}
	//fmt.Printf("min=%d(%b) max=%d(%b),max-min=%b\n", min, min, max, max, max-min)
	//for idx := 0; idx < len(obj.mphs); idx++ {
	//	fmt.Printf("range=[%d,%d)\n", obj.mphsRange[idx][0]-min, obj.mphsRange[idx][1]-min)
	//}
	return

}

//Encode 编码，num范围必须在Init中[min,max)之内
func (obj *ShuffleType) Encode(num uint64) (cipherNum uint64, err error) {
	for idx := 0; idx < len(obj.mphs); idx++ {
		min, max := obj.mphsRange[idx][0], obj.mphsRange[idx][1]
		//fmt.Printf("idx=%d,min=%d,max=%d\n", idx, min, max)
		if num >= min && num < max {
			mph := obj.mphs[idx]
			cipherNum, err = mph.Encode(num - min)
			if err != nil {
				return
			}
			cipherNum += min
			return
		}
	}
	err = fmt.Errorf("encode not found mph for num=%d", num)
	return
}

//Encode 解码，num范围必须在Init中[min,max)之内
func (obj *ShuffleType) Decode(num uint64) (plainNum uint64, err error) {
	for idx := 0; idx < len(obj.mphs); idx++ {
		min, max := obj.mphsRange[idx][0], obj.mphsRange[idx][1]
		if num >= min && num < max {
			mph := obj.mphs[idx]
			plainNum, err = mph.Decode(num - min)
			if err != nil {
				return
			}
			plainNum += min
			return
		}
	}
	err = fmt.Errorf("decode not found mph for num=%d", num)
	return
}
