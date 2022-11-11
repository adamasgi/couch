package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
)

type Couch struct {
	Client  http.Client
	Address string `default:"http://100.102.100.49:5984"`
}

type DB struct {
	C    Couch
	Name string
}

func (c Couch) Login() {
	user := os.Getenv("couch_user")
	password := os.Getenv("couch_password")
	body, _ := json.Marshal(map[string]string{
		"name":     user,
		"password": password,
	})
	add := c.Address + "/_session"
	res, err := c.Client.Post(add, "application/json", bytes.NewBuffer(body))

	if err != nil {
		fmt.Printf("%s", err)
	}
	reqBody, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	fmt.Printf("%s", reqBody)
}

func (c Couch) Get(path string) []byte {
	res, _ := c.Client.Get(c.Address + "/" + path)
	reqBody, _ := ioutil.ReadAll(res.Body)
	return reqBody
}

func (c Couch) NewDB(DBname string) []byte {
	req, _ := http.NewRequest("PUT", c.Address+"/"+DBname, nil)
	res, _ := c.Client.Do(req)
	reqBody, _ := ioutil.ReadAll(res.Body)
	return reqBody
}

func (d DB) Doc(DB string, ID string, body map[string]interface{}) []byte {
	docBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", d.C.Address+"/"+DB+"/"+ID, bytes.NewBuffer(docBody))
	res, _ := d.C.Client.Do(req)
	reqBody, _ := ioutil.ReadAll(res.Body)
	return reqBody
}

func InitCouch() Couch {
	c := new(Couch)
	Jar, _ := cookiejar.New(nil)
	c.Client = http.Client{
		Jar: Jar,
	}
	c.Address = "http://100.102.100.49:5984"
	c.Login()
	return *c
}

func loadFile(filepath string) map[string]interface{} {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}
	var payload map[string]interface{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	id, ok := payload["_id"]
	if ok {
		fmt.Printf("%s ID", id)
	} else {
		log.Fatal("No id field")
	}
	return payload
}

func main() {
	couch := InitCouch()
	//couch.NewDB("testes")
	//fmt.Printf("%s", db)
	testesDB := DB{couch, "testes"}
	body := loadFile("./test.json")
	idd := fmt.Sprintf("%s", body["_id"])
	doc := testesDB.Doc("testes", idd, body)
	fmt.Printf("%s", doc)
	ge := couch.Get("testes/roa")
	fmt.Printf("%s", ge)
}
