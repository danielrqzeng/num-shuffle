package shuffle

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

//NumMPHType 完美哈希一段范围的整数（即将这段范围内的整数，做一个shuffle，由于数据空间太大，是以不能用类似52张扑克牌的方式做shuffle）
type NumMPHType struct {
	max       uint64 //使用者想要的最大值
	realMax   uint64 //程序能达到的最大值，取值范围为[min，2^(first_1_bit_idx_from_left-1))
	secretKey string //密钥

	first1bitIdx uint64   //max的二进制值，最左边的1的位置（位置为从1开始，从右往左计算）
	mask         uint64   //最大范围的二进制掩码
	secret       []uint64 //密钥队列，通过用户传过来的密钥生成，其是一个队列，将会在encode时候次第做异或
}

//String ...
func (obj *NumMPHType) String() string {
	str := ""
	str += fmt.Sprintf("max=%d,realMax=%d,mask=%x", obj.max, obj.realMax, obj.mask)
	return str
}

//Init 初始化
//取值范围为[min,max),其为用户想 要的取值范围，但是这里要打个折扣，取值范围将会是[min，2^(first_1_bit_idx_from_left-1))
func (obj *NumMPHType) Init(max uint64, secretKey string) {
	//校验参数
	if max == 0 || secretKey == "" {
		panic(fmt.Errorf("param is unvalid"))
	}

	obj.max = max
	obj.secretKey = secretKey

	//计算max的第一个二进制的下标值
	//n := leftBitIdx(obj.max)
	n := FirstNoneZeroBit(obj.max)
	obj.first1bitIdx = n //最大非零下标，另最大次位下标=最大非零下标-1，即obj.first1bitIdx-1

	//NOTE: 1<<0=1
	for i := uint64(0); i < obj.first1bitIdx-1; i++ {
		obj.mask = obj.mask | (1 << i)
	}
	obj.realMax = obj.mask

	//计算做异或运算的密钥
	sha1Key := sha1.Sum([]byte(secretKey))
	ms5Key := md5.Sum([]byte(secretKey))
	sha256Key := sha256.Sum256([]byte(secretKey))
	obj.secret = []uint64{
		obj.mask & binary.LittleEndian.Uint64(sha1Key[:]),
		obj.mask & binary.LittleEndian.Uint64(ms5Key[:]),
		obj.mask & binary.LittleEndian.Uint64(sha256Key[:]),
	}
}

//GetRange 获取能完美hash的最大值
//取值范围为[min,max),其为用户想要的取值范围，但是这里要打个折扣，取值范围将会是[min，2^(first_1_bit_idx_from_left-1))
func (obj *NumMPHType) GetRange() (max uint64) {
	max = obj.realMax
	return
}

//Encode 整数加密，输入明文整数，输出密文整数
func (obj *NumMPHType) Encode(plainNum uint64) (cipherNum uint64, err error) {
	if plainNum > obj.realMax {
		err = fmt.Errorf("num=%d not in [0,%d)", plainNum, obj.realMax)
		return
	}

	//在范围值内，坐反转
	num := RevBit(plainNum, obj.first1bitIdx-1)
	cipherNum = num

	for _, s := range obj.secret {
		cipherNum = cipherNum ^ s
		//fmt.Printf("idx=%d,s=%d,plain=%d,cipher=%d\n",idx,s,plainNum,cipherNum)
	}

	return
}

//Encode 整数解密，输入密文整数，输出明文整数
func (obj *NumMPHType) Decode(cipherNum uint64) (plainNum uint64, err error) {
	if plainNum > obj.realMax {
		err = fmt.Errorf("num=%d not in [0,%d)", cipherNum, obj.realMax)
		return
	}

	plainNum = cipherNum
	for i := len(obj.secret) - 1; i >= 0; i-- {
		s := obj.secret[i]
		plainNum = plainNum ^ s
	}

	plainNum = RevBit(plainNum, obj.first1bitIdx-1)
	return
}

/*
doc
mph:完美哈希映射
此类是为了在一个整数某段空间内，做一个整数映射
我们熟悉的对52张牌做一个shuffle，乱序后，可以得到一个下标对应一张卡，此处不能用shuffle，因为需要映射的空间太大(比如500亿)
数学表达如下：
对y=encode(x)，x=decode(y)，其中encode为加密函数，decode为解密函数，满足，
1.当x1!=x2时候，f(x1)!=f(x2)
2.x和y同属某个范围，比如[min,max)

举例：在某个空间[0,50000)
y=f(1000)=2001
x=f(2001)=1000


解决方案
* 用户决定其想要映射的范围，比如[0,50000)
* 调用函数初始化Init
	* 将50000计算为二进制即为1100_0011_0101_0000
	* 最大次位下标：一个整数二进制的首位为一的下标称为最大非零下标（从右到左，首位下标为1），其再减一即为最大次位下标
		> 比如1100_0011_0101_0000的最大非零下标为16，其最大次位下标为15
	* 计算最大此位和最大映射范围
	* 最大次位下标为15，是以其能映射的范围为0~7fff，这样能保证其无论怎么操作15位，均不会超过16位，限制了其范围
		> @TODO 此处注意，和用户想要的范围[0,50000)不一样，范围小了，只到[0,0x7FFF)，是其子集
	* 计算密钥，此处使用了多个摘要算法，避免被猜出来
	* 对密钥进行掩码，只允许其最大的值为上步骤中的掩码范围
* 调用函数Encode
	* 先做反转，但是只做最大次位下标内的二进制反转，此步骤是为了混淆规避数字相邻带来的规律
	* 使用密钥数组顺序做异或（因为异或两次可得同一个数字）
	> 加密的数字，必须在映射范围之内
* 调用函数Decode
	* 使用密钥数组倒序做异或（因为异或两次可得同一个数字）
	* 最大次位下标内的二进制反转
	> 加密的数字，必须在映射范围之内

*/
