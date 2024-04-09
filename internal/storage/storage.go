package storage

import "github.com/Grey0520/s3proxy/internal/config"

type StorageProvider interface {
	CreateBucket(bucketName string) error
	DeleteBucket(bucketName string) error
	ListBucket(bucketName string) (*ListBucketResult, error)
	ListAllMyBuckets() (*ListAllMyBucketsResult, error)
	GetBucketAcl(bucketName string) (*AccessControlPolicy, error)

	PutObject(bucketName, objectKey string, data *Object) error
	GetObject(bucketName, objectKey string) (*Object, error)
	DeleteObject(bucketName, objectKey string) error
	// ListObjects(bucketName string, prefix string, recursive bool) ([]*Object, error)
	CopyObject(srcBucketName, srcObjectKey, destBucketName, destObjectKey string) error
	MoveObject(srcBucketName, srcObjectKey, destBucketName, destObjectKey string) error
	HeadObject(bucketName, objectKey string) (map[string]string, error)
}

func NewStorageProvider(cfg config.Config) (StorageProvider, error) {
	// Stub
	switch cfg.Cloud.Provider {
	case "aws":
		return NewAWSStore(cfg.Cloud.Identity, cfg.Cloud.Key, cfg.Cloud.Region)
	case "local":
		return NewLFSStore(cfg.Cloud.Filesystem.Basedir)
	default:
		return nil, nil
	}
}
