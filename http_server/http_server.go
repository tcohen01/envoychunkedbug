package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/chunked", HandleChunked)
	http.HandleFunc("/normal", HandleNormal)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

const DefaultSize = 1 * 1024     // 1 KB
const MaxSize = 1024 * 1024      // 1 MiB
const DefaultChunkSize = 1024    // 1 KB
const MaxChunkSize = 1024 * 1024 // 1 MB

func HandleChunked(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Panic("expected ResponseWriter to be Flusher")
	}

	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	sizeParam := r.URL.Query().Get("size")
	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		size = DefaultSize
	}
	if size > MaxSize {
		size = MaxSize
	}
	chunkSizeParam := r.URL.Query().Get("chunksize")
	chunkSize, err := strconv.Atoi(chunkSizeParam)
	if err != nil {
		chunkSize = DefaultChunkSize
	}
	if chunkSize > MaxChunkSize {
		chunkSize = MaxChunkSize
	}
	log.Printf("sending chunked response with %d bytes and %d bytes per chunk", size, chunkSize)
	chunks := 0
	chunkBytes := make([]byte, chunkSize)
	for size > 0 {
		randomBase64Bytes(chunkBytes)
		writeSize := chunkSize
		if size < chunkSize {
			writeSize = size
		}
		_, err := w.Write(chunkBytes[:writeSize])
		if err != nil {
			log.Panicf("failed to write chunk %d bytes: %v", writeSize, err)
		}
		flusher.Flush()
		chunks++
		size -= writeSize
	}
	log.Printf("%d chunks sent", chunks)
}

func HandleNormal(w http.ResponseWriter, r *http.Request) {
	sizeParam := r.URL.Query().Get("size")
	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		size = DefaultSize
	}
	if size > MaxSize {
		size = MaxSize
	}
	log.Printf("sending response with %d bytes", size)
	responseBytes := make([]byte, size)
	randomBase64Bytes(responseBytes)
	_, err = w.Write(responseBytes)
	if err != nil {
		log.Panicf("failed to write normal response of %d bytes: %v", size, err)
	}
}

var base64Charset = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/")

func randomBase64Bytes(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = base64Charset[rand.Intn(len(base64Charset))]
	}
}
