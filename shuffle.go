package shuffle

import (
	"fmt"
	"github.com/spaolacci/murmur3"
	"math/rand"
)

func MurMurHash64(str string) (val uint64) {
	m := murmur3.New64()
	_, _ = m.Write([]byte(str))
	val = m.Sum64()
	return
}

func shuffle(seed uint64, nums []uint64) []uint64 {
	r := rand.New(rand.NewSource(int64(seed)))
	for i := len(nums); i > 0; i-- {
		last := i - 1
		idx := r.Intn(i)
		nums[last], nums[idx] = nums[idx], nums[last]
	}
	return nums
}

type ShuffleType struct {
	min uint64
	max uint64

	mphs      []*NumMPHType //idx=mph
	mphsRange [][]uint64    //idx=[min,max)

	ascMphs      []*NumMPHType //idx=mph,二进制从小到大
	ascMphsRange [][]uint64    //idx=[min,max)

	descMphs      []*NumMPHType //idx=mph，二进制从大到小
	descMphsRange [][]uint64    //idx=[min,max)

	mphList   [][]*NumMPHType
	mphRanges [][][]uint64
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
		maxVal := uint64(1) << uint64(i)
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
	/*
		//二进制从大到小
		baseMin = obj.min
		pos = uint64(0)
		for i := int64(63); i >= 0; i-- {
			maxVal := uint64(1) << uint64(i)
			if rangeMax&maxVal == 0 {
				continue
			}
			mph := &NumMPHType{}
			mph.Init(maxVal, fmt.Sprintf("%s-%d", obj.secretKey, maxVal))
			obj.descMphs = append(obj.descMphs, mph)

			//计算mph的范围
			base := baseMin + pos
			incr := mph.GetRange() + 1
			obj.descMphsRange = append(obj.descMphsRange, []uint64{base, base + incr})
			pos += incr
		}

		//二进制从小到大
		baseMin = obj.min
		pos = uint64(0)
		for i := int64(0); i < 64; i++ {
			maxVal := uint64(1) << uint64(i)
			if rangeMax&maxVal == 0 {
				continue
			}
			mph := &NumMPHType{}
			mph.Init(maxVal, fmt.Sprintf("%s-%d", obj.secretKey, maxVal))
			obj.ascMphs = append(obj.ascMphs, mph)

			//计算mph的范围
			base := baseMin + pos
			incr := mph.GetRange() + 1
			obj.ascMphsRange = append(obj.ascMphsRange, []uint64{base, base + incr})
			pos += incr
		}

		const listLen = 2
		obj.mphList = make([][]*NumMPHType, listLen)
		obj.mphRanges = make([][][]uint64, listLen)
		uint64Shuffle := make([][]uint64, listLen)
		for i := 0; i < listLen; i++ {
			sf := make([]uint64, 0)
			for j := uint64(0); j < 64; j++ {
				maxVal := uint64(1) << uint64(j)
				if rangeMax&maxVal == 0 {
					continue
				}
				sf = append(sf, j)
			}
			seed := MurMurHash64(fmt.Sprintf("%s-%d1", obj.secretKey, i))
			uint64Shuffle[i] = shuffle(seed, sf[:])
			fmt.Println("idx=", i, ",seed=", seed)
			fmt.Println("uint64Shuffle[i] =", uint64Shuffle[i])
		}
		for idx := 0; idx < listLen; idx++ {
			baseMin = obj.min
			pos = uint64(0)
			for _, i := range uint64Shuffle[idx] {
				maxVal := uint64(1) << uint64(i)
				if rangeMax&maxVal == 0 {
					continue
				}
				mph := &NumMPHType{}
				mph.Init(maxVal, fmt.Sprintf("%s-%d", obj.secretKey, maxVal))
				obj.mphList[idx] = append(obj.mphList[idx], mph)

				//计算mph的范围
				base := baseMin + pos
				incr := mph.GetRange() + 1
				obj.mphRanges[idx] = append(obj.mphRanges[idx], []uint64{base, base + incr})
				pos += incr
				fmt.Println("idx=", idx, ",bit=", i, ",range=[", base, ",", base+incr, ")")
			}
		}*/
	return
}

