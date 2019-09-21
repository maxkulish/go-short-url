package cache

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

const (
	localFile = "./cache/local.db"
	bucket    = "shortLinks"
)

type BoltDB struct {
	Conn *bolt.DB
}

func NewBoltDB() *BoltDB {

	db, err := bolt.Open(localFile, 0600, &bolt.Options{Timeout: 1 * time.Second})

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("create bucket: %+v", err)
		}

		return nil
	})
	if err != nil {
		log.Fatalf("can't open local BoltDB file: %s", localFile)
	}

	return &BoltDB{
		Conn: db,
	}
}

// Returns one kye/value pair from BoltDB
func (d *BoltDB) GetFullURL(short string) string {

	var fullURL []byte
	err := d.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		fullURL = b.Get([]byte(short))
		return nil
	})
	if err != nil {

	}

	return string(fullURL)
}

// Returns one kye/value pair from BoltDB
func (d *BoltDB) GetALL() map[string][]byte {

	var allInCache = map[string][]byte{}
	err := d.Conn.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		err := b.ForEach(func(k, v []byte) error {
			allInCache[string(k)] = v
			return nil
		})

		if err != nil {
			return fmt.Errorf("get all error: %+v", err)
		}

		return nil

	})

	if err != nil {
		return nil
	}

	return allInCache
}

// Save a bunch of keys/values in one bulk operation
func (d *BoltDB) Put(short, fullURL string) error {

	err := d.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put([]byte(short), []byte(fullURL))
		return err
	})
	if err != nil {
		return fmt.Errorf("put error: %+v", err)
	}

	return nil
}

// Close connection to the BoltDB
func (d *BoltDB) Close() {
	d.Close()
}
