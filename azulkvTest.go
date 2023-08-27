package main


import (
	"fmt"
//	"log"
//	"os"
	"db/azulkv/azulkvLib"
)

func main() {

	kvMap, err := azulkv.InitKV()
	if err != nil {panic(err)}

//	fmt.Printf("kvmap actual %d capacity: %d\n", kvMap.NumEntries, kvMap.Cap)
//	fmt.Printf("hash len: %d\n", len(*kvMap.Hash))

	kvMap.FillRan(5)

	kvMap.PrintKV(0,5)

	fmt.Printf("success\n")

}

