package main

import (
	"fmt"
	"github.com/RQZeng/num-shuffle"
)

func main() {
	//ExtendEncode()
	//fmt.Println("------Encode----------")

	numShuffle := shuffle.ShuffleType{}
	//对[1000,100000)范围内的整数做shuffle
	min, max := uint64(0), uint64(8)
	//min, max = uint64(0), uint64(35)
	err := numShuffle.Init(min, max, "test")
	if err != nil {
		panic(err)
	}
	uniq := make(map[uint64]bool)
	for i := uint64(min); i < max; i++ {
		encodeNum, err := numShuffle.Encode(i)
		if err != nil {
			panic(err)
		}
		decodeNum, err := numShuffle.Decode(encodeNum)
		if err != nil {
			panic(err)
		}
		if decodeNum > max {
			panic(fmt.Sprintf("mph i=%d to %d err,cuz range", i, encodeNum))
		}

		if decodeNum != i {
			panic(fmt.Sprintf("mph i=%d to %d err", i, decodeNum))
		}
		if _, ok := uniq[encodeNum]; ok {
			panic(fmt.Sprintf("mph i=%d to %d err,cuz decode exist", i, encodeNum))
		}
		uniq[encodeNum] = true
		fmt.Printf("mph i=%d to %d\n", i, encodeNum)
	}
	if uint64(len(uniq)) != max-min {
		panic(fmt.Sprintf("mph err,cuz got=%d,expect=%d", uint64(len(uniq)), max))
	}
}

func ExtendEncode() {
	fmt.Println("------ExtendEncode----------")
	numShuffle := shuffle.ShuffleType{}
	//对[1000,100000)范围内的整数做shuffle
	min, max := uint64(100), uint64(200)
	//min, max = uint64(0), uint64(35)
	err := numShuffle.Init(min, max, "test")
	if err != nil {
		panic(err)
	}
	uniq := make(map[uint64]bool)
	for i := uint64(min); i < max; i++ {
		encodeNum, err := numShuffle.ExtendDecode(i)
		if err != nil {
			panic(err)
		}
		decodeNum, err := numShuffle.ExtendDecode(encodeNum)
		if err != nil {
			panic(err)
		}
		if decodeNum > max {
			panic(fmt.Sprintf("mph i=%d to %d err,cuz range", i, encodeNum))
		}

		if decodeNum != i {
			panic(fmt.Sprintf("mph decode from %d to i=%d err,got %d", encodeNum, i, decodeNum))
		}
		if _, ok := uniq[encodeNum]; ok {
			panic(fmt.Sprintf("mph i=%d to %d err,cuz decode exist", i, encodeNum))
		}
		uniq[encodeNum] = true
		fmt.Printf("mph i=%d to %d\n", i, encodeNum)
	}
	if uint64(len(uniq)) != max-min {
		panic(fmt.Sprintf("mph err,cuz got=%d,expect=%d", uint64(len(uniq)), max))
	}
}
