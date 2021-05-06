package shuffle

import "math/bits"

/*
RevBit 反转二进制，val为想要反转的值，num为反转val几位
举例：
	* val=0b1001,num=2  =>  revVal=0b1010 (从右数起，反转两位
	* val=0b1011,num=3  =>  revVal=0b1110 (从右数起，反转三位
	* val=0b1100_1011,num=5  =>  revVal=0b1101_1010
*/
func RevBit(val, num uint64) (revVal uint64) {
	var bitArr []int

	nn := val
	//记录[0,bitNum)的二进制序列
	for i := uint64(0); i < 64; i++ {
		if nn&(1<<i) != 0 {
			bitArr = append(bitArr, 1)
		} else {
			bitArr = append(bitArr, 0)
		}
	}

	revArr := bitArr[:num]

	for i, j := 0, len(revArr)-1; i < j; i, j = i+1, j-1 {
		revArr[i], revArr[j] = revArr[j], revArr[i]
	}

	revVal = 0
	for i := uint64(0); i < 64; i++ {
		bitVal := bitArr[i]
		if i < num {
			bitVal = revArr[i]
		}
		if bitVal == 1 {
			revVal = revVal | (1 << i)
		}
	}
	//fmt.Printf("val=%d(%x),bitVal=%b,num=%d,result=%d(%x)\n", val, val, val, num, revVal, revVal)
	return
}

/*
FirstNoneZeroBit 计算n的第一个二进制的1的下标值（小端序，从右往左,首位下标为1）
举例:
	* n=0,其为0，是以没有下标为1 ==> ret=0
	* n=0b0010_1010,其第一个1的下标值为6 ==> ret=6
	* n=0b0101,其第一个1的下标值为3 ==> ret=3
	* n=0xFFFF_FFFF_FFFF_FFFF,其第一个1的下标值为64 ==> ret=64
*/
func FirstNoneZeroBit(n uint64) uint64 {
	//revN := bits.Reverse64(n)
	//zeroNumToLeft := bits.TrailingZeros64(revN)
	//return 64 - uint64(zeroNumToLeft)
	return 64 - uint64(bits.LeadingZeros64(n))
}
