package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func readDir(conn net.Conn)  {
	files, err := os.ReadDir(ROOTDIR)
	if err != nil{
		conn.Write([]byte(err.Error()))
		log.Println(err)
		return
	}
	fileInfo := ""
	for _, file := range files{
		fi, _ := file.Info()
		fileInfo += fi.Mode().String() + " " + fmt.Sprintf("%d", fi.Size()) + " " + file.Name() + "\n"
	}
	conn.Write([]byte(fileInfo))
}

func chDir(conn net.Conn, dir string)  {
	const BASEDIR = "/media/debarghya/927268E17268CC13/Home/Desktop/Labs/Computer-Networks/Assn7/goftpd/ftp"
	rel, err := filepath.Rel(BASEDIR, ROOTDIR+"/"+dir)
	if err != nil{
		conn.Write([]byte(err.Error()))
		log.Println(err)
		return
	}
	if strings.Contains(rel, ".."){
		log.Println("cd ", rel)
		conn.Write([]byte("Changed directory successfuly!"))
		return
	}else{
		ROOTDIR,_ = filepath.Abs(ROOTDIR + "/" + dir)
		log.Println("cd ", ROOTDIR)
		conn.Write([]byte("Changed directory successfuly!"))
		return
	}
}

func getFile(conn net.Conn, name string)  {
	inputFile, err := os.Open(ROOTDIR + "/" + name)
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}else {
		stats,_ := inputFile.Stat()
		//send file Size
		conn.Write([]byte(strconv.FormatInt(stats.Size(), 10)))
	}
	defer inputFile.Close()
	buffer := make([]byte, BUFFSIZE)
	for {
		_, err := inputFile.Read(buffer)
		if err == io.EOF{
			break
		}
		conn.Write(buffer)
	}
	log.Println("File Transfered Successfully!")
}

func putFile(conn net.Conn, name string, filesize int64)  {
	outputFile, err := os.Create(ROOTDIR + "/" + name)
	if err != nil {
		log.Println(err)
	}
	defer outputFile.Close()
	var fileSizeReceived int64
	for {
		if (filesize - fileSizeReceived) < BUFFSIZE {
			io.CopyN(outputFile, conn, (filesize - fileSizeReceived))
			conn.Read(make([]byte, (fileSizeReceived+BUFFSIZE) - filesize))
			break
		}
		io.CopyN(outputFile, conn, BUFFSIZE)
		fileSizeReceived += BUFFSIZE
	}
	log.Println("File Received Successfully, saved at:", outputFile)
}

func delDir(conn net.Conn, name string)	{
	log.Println(ROOTDIR + "/" +name)
	err := os.Remove(ROOTDIR + "/" +name)
	if err != nil{
		conn.Write([]byte(err.Error()))
		log.Println(err)
		return
	}
	conn.Write([]byte("File Deleted Successfully"))
}

func showFileContents(conn net.Conn, name string)  {
	filename, err := os.Open(ROOTDIR + "/" + name)
	log.Println("cat ", name)
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}else {
		scanner := bufio.NewScanner(filename)
		for scanner.Scan() {
			line := scanner.Bytes()
			//log.Println(string(line))
			conn.Write(line)
			time.Sleep(200*time.Millisecond)
		}
		//log.Println("Out of loop")
		conn.Write([]byte("EOF"))
	}
	defer filename.Close()
}