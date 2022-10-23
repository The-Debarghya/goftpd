package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	//"syscall"
	//"golang.org/x/crypto/ssh/terminal"
)


const (
	HOST = "localhost"
	PORT = "9999"
	BUFFSIZE = 5120
)
var CLIENTROOT = "/ftpfiles/"

func init(){
	cdir, _ := os.Getwd()
	splits := strings.Split(cdir, "/")
	CLIENTROOT = strings.Join(splits[:len(splits)-1],"/") + CLIENTROOT
}

func main()  {
	addr, err := net.ResolveTCPAddr("tcp", HOST+":"+PORT)
	if err != nil {
    	log.Fatal(err)
	}
	server, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
    	log.Fatal(err)
	}
	log.Println("FTP client is connected to", server.RemoteAddr())
	authorizeClient(server)
	clientHandler(server)
}

func clientHandler(conn net.Conn)  {
	stdreader := bufio.NewReader(os.Stdin)
	buffer := make([]byte, BUFFSIZE)
	for{
		fmt.Printf("ftp> ")
		cmd, _  := stdreader.ReadString('\n')
		cmd = strings.Trim(cmd, "\n")
		cmdList := strings.Split(cmd, " ")
		switch strings.ToLower(cmdList[0]) {
		case "get":
			if len(cmdList) == 1 {
				fmt.Println("Invalid command\nType 'help' for more info.")
				continue
			}
			conn.Write([]byte(cmd))
			n, _ := conn.Read(buffer)
			fileSize , err := strconv.ParseInt(string(buffer[:n]), 10, 64)
			if err != nil{
				fmt.Println("ERROR: ", string(buffer[:n]))
				continue
			}
			saveFile(conn, cmdList[1], fileSize)
		case "put":
			if len(cmdList) == 1 {
				fmt.Println("Invalid command\nType 'help' for more info.")
				continue
			}
			conn.Write([]byte(cmd))
			sendFile(conn, cmdList[1])
		case "del":
			if len(cmdList) == 1 {
				fmt.Println("Invalid command\nType 'help' for more info.")
				continue
			}
			conn.Write([]byte(cmd))
			n, _ := conn.Read(buffer)
			fmt.Println(string(buffer[:n]))
		case "cd":
			if len(cmdList) == 1 {
				fmt.Println("Invalid command\nType 'help' for more info.")
				continue
			}
			conn.Write([]byte(cmd))
			n, _ := conn.Read(buffer)
			fmt.Println(string(buffer[:n]))
		case "pwd":
			conn.Write([]byte(cmd))
			n, _ := conn.Read(buffer)
			fmt.Println(string(buffer[:n]))
		case "ls":
			conn.Write([]byte(cmd))
			n, _ := conn.Read(buffer)
			fmt.Println(string(buffer[:n]))
		case "cat":
			if len(cmdList) == 1 {
				fmt.Println("Invalid command\nType 'help' for more info.")
				continue
			}
			conn.Write([]byte(cmd))
			tmp := make([]byte, BUFFSIZE)     
    		for {
        		n, err := conn.Read(tmp)
        		if err != nil {
            		if err != io.EOF {
                		fmt.Println("Read Error:", err)
            		}
            		break
       			}
				if string(tmp[:n]) == "EOF"{
					break
				}
        		fmt.Println(string(tmp[:n]))
        		//buffer = append(buffer, tmp[:n]...)
    		}
			//fmt.Println(string(buffer))
		case "quit":
			fmt.Println("Logging out\nGoodbye!")
			conn.Write([]byte(cmd))
			return
		case "help":
			conn.Write([]byte(cmd))
			fmt.Print("Supported Commands:\n\nget <filename>:\tDownload remote file\nput <filename>:\tUpload local file\ndel <filename/dir>:\tDelete remote dir/file\ncd <dirname>:\tChange directory\nls:\tList all file(s)/directories in current remote directory\npwd:\tGet the present working directory name\ncat <filename>:\tDisplay file contents\nquit:\tClose the connection\nhelp:\tShow this menu.\n\n")
		case "":
			continue
		default:
			fmt.Println("Invalid command! To know all supported commands use 'help'")			
		}
	}
}

func authorizeClient(conn net.Conn)  {
	stdreader := bufio.NewReader(os.Stdin)
	buffer := make([]byte, BUFFSIZE)
	n, _ := conn.Read(buffer)

	fmt.Println(string(buffer[:n]))
	fmt.Println("Authentication Required To transfer File(s)")

	auth := false
	for !auth{
		fmt.Printf("Username: ")
		user, _ := stdreader.ReadString('\n')
		fmt.Printf("Password: ")
		pass, _ := stdreader.ReadString('\n')
		//strings.TrimSpace(pass)
		conn.Write([]byte(user))
		conn.Write([]byte(pass))
		n, _ := conn.Read(buffer)
		if string(buffer[:n]) == "1"{
			fmt.Println("Authentication Successful!")
			auth = true
			break
		}
		fmt.Println("Invalid Credentials!")
	}
}