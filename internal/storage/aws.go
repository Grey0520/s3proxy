package storage

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"
)

type AWSStore struct {
	Session *session.Session
	ctx     context.Context
}

func NewAWSStore(accessKeyID, secretAccessKey, region string) (*AWSStore, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	return &AWSStore{
		Session: sess,
		ctx:     context.Background(),
	}, nil
}

func (store *AWSStore) CreateBucket(bucketName string) error {
	s3Client := s3.New(store.Session)

	_, err := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %v", err)
	}

	return nil
}

// DeleteBucket 使用存储在结构体中的session删除一个指定的S3 Bucket。
func (store *AWSStore) DeleteBucket(bucketName string) error {
	// 使用结构体中的Session创建一个S3服务客户端
	s3Client := s3.New(store.Session)

	// 调用DeleteBucket API
	_, err := s3Client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName), // 指定要删除的bucket名称
	})
	if err != nil {
		// 如果有错误发生，返回错误
		return fmt.Errorf("failed to delete bucket: %v", err)
	}

	// 如果删除成功，返回nil
	return nil
}

func (store *AWSStore) ListBucket(bucketName string) (*ListBucketResult, error) {
	// 使用结构体中的Session创建一个S3服务客户端
	s3Client := s3.New(store.Session)

	// 调用ListObjectsV2 API
	output, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:  aws.String(bucketName),
		MaxKeys: aws.Int64(1000),
	})
	if err != nil {
		// 如果有错误发生，返回错误
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	// 创建一个ListBucketResult类型的实例
	result := &ListBucketResult{
		XMLName: xml.Name{Local: "ListBucketResult"},
	}

	// 遍历ListObjectsV2Output中的Contents，将每个对象的Key和LastModified添加到ListBucketResult中
	for _, object := range output.Contents {
		result.Contents = append(result.Contents, Content{
			Key:          *object.Key,
			LastModified: *object.LastModified,
			ETag:         *object.ETag,
			Size:         *object.Size,
			StorageClass: *object.StorageClass,
		})
	}

	// 返回ListBucketResult实例
	return result, nil
}

func (store *AWSStore) ListAllMyBuckets() (*ListAllMyBucketsResult, error) {
	// 使用结构体中的Session创建一个S3服务客户端
	s3Client := s3.New(store.Session)

	// 调用ListBuckets API
	output, err := s3Client.ListBuckets(nil)
	if err != nil {
		// 如果有错误发生，返回错误
		return nil, fmt.Errorf("failed to list buckets: %v", err)
	}
	var buckets Buckets
	for _, s3Bucket := range output.Buckets {
		buckets.Bucket = append(buckets.Bucket, Bucket{
			Name:         *s3Bucket.Name,
			CreationDate: *s3Bucket.CreationDate,
		})
	}

	return &ListAllMyBucketsResult{
		XMLName: xml.Name{Local: "ListAllMyBucketsResult"},
		Xmlns:   "http://s3.amazonaws.com/doc/2006-03-01/",
		Buckets: buckets,
	}, nil
}

func (store *AWSStore) GetBucketAcl(bucketName string) (*AccessControlPolicy, error) {
	// 使用结构体中的Session创建一个S3服务客户端
	s3Client := s3.New(store.Session)

	// 调用GetBucketAcl API
	output, err := s3Client.GetBucketAcl(&s3.GetBucketAclInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		// 如果有错误发生，返回错误
		return nil, fmt.Errorf("failed to get bucket ACL: %v", err)
	}

	var owner Owner
	if output.Owner != nil {
		ownerId := ""
		ownerDisplayName := ""

		if output.Owner.ID != nil {
			ownerId = *output.Owner.ID
		}
		if output.Owner.DisplayName != nil {
			ownerDisplayName = *output.Owner.DisplayName
		}
		owner = Owner{
			ID:          ownerId,
			DisplayName: ownerDisplayName,
		}
	}
	// 创建一个AccessControlPolicy类型的实例
	policy := AccessControlPolicy{
		Owner: owner,
	}

	// 遍历GetBucketAclOutput中的Grants，将每个Grant的Grantee和Permission添加到AccessControlPolicy中
	for _, grant := range output.Grants {
		// 初始化需要的变量，用于避免nil解引用
		granteeID := ""
		granteeDisplayName := ""
		permission := ""

		// 检查Grantee是否为nil
		if grant.Grantee != nil {
			// 检查Grantee.ID是否为nil
			if grant.Grantee.ID != nil {
				granteeID = *grant.Grantee.ID
			}
			// 检查Grantee.DisplayName是否为nil
			if grant.Grantee.DisplayName != nil {
				granteeDisplayName = *grant.Grantee.DisplayName
			}
		}

		// 检查Permission是否为nil
		if grant.Permission != nil {
			permission = *grant.Permission
		}

		// 使用安全获取的值构造Grant
		policy.AccessControlList.Grant = append(policy.AccessControlList.Grant, Grant{
			Grantee: Grantee{
				XmlnsXsi:    "http://www.w3.org/2001/XMLSchema-instance",
				ID:          granteeID,
				DisplayName: granteeDisplayName,
				XsiType:     "CanonicalUser", // 假设XsiType是已知的固定值
			},
			Permission: permission,
		})
	}

	// 返回AccessControlPolicy实例
	return &policy, nil
}

