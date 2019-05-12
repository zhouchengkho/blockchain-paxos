package blockchain

import "encoding/json"

type Transaction struct {
	From string
	To string
	Amount float32
}

type Block struct {
	Transactions []Transaction
}

type Partition struct {
	Partition1 []int
	Partition2 []int
	Partition3 []int
}

type BlockInfo struct {
	BlockData Block
	Index int
}


func BlockToStr(block Block) string {
	res, _ := json.Marshal(block)
	return string(res)
}

func StrToBlock(s string) Block {
	res := Block{}
	json.Unmarshal([]byte(s), &res)
	return res
}

func BlocksToStr(blocks []Block) string {
	res, _ := json.Marshal(blocks)
	return string(res)
}

func BytesToBlock(bytes []byte) Block {
	res := Block{}
	json.Unmarshal(bytes, &res)
	return res
}

func BytesToPartition(bytes []byte) Partition {
	res := Partition{}
	json.Unmarshal(bytes, &res)
	return res
}
