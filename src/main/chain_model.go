package main

import (
	"chainpaxos"
	"net/http"
	"fmt"
	"strconv"
	"io"
	"io/ioutil"
	"github.com/julienschmidt/httprouter"
	"blockchain"
	"os"
	"encoding/json"
)

func cleanup(kva []*KVPaxos) {
	for i := 0; i < len(kva); i++ {
		if kva[i] != nil {
			kva[i].Kill()
		}
	}
}

func pp(tag string, src int, dst int) string {
	s := "/var/tmp/824-"
	s += strconv.Itoa(os.Getuid()) + "/"
	s += "kv-" + tag + "-"
	s += strconv.Itoa(os.Getpid()) + "-"
	s += strconv.Itoa(src) + "-"
	s += strconv.Itoa(dst)
	return s
}

func cleanpp(tag string, n int) {
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			ij := pp(tag, i, j)
			os.Remove(ij)
		}
	}
}

func part(tag string, npaxos int, p1 []int, p2 []int, p3 []int) {
	cleanpp(tag, npaxos)

	pa := [][]int{p1, p2, p3}
	for pi := 0; pi < len(pa); pi++ {
		p := pa[pi]
		for i := 0; i < len(p); i++ {
			for j := 0; j < len(p); j++ {
				ij := pp(tag, p[i], p[j])
				pj := port(tag, p[j])
				err := os.Link(pj, ij)
				if err != nil {
					fmt.Println("partition failure")
					fmt.Println(err)
				}
			}
		}
	}
	model.currentPartition.Partition1 = p1
	model.currentPartition.Partition2 = p2
	model.currentPartition.Partition3 = p3
}

func port(tag string, host int) string {
	s := "/var/tmp/824-"
	s += strconv.Itoa(os.Getuid()) + "/"
	os.Mkdir(s, 0777)
	s += "kv-"
	s += strconv.Itoa(os.Getpid()) + "-"
	s += tag + "-"
	s += strconv.Itoa(host)
	return s
}

type ChainModel struct {
	kva []*KVPaxos
	tag string
	cka [NSERVERS]*Clerk
	currentPartition blockchain.Partition
}

func NewChainModel() *ChainModel {
	model := new(ChainModel)
	model.kva = make([]*KVPaxos, NSERVERS)

	for i := 0; i < NSERVERS; i++ {
		var kvh []string = make([]string, NSERVERS)
		for j := 0; j < NSERVERS; j++ {
			if j == i {
				kvh[j] = port(TAG, i)
			} else {
				kvh[j] = pp(TAG, i, j)
			}
		}
		model.kva[i] = chainpaxos.StartServer(kvh, i)
	}
	for i := 0; i < NSERVERS; i++ {
		model.cka[i] = chainpaxos.MakeClerk([]string{port(TAG, i)})
	}
	return model
}

func bodyToBytes(body io.ReadCloser) []byte {
	bytes, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("error parsing body!!")
		return make([]byte, 0)
	}
	return bytes
}

func GetBlocks(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// client i is connected to server i
	blockPosition, err := strconv.Atoi(ps.ByName("block"))
	if err != nil {
		blockPosition = chainpaxos.LastBlock
	}
	clientIndex, err := strconv.Atoi(ps.ByName("client"))
	if err != nil {
		clientIndex = 0
	}
	fmt.Fprint(w, model.cka[clientIndex].Get(blockPosition))
}

func CreateBlock(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// client i is connected to server i
	client := ps.ByName("client")
	// block := ps.ByName("name")
	clientIndex, err := strconv.Atoi(client)
	if err != nil {
		clientIndex = 0
	}
	block := blockchain.BytesToBlock(bodyToBytes(r.Body))
	model.cka[clientIndex].Append(chainpaxos.LastBlock, blockchain.BlockToStr(block))
	fmt.Fprint(w, "OK")
}

func CreateNetworkPartition(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	partition := blockchain.BytesToPartition(bodyToBytes(r.Body))
	part(TAG, NSERVERS, partition.Partition1, partition.Partition2, partition.Partition3)
	fmt.Fprint(w, "OK")
}

func GetNetworkPartition(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	res, _ := json.Marshal(model.currentPartition)
	fmt.Fprint(w, string(res))
}