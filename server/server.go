package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/The-Debarghya/goftpd/server/creds"
)

const (
	PORT = "9999"
	BUFFSIZE = 5120
)
var ROOTDIR = "/ftp"

func init()  {
	currdir, _ := os.Getwd()
	split := strings.Split(currdir, "/")
	ROOTDIR = strings.Join(split[:len(split)-1], "/") + ROOTDIR
}

func main()  {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:"+PORT)
	if err != nil {
    	log.Fatal(err)
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil{
		log.Fatal(err)
	}
	defer server.Close()
	log.Println("FTP server is up and running at", addr)
	for{
		conn, err := server.AcceptTCP()
		if err != nil {
    		log.Fatal("Client connection dropped!")
		}
		go serverHandler(conn)
	}

}

func getCredentialsList() []creds.Credentials {
	//var data creds.Data
	data := creds.ImportFromJSON()
	return data.Credslist
}

func getKeyStr() string {
	//var data creds.Data
	data := creds.ImportFromJSON()
	return data.Key
}

func decrypt(encryptedpass string, Key string) (decryptedpass string){
	key, _ := hex.DecodeString(Key)
	enc, _ := hex.DecodeString(encryptedpass)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Error initializing AES cipher!")
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("Error initializing AES GCM!")
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println("Error in decryption!")
	}
	return string(plaintext)
}

func authenticateClient(conn net.Conn){
	reader := bufio.NewScanner(conn)
	auth := false
	conn.Write([]byte("Connection with Server Established!"))
	Creds := getCredentialsList()
	Key := getKeyStr()
	for !auth{
		reader.Scan()
		user := reader.Text()
		reader.Scan()
		pass := reader.Text()
		for _, credential := range Creds{
			if credential.Username == user && pass == decrypt(credential.Password, Key){
				auth = true
				log.Println("Client verified successfully, as user", user)
				break
			}
		}
		if auth{
			conn.Write([]byte("1"))
			break
		}
		conn.Write([]byte("0"))
	}
}

func serverHandler(conn net.Conn)  {
	defer conn.Close()
	authenticateClient(conn)
	buffer := make([]byte, BUFFSIZE)
	for{
		n, _ := conn.Read(buffer)
		cmd := string(buffer[:n])
		cmdList := strings.Split(cmd, " ")
		switch strings.ToLower(cmdList[0]){
		case "put":
			log.Println("Put File", cmdList[1])
			n, _ := conn.Read(buffer)
			fileSize, err := strconv.ParseInt(string(buffer[:n]), 10, 64)
			log.Println(fileSize)
			if err != nil || fileSize == -1{
				log.Println("Error while parsing file")
				continue
			}
			putFile(conn, cmdList[1], fileSize)

		case "get":
			log.Println("Get File", cmdList[1])
			getFile(conn, cmdList[1])
		case "ls":
			log.Println("ls")
			readDir(conn)
		case "pwd":
			log.Println("pwd")
			conn.Write([]byte(ROOTDIR))
		case "cd":
			chDir(conn, cmdList[1])
		case "cat":
			showFileContents(conn, cmdList[1])
		case "del":
			delDir(conn, cmdList[1])
		case "help":
			log.Println("Help Menu Display")
		case "quit":
			log.Println("Client Closed the Connection")
			return
		case "":
			log.Println("Error in buffer read!")
			return
		default:
			log.Println("Invalid command", cmd)
			continue
		}
	}

}