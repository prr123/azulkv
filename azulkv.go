// azulkv.go
// cli program to manipulate the azulkv db
// author: prr azul software
// date: 30 Aug 2023
// copyright 2023 prr, azulsoftware
//

package main

import (
    "log"
    "fmt"
    "os"
	"strings"
    util "github.com/prr123/utility/utilLib"
)

func main() {

    numarg := len(os.Args)
    dbg := false

    useStr := "./azulkv cmd [/kv=key:value]or[/key=key] [/dbg]"
    helpStr := "program to manipulate the azulkv db\ncommands are:\n"
	helpStr = helpStr + "  add /kv=key:value\n"
	helpStr = helpStr + "  del /key=key\n"
	helpStr = helpStr + "  get [/key=key]\n"
	helpStr = helpStr + "  upd /kv=key:value\n"
	helpStr = helpStr + "  entries\n"
	helpStr = helpStr + "  help\n"

    if numarg > 4 {
        fmt.Println("too many arguments in cli!")
		fmt.Printf("usage is: %s\n", useStr)
        os.Exit(-1)
    }

    if numarg == 1 {
		fmt.Printf("insufficient arguments!\n")
		fmt.Printf("usage is: %s\n", useStr)
		os.Exit(1)
	}


	cmdStr := os.Args[1]
    flags:=[]string{"dbg"}

	switch cmdStr {
	case "add":
		if dbg {fmt.Println("add")}
		flags = []string{"kv", "dbg"}
	case "del":
		if dbg {fmt.Println("del")}
		flags = []string{"key", "dbg"}
	case "upd":
		if dbg {fmt.Println("upd")}
		flags = []string{"kv", "dbg"}
	case "get":
		if dbg {fmt.Println("get")}
		flags = []string{"key", "dbg"}
	case "entries":
		if dbg {fmt.Println("entries")}
	case "help":
		fmt.Printf("%s", helpStr)
		fmt.Printf("usage is: %s\n", useStr)
		os.Exit(1)

	default:
		fmt.Printf(" command %s is not vald!\n For more information, see: azulkv help", cmdStr)
		fmt.Printf("usage is: %s\n", useStr)
		os.Exit(1)
	}

	if dbg {
		fmt.Println("dbg -- flags:")
		for i:=0; i<len(flags); i++ {fmt.Printf("  %s\n",flags[i])}
	}

    // default file
    flagMap, err := util.ParseFlagsStart(os.Args, flags, 2)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

    _, ok := flagMap["dbg"]
    if ok {dbg = true}
    if dbg {
		fmt.Printf("dbg -- flag list:\n")
        for k, v :=range flagMap {
            fmt.Printf("  flag: /%s value: %s\n", k, v)
        }
    }

	keyStr := ""
	valStr := ""

	switch cmdStr {
	case "add":
		kval, ok := flagMap["kv"]
    	if !ok {
        	fmt.Printf("cli add error: no kv flag\n",)
			fmt.Printf("usage is: %s\n", useStr)
			os.Exit(-1)
		} else {
			if kval.(string) == "none" {log.Fatalf("cli add error: no key:val string provided with kv flag!")}
			kvStr := kval.(string)
			idx := strings.Index(kvStr, ":")
			if idx == -1 {log.Fatalf("cli add error: no key:val seperator provided win kv value string!")}
			keyStr = kvStr[:idx]
			valStr = kvStr[idx+1:]
			if dbg {fmt.Printf("-- add key: %s value %s\n", keyStr, valStr)}
		}
		// process add
		log.Printf("processing add key: %s value: %s\n", keyStr, valStr)

	case "upd":
		kval, ok := flagMap["kv"]
    	if !ok {
        	fmt.Printf("cli upd error: no kv flag\n",)
			fmt.Printf("usage is: %s\n", useStr)
			os.Exit(-1)
		} else {
			if kval.(string) == "none" {log.Fatalf("cli upd error: no key:val string provided with kv flag!")}
			kvStr := kval.(string)
			idx := strings.Index(kvStr, ":")
			if idx == -1 {log.Fatalf("cli upd error: no key:val seperator provided win kv value string!")}
			keyStr = kvStr[:idx]
			valStr = kvStr[idx+1:]
			if dbg {fmt.Printf("-- upd key: %s value %s\n", keyStr, valStr)}
		}
		// process upd
		log.Printf("processing upd key: %s value: %s\n", keyStr, valStr)

	case "del":
		kval, ok := flagMap["key"]
    	if !ok {
        	fmt.Printf("cli del error: no key flag\n",)
			fmt.Printf("usage is: %s\n", useStr)
			os.Exit(-1)
		} else {
			if kval.(string) == "none" {log.Fatalf("cli del error: no key string provided with key flag!")}
			keyStr = kval.(string)
			if dbg {fmt.Printf("-- del key: %s\n", keyStr)}
		}

		// process del
		log.Printf("processing del key: %s\n", keyStr)

	case "get":
		kval, ok := flagMap["key"]
    	if !ok {
        	fmt.Printf("cli get error: no key flag\n",)
			fmt.Printf("usage is: %s\n", useStr)
			os.Exit(-1)
		} else {
			if kval.(string) == "none" {log.Fatalf("cli get error: no key string provided with key flag!")}
			keyStr = kval.(string)
			if dbg {fmt.Printf("-- get key: %s\n", keyStr)}
		}

		// process get
		log.Printf("processing get key: %s\n", keyStr)

	case "entries":
		// process entries
		log.Printf("processing entries\n")

	default:
		if dbg {fmt.Printf("default cmd: %s\n", cmdStr)}
		for k, _ :=range flagMap {
            if k != "dbg" {
				fmt.Printf("cli error: invalid flag: %s\n",k)
				fmt.Printf("usage is: %s\n", useStr)
				os.Exit(-1)
			}
        }

	}

	log.Println("success end template!")
}
