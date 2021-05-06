# num-shuffle
将某段范围内的整数打乱
> 其数学表达为，对于`y=encode(x)`,`x=decode(y)`,x,y∈`[min,max)`,且当x1!=x2时候，y1!=y2

## 快速开始
```golang
package main

import (
	"fmt"
	"github.com/RQZeng/num-shuffle"
)

func main() {
	numShuffle := shuffle.ShuffleType{}
	//对[1000,100000)范围内的整数做shuffle
	min, max := uint64(1000), uint64(100000)
	err := numShuffle.Init(min, max, "test")
	if err != nil {
		panic(err)
	}
	
	num := uint64(1001)
	//编码
	encodeNum, err := numShuffle.Encode(num)
	if err != nil {
		panic(err)
	}
	//解码
	decodeNum, err := numShuffle.Decode(encodeNum)
	if err != nil {
		panic(err)
	}
}
```

## 原理
* skip32是一个将32为正整数打乱的算法，本算法不同于skip32,可以控制对某段范围内的整数做打乱
* 对于某范围内的整数空间的打乱，是利用加密算法里面普遍使用的异或来实现
* 具体做法
    * 先指定范围和密钥，初始化
        * 范围整数拆解为二进制幂和的形式
        * 进行二进制幂的整数完美哈希组合
    * 编码
        * 判断要编码整数范围，使用以之对应的二进制幂完美哈希函数进行编码
    * 解码
        * 判断要编码整数范围，使用以之对应的二进制幂完美哈希函数进行解码
> 对于同一个范围的整数，密钥不同的情况下，其编码解码出来的输出都是不一致的，外部难以猜测

## 适用场景
* 对于订单id，有内部外部id之分，内部和外部id可以相互转换，内部id希望可以递增以方便业务控制，外部id不希望有可被猜测的规则
                
## 一些限制
* 编码范围为uint64范围，即[0,0xFFFFFFFFFFFFFFFF),注意，其不包含0xFFFFFFFFFFFFFFFF这个数字的
* 由于整数需要分解为二进制幂和，是以其编码范围会有一个二进制幂范围层，举例
    * 对于[0,35),其最大值为1000=0b10_0011
    * 将会有3个二进制幂的完美哈希实例对应，分别为[0,1),[1,3),[4,35)
    * 0∈[0,1),是以其必定被映射到0
    * 2∈[1,3),是以其也会被映射到[1,3)范围内
  