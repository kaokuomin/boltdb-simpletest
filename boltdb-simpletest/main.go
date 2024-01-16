package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/boltdb/bolt"
)

func main() {
	// 调用creatdb函数创建数据库文件
	nodeID := "N1"
	db, err := creatdb(nodeID)

	// //如果已生成一个数据库文件，那就从绝对路径打开这个.db文件
	// dbPath := "E:/GOPATH/goworkplace/BFT_Rpc4/pbft/N0.db"
	// db, err := bolt.Open(dbPath, 0600, nil)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // 在main函数中确保在所有操作完成后关闭数据库

	//根据自己是否需要执行以下功能来选择注释/不注释某段功能代码
	// 写入数据
	err = writeData(db, strconv.Itoa(1), "myValue")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("写入数据成功")
	}

	// 读取数据
	value, err := readData(db, strconv.Itoa(1))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Read data:", value)

	// // 删除数据
	// err = deleteData(db, "myKey")
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println("删除数据成功")
	// }

}

// 创建数据库文件
func creatdb(nodeID string) (*bolt.DB, error) {
	//比如nodeID=N1,那么创建一个名为N1.db的数据库文件
	db, err := bolt.Open(nodeID+".db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err // 返回 nil 和错误，而不是错误本身
	}
	// defer db.Close() // 确保在函数返回前关闭数据库
	return db, nil // 成功创建数据库后返回数据库实例和nil错误
}

// 写入数据到数据库
func writeData(db *bolt.DB, key string, value string) error {
	return db.Update(func(tx *bolt.Tx) error {
		// 创建或获取名为 "bucket1" 的 bucket
		b, err := tx.CreateBucketIfNotExists([]byte("bucket1"))
		if err != nil {
			return err
		}

		// 将键值对写入到 bucket 中
		return b.Put([]byte(key), []byte(value))
	})
}

// 从数据库读取数据
func readData(db *bolt.DB, key string) (string, error) {
	var value string

	err := db.View(func(tx *bolt.Tx) error {
		// 获取名为 "bucket1" 的 bucket
		b := tx.Bucket([]byte("bucket1"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		// 从 bucket 中读取键对应的值
		v := b.Get([]byte(key))
		if v == nil {
			return fmt.Errorf("key not found")
		}

		// 将字节切片转换为字符串
		value = string(v)
		return nil
	})

	return value, err
}

// 从数据库删除数据
func deleteData(db *bolt.DB, key string) error {
	return db.Update(func(tx *bolt.Tx) error {
		// 获取名为 "bucket1" 的 bucket
		b := tx.Bucket([]byte("bucket1"))
		if b == nil {
			return fmt.Errorf("bucket not found")
		}

		// 删除键值对
		return b.Delete([]byte(key))
	})
}
