package main

import (
	"fmt"
	"log"
	"v3/bolt"
)
func main() {

	//打开数据库, 第一参数：数据库的名字，第二个参数：权限
	db, err := bolt.Open("test.db", 0600, nil)
	defer db.Close()
	if err != nil{
		log.Panic(err)
	}

	//写数据
	db.Update(func(tx *bolt.Tx) error {
		//找到抽屉bucket,(如果没有就创建)
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil {
			//没有抽屉，就需要我们创建
			bucket, err = tx.CreateBucket([]byte("b1"))
			if err != nil {
				log.Panic("创建bucket(b1)失败")
			}
		}
		//写入数据
		bucket.Put([]byte("1"),[]byte("hello"))

		return nil


	})

	//读数据
	db.View(func(tx *bolt.Tx) error {
		//找到抽屉，没有的话直接报错退出
		bucket :=  tx.Bucket([]byte("b1"))
		if bucket == nil{
			log.Panic("bucket 不应该为空，请检查！！！")
		}

		//读取value
		res := bucket.Get([]byte("1"))

		fmt.Printf("1: %s\n", res)

		return nil
	})
}