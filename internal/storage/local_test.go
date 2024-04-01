package storage

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

var baseDir = "/tmp/buckets"

func TestNewLFSStore(t *testing.T) {
	dir := baseDir
	_, err := NewLFSStore(dir)
	if err != nil {
		t.Errorf("failed to create new local storage provider: %v", err)
	}
}

func TestLFSStoreCreateBucket(t *testing.T) {
	dir := baseDir
	store, _ := NewLFSStore(dir)

	bucketName := "test_bucket_create"
	_ = store.CreateBucket(bucketName)
}

func TestLFSStoreDeleteBucket(t *testing.T) {
	dir := baseDir
	store, _ := NewLFSStore(dir)

	bucketName := "test_bucket_delete"
	err := store.CreateBucket(bucketName)
	if err != nil {
		t.Fatalf("Failed to create bucket: %v", err)
	}

	err = store.DeleteBucket(bucketName)
	if err != nil {
		t.Fatalf("Failed to delete bucket: %v", err)
	}
}

func TestLFSStoreListBucket(t *testing.T) {
	dir := baseDir
	store, _ := NewLFSStore(dir)

	bucketName := "test-bucket-listbucket"
	err := store.CreateBucket(bucketName)

	result, err := store.ListBucket(bucketName)
	if err != nil {
		t.Fatalf("Failed to list bucket: %v", err)
	}
	resultBucket := (*result).Name
	if resultBucket != bucketName {
		t.Errorf("Expected bucket name %s, got %s", bucketName, (*result).Name)
	}
}

func TestLFSStoreListAllMyBuckets(t *testing.T) {
	dir := baseDir
	store, _ := NewLFSStore(dir)

	bucketName := "test-bucket-listallmybuckets"
	err := store.CreateBucket(bucketName)

	result, err := store.ListAllMyBuckets()
	if err != nil {
		t.Fatalf("Failed to list all my buckets: %v", err)
	}
	if result.Buckets.Bucket == nil {
		t.Fatalf("Failed to list all my buckets: %v", err)
	}
	found := false
	for _, b := range result.Buckets.Bucket {
		if b.Name == bucketName {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected bucket name %s, got %v", bucketName, result.Buckets)
	}
}

func TestLFSStoreGetBucketAcl(t *testing.T) {
	dir := baseDir
	store, _ := NewLFSStore(dir)

	bucketName := "test-bucket-getbucketacl"
	err := store.CreateBucket(bucketName)

	result, err := store.GetBucketAcl(bucketName)
	if err != nil {
		t.Fatalf("Failed to get bucket acl: %v", err)
	}
	if (*result).Owner.ID == "" {
		t.Fatalf("Failed to get bucket acl: %v", err)
	}
	fmt.Println(result)
}

func TestLFSStorePutObject(t *testing.T) {
	dir := baseDir
	store, _ := NewLFSStore(dir)

	bucketName := "test-bucket-putobject"
	store.CreateBucket(bucketName)

	objectKey := "test-object-key"
	objectData := &Object{
		Key:         objectKey,
		Data:        bytes.NewReader([]byte("test content")),
		ContentType: "text/plain",
	}
	err := store.PutObject(bucketName, objectKey, objectData)
	if err != nil {
		t.Fatalf("Failed to put object: %v", err)
	}
}

func TestLFSStoreGetObject(t *testing.T) {
	dir := baseDir
	store, _ := NewLFSStore(dir)

	bucketName := "test-bucket-getobject"
	store.CreateBucket(bucketName)

	objectKey := "test-object-key"
	objectData := &Object{
		Key:         objectKey,
		Data:        bytes.NewReader([]byte("test content")),
		ContentType: "text/plain",
	}
	store.PutObject(bucketName, objectKey, objectData)

	result, err := store.GetObject(bucketName, objectKey)
	if err != nil {
		t.Fatalf("Failed to get object: %v", err)
	}
	retrievedData, err := io.ReadAll(result.Data)
	if err != nil {
		t.Fatalf("Failed to read object data: %v", err)
	}
	if !bytes.Equal(retrievedData, []byte("test content")) {
		t.Errorf("Expected object data %s, got %s", "test content", string(retrievedData))
	}
}

func TestLFSStoreDeleteObject(t *testing.T) {
	dir := baseDir
	store, _ := NewLFSStore(dir)

	bucketName := "test-bucket-deleteobject"
	store.CreateBucket(bucketName)

	objectKey := "test-object-key"
	objectData := &Object{
		Key:         objectKey,
		Data:        bytes.NewReader([]byte("test content")),
		ContentType: "text/plain",
	}
	store.PutObject(bucketName, objectKey, objectData)

	err := store.DeleteObject(bucketName, objectKey)
	if err != nil {
		t.Fatalf("Failed to delete object: %v", err)
	}
}

func TestLFSStoreCopyObject(t *testing.T) {
	dir := baseDir
	store, _ := NewLFSStore(dir)

	bucketName := "test-bucket-copyobject"
	store.CreateBucket(bucketName)

	objectKey := "test-object-key"
	objectData := &Object{
		Key:         objectKey,
		Data:        bytes.NewReader([]byte("test content")),
		ContentType: "text/plain",
	}
	store.PutObject(bucketName, objectKey, objectData)

	destObjectKey := "test-object-key-copy"
	err := store.CopyObject(bucketName, objectKey, bucketName, destObjectKey)
	if err != nil {
		t.Fatalf("Failed to copy object: %v", err)
	}
}

func TestLFSStoreMoveObject(t *testing.T) {
	dir := baseDir
	store, _ := NewLFSStore(dir)

	bucketName := "test-bucket-moveobject"
	store.CreateBucket(bucketName)

	objectKey := "test-object-key"
	objectData := &Object{
		Key:         objectKey,
		Data:        bytes.NewReader([]byte("test content")),
		ContentType: "text/plain",
	}
	store.PutObject(bucketName, objectKey, objectData)

	destObjectKey := "test-object-key-move"
	err := store.MoveObject(bucketName, objectKey, bucketName, destObjectKey)
	if err != nil {
		t.Fatalf("Failed to move object: %v", err)
	}
}
