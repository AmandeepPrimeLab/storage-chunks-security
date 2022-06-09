# BASIC FLOW

This repo is to create shards of files and save its information into Manifest file

## INSTALLATION SETUP

- Just install the Go https://go.dev/dl/ link as per your OS
- Then clone the project and go insude the golang-blockchain folder
- Then delele the go.sum file and run 'go get' command
- Then run 'go run main.go' command

- Then we have a list of endpoints that 
    - "/enryptFile/{filename}" Methods("GET")
	- "/getChunkByKey/{key}" .Methods("GET")
	- "/getChunkByFilename/{filename}" Methods("GET")
	- "/deryptFile/{filename}" .Methods("GET")
