package main

import (
	//"testing"
	//"runtime"
	//"time"
	//"strings"
	//"sync/atomic"
	"chainpaxos"
	"runtime"
	// "bufio"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"blockchain"
	"github.com/rs/cors"
)

type Clerk = chainpaxos.Clerk
type KVPaxos = chainpaxos.KVPaxos
type Block = blockchain.Block
type Transaction = blockchain.Transaction

const NSERVERS = 5

const TAG = "partition"

var model = NewChainModel()

func main() {

	// init
	runtime.GOMAXPROCS(4)

	// model defined in chain_model
	defer cleanup(model.kva)
	defer cleanpp(TAG, NSERVERS)


	defer part(TAG, NSERVERS, []int{}, []int{}, []int{})


	// default status: start all server without partition
	part(TAG, NSERVERS, []int{0, 1, 2, 3, 4}, []int{}, []int{})

	router := httprouter.New()
	router.GET("/client/:client/block/:block", GetBlocks)
	router.POST("/client/:client/block", CreateBlock)
	router.POST("/partition", CreateNetworkPartition)
	router.GET("/partition", GetNetworkPartition)


	handler := cors.AllowAll().Handler(router)

	http.ListenAndServe(":3000", handler)

}