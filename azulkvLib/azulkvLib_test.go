
package azulkv

import (
//	"log"
	"testing"
//	"db/azulkv/azulkvLib"
)

func TestAddEntry(t *testing.T) {
	kv, err := InitKV("testDb", true)
	if err != nil {t.Errorf("error -- InitKV: %v", err)}

	err = kv.AddEntry("key1", "val1")
	if err != nil {t.Errorf("error -- AddEntry: %v", err)}

	if (*kv.Keys)[0] != "key1" {t.Errorf("keys do not agree: %s is not %s!", (*kv.Keys)[0], "key1")}

}

func TestDelEntry(t *testing.T) {

	kv, err := InitKV("testDb", true)
	if err != nil {t.Errorf("error -- InitKV: %v", err)}

	err = kv.AddEntry("key1", "val1")
	if err != nil {t.Errorf("error -- AddEntry: %v", err)}

	if (*kv.Keys)[0] != "key1" {t.Errorf("keys do not agree: %s is not %s!", (*kv.Keys)[0], "key1")}

	idx := kv.FindKey("key1")
	if idx == -1 {t.Errorf("error -- FindKey: %d key1 not found!", idx)}
	
	err = kv.DelEntry(idx)
	if err != nil {t.Errorf("error -- DelEntry: %v", err)}

}
