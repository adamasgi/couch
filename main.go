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

type CouchDB struct {
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
	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}

	defer res.Body.Close()

	fmt.Printf("%s", reqBody)
}

func (c Couch) Get(path string) []byte {
	res, err := c.Client.Get(c.Address + "/" + path)
	if err != nil {
		fmt.Printf("%s", err)
	}
	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}

	return reqBody
}

func (c Couch) NewDB(DBname string) []byte {
	req, err := http.NewRequest("PUT", c.Address+"/"+DBname, nil)
	if err != nil {
		fmt.Printf("%s", err)
	}

	res, err := c.Client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
	}

	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}

	return reqBody
}

func (d CouchDB) Add(ID string, body map[string]interface{}) []byte {
	docBody, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("%s", err)
	}

	req, err := http.NewRequest("PUT", d.C.Address+"/"+d.Name+"/"+ID, bytes.NewBuffer(docBody))
	if err != nil {
		fmt.Printf("%s", err)
	}

	res, err := d.C.Client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
	}

	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}

	return reqBody
}

func (d CouchDB) Delete(ID string) []byte {
	req, err := http.NewRequest("DELETE", d.C.Address+"/"+d.Name+"/"+ID, nil)
	if err != nil {
		fmt.Printf("%s", err)
	}

	res, err := d.C.Client.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
	}

	reqBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}

	return reqBody
}

func InitCouch() Couch {
	c := new(Couch)
	Jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Printf("%s", err)
	}

	c.Client = http.Client{
		Jar: Jar,
	}
	c.Address = "http://100.102.100.49:5984"
	c.Login()
	return *c
}

func InitCouchDB(name string) CouchDB {
	c := InitCouch()
	d := new(CouchDB)
	d.C = c
	d.Name = name
	return *d
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
	testesDB := CouchDB{couch, "testes"}
	body := loadFile("./test.json")
	idd := fmt.Sprintf("%s", body["_id"])
	doc := testesDB.Add(idd, body)
	fmt.Printf("%s", doc)
	ge := couch.Get("testes/roa")
	fmt.Printf("%s", ge)
}
