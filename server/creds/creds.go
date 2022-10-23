package creds

import (
	"encoding/json"
	"os"
	"log"
)

type Credentials struct{
	Username string `json:"Username"`
	Password string	`json:"Password"`
}

//var CredentialsList []Credentials

type Data struct{
	Credslist []Credentials
	Key string `json:"Key"`
}

func ImportFromJSON() Data {
	var payload Data
	content, err1 := os.ReadFile("./creds/creds.json")
	if err1 != nil {
        log.Fatal("Error when opening file: ", err1)
    }
	err := json.Unmarshal(content, &payload)
	if err != nil {
        log.Fatal("Error when opening file: ", err)
    }
	return payload
}