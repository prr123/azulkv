// azulkv
// library of simple kv
// Author: prr azulsoftware
// Date: 27. Aug 2023
// copyright 2027 prr azul software
//

package azulkv

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"os"
	"github.com/dgryski/go-t1ha"
)

type kvObj struct {
	DirPath string
	Dbg bool
	Cap int
	Num int
	Entries *int
	Hash *[]uint64
	Keys *[]string
	Vals *[]string
	TabNam string
	Tab *os.File
}


func InitKV(dirPath string, dbg bool) (dbpt *kvObj, err error){

	db := kvObj {
		Cap: 500,
		Num: 1,
		Dbg: dbg,
	}

    // find dir
    _, err = os.Stat(dirPath)
    if err != nil {
        if os.IsNotExist(err) {
            if dbg {log.Printf("db dir does not exist!\n")}

            //create directory
           	if  err1 := os.Mkdir(dirPath, 0755); err1 != nil {return &db, fmt.Errorf("could not create dir: %v", err1)}

            db.DirPath = dirPath

            //create files
            tabNam := dirPath + "/azulkvBase.dat"
            fil, err1:= os.Create(tabNam)
            if err1 != nil {return &db, fmt.Errorf("could not create table: %v", err1)}
            db.DirPath = dirPath
            db.Tab=fil
            db.TabNam = tabNam
            return &db, nil
        } else {
            return &db, fmt.Errorf("could not open dir: %v", err)
        }
    }
    log.Printf("azulkv dir exists!\n")

    db.DirPath = dirPath

	tabNam := dirPath + "/azulkvBase.dat"
	fil, err2 := os.Open(tabNam)
	if err2 != nil {
		return &db, fmt.Errorf("could not open table: %v", err2)
	}
	db.Tab = fil
	db.TabNam = tabNam

	fill :=0
	db.Entries = &fill
	capacity := db.Cap
	hash := make([]uint64, capacity)
	db.Hash = &hash
	keys := make([]string, capacity)
	db.Keys = &keys
	vals := make([]string, capacity)
	db.Vals = &vals

	err = db.Load()
	if err != nil {return &db, fmt.Errorf(" could not load table! %v", err)}
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


func (dbpt *kvObj) FillRan (level int) (err error){

	db := *dbpt
	for i:=0; i<level; i++ {
		bdat := GenRanData()
		hashval := GetHash(bdat)
		valstr := fmt.Sprintf("val_%d",i)
		valb := []byte(valstr)
		(*db.Keys)[i] = string(bdat)
		(*db.Hash)[i] = hashval
		(*db.Vals)[i] = string(valb)
//		fmt.Printf(" %d: %d %s %s\n", i, (*db.Hash)[i], (*db.Keys)[i], (*db.Vals)[i])
	}
	db.Num = level
	(*db.Entries) = level
	dbpt = &db
//fmt.Printf("fil db: %v\n", dbpt)
	return nil
}

func (dbpt *kvObj) PrintKV (idx int, num int) {

	db := *dbpt
	fmt.Printf("Entries: %d\n", (*db.Entries))
	if idx+num > (*db.Entries) {
		fmt.Printf("invalid idx; idx + num > %d!\n", db.Entries)
		return
	}
	for i:=idx; i<idx + num; i++ {
		fmt.Printf("  [%2d]: %d %20s %s\n", i, (*db.Hash)[i], (*db.Keys)[i], (*db.Vals)[i])
	}
	return
}


func (dbp *kvObj) AddEntry (key, val string) (err error){

	db := *dbp
	idx := (*db.Entries)
	if idx > db.Cap-2 {return fmt.Errorf("entry exceeds limits")}

	hashval := GetHash([]byte(key))
	(*db.Hash)[idx] = hashval
	(*db.Keys)[idx] = key
	(*db.Vals)[idx] = val

	(*db.Entries)++
	dbp = &db

	return nil
}

func (dbp *kvObj) UpdEntry (idx int, val string) (err error){

	db := *dbp
	if idx > db.Cap {return fmt.Errorf("invalid index")}
	(*db.Vals)[idx] = val
	dbp = &db
	return nil
}

func (dbp *kvObj) DelEntry (idx int) (err error){

	db := *dbp
	if idx > db.Cap {return fmt.Errorf("invalid index")}
	(*db.Hash)[idx] = 0
	(*db.Keys)[idx] = ""
	(*db.Vals)[idx] = ""
	dbp = &db

	return nil
}

func (dbp *kvObj) GetVal (keyStr string) (idx int, valstr string){

	db := *dbp
	idx = -1
	for i:=0; i< (*db.Entries); i++ {
		if (*db.Keys)[i] == keyStr {
			idx = i
			valstr = (*db.Vals)[i]
			return idx, valstr
		}
	}
	return idx, ""
}

func (dbp *kvObj) GetValByHash (key string) (idx int, valstr string){

	db := *dbp
	hashval := GetHash([]byte(key))

	for i:=0; i< (*db.Entries); i++ {
		if (*db.Hash)[i] == hashval {
			idx = i
			valstr = (*db.Vals)[i]
			return idx, valstr
		}
	}
	return idx, ""
}

func (dbp *kvObj) FindKeyByHash (key string) (idx int){
	db := *dbp
	hashval := GetHash([]byte(key))

	for i:=0; i< (*db.Entries); i++ {
		if (*db.Hash)[i] == hashval {
			idx = i
			return idx
		}
	}
	return -1
}

func (dbp *kvObj) FindKey (keyStr string) (idx int) {

	db := *dbp
	for i:=0; i< (*db.Entries); i++ {
		if (*db.Keys)[i] == keyStr {
			idx = i
			return idx
		}
	}
	return -1

}

func (dbp *kvObj) GetKeyByIdx (idx int) (key string) {

	db := *dbp
	if idx > (*db.Entries) {return ""}

	key = (*db.Keys)[idx]
	return key
}

func (db *kvObj) Clean () (err error){

	return err
}

func (db *kvObj) Backup () (err error){

	return err
}

func (db *kvObj) Load () (err error){

	return nil
}

func PrintDb(dbp *kvObj) {

    db := *dbp
//  dbg := db.Dbg

    fmt.Printf("******* AzulKV: %s *******\n", db.DirPath)
    fmt.Printf("Dir:    %s\n",db.DirPath)
//    fmt.Printf("New DB: %t\n",db.Ndb)
    fmt.Printf("  table: %s\n", db.TabNam)
    fmt.Printf("********* End AzulKV: *******\n")
    return
}
