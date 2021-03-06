package crawler

import (
	"github.com/leona/domain-crawler/src/crawler/utilities"
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"strings"
	"errors"
)

type StoreType struct {
	db *bolt.DB
}

func (self *StoreType) init() {
	var err error
	self.db, err = bolt.Open(*utilities.InputOptions.Db + ".db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}
}

func (self *StoreType) put(root StoreKey, value string) (error) {
	tx, err := self.db.Begin(true)

	if err != nil {
		fmt.Println("Error starting db for put", err)
		return nil
	}
	defer tx.Rollback()

	bucket, _ := tx.CreateBucketIfNotExists([]byte(root))
	err = bucket.Put([]byte(value), []byte(""))

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		fmt.Println("Transaction error")
		return err
	}

	return nil
}

func (self *StoreType) Pop(root StoreKey, count int) []string {
	utilities.Info(3, "Popping:", root)
	tx, err := self.db.Begin(true)

	if err != nil {
		fmt.Println("Error starting db for pop", err)
		return nil
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte(root))

	if bucket == nil {
		utilities.Info(3, "Failed to pop. Could not find root bucket.")
	}

	cursor := bucket.Cursor()
	output := []string{}

	for k, _ := cursor.First(); k != nil; k, _ = cursor.Next() {
		output = append(output, string(k))
		bucket.Delete(k)

		if len(output) >= count {
			break
		}
	}

	if err = tx.Commit(); err != nil {
		fmt.Println("Transaction error", err)
		return nil
	}

	return output
}

func (self *StoreType) getNestedBuckets(root StoreKey, path []string, limit int) ([]string) {
	utilities.Info(3, "getNestedBuckets in root:", root, "- path:", path)
	output := []string{}

	self.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(root))

		if len(path) > 0 {
			for _, item := range path {
				if len(item) == 0 {
					return nil
				}
				bucket = bucket.Bucket([]byte(item))
			}
		}

		if bucket == nil {
			return nil
		}

		var recurse func(*bolt.Bucket, []string) 

		recurse = func(_bucket *bolt.Bucket, base []string) {
			cursor := _bucket.Cursor()

			for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
				if limit > 0 && len(output) > limit {
					return
				}

				component := append(base, string(key))
				output = append(output, strings.Join(utilities.Reverse(component), "."))
				
				_bucket := _bucket.Bucket(key)

				if _bucket != nil {
					recurse(_bucket, component)
				}
			}
		}

		recurse(bucket, path)
		return nil
	})

	return output
}

func (self *StoreType) createNestedBucket(root StoreKey, path []string) error {
	utilities.Info(3, "createNestedBucket in root:", root, "- path:", path)
	tx, err := self.db.Begin(true)

	if err != nil {
		fmt.Println("Error starting db for createNestedBucket", err)
		return errors.New("Error starting db transaction")
	}
	defer tx.Rollback()

	var bucket *bolt.Bucket
	bucket, err = tx.CreateBucketIfNotExists([]byte(root))

	for _, item := range path {
		if len(item) == 0 {
			fmt.Println("Error creating nested bucket. Path empty.")
			return errors.New("Error created nested bucket. Path emtpy.")
		}
		
		bucket, err = bucket.CreateBucketIfNotExists([]byte(item))

		if err != nil {
			fmt.Println("Error in createNestedBucket", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		fmt.Println("Transaction error", err)
		return err
	}

	return nil
}

func (self *StoreType) getNestedBucket(root StoreKey, path []string) (*bolt.Bucket) {
	utilities.Info(3, "getNestedBucket in root:", root, "- path:", path)
	tx, err := self.db.Begin(true)

	if err != nil {
		fmt.Println("Error starting db for getNestedBucket", err)
		return nil
	}
	defer tx.Rollback()

	var bucket *bolt.Bucket
	bucket = tx.Bucket([]byte(root))

	if bucket == nil {
		return nil
	}

	for _, item := range path {
		bucket = bucket.Bucket([]byte(item))

		if bucket == nil {
			break
		}
	}

	return bucket
}