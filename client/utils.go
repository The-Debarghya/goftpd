package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)


func sendFile(conn net.Conn, filename string) {
	inputFile, err := os.Open(CLIENTROOT + "/" + filename)
	if err != nil {
		fmt.Println("ERROR: ", err.Error())
		conn.Write([]byte("-1"))
		return
	} else {
		stats,_ := inputFile.Stat()
		//send file Size
		conn.Write([]byte(strconv.FormatInt(stats.Size(),10)))
	}
	defer inputFile.Close()
	//time.Sleep(100*time.Microsecond)
	buffer := make([]byte, BUFFSIZE)
	for {
		_, err := inputFile.Read(buffer)
		if err == io.EOF{
			break
		}
		conn.Write(buffer)
	}
	fmt.Println("File Uploaded Successfully")
}

func saveFile(conn net.Conn, filename string, fileSize int64) {
	outputFile, err := os.Create(CLIENTROOT + filename)
	if err != nil {
		fmt.Println(err)
	}
	defer outputFile.Close()
	var fileSizeReceived int64
	for {
		if (fileSize - fileSizeReceived) < BUFFSIZE {
			io.CopyN(outputFile, conn, (fileSize - fileSizeReceived))
			conn.Read(make([]byte, (fileSizeReceived + BUFFSIZE) - fileSize))
			break
		}
		io.CopyN(outputFile, conn, BUFFSIZE)
		fileSizeReceived += BUFFSIZE
	}
	fmt.Println("File Downloaded Successfully")
}