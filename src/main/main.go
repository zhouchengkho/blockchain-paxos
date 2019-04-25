package main

import (
	//"testing"
	"strconv"
	"os"
	//"runtime"
	"fmt"
	//"time"
	//"strings"
	//"sync/atomic"
	"chainpaxos"
	"encoding/json"
	"runtime"
	"time"
)

type Clerk = chainpaxos.Clerk
type KVPaxos = chainpaxos.KVPaxos


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

func TestBasic() {
	const nservers = 3
	const nclients = 2
	var kva []*KVPaxos = make([]*KVPaxos, nservers)
	var kvh []string = make([]string, nservers)
	defer cleanup(kva)

	for i := 0; i < nservers; i++ {
		kvh[i] = port("basic", i)
	}
	for i := 0; i < nservers; i++ {
		kva[i] = chainpaxos.StartServer(kvh, i)
	}

	var cka [nclients]*Clerk
	for i := 0; i < nclients; i++ {
		cka[i] = chainpaxos.MakeClerk(kvh)
	}
	block0 := Block{}
	block0.Transactions = make([]Transaction, 3)
	for i := 0; i < 3; i++ {
		block0.Transactions[i] = Transaction{"Cheng", "Lingxue", 33.33 }
	}
	cka[0].Append(chainpaxos.LastBlock, createValue(block0))
	gotBlock := cka[1].Get(chainpaxos.LastBlock)
	fmt.Println(gotBlock)
	parsedBlock := parseBlock(gotBlock)
	fmt.Println(parsedBlock)
}

func TestPartition() {
	runtime.GOMAXPROCS(4)

	tag := "partition"
	const nservers= 5
	var kva []*KVPaxos = make([]*KVPaxos, nservers)
	defer cleanup(kva)
	defer cleanpp(tag, nservers)

	for i := 0; i < nservers; i++ {
		var kvh []string = make([]string, nservers)
		for j := 0; j < nservers; j++ {
			if j == i {
				kvh[j] = port(tag, i)
			} else {
				kvh[j] = pp(tag, i, j)
			}
		}
		kva[i] = chainpaxos.StartServer(kvh, i)
	}
	defer part(tag, nservers, []int{}, []int{}, []int{})

	var cka [nservers]*Clerk
	for i := 0; i < nservers; i++ {
		cka[i] = chainpaxos.MakeClerk([]string{port(tag, i)})
	}

	fmt.Printf("Test: No partition ...\n")

	part(tag, nservers, []int{0, 1, 2, 3, 4}, []int{}, []int{})
	cka[0].Put(chainpaxos.LastBlock, "12")
	cka[2].Put(chainpaxos.LastBlock, "13")

	fmt.Printf("Test: Progress in majority ...\n")

	part(tag, nservers, []int{2, 3, 4}, []int{0, 1}, []int{})

	done := make(chan bool)
	go func() {
		cka[0].Put(chainpaxos.LastBlock, "14")
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("fatal: put in minority succeeded")
	case <-time.After(time.Second):
	}

	fmt.Println("Put in minority will not succeed until heal")


	fmt.Printf("Test: Completion after heal ...\n")

	part(tag, nservers, []int{0, 2, 3, 4}, []int{1}, []int{})

	select {
	case <-done:
	case <-time.After(30 * 100 * time.Millisecond):
		fmt.Println("fatal: put did not complete after heal")
	}
	fmt.Println(cka[0].Get(chainpaxos.LastBlock))
	// fmt.Println(cka[1].Get("1"))



}


func createValue(block Block) string {
	res, _ := json.Marshal(block)
	return string(res)
}

func parseBlock(s string) Block {
	res := Block{}
	json.Unmarshal([]byte(s), &res)
	return res
}

func main() {
	// TestBasic()
	TestPartition()
}