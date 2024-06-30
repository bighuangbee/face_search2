package storage

import (
	"fmt"
	"testing"
)

func TestBoltDb(t *testing.T) {
	// 创建一个新的数据库
	db, err := NewBoltDb("test.db")
	if err != nil {
		t.Fatalf("Failed to create BoltDb: %v", err)
	}
	defer db.db.Close()

	// 创建测试数据
	testKey := "testKey"
	testValue := &RegisteInfo{Filename: "testName"}

	// 测试 Update 方法
	err = db.Update(testKey, testValue)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	err = db.Update("testKey1", testValue)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	err = db.Update("testKey2", testValue)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// 测试 Read 方法
	value, ok := db.Read(testKey)
	if ok {
		t.Fatalf("Read failed: ")
	}
	if value.Filename != testValue.Filename {
		t.Fatalf("Read returned incorrect value: got %v, want %v", value.Filename, testValue.Filename)
	}

	// 测试 ReadBatch 方法
	values, err := db.ReadBatch()
	if err != nil {
		t.Fatalf("ReadBatch failed: %v", err)
	}
	for i, info := range values {
		fmt.Println("ReadBatch", i, *info)
	}

	// 测试 Delete 方法
	err = db.Delete(testKey)
	if err != nil {
		fmt.Println("Delete failed: ", err)
	}

	// 确认键值对已经被删除
	value, ok = db.Read(testKey)
	fmt.Println("read , ", testKey, err, value)
	if err == nil {
		fmt.Println("Expected error, got nil")
	}
	if value != nil {
		fmt.Println("Expected nil value, got ", value)
	}
}
