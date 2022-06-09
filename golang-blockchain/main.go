package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"golang-blockchain/blockchain"

	"github.com/gorilla/mux"
)

type Message struct {
	Status string
}

func checkServer(w http.ResponseWriter, r *http.Request) {

	var newEmployee = Message{Status: "Server is in running state"}
	setupHeader(w)
	json.NewEncoder(w).Encode(newEmployee)
}

func setupHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
}

func enryptFile(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]
	setupHeader(w)
	blockchain.CreateChunksAndEncrypt(filename)
	inst := blockchain.GetDBinstacnce()
	data := inst.GetEncryptedFiles(filename)
	inst.Database.Close()
	json.NewEncoder(w).Encode(data)
}

func getChunkByKey(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	setupHeader(w)
	inst := blockchain.GetDBinstacnce()
	data := inst.GetChunkByKey(key)
	inst.Database.Close()
	w.Write([]byte(data))
}

func getChunkByFilename(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]
	setupHeader(w)
	inst := blockchain.GetDBinstacnce()
	data := inst.GetChunksByPrefix(filename)
	inst.Database.Close()
	json.NewEncoder(w).Encode(data)
}

func deryptFile(w http.ResponseWriter, r *http.Request) {

	filename := mux.Vars(r)["filename"]
	setupHeader(w)
	inst := blockchain.GetDBinstacnce()
	inst.ConvertDecryptFiles(filename)
	files := blockchain.ReadFile(blockchain.DecryptedLoc + "final.txt")
	inst.Database.Close()
	w.Write(files)
	os.Remove(blockchain.DecryptedLoc + "final.txt")
}

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", checkServer)
	router.HandleFunc("/enryptFile/{filename}", enryptFile).Methods("GET")
	router.HandleFunc("/getChunkByKey/{key}", getChunkByKey).Methods("GET")
	router.HandleFunc("/getChunkByFilename/{filename}", getChunkByFilename).Methods("GET")
	router.HandleFunc("/deryptFile/{filename}", deryptFile).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
