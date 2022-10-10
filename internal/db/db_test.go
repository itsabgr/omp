package db

import (
	"encoding/hex"
	"github.com/dgraph-io/badger/v3"
	"github.com/itsabgr/omp/internal/model"
	"github.com/itsabgr/omp/internal/utils"
	"strings"
	"testing"
)

func TestCorrectness(t *testing.T) {
	db := New(utils.Must(badger.Open(badger.DefaultOptions("").WithLogger(nil).WithInMemory(true))))
	defer db.badger.Close()
	data := hex.EncodeToString(utils.RandBytes(100))
	chunkSize := uint(23)
	sum := utils.Sum([]byte(data))
	utils.Throw(db.CreateImage(&model.Image{ChunkSize: chunkSize, Size: uint(len(data)), Sha256: sum}))
	for i := uint(0); i <= uint(len(data))/chunkSize; i++ {
		offset := i * chunkSize
		chunk := data[offset:utils.Min(int(offset+chunkSize), len(data))]
		utils.Throw(db.PutChunk(&model.Chunk{Id: i, Size: chunkSize, Data: chunk, Image: sum}))
	}
	builder := new(strings.Builder)
	utils.Throw(db.WriteImageTo(sum, builder))
	if sum != utils.Sum([]byte(builder.String())) {
		t.FailNow()
	}
}
