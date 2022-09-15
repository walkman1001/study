package sample

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

const dbName = "my.db"
const bucketName = "bucket01"

func createDB() {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

func insertRecord(k string, v string) {

	db, err := bolt.Open(dbName, 0600, nil)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		fmt.Println("b=", b)
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte(k), []byte(v))
		return nil
	})
}

func searchRecordByKey(k string) {
	db, err := bolt.Open(dbName, 0600, nil)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		v := b.Get([]byte(k))
		fmt.Printf("The answer is: %s\n", v)
		return nil
	})
}

func boltBbTest(i int) {

	createDB()
	insertRecord(string(time.Now().Second()), string(time.Now().Nanosecond()))
	searchRecordByKey("1")
	fmt.Println("i=", i)
}

func mainTest() {
	for i := 0; i < 10000; i++ {
		//fmt.Println("i=", i)
		go boltBbTest(i)
	}
	time.Sleep(200 * time.Second)
}
