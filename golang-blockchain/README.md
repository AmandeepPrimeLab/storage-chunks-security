# BASIC FLOW

This repo is to create shards of files and save its information into Manifest file

## INSTALLATION SETUP

- Just install the Go https://go.dev/dl/ link as per your OS
- Then clone the project and go insude the golang-blockchain folder
- Then delele the go.sum file and run 'go get' command

## How to use
Go inside the golang-blockchain folder then running ```go run main.go``` 

## Create Chunks and encrypt the file 
```
/enryptFile/{filename}: Method is GET

```
## Get Chunks By Key 
```
/getChunkByKey/{key}: Method is GET

```
## Get Chunk By Filename 
```
/getChunkByFilename/{filename}: Method is GET

```
## Get Chunks and decrypt the file 
```
/deryptFile/{filename}: Method is GET

```

#### We have 2 files right now inside this repo:
```
demo.txt
test.txt

```