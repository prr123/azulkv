package main


import (
	"fmt"
	"math/rand"
	"time"
//	"os"
	"github.com/dgryski/go-t1ha"
)

type kvObj struct {
	Cap int
	NumEntries int
	Hash *[]uint64
	Keys *[]string
	Vals *[]string
}


func main() {

	kvMap, err := InitKV()
	if err != nil {panic(err)}

	fmt.Printf("kvmap actual %d capacity: %d\n", kvMap.Cap, kvMap.NumEntries)


	kvMap.FillRan(5)

	kvMap.PrintKV(5)

	fmt.Printf("success\n")

}

func InitKV() (dbpt *kvObj, err error){

	db := kvObj {
		Cap: 500,
		NumEntries: 0,
	}

	capacity := db.Cap
	hash := make([]uint64, capacity)
	db.Hash = &hash
	keys := make([]string, capacity)
	db.Keys = &keys
	vals := make([]string, capacity)
	db.Vals = &vals

	return &db, nil
}

func GetHash(bdat []byte) (hash uint64) {

	seed :=uint64(0)
	hash = t1ha.Sum64(bdat, seed)

	return hash
}

func GenRanData () (bdat []byte) {

	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

    rangeStart := 5
    rangeEnd := 25
    offset := rangeEnd - rangeStart

    randLength := seededRand.Intn(offset) + rangeStart
    bdat = make([]byte, randLength)

    charset := "abcdefghijklmnopqrstuvw0123456789"
    for i := range bdat {
        bdat[i] = charset[seededRand.Intn(len(charset)-1)]
    }
	return bdat
}


func (dbpt *kvObj) FillRan (level int) {

	db := *dbpt
	for i:=0; i<level; i++ {
		bdat := GenRanData()
		hashval := GetHash(bdat)
		valstr := fmt.Sprintf("val_%d",i)
		valb := []byte(valstr)

		fmt.Printf("str{%d]: %s %d\n", len(bdat), string(bdat), hashval)

		(*db.Keys)[i] = string(bdat)
		(*db.Hash)[i] = hashval
		(*db.Vals)[i] = string(valb)
	}
	db.NumEntries = level
	dbpt = &db
	return
}

func (db *kvObj) PrintKV (level int) {

//	level := db.NumEntries
	fmt.Printf("Entries: %d\n", db.NumEntries)

	for i:=0; i<level; i++ {
		fmt.Printf("  [%2d]: %d %20s %s\n", i, (*db.Hash)[i], (*db.Keys)[i], (*db.Vals)[i])
	}
}


func (db *kvObj) AddEntry (key, val string) (err error){

	return nil
}

func (db *kvObj) UpdEntry (key, val string) (err error){

	return nil
}

func (db *kvObj) DelEntry (key string) (err error){

	return err
}

func (db *kvObj) FindKey (keystring) (bool) {

	return false
}

func (db *kvObj) Clean () (err error){

	return err
}

func (db *kvObj) Backup () (err error){

	return err
}

func (db *kvObj) KVLoad () (err error){

	return err
}


