package client

import (
	"database/sql"
	// "fmt"
	_ "github.com/mattn/go-sqlite3"
)

type ChunkStore struct {
	db *sql.DB
}

// Create a new sqlite database to store file chunks
// The database will be created at the specified path
func NewChunkStore(dbPath string) (*ChunkStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

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
	_, err := cs.db.Exec(
		"INSERT INTO file_chunks (id, index_num, data) VALUES (?, ?, ?)",
		chunk.id, chunk.index, chunk.data,
	)
	return err
}

func Store(chunk *FileChunk) error {
	// fmt.Println("Storing chunk:", chunk.id)
	store, err := NewChunkStore("file_chunks.db")
	if err != nil {
		return err
	}
	defer store.Close()

	err = store.StoreChunk(chunk)
	if err != nil {
		return err
	}
	return nil
}

func (cs *ChunkStore) Retrieve(id string) (FileChunk, error) {
	var chunk FileChunk
	err := cs.db.QueryRow(
		"SELECT id, index_num, data FROM file_chunks WHERE id = ?", id,
	).Scan(&chunk.id, &chunk.index, &chunk.data)
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
		if err := rows.Scan(&chunk.id, &chunk.index, &chunk.data); err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}
	return chunks, nil
}

func (cs *ChunkStore) Close() error {
	return cs.db.Close()
}

/* func main() {
    store, err := NewChunkStore("file_chunks.db")
    if err != nil {
        fmt.Println("Error creating store:", err)
        return
    }
    defer store.Close()

    // Example usage
    chunk := FileChunk{
        id:    "chunk1",
        index: 1,
        data:  []byte("Hello, World!"),
    }

    err = store.Store(chunk)
    if err != nil {
        fmt.Println("Error storing chunk:", err)
        return
    }

    retrievedChunk, err := store.Retrieve("chunk1")
    if err != nil {
        fmt.Println("Error retrieving chunk:", err)
        return
    }

    fmt.Printf("Retrieved chunk: id=%s, index=%d, data=%s\n",
               retrievedChunk.id, retrievedChunk.index, string(retrievedChunk.data))

    chunks, err := store.RetrieveChunksByIndex(1, 10)
    if err != nil {
        fmt.Println("Error retrieving chunks by index:", err)
        return
    }

    fmt.Printf("Retrieved %d chunks\n", len(chunks))
} */
