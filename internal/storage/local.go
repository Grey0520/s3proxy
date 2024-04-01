package storage

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
)

// Local File System (LFS) Store
type LFSStore struct {
	Bucket   *blob.Bucket
	ctx      context.Context
	basePath string
}

func NewLFSStore(basePath string) (*LFSStore, error) {
	err := createDirIfNotExist(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create dir: %v", err)
	}

	b, err := fileblob.OpenBucket(basePath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open bucket: %v", err)
	}
	return &LFSStore{
		Bucket:   b,
		ctx:      context.Background(),
		basePath: basePath,
	}, nil
}

func (local *LFSStore) CreateBucket(bucketName string) error {
	return local.createBucket(bucketName)
}

func (local *LFSStore) DeleteBucket(bucketName string) error {
	return local.deleteBucket(bucketName)
}

func (local *LFSStore) ListBucket(bucketName string) (*ListBucketResult, error) {
	return local.listBucket(bucketName)
}

func (local *LFSStore) ListAllMyBuckets() (*ListAllMyBucketsResult, error) {
	return local.listAllMyBuckets()
}

func (local *LFSStore) GetBucketAcl(bucketName string) (*AccessControlPolicy, error) {
	return local.getBucketAcl(bucketName)
}

func (local *LFSStore) PutObject(bucketName, objectKey string, data *Object) error {
	return local.putObject(bucketName, objectKey, data)
}

func (local *LFSStore) GetObject(bucketName, objectKey string) (*Object, error) {
	return local.getObject(bucketName, objectKey)
}

func (local *LFSStore) DeleteObject(bucketName, objectKey string) error {
	return local.deleteObject(bucketName, objectKey)
}

func (local *LFSStore) CopyObject(srcBucketName, srcObjectKey, destBucketName, destObjectKey string) error {
	return local.copyObject(srcBucketName, srcObjectKey, destBucketName, destObjectKey)
}

func (local *LFSStore) MoveObject(srcBucketName, srcObjectKey, destBucketName, destObjectKey string) error {
	return local.moveObject(srcBucketName, srcObjectKey, destBucketName, destObjectKey)
}

func (local *LFSStore) HeadObject(bucketName, objectKey string) (map[string]string, error) {
	return local.headObject(bucketName, objectKey)
}

func (local *LFSStore) createBucket(bucketName string) error {
	dir := fmt.Sprintf("%s/%s", local.basePath, bucketName)

	// 确保不存在同名的目录
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return fmt.Errorf("bucket %s already exists", bucketName)
	}

	if err := createDir(dir); err != nil {
		return err
	}

	b, err := fileblob.OpenBucket(dir, nil)
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %v", bucketName, err)
	}
	local.Bucket = b
	defer b.Close()

	return nil
}

func (local *LFSStore) deleteBucket(bucketName string) error {
	dir := fmt.Sprintf("%s/%s", local.basePath, bucketName)

	// 桶内有对象则不允许删除
	result, err := local.ListBucket(bucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket %s: %v", bucketName, err)
	}
	if len(result.Contents) > 0 {
		return fmt.Errorf("bucket %s is not empty", bucketName)
	}

	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to delete bucket %s: %v", bucketName, err)
	}
	return nil
}

func (local *LFSStore) listBucket(bucketName string) (*ListBucketResult, error) {
	if local.Bucket == nil {
		return nil, fmt.Errorf("bucket %s is not initialized", bucketName)
	}

	if err := local.checkoutBucket(bucketName); err != nil {
		return nil, err
	}
	iter := local.Bucket.List(&blob.ListOptions{})

	result := &ListBucketResult{
		XMLName: xml.Name{Local: "ListBucketResult"},
		Name:    bucketName,
	}
	for {
		obj, err := iter.Next(local.ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %v", err)
		}
		content := Content{
			Key:          obj.Key,
			LastModified: obj.ModTime,
			ETag:         string(obj.MD5),
			Size:         obj.Size,
			StorageClass: "STANDARD",
			Owner:        newFakeOwner(),
		}
		result.Contents = append(result.Contents, content)
	}

	return result, nil
}

