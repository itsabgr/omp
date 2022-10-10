package db

import (
	"errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/itsabgr/omp/internal/model"
	"github.com/itsabgr/omp/internal/utils"
	"io"
)

var ErrExists = errors.New("db: exists")
var ErrInvalidChunkID = errors.New("db: invalid chunk id")
var ErrInvalidChunkSize = errors.New("db: invalid chunk-size")
var ErrNotFound = errors.New("db: not found")

type DB struct {
	badger *badger.DB
}

func New(badger *badger.DB) *DB {
	return &DB{badger: badger}
}

var (
	tableImages = []byte{'I'}
	tableChunks = []byte{'C'}
)

func getRecord(tx *badger.Txn, key []byte) (*badger.Item, error) {
	item, err := tx.Get(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if item.IsDeletedOrExpired() {
		return nil, ErrNotFound
	}
	return item, nil

}
func (db *DB) CreateImage(image *model.Image) error {
	tx := db.badger.NewTransaction(true)
	defer tx.Discard()
	key := utils.Concat(tableImages, []byte(image.Sha256))
	if _, err := getRecord(tx, key); err == nil {
		return ErrExists

	}
	if err := tx.Set(
		key,
		utils.Must(image.Marshal()),
	); err != nil {
		return err
	}
	return tx.Commit()
}
func (db *DB) PutChunk(chunk *model.Chunk) error {
	tx := db.badger.NewTransaction(true)
	defer tx.Discard()
	key := utils.Concat(tableImages, []byte(chunk.Image))
	{ //check image
		image := new(model.Image)
		if item, err := getRecord(tx, key); err != nil {
			return err
		} else {
			if err := item.Value(func(val []byte) error {
				return image.Unmarshal(val)
			}); err != nil {
				return err
			}
		}
		if chunk.Id > image.Size/image.ChunkSize {
			return ErrInvalidChunkID
		}
		if chunk.Size > image.ChunkSize {
			return ErrInvalidChunkSize
		}
	}
	if err := tx.Set(
		utils.Concat(tableChunks, []byte(chunk.Image), utils.Int64ToBE(uint64(chunk.Id))),
		[]byte(chunk.Data),
	); err != nil {
		return err
	}
	return tx.Commit()
}
func (db *DB) WriteImageTo(sha256 string, dst io.Writer) error {
	tx := db.badger.NewTransaction(false)
	defer tx.Discard()
	iterator := tx.NewIterator(badger.IteratorOptions{
		Prefix: utils.Concat(tableChunks, []byte(sha256)),
	})
	defer iterator.Close()
	anyWrite := false
	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		anyWrite = true
		if err := iterator.Item().Value(func(val []byte) error {
			_, err := dst.Write(val)
			return err
		}); err != nil {
			return err
		}
	}
	if !anyWrite {
		return ErrNotFound
	}
	return nil
}
