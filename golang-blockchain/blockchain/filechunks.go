package blockchain

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger"
	"github.com/google/uuid"
)

const (
	dbPath       = "./database"
	cryptoKey    = "teteteteteetesdsdsdsdsdt"
	EncryptedLoc = "./chunks/encrypted/"
	DecryptedLoc = "./chunks/decrypted/"
)

func ReadDir(dirname string) []os.FileInfo {

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func ReadFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func (chain *FileDB) GetEncryptedFiles(fileName string) []string {

	var chunks []string

	chain.Database.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(fileName)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			fmt.Println("K >", string(k))
			chunks = append(chunks, string(k))
		}
		return nil
	})
	return chunks
}

func (chain *FileDB) ConvertDecryptFiles(fileName string) {

	chunks := chain.GetEncryptedFiles(fileName)

	filename := DecryptedLoc + "final.txt"

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range chunks {

		databyte := ReadFile(EncryptedLoc + f)
		data := DecryptFile(string(databyte))
		length, err := io.WriteString(file, data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("file length is ", length, data)
	}
	defer file.Close()
	databyte := ReadFile(filename)
	fmt.Println("Actual data of saved file is ", string(databyte))
}

func CreateChunksAndEncrypt(fileNameE string) {

	split := 4
	file, err := os.Open(fileNameE)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	texts := make([]string, 0)
	for scanner.Scan() {
		text := scanner.Text()
		texts = append(texts, text)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	lengthPerSplit := len(texts) / split
	for i := 0; i < split; i++ {
		if i+1 == split {
			chunkTexts := texts[i*lengthPerSplit:]
			fmt.Println(chunkTexts)
			writefile(EncryptFile(strings.Join(chunkTexts, "\n")), i+1, fileNameE)
		} else {
			chunkTexts := texts[i*lengthPerSplit : (i+1)*lengthPerSplit]
			fmt.Println(chunkTexts)
			writefile(EncryptFile(strings.Join(chunkTexts, "\n")), i+1, fileNameE)
		}
	}
}

func writefile(data string, index int, fileNameE string) {

	fileChunk := fileNameE + "-" + strconv.Itoa(index) + "-chunks-" + uuid.New().String() + ".txt"
	file, err := os.Create("./chunks/encrypted/" + fileChunk)

	var chunks File

	chunks.Chunkname = fileChunk
	chunks.Filename = fileNameE // filename should be unique
	chunks.Ownername = "Amandeep"
	chunks.NodeAddress = "node1"
	chunks.BlockHash = []byte("SomeHash")
	chunks.ChuckIndex = index

	inst := SaveFileInfo(chunks)
	inst.Database.Close()

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	file.WriteString(data)
}

type File struct {
	Chunkname   string
	Filename    string
	Ownername   string
	NodeAddress string
	BlockHash   []byte
	ChuckIndex  int
}

type FileDB struct {
	Database *badger.DB
}

func GetDBinstacnce() *FileDB {
	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	blockchain := FileDB{Database: db}
	return &blockchain
}

func SaveFileInfo(chunk File) *FileDB {
	db := GetDBinstacnce()

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(chunk)
	key := chunk.Chunkname
	fmt.Println("key is ", key)
	db.Database.Update(func(txn *badger.Txn) error {
		err2 := txn.Set([]byte(key), reqBodyBytes.Bytes())
		return err2
	})

	blockchain := FileDB{Database: db.Database}
	return &blockchain
}

func (chain *FileDB) GetChunkByKey(key string) string {

	var file []byte

	chain.Database.View(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(key))
		if err != nil {
			fmt.Println("Key not found.")
			return err
		}
		file, err = item.Value()
		fmt.Println("Item: ", string(file))
		Handle(err)
		return err
	})

	return string(file)
}

func (chain *FileDB) GetChunksByPrefix(prefix string) []string {

	var chunks []string

	chain.Database.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(prefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			fmt.Println("K >", string(k))
			chunks = append(chunks, string(k))
		}
		return nil
	})
	return chunks
}
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
