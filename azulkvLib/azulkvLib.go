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
	"unsafe"
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

	fill :=0
	db.Entries = &fill
	capacity := db.Cap
	hash := make([]uint64, capacity)
	db.Hash = &hash
	keys := make([]string, capacity)
	db.Keys = &keys
	vals := make([]string, capacity)
	db.Vals = &vals

    // find dir
    _, err = os.Stat(dirPath)
    if err != nil {
        if os.IsNotExist(err) {
            if dbg {log.Printf("db dir does not exist!\n")}

            //create directory
           	if  err1 := os.Mkdir(dirPath, 0755); err1 != nil {return &db, fmt.Errorf("could not create dir: %v", err1)}

            db.DirPath = dirPath

            //create files
/*
            tabNam := dirPath + "/azulkvBase.dat"
            fil, err1:= os.Create(tabNam)
            if err1 != nil {return &db, fmt.Errorf("could not create table: %v", err1)}
            db.Tab=fil
            db.TabNam = tabNam
*/  
          return &db, nil
        } else {
            return &db, fmt.Errorf("could not open dir: %v", err)
        }
    }
    log.Printf("azulkv dir exists!\n")

    db.DirPath = dirPath

	return &db, nil
}

func GetHash(bdat []byte) (hash uint64) {

	seed :=uint64(0)
	hash = t1ha.Sum64(bdat, seed)

	return hash
}

func GenRanData (rangeStart, rangeEnd int) (bdat []byte) {

	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

//    rangeStart := 5
//    rangeEnd := 25
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
		bdat := GenRanData(5, 25)
		hashval := GetHash(bdat)
		valdat := GenRanData(5, 40)
		valstr := fmt.Sprintf("val-%d_%s",i,string(valdat))
//		valb := []byte(valstr)
		(*db.Keys)[i] = string(bdat)
		(*db.Hash)[i] = hashval
		(*db.Vals)[i] = valstr
//		fmt.Printf(" %d: %d %s %s\n", i, (*db.Hash)[i], (*db.Keys)[i], (*db.Vals)[i])
	}
	db.Num = level
	(*db.Entries) = level
	dbpt = &db
//fmt.Printf("fil db: %v\n", dbpt)
	return nil
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

func (dbP *kvObj) Clean () (err error){

	return err
}

func (dbp *kvObj) Backup (tabNam string) (err error){

    db := *dbp
//    kvMap, err := InitKV("testDb", true)
//    if err != nil {log.Fatalf("error -- InitKV: %v", err)}
	numEntries := *db.Entries
	dirPath := db.DirPath
	filPath := dirPath + "/" + tabNam
	if len(dirPath) == 0 {return fmt.Errorf("DirPath not found!")}
    _, err = os.Stat(filPath)
    if err == nil {
		return fmt.Errorf("table %s already exists!: %v", tabNam, err)
    }

	//create table
	outfil, err:= os.Create(filPath)
	defer outfil.Close()
	if err != nil {return fmt.Errorf("could not create table: %v", err)}

	numEnt := uint32(numEntries)
	backSize := 4 + int(unsafe.Sizeof(numEnt))*numEntries *2

	bck := make([]byte, backSize, 4096)


	pt := (*[4]byte)(unsafe.Pointer(&numEnt))[:]
	copy(bck[:4], pt)
	for i:=0; i<3; i++ {
		fmt.Printf("%d:", bck[i])
	}
	fmt.Printf("%d\n", bck[4])

	start := 4
	for i:=0; i<numEntries; i++ {
		entry := uint32(i)
		pt := (*[4]byte)(unsafe.Pointer(&entry))[:]
		copy(bck[start+i*4:start+(i+1)*4], pt)
	}

	start = numEntries*4 + 4
	for i:=0; i<numEntries; i++ {
		klen := uint16(len((*db.Keys)[i]))
		pt := (*[2]byte)(unsafe.Pointer(&klen))[:]
		copy(bck[start:start+2], pt)
		vlen := uint16(len((*db.Vals)[i]))
		pt2 := (*[2]byte)(unsafe.Pointer(&vlen))[:]
		copy(bck[start+2:start+4], pt2)
		start = start + 4
		fmt.Printf("  %d: kl %d vl %d\n",i, klen, vlen)
		key := []byte((*db.Keys)[i])
		copy(bck[start:start+int(klen)],key)
		val := []byte((*db.Vals)[i])
		copy(bck[start +int(klen):start+int(klen)+int(vlen)],val)
		fmt.Printf("klen: %d vlen: %d key: %s val %s\n", klen, vlen, string(key), string(val))
		start = start + int(klen) + int(vlen)
	}
	endpt := start
	fmt.Printf("endpt: %d\n",endpt)
	_, err = outfil.Write(bck[:endpt])
	if err !=nil {return fmt.Errorf("backup write: %v", err)}
	return nil
}


func (dbp *kvObj) Load(tabNam string) (err error){
	var numEntries uint32

    db := *dbp
//	capacity := db.Cap

	dirPath := db.DirPath
	filPath := dirPath + "/" + tabNam

	bckup, err := os.ReadFile(filPath)
	if err != nil {return fmt.Errorf("could not read table: %v", err)}

	siz := len(bckup)

	fmt.Printf("backup: %d\n",siz)

	if siz < 4 {return fmt.Errorf("no valid numEntries found!")}

	numEntries = *(*uint32)(unsafe.Pointer(&bckup[0]))
	numKeys := int(numEntries)
	(*db.Entries) = numKeys

	// fill kvObj
//	Keys := make([]string,numKeys, capacity)

/*
	for i:=0; i<3; i++ {
		fmt.Printf("%d:", bckup[i])
	}
	fmt.Printf("%d\n", bckup[4])
*/
	fmt.Printf("numEntries: %d\n", numEntries)
	entries := make([]uint32, numKeys)

	for i:=0; i< numKeys; i++ {
		entries[i] = *(*uint32)(unsafe.Pointer(&bckup[4+i*4]))
	}

	for i:=0; i< numKeys; i++ {
		fmt.Printf("entry[%d]: %d\n", i, entries[i])
	}

	start := 4 + numKeys*4
	for i:=0; i< numKeys; i++ {
		klen := *(*uint16)(unsafe.Pointer(&bckup[start]))
		vlen := *(*uint16)(unsafe.Pointer(&bckup[start +2]))
		fmt.Printf("  %d: klen %d vlen %d\n", i, klen, vlen)
		key := bckup[start +4: start+4+int(klen)]
		val := bckup[start +4 + int(klen): start+4+int(klen)+int(vlen)]
		start = start + 4 + int(klen) + int(vlen)
		(*db.Keys)[i] = string(key)
		(*db.Vals)[i] = string(val)
		(*db.Hash)[i] = GetHash(key)

		fmt.Printf("klen: %d vlen: %d key: %s val %s\n", klen, vlen, string(key), string(val))
	}
    dbp = &db
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

func (dbpt *kvObj) PrintKV (idx int, num int) {

	db := *dbpt
	fmt.Printf("********* Entries: %d *************\n", (*db.Entries))
	if idx+num > (*db.Entries) {
		fmt.Printf("invalid idx; idx + num > %d!\n", db.Entries)
		return
	}
	for i:=idx; i<idx + num; i++ {
		fmt.Printf("  [%2d]: %d %20s %s\n", i, (*db.Hash)[i], (*db.Keys)[i], (*db.Vals)[i])
	}
	fmt.Printf("********* End Entries *************\n")
	return
}

