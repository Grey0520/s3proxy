package storage

import (
	"encoding/xml"
	"io"
	"time"
)

// Object 表示S3中存储的一个对象，包括其元数据和数据内容。
type Object struct {
	Key          string    // 对象的键（文件名）
	Size         int64     // 对象的大小，以字节为单位
	LastModified time.Time // 对象最后被修改的时间
	ContentType  string    // 对象的MIME类型
	Data         io.Reader // 对象的数据流
}

// ListAllMyBucketsResult 是 GET / 的根 xml 元素
type ListAllMyBucketsResult struct {
	XMLName xml.Name `xml:"ListAllMyBucketsResult"`
	Xmlns   string   `xml:"xmlns,attr"`
	Owner   Owner    `xml:"Owner"`
	Buckets Buckets  `xml:"Buckets"`
}

// ListBucketResult 是 GET /BUCKETNAME 的根 xml 元素
type ListBucketResult struct {
	XMLName     xml.Name  `xml:"ListBucketResult"`
	Xmlns       string    `xml:"xmlns,attr"`
	Name        string    `xml:"Name"`
	Prefix      string    `xml:"Prefix"`
	MaxKeys     int       `xml:"MaxKeys"`
	Marker      string    `xml:"Marker"`
	IsTruncated bool      `xml:"IsTruncated"`
	Contents    []Content `xml:"Contents"`
}

// AccessControlList 是 GET /BUCKETNAME?acl 的根 xml 元素
type AccessControlPolicy struct {
	XMLName           xml.Name          `xml:"AccessControlPolicy"`
	Xmlns             string            `xml:"xmlns,attr"`
	Owner             Owner             `xml:"Owner"`
	AccessControlList AccessControlList `xml:"AccessControlList"`
}

// ListPartsResult 是 GET /BUCKETNAME/OBJECTNAME?uploadId=UPLOADID 的根 xml 元素
type ListPartsResult struct {
	XMLName      xml.Name  `xml:"ListPartsResult"`
	Xmlns        string    `xml:"xmlns,attr"`
	Bucket       string    `xml:"Bucket"`
	Key          string    `xml:"Key"`
	UploadId     string    `xml:"UploadId"`
	Initiator    Initiator `xml:"Initiator"`
	Owner        Owner     `xml:"Owner"`
	StorageClass string    `xml:"StorageClass"`
	Part         Part      `xml:"Part"`
}

// Buckets 与 ListAllMyBucketsResult.Buckets 相对应
type Buckets struct {
	Bucket []Bucket `xml:"Bucket"`
}

// Bucket 与 Buckets.Bucket 相对应
type Bucket struct {
	Name         string    `xml:"Name"`
	CreationDate time.Time `xml:"CreationDate"`
}

// Content 与 ListBucketResult.Contents 相对应
type Content struct {
	Key          string    `xml:"Key"`
	LastModified time.Time `xml:"LastModified"`
	ETag         string    `xml:"ETag"`
	Size         int64     `xml:"Size"`
	StorageClass string    `xml:"StorageClass"`
	Owner        Owner     `xml:"Owner"`
}

// Owner 与 ListAllMyBucketsResult.Owner 相对应
type Owner struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

type AccessControlList struct {
	Grant []Grant `xml:"Grant"`
}

type Grant struct {
	Grantee    Grantee `xml:"Grantee"`
	Permission string  `xml:"Permission"`
}

type Grantee struct {
	XMLName     xml.Name `xml:"Grantee"`
	XmlnsXsi    string   `xml:"xmlns:xsi,attr"`
	XsiType     string   `xml:"xsi:type,attr"`
	ID          string   `xml:"ID"`
	DisplayName string   `xml:"DisplayName"`
}

// Initiator 与 ListPartsResult.Initiator 相对应
type Initiator struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

// Part 与 ListPartsResult.Part 相对应
type Part struct {
	PartNumber   int       `xml:"PartNumber"`
	LastModified time.Time `xml:"LastModified"`
	ETag         string    `xml:"ETag"`
	Size         int       `xml:"Size"`
}