func (local *LFSStore) listAllMyBuckets() (*ListAllMyBucketsResult, error) {
	var buckets []Bucket
	err := filepath.Walk(local.basePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			buckets = append(buckets, Bucket{
				Name:         info.Name(),
				CreationDate: info.ModTime(),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ListAllMyBucketsResult{
		XMLName: xml.Name{Local: "ListAllMyBucketsResult"},
		Owner:   newFakeOwner(),
		Buckets: Buckets{
			Bucket: buckets,
		},
	}, nil
}

func (local *LFSStore) getBucketAcl(bucketName string) (*AccessControlPolicy, error) {
	acp := &AccessControlPolicy{
		Owner: newFakeOwner(),
		AccessControlList: AccessControlList{
			Grant: []Grant{
				{
					Grantee: Grantee{
						ID:          "fake-owner-id",
						DisplayName: "fake-owner-name",
					},
					Permission: "FULL_CONTROL",
				},
			},
		},
	}
	return acp, nil
}

func (local *LFSStore) putObject(bucketName, objectKey string, data *Object) error {
	if err := local.checkoutBucket(bucketName); err != nil {
		return err
	}

	opts := &blob.WriterOptions{
		ContentType: data.ContentType,
	}
	writer, err := local.Bucket.NewWriter(local.ctx, objectKey, opts)
	if err != nil {
		return fmt.Errorf("failed to create object %s: %v", objectKey, err)
	}

	// if data.ContentType != "" {
	// 	writer.ContentType = data.ContentType
	// }
	if _, err := io.Copy(writer, data.Data); err != nil {
		writer.Close()
		return err
	}

	if err = writer.Close(); err != nil {
		return err
	}
	return nil
}

func (local *LFSStore) getObject(bucketName, objectKey string) (*Object, error) {
	if err := local.checkoutBucket(bucketName); err != nil {
		return nil, err
	}

	reader, err := local.Bucket.NewReader(local.ctx, objectKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get object %s: %v", objectKey, err)
	}
	// defer reader.Close()

	contentType := reader.ContentType()

	return &Object{
		Key:         objectKey,
		ContentType: contentType,
		Data:        reader,
	}, nil
}

func (local *LFSStore) deleteObject(bucketName, objectKey string) error {
	if err := local.checkoutBucket(bucketName); err != nil {
		return err
	}

	if err := local.Bucket.Delete(local.ctx, objectKey); err != nil {
		return fmt.Errorf("failed to delete object %s: %v", objectKey, err)
	}
	return nil
}

func (local *LFSStore) copyObject(srcBucket, srcObject, dstBucket, dstObject string) error {
	if err := local.checkoutBucket(srcBucket); err != nil {
		return err
	}
	if err := local.checkoutBucket(dstBucket); err != nil {
		return err
	}

	srcData, err := local.getObject(srcBucket, srcObject)
	if err != nil {
		return fmt.Errorf("failed to get object %s: %v", srcObject, err)
	}
	dstData := &Object{
		Key:         dstObject,
		ContentType: srcData.ContentType,
		Data:        srcData.Data,
	}
	if err := local.putObject(dstBucket, dstObject, dstData); err != nil {
		return fmt.Errorf("failed to put object %s: %v", dstObject, err)
	}
	return nil
}

func (local *LFSStore) moveObject(srcBucketName, srcObjectKey, destBucketName, destObjectKey string) error {
	err := local.copyObject(srcBucketName, srcObjectKey, destBucketName, destObjectKey)
	if err != nil {
		return fmt.Errorf("failed to copy object %s to %s: %v", srcObjectKey, destObjectKey, err)
	}

	err = local.deleteObject(srcBucketName, srcObjectKey)
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %v", srcObjectKey, err)
	}

	return nil
}

func (local *LFSStore) headObject(bucketName, objectKey string) (map[string]string, error) {
	// if err := local.checkoutBucket(bucketName); err != nil {
	// 	return nil, err
	// }

	// attrs, err := local.Bucket.Attrs(local.ctx, objectKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get object %s: %v", objectKey, err)
	// }

	// return map[string]string{
	// 	"Content-Length": strconv.FormatInt(attrs.Size, 10),
	// 	"Content-Type":   attrs.ContentType,
	// }, nil
	return nil, nil
}

// Some Utils
func createDir(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}
	return nil
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0o755)
	}
	return nil
}

func (local *LFSStore) checkoutBucket(bucketName string) error {
	if local.basePath == "" {
		return fmt.Errorf("basePath cannot be empty")
	}

	//假设不存在这个桶
	exist := false
	filepath.Walk(local.basePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == bucketName {
			exist = true
		}
		return nil
	})

	if exist {
		dir := fmt.Sprintf("%s/%s", local.basePath, bucketName)
		b, err := fileblob.OpenBucket(dir, nil)
		if err != nil {
			return fmt.Errorf("failed to open bucket %s: %v", bucketName, err)
		}
		local.Bucket = b
		return nil
	}
	return fmt.Errorf("bucket %s does not exist", bucketName)
}

func newFakeOwner() Owner {
	return Owner{
		ID:          "capgrry",
		DisplayName: "local",
	}
}
