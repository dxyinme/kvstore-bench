package kv

import (
	"github.com/nutsdb/nutsdb"
)

type nutsdbStore struct {
	db *nutsdb.DB
}

const Bucket = "Bucket"

func newNutsDBCommon(path string, options nutsdb.Options) (Store, error) {
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

	return &nutsdbStore{db: db}, nil
}

func newNutsDB(path string) (Store, error) {
	options := nutsdb.DefaultOptions
	options.Dir = path
	options.EntryIdxMode = nutsdb.HintKeyAndRAMIdxMode
	options.SyncEnable = false
	options.HintKeyAndRAMIdxCacheSize = 0
	return newNutsDBCommon(path, options)
}

func newNutsDBMmap(path string) (Store, error) {
	options := nutsdb.DefaultOptions
	options.Dir = path
	options.EntryIdxMode = nutsdb.HintKeyAndRAMIdxMode
	options.SyncEnable = false
	options.HintKeyAndRAMIdxCacheSize = 0
	options.RWMode = nutsdb.MMap
	return newNutsDBCommon(path, options)
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
	err := n.db.View(func(tx *nutsdb.Tx) error {
		value, _ = tx.Get(Bucket, key)
		return nil
	})
	return value, err
}

func (n nutsdbStore) Delete(key []byte) error {
	return n.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(Bucket, key)
	})
}

func (n nutsdbStore) Close() error {
	return n.db.Close()
}
