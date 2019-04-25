package main

type Transaction struct {
	From string
	To string
	Amount float32
}

type Block struct {
	Transactions []Transaction
}