//Encode 编码，num范围必须在Init中[min,max)之内
func (obj *ShuffleType) encode(num uint64, mphs []*NumMPHType, mphsRange [][]uint64) (cipherNum uint64, err error) {
	for idx := 0; idx < len(mphs); idx++ {
		min, max := mphsRange[idx][0], mphsRange[idx][1]
		//fmt.Printf("encode idx=%d,min=%d,max=%d\n", idx, min, max)
		if num >= min && num < max {
			mph := mphs[idx]
			cipherNum, err = mph.Encode(num - min)
			if err != nil {
				return
			}
			cipherNum += min
			//fmt.Printf("encode idx=%d,min=%d,max=%d from %d to %d\n", idx, min, max, num, cipherNum)

			return
		}
	}
	err = fmt.Errorf("encode not found mph for num=%d", num)
	return
}

//Encode 解码，num范围必须在Init中[min,max)之内
func (obj *ShuffleType) decode(num uint64, mphs []*NumMPHType, mphsRange [][]uint64) (plainNum uint64, err error) {
	for idx := 0; idx < len(mphs); idx++ {
		min, max := mphsRange[idx][0], mphsRange[idx][1]
		if num >= min && num < max {
			mph := mphs[idx]
			plainNum, err = mph.Decode(num - min)
			if err != nil {
				return
			}
			plainNum += min
			//fmt.Printf("decode idx=%d,min=%d,max=%d from %d to %d\n", idx, min, max, num, plainNum)
			return
		}
	}
	err = fmt.Errorf("decode not found mph for num=%d", num)
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

func (obj *ShuffleType) ExtendEncode(num uint64) (cipherNum uint64, err error) {
	//c1, err := obj.encode(num, obj.descMphs, obj.descMphsRange)
	//if err != nil {
	//	return
	//}
	//c2, err := obj.decode(c1, obj.ascMphs, obj.ascMphsRange)
	//if err != nil {
	//	return
	//}
	//cipherNum = c2
	//fmt.Println("DoubleEncode ", num, "->", c1, "->", c2, "->", cipherNum)
	//fmt.Println(obj.encode(125, obj.mphList[2], obj.mphRanges[2]))
	//fmt.Println(obj.decode(132, obj.mphList[2], obj.mphRanges[2]))

	c := make([]uint64, 0)
	cipherNum = num
	c = append(c, cipherNum)
	for i := 0; i < len(obj.mphList); i++ {
		cipherNum, err = obj.encode(cipherNum, obj.mphList[i], obj.mphRanges[i])
		//if i%2 == 0 {
		//	cipherNum, err = obj.decode(cipherNum, obj.mphList[i], obj.mphRanges[i])
		//} else {
		//	cipherNum, err = obj.encode(cipherNum, obj.mphList[i], obj.mphRanges[i])
		//}
		if err != nil {
			return
		}
		c = append(c, cipherNum)
	}
	//fmt.Println("DoubleEncode ", num, ",chain=", c)

	return
}

func (obj *ShuffleType) ExtendDecode(num uint64) (plainNum uint64, err error) {
	//c1, err := obj.encode(num, obj.ascMphs, obj.ascMphsRange)
	//if err != nil {
	//	return
	//}
	//c2, err := obj.decode(c1, obj.descMphs, obj.descMphsRange)
	//if err != nil {
	//	return
	//}
	//plainNum = c2
	//fmt.Println("DoubleDecode ", num, "->", c1, "->", c2, "->", plainNum)

	c := make([]uint64, 0)
	plainNum = num
	c = append(c, plainNum)
	for i := len(obj.mphList) - 1; i >= 0; i-- {
		plainNum, err = obj.decode(plainNum, obj.mphList[i], obj.mphRanges[i])
		//if i%2 == 0 {
		//	plainNum, err = obj.encode(plainNum, obj.mphList[i], obj.mphRanges[i])
		//} else {
		//	plainNum, err = obj.decode(plainNum, obj.mphList[i], obj.mphRanges[i])
		//}
		if err != nil {
			return
		}
		c = append(c, plainNum)
	}
	fmt.Println("DoubleDecode ", num, ",chain=", c)
	return
}
