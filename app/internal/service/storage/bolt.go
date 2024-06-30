package storage

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
)

type BoltDb struct {
	dbname string
	db     *bolt.DB
	Bucket string
}

func NewBoltDb(dbname string) (*BoltDb, error) {
	db, err := bolt.Open(dbname, 0600, nil)
	if err != nil {
		return nil, err

	}

	return &BoltDb{dbname: dbname, db: db, Bucket: "face"}, nil
}

func (b *BoltDb) Update(key string, value *RegisteInfo) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(b.Bucket))
		if err != nil {
			return err
		}
		item, _ := json.Marshal(value)
		return bucket.Put([]byte(key), item)
	})
}

func (b *BoltDb) Read(key string) (value *RegisteInfo, ok bool) {
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.Bucket))
		if bucket == nil {
			return fmt.Errorf("Bucket not found")
		}

		item := bucket.Get([]byte(key))
		if item == nil {
			return fmt.Errorf("Key not found")
		}

		err := json.Unmarshal(item, &value)
		if err != nil {
			return err
		}

		return nil
	})

	if err == nil {
		ok = true
	}

	return value, ok
}

func (b *BoltDb) ReadBatch() (values []*RegisteInfo, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.Bucket))
		if bucket == nil {
			return fmt.Errorf("Bucket not found")
		}

		return bucket.ForEach(func(k, v []byte) error {
			info := RegisteInfo{}
			err := json.Unmarshal(v, &info)
			if err != nil {
				return err
			}
			values = append(values, &info)
			return nil
		})
	})
	return values, err
}

func (b *BoltDb) Delete(key string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(b.Bucket))
		if bucket == nil {
			return fmt.Errorf("Bucket not found")
		}
		return bucket.Delete([]byte(key))
	})
}
