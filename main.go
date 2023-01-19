package main

import (
	"candy_db/client"
	"encoding/json"
	"fmt"
)

func main() {

	// test i64
	db := client.NewCandyDBClient[int64]()
	key, _ := json.Marshal("key")
	db.Add(key, int64(1))
	ret := db.Find(key)
	fmt.Println(ret.Assert())

	// test string
	db2 := client.NewCandyDBClient[string]()
	db2.Add(key, "test candy db")
	ret2 := db2.Find(key)
	fmt.Println(ret2.Assert())

	// test struct
	type Test struct {
		name string
		age  int
	}

	db3 := client.NewCandyDBClient[Test]()
	db3.Add(key, Test{
		name: "test db",
		age:  123,
	})
	ret3 := db3.Find(key)
	fmt.Println(ret3.Assert())

}
