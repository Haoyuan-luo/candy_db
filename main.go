package main

import (
	"candy_db/common"
	"fmt"
)

func main() {
	nodeList := common.NodeList{
		common.NewNode(1, []byte("12345678")),
		common.NewNode(2, []byte("23456789")),
		common.NewNode(3, []byte("34567890")),
		common.NewNode(4, []byte("45678901")),
	}

	s := common.NewSkipList(nodeList)
	fmt.Println(s)
}
