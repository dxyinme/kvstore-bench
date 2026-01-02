package kv

import (
	"log"
	"time"

	"github.com/nutsdb/nutsdb"
)

type nutsdbStore struct {
	db       *nutsdb.DB
	needWait bool
}

const Bucket = "Bucket"

func newNutsDBCommon(_ string, options nutsdb.Options, needWait bool) (Store, error) {
	db, err := nutsdb.Open(options)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *nutsdb.Tx) error {
		return tx.NewKVBucket(Bucket)
	})
	if err != nil {
		return nil, err
	}

	return &nutsdbStore{db: db, needWait: needWait}, nil
}

func newNutsDB(path string) (Store, error) {
	options := nutsdb.DefaultOptions
	options.Dir = path
	options.EntryIdxMode = nutsdb.HintKeyAndRAMIdxMode
	options.SyncEnable = false
	options.HintKeyAndRAMIdxCacheSize = 0
	return newNutsDBCommon(path, options, false)
}

func newNutsDBMerge(path string) (Store, error) {
	options := nutsdb.DefaultOptions
	options.Dir = path
	options.EntryIdxMode = nutsdb.HintKeyAndRAMIdxMode
	options.SyncEnable = false
	options.HintKeyAndRAMIdxCacheSize = 0
	options.SegmentSize = 4 * nutsdb.MB
	options.MergeInterval = 2 * time.Minute
	return newNutsDBCommon(path, options, true)
}

func newNutsDBMergeV2(path string) (Store, error) {
	options := nutsdb.DefaultOptions
	options.Dir = path
	options.EntryIdxMode = nutsdb.HintKeyAndRAMIdxMode
	options.SyncEnable = false
	options.HintKeyAndRAMIdxCacheSize = 0
	options.SegmentSize = 4 * nutsdb.MB
	options.MergeInterval = 2 * time.Minute
	options.EnableMergeV2 = true
	return newNutsDBCommon(path, options, true)
}

func newNutsDBMmap(path string) (Store, error) {
	options := nutsdb.DefaultOptions
	options.Dir = path
	options.EntryIdxMode = nutsdb.HintKeyAndRAMIdxMode
	options.SyncEnable = false
	options.HintKeyAndRAMIdxCacheSize = 0
	options.RWMode = nutsdb.MMap
	return newNutsDBCommon(path, options, false)
}

func (n nutsdbStore) Put(key []byte, value []byte) error {
	return n.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(Bucket, key, value, nutsdb.Persistent)
	})
}

func (n nutsdbStore) Get(key []byte) ([]byte, error) {
	var (
		value []byte
	)
	err := n.db.View(func(tx *nutsdb.Tx) (e error) {
		value, e = tx.Get(Bucket, key)
		return
	})
	return value, err
}

func (n nutsdbStore) Delete(key []byte) error {
	return n.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(Bucket, key)
	})
}

func (n nutsdbStore) Close() error {
	if n.needWait {
		log.Println("wait for 10 minutes for merging")
		<-time.After(10 * time.Minute)
	}
	return n.db.Close()
}
