package client

import (
	"database/sql"
	"log"
	"strings"
	"time"

	_ "modernc.org/sqlite" // Change to the new SQLite driver
)

type ChunkStore struct {
	db *sql.DB
}

var store *ChunkStore

func init() {
	var err error
	store, err = NewChunkStore("file_chunks.db")
	if err != nil {
		log.Fatal(err)
	}
}

// Create a new sqlite database to store file chunks
// The database will be created at the specified path
func NewChunkStore(dbPath string) (*ChunkStore, error) {
	db, err := sql.Open("sqlite", dbPath) // Use "sqlite" for the modernc driver
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA busy_timeout = 5000;") // Set timeout to 5 seconds
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(1)

	// Create table if it doesn't exist
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS file_chunks (
            id TEXT PRIMARY KEY,
            index_num INTEGER,
            data BLOB
        )
    `)
	if err != nil {
		return nil, err
	}

	return &ChunkStore{db: db}, nil
}

func (cs *ChunkStore) StoreChunk(chunk *FileChunk) error {
	var existingID string
	err := cs.db.QueryRow("SELECT id FROM file_chunks WHERE id = ?", chunk.Id).Scan(&existingID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existingID == "" {
		_, err = cs.db.Exec(
			"INSERT INTO file_chunks (id, index_num, data) VALUES (?, ?, ?)",
			chunk.Id, chunk.Index, chunk.Data,
		)
	} else {
		_, err = cs.db.Exec(
			"UPDATE file_chunks SET index_num = ?, data = ? WHERE id = ?",
			chunk.Index, chunk.Data, chunk.Id,
		)
	}
	return err
}

func Store(chunk *FileChunk) error {

	err := store.StoreChunk(chunk)
	if err != nil {
		store.StoreChunkWithRetry(chunk, 5)
	}
	return nil
}

func (cs *ChunkStore) Retrieve(id string) (FileChunk, error) {
	var chunk FileChunk
	err := cs.db.QueryRow(
		"SELECT id, index_num, data FROM file_chunks WHERE id = ?", id,
	).Scan(&chunk.Id, &chunk.Index, &chunk.Data)
	return chunk, err
}

func (cs *ChunkStore) RetrieveChunksByIndex(startIndex, endIndex int) ([]FileChunk, error) {
	rows, err := cs.db.Query(
		"SELECT id, index_num, data FROM file_chunks WHERE index_num BETWEEN ? AND ? ORDER BY index_num",
		startIndex, endIndex,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []FileChunk
	for rows.Next() {
		var chunk FileChunk
		if err := rows.Scan(&chunk.Id, &chunk.Index, &chunk.Data); err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}
	return chunks, nil
}

func (cs *ChunkStore) Close() error {
	return cs.db.Close()
}

func (cs *ChunkStore) StoreChunkWithRetry(chunk *FileChunk, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = cs.StoreChunk(chunk)
		if err == nil {
			return nil
		}
		if strings.Contains(err.Error(), "database is locked") {
			time.Sleep(time.Millisecond * 100 * time.Duration(i+1))
			continue
		}
		return err
	}
	return err
}
