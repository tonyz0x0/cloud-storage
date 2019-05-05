package ceph

import (
	"gopkg.in/amz.v1/aws"
	"gopkg.in/amz.v1/s3"

	cfg "cloud-storage/config"
)

var cephConn *s3.S3

// GetCephConnection : Get Ceph Connection
func GetCephConnection() *s3.S3 {
	if cephConn != nil {
		return cephConn
	}

	// Config Initialization
	auth := aws.Auth{
		AccessKey: cfg.CephAccessKey,
		SecretKey: cfg.CephSecretKey,
	}

	curRegion := aws.Region{
		Name:                 "default",
		EC2Endpoint:          cfg.CephGWEndpoint,
		S3Endpoint:           cfg.CephGWEndpoint,
		S3BucketEndpoint:     "",
		S3LocationConstraint: false,
		S3LowercaseBucket:    false,
		Sign:                 aws.SignV2,
	}

	// 2. Create S3 Connection
	return s3.New(auth, curRegion)
}

// GetCephBucket : Get bucket
func GetCephBucket(bucket string) *s3.Bucket {
	conn := GetCephConnection()
	return conn.Bucket(bucket)
}

// PutObject : Upload to Ceph cluster
func PutObject(bucket string, path string, data []byte) error {
	return GetCephBucket(bucket).Put(path, data, "octet-stream", s3.PublicRead)
}
