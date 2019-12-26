package fdb

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"path/filepath"
	"util"
)

var (
	log = logrus.WithFields(logrus.Fields{"pkg": "fdb"})
)

type FalDB struct {
	DB *bolt.DB
}

func NewDB(path string) *FalDB {
	util.MkdirIfNotExists(filepath.Dir(path))
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}
	return &FalDB{DB:db}
}

func (db *FalDB) CreateBucket(bucketName string) error {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (db *FalDB) DeleteBucket(bucketName string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketName))
	})
}

func (db *FalDB) GetAllBucket() ([]string, error) {
	var err error
	bucketList := make([]string, 0)
	err = db.DB.Update(func(tx *bolt.Tx) error {
		err := tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			bucketList = append(bucketList, string(name))
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return bucketList, nil
}

func (db *FalDB) GetAllKV(bucketName string) (map[string]string, error) {

	var err error
	keyList := make(map[string]string)

	err = db.ForEach(bucketName, func(k, v []byte) error {
		keyList[string(k)] = string(v)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return keyList, nil
}

func (db *FalDB) ForEach(bucketName string, fn func(k, v []byte) error) error {
	log.Debug("enter BoltDB ForEach func")
	err := db.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			log.Errorf("boltdb bucket:%s doesn't exist", bucketName)
			return fmt.Errorf("boltdb bucket:%s doesn't exist", bucketName)
		}
		return bucket.ForEach(fn)
	})
	return err
}


func (db *FalDB) Put (key, value, bucket string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		//根据name找到对应bucket
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("the bucket does not exist")
		}

		//将key/value对加入bucket
		return b.Put([]byte(key), []byte(value))
	})
}

func (db *FalDB) Get (key, bucket string) (string, error) {
	var value []byte
	err := db.DB.View(func(tx *bolt.Tx) error {
		//根据name找到对应bucket
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("the bucket does not exist")
		}

		//获得key对应value
		value = b.Get([]byte(key))
		if value == nil {
			return errors.New("the key does not exist")
		}
		return nil
	})
	return string(value), err
}

func (db *FalDB) Delete (key, bucket string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		//根据name找到对应bucket
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("the bucket does not exist")
		}

		//删除key/value
		return b.Delete([]byte(key))
	})
}

func (db *FalDB) PutBatch(data map[string]map[string]string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		for b, kvs := range data {
			bucket, err := tx.CreateBucketIfNotExists([]byte(b))
			if err != nil {
				log.Error(err)
				return err
			}
			for k, v := range kvs {
				err = bucket.Put([]byte(k), []byte(v))
				if err != nil {
					log.Error(err)
					return err
				}
			}
		}
		return nil
	})
}
