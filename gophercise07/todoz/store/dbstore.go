package store

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

var DbCon *bolt.DB

// Tablename is Bucket name of the database.
var Tablename = []byte("todo")

// Todoz is the datastructure which represents database values for operations.
type Todoz struct {
	Id   int
	Task string
}

func itob(i int) []byte {
	bt := make([]byte, 8)
	binary.BigEndian.PutUint64(bt, uint64(i))
	return bt
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

// InsertTask is used to add new entry in database
func InsertTask(task string) (int, error) {
	var id int
	err := DbCon.Update(func(t *bolt.Tx) error {
		buc := t.Bucket(Tablename)

		id64, _ := buc.NextSequence()
		id := int(id64)
		key := itob(id)
		return buc.Put(key, []byte(task))
	})

	return id, err
}

// RemoveTasks deletes the entry by selecting the parameter provided.
func RemoveTasks(id int) error {
	return DbCon.Update(func(t *bolt.Tx) error {
		buc := t.Bucket(Tablename)
		return buc.Delete(itob(id))
	})
}

// GetAll gets all the entries from the database
func GetAll() ([]Todoz, error) {
	var todolist []Todoz
	err := DbCon.View(func(t *bolt.Tx) error {
		buc := t.Bucket(Tablename)
		cur := buc.Cursor()
		for key, val := cur.First(); key != nil; key, val = cur.Next() {
			todolist = append(todolist, Todoz{
				Id:   btoi(key),
				Task: string(val),
			})
		}
		return nil
	})
	return todolist, err
}

// Init initializes the database and connects to it.
// In case the database do not exists. It creates one
func Init(dbpath string) error {
	var err error
	DbCon, err = bolt.Open(dbpath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return DbCon.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucketIfNotExists(Tablename)
		return err
	})

}
