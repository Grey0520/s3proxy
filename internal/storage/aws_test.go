package storage

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

// 注意：当仓库公开时，这个测试会失效
const (
	accessKey = "AKIAVI2TJMRA4TBYJRFL"
	secretKey = "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"
	region    = "ap-east-1"
)

func TestNewAWSStore(t *testing.T) {

	region := "ap-east-1"

	store, err := NewAWSStore(accessKey, secretKey, region)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if store == nil {
		t.Fatal("Expected non-nil AWSStore instance")
	}

	// 尝试使用 store.Session 创建一个 S3 服务客户端，并执行一个无副作用的操作来验证会话是否有效。
	s3svc := s3.New(store.Session)
	_, err = s3svc.ListBuckets(nil)
	if err != nil {
		t.Fatalf("Expected valid session, got error %v", err)
	}
}

func TestAWSStore_CreateBucket(t *testing.T) {
	store, err := NewAWSStore(accessKey, secretKey, region)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	bucketName := fmt.Sprintf("s3proxy-test-%d", time.Now().UnixNano())
	if err := store.CreateBucket(bucketName); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestAWSStore_DeleteBucket(t *testing.T) {
	store, err := NewAWSStore(accessKey, secretKey, region)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	bucketName := fmt.Sprintf("s3proxy-test-%d", time.Now().UnixNano())
	if err := store.CreateBucket(bucketName); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if err := store.DeleteBucket(bucketName); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestListBucket(t *testing.T) {
	store, err := NewAWSStore(accessKey, secretKey, region)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	objects, err := store.ListBucket("s3proxy-reserved")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Log(objects)
	// fmt.Println(objects)
}

func TestListAllMyBuckets(t *testing.T) {
	store, err := NewAWSStore(accessKey, secretKey, region)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	buckets, err := store.ListAllMyBuckets()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Log(buckets)
	// fmt.Println(buckets)
}

func TestGetBucketAcl(t *testing.T) {
	store, err := NewAWSStore(accessKey, secretKey, region)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	acl, err := store.GetBucketAcl("s3proxy-reserved")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Log(acl)
	// fmt.Println(acl)
}

// 以下是 Object relative test
func TestPutObject(t *testing.T) {
	store, err := NewAWSStore(accessKey, secretKey, region)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	content := []byte("Hello, world!")
	object := &Object{
		Data: bytes.NewReader(content),
	}
	if err := store.PutObject("s3proxy-reserved", "test.txt", object); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestGetObject(t *testing.T) {
	store, err := NewAWSStore(accessKey, secretKey, region)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	object, err := store.GetObject("s3proxy-reserved", "test.txt")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Log(object)
	fmt.Println(object)
}

func TestDeleteObject(t *testing.T) {
	store, err := NewAWSStore(accessKey, secretKey, region)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	store.PutObject("s3proxy-reserved", "test-for-delete.txt", &Object{
		Data: bytes.NewReader([]byte("Hello, world!")),
	})
	if err := store.DeleteObject("s3proxy-reserved", "test-for-delete.txt"); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestCopyObject(t *testing.T) {
	store, err := NewAWSStore(accessKey, secretKey, region)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if err := store.CopyObject("s3proxy-reserved", "test.txt", "s3proxy-copy", "test-copy.txt"); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestMoveObject(t *testing.T) {
	store, err := NewAWSStore(accessKey, secretKey, region)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	store.PutObject("s3proxy-reserved", "test-for-move.txt", &Object{
		Data: bytes.NewReader([]byte("Hello, world!")),
	})
	if err := store.MoveObject("s3proxy-reserved", "test-for-move.txt", "s3proxy-copy", "test-for-move.txt"); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestHeadObject(t *testing.T) {
	store, err := NewAWSStore(accessKey, secretKey, region)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	object, err := store.HeadObject("s3proxy-reserved", "test.txt")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Log(object)
	fmt.Println(object)
}
