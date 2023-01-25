package main

import (
	"candy_db/client"
	"candy_db/common"
	"candy_db/common/util"
	"context"
	"fmt"
)

func after(ctx context.Context, req client.CallReq) error {
	fmt.Println("after")
	return nil
}

func before(ctx context.Context, req client.CallReq) error {
	fmt.Println("before")
	return nil
}

func main() {
	ctx := context.Background()
	// test i64
	db, err := client.NewCandyDBClient[string, int64](
		client.WithBeforeCall[string, int64](before),
		client.WithAfterCall[string, int64](after),
		client.WithMemTable[string, int64](
			common.CustomCache[string](util.NewCacheService[string](util.LFU)),
		),
	)
	if err != nil {
		panic(err)
	}
	err = db.Add(ctx, "key", int64(123))
	if err != nil {
		panic(err)
	}
	ret, err := db.Find(ctx, "key")
	fmt.Println(ret.Value)
}