func (store *AWSStore) PutObject(bucketName, objectKey string, data *Object) error {
	bucket, err := blob.OpenBucket(store.ctx, fmt.Sprintf("s3://%s", bucketName))
	if err != nil {
		return fmt.Errorf("failed to open bucket: %v", err)
	}
	defer bucket.Close()

	currentDate := time.Now().Format(time.RFC3339)
	opts := &blob.WriterOptions{
		Metadata: map[string]string{
			"creation-date": currentDate,
		},
		ContentType: data.ContentType,
	}

	w, err := bucket.NewWriter(store.ctx, objectKey, opts)
	if err != nil {
		return fmt.Errorf("failed to obtain writer: %v", err)
	}

	if _, err := io.Copy(w, data.Data); err != nil {
		return fmt.Errorf("failed to write object: %v", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	return nil
}

func (store *AWSStore) GetObject(bucketName, objectKey string) (*Object, error) {
	bucket, err := blob.OpenBucket(store.ctx, fmt.Sprintf("s3://%s", bucketName))
	if err != nil {
		return nil, fmt.Errorf("failed to open bucket: %v", err)
	}
	defer bucket.Close()

	r, err := bucket.NewReader(store.ctx, objectKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain reader: %v", err)
	}

	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %v", err)
	}

	return &Object{
		Size:         r.Size(),
		LastModified: r.ModTime(),
		ContentType:  r.ContentType(),
		Data:         bytes.NewReader(data),
	}, nil
}

func (store *AWSStore) DeleteObject(bucketName, objectKey string) error {
	bucket, err := blob.OpenBucket(store.ctx, fmt.Sprintf("s3://%s", bucketName))
	if err != nil {
		return fmt.Errorf("failed to open bucket: %v", err)
	}
	defer bucket.Close()

	if err := bucket.Delete(store.ctx, objectKey); err != nil {
		return fmt.Errorf("failed to delete object: %v", err)
	}

	return nil
}

func (store *AWSStore) CopyObject(srcBucketName, srcObjectKey, destBucketName, destObjcetKey string) error {
	bucket, err := blob.OpenBucket(store.ctx, fmt.Sprintf("s3://%s", srcBucketName))
	if err != nil {
		return fmt.Errorf("failed to open source bucket: %v", err)
	}
	defer bucket.Close()

	r, err := bucket.NewReader(store.ctx, srcObjectKey, nil)
	if err != nil {
		return fmt.Errorf("failed to obtain reader: %v", err)
	}
	defer r.Close()

	destBucket, err := blob.OpenBucket(store.ctx, fmt.Sprintf("s3://%s", destBucketName))
	if err != nil {
		return fmt.Errorf("failed to open destination bucket: %v", err)
	}
	defer destBucket.Close()

	w, err := destBucket.NewWriter(store.ctx, destObjcetKey, nil)
	if err != nil {
		return fmt.Errorf("failed to obtain writer: %v", err)
	}

	if _, err := io.Copy(w, r); err != nil {
		return fmt.Errorf("failed to copy object: %v", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	return nil
}

func (store *AWSStore) MoveObject(srcBucketName, srcObjectKey, destBucketName, destObjcetKey string) error {
	if err := store.CopyObject(srcBucketName, srcObjectKey, destBucketName, destObjcetKey); err != nil {
		return fmt.Errorf("failed to copy object: %v", err)
	}

	if err := store.DeleteObject(srcBucketName, srcObjectKey); err != nil {
		return fmt.Errorf("failed to delete object: %v", err)
	}

	return nil
}

func (store *AWSStore) HeadObject(bucketName, objectKey string) (map[string]string, error) {
	// Open the bucket
	bucket, err := blob.OpenBucket(store.ctx, fmt.Sprintf("s3://%s", bucketName))
	if err != nil {
		return nil, fmt.Errorf("failed to open bucket: %v", err)
	}
	defer bucket.Close()

	attrs, err := bucket.Attributes(store.ctx, objectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get object attributes: %v", err)
	}

	metadata := make(map[string]string)
	metadata["Size"] = fmt.Sprintf("%d", attrs.Size)
	metadata["LastModified"] = attrs.ModTime.Format(time.RFC3339)
	metadata["ContentType"] = attrs.ContentType

	// 其他的信息也加进去吧
	for key, value := range attrs.Metadata {
		metadata[key] = value
	}

	return metadata, nil
}
