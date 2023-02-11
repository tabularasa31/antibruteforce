//package boltdb
//
//import (
//	"fmt"
//	"github.com/boltdb/bolt"
//	"github.com/tabularasa31/antibruteforce/config"
//	"log"
//)
//
//func InitDB(cfg *config.Config) (*bolt.DB, error) {
//	fmt.Println("---------Starting BoltDB Init......")
//	db, err := bolt.Open(cfg.DB, 0600, nil)
//	if err != nil {
//		return nil, err
//	}
//
//	if err := db.Update(func(tx *bolt.Tx) error {
//		_, err := tx.CreateBucketIfNotExists([]byte("lists"))
//		if err != nil {
//			return err
//		}
//		return nil
//	}); err != nil {
//		return nil, err
//	}
//
//	fmt.Println("---------Ending Init BoltDB......")
//
//	return db, nil
//}
//

package boltdb

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

var db *bolt.DB

// Init opens a connection to the BoltDB database.
func Init(dbPath string) error {
	fmt.Println("---------trying to init boltdb")
	var err error
	db, err = bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return err
	}
	return nil
}

// Close closes the connection to the BoltDB database.
func Close() {
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// GetDB returns the BoltDB instance.
func GetDB() *bolt.DB {
	return db
}
