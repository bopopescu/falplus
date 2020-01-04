package fdb

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"path/filepath"
	"sync"
	"time"
	"util"
)

var (
	log = logrus.WithFields(logrus.Fields{"pkg": "fdb"})
)

type FalDB struct {
	DB      *bolt.DB
	setchan chan *kvSetBatch
	delchan chan *kvDelBatch
	closet  chan struct{}
	clodel  chan struct{}
}

type kvSetBatch struct {
	key    string
	value  string
	bucket string
	err    error
	wg     sync.WaitGroup
}

type kvDelBatch struct {
	key    string
	bucket string
	err    error
	wg     sync.WaitGroup
}

func NewDB(path string) *FalDB {
	util.MkdirIfNotExists(filepath.Dir(path))
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}
	b := &FalDB{
		DB:      db,
		setchan: make(chan *kvSetBatch, 8191),
		delchan: make(chan *kvDelBatch, 8191),
		closet:  make(chan struct{}),
		clodel:  make(chan struct{}),
	}

	go b.goSetKVBatch()
	go b.goDelKVBatch()
	return b
}

func (db *FalDB) Close() {
	if db == nil {
		return
	}
	for {
		log.Debugf("closing db len(b.setchan)=%d, len(b.delchan)=%d", len(db.setchan), len(db.delchan))
		time.Sleep(time.Second)
		if len(db.setchan) == 0 && len(db.delchan) == 0 {
			break
		}
	}
	close(db.clodel)
	close(db.closet)
	close(db.setchan)
	close(db.delchan)
	db.DB.Close()
	return
}

func (db *FalDB) CreateBucket(bucket string) error {
	err := db.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (db *FalDB) DeleteBucket(bucket string) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucket))
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

func (db *FalDB) GetAllKV(bucket string) (map[string]string, error) {
	var err error
	keyList := make(map[string]string)

	err = db.ForEach(bucket, func(k, v []byte) error {
		keyList[string(k)] = string(v)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return keyList, nil
}

func (db *FalDB) ForEach(bucket string, fn func(k, v []byte) error) error {
	log.Debug("enter BoltDB ForEach func")
	err := db.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			log.Errorf("bucket:%s doesn't exist", bucket)
			return fmt.Errorf("bucket:%s doesn't exist", bucket)
		}
		return b.ForEach(fn)
	})
	return err
}

func (db *FalDB) Get(key, bucket string) (string, error) {
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

func (db *FalDB) Put(key, value, bucket string) error {
	kv := &kvSetBatch{
		bucket: bucket,
		key:    key,
		value:  value,
	}
	kv.wg.Add(1)
	db.setchan <- kv
	kv.wg.Wait()
	return kv.err
}

func (db *FalDB) PutBatch(data map[string]map[string]string) error {
	wg := &sync.WaitGroup{}
	var errs error
	for bucket, value := range data {
		for k, v := range value {
			wg.Add(1)
			go func(k, v, bucket string) {
				err := db.Put(k, v, bucket)
				if err != nil {
					errs = err
				}
				wg.Done()
			}(k, v, bucket)
		}
	}
	wg.Wait()
	return errs
}

func (db *FalDB) putBatch(data map[string]map[string]string) error {
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

func (db *FalDB) Delete(key, bucket string) error {
	kv := &kvDelBatch{
		bucket: bucket,
		key:    key,
	}
	kv.wg.Add(1)
	db.delchan <- kv
	kv.wg.Wait()
	return kv.err
}

func (db *FalDB) DeleteBatch(data []string, bucket string) error {
	wg := &sync.WaitGroup{}
	var errs error
	for _, v := range data {
		wg.Add(1)
		go func(v string) {
			err := db.Delete(v, bucket)
			if err != nil {
				errs = err
			}
			wg.Done()
		}(v)
	}
	wg.Wait()
	return errs
}

func (db *FalDB) delBatch(data map[string][]string) error {
	var err error
	err = db.DB.Update(func(tx *bolt.Tx) error {
		for k, v := range data {
			if v == nil {
				err = tx.DeleteBucket([]byte(k))
				if err != nil {
					log.Errorf("DeleteBucket %s error:%s", k, err)
					return err
				}
			} else {
				bucket := tx.Bucket([]byte(k))
				if bucket == nil {
					log.Errorf("bucket:%s doesn't exist", k)
					return fmt.Errorf("bucket:%s doesn't exist", k)
				}
				for _, key := range v {
					err = bucket.Delete([]byte(key))
					if err != nil {
						log.Errorf("Delete key %s Bucket %s error", key, k)
						return err
					}
				}
			}
		}
		return nil
	})
	return err
}

func (db *FalDB) goSetKVBatch() {
	kvs := make([]*kvSetBatch, 8192)
	for {
		select {
		case v := <-db.setchan:
			if v == nil {
				continue
			}
			length := len(db.setchan)
			log.Debugf("goSetKVBatch length %d", length+1)
			data := make(map[string]map[string]string)
			bucket := make(map[string]string)
			bucket[v.key] = v.value
			data[v.bucket] = bucket
			kvs[0] = v
			for i := 1; i < length+1; i++ {
				v = <-db.setchan
				if v == nil {
					continue
				}
				if _, exist := data[v.bucket]; exist {
					data[v.bucket][v.key] = v.value
				} else {
					bucket := make(map[string]string)
					bucket[v.key] = v.value
					data[v.bucket] = bucket
				}
				kvs[i] = v
			}
			err := db.putBatch(data)
			for i := 0; i < length+1; i++ {
				kvs[i].wg.Done()
				kvs[i].err = err
			}
		case <-db.closet:
			log.Debugf("close db setkv")
			return
		}
	}
}

func (db *FalDB) goDelKVBatch() {
	kvs := make([]*kvDelBatch, 8192)
	for {
		select {
		case v := <-db.delchan:
			if v == nil {
				continue
			}
			length := len(db.delchan)
			log.Debugf("goDelKVBatch length %d", length+1)
			data := make(map[string][]string)
			key := make([]string, 0)
			if v.key != "" {
				key = append(key, v.key)
				data[v.bucket] = key
			} else {
				data[v.bucket] = nil
			}
			kvs[0] = v
			for i := 1; i < length+1; i++ {
				v = <-db.delchan
				if v == nil {
					continue
				}
				if v.key == "" {
					data[v.bucket] = nil
				} else if keys, exist := data[v.bucket]; !exist || exist && keys != nil {
					data[v.bucket] = append(data[v.bucket], v.key)
				}
				kvs[i] = v
			}
			err := db.delBatch(data)
			for i := 0; i < length+1; i++ {
				kvs[i].wg.Done()
				kvs[i].err = err
			}
		case <-db.clodel:
			log.Debugf("close db delkv")
			return
		}
	}
}
