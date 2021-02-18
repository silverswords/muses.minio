package driver

import (
	"encoding/xml"
	"github.com/minio/minio-go/v7/pkg/encrypt"
)

type MakeBucketOptions struct {
	Region        string
	ObjectLocking bool
}

type OtherMakeBucketOption func(o *MakeBucketOptions)

func WithRegion(region string) OtherMakeBucketOption {
	return func(o *MakeBucketOptions) {
		o.Region = region
	}
}

func WithObjectLocking(objectLocking bool) OtherMakeBucketOption {
	return func(o *MakeBucketOptions) {
		o.ObjectLocking = objectLocking
	}
}

type SetBucketVersioningOptions struct {
	XMLName   xml.Name
	Status    string
	MFADelete string
}

type OtherSetBucketVersioningOption func(o *SetBucketVersioningOptions)

func WithXMLName(xmlName xml.Name) OtherSetBucketVersioningOption {
	return func(o *SetBucketVersioningOptions) {
		o.XMLName = xmlName
	}
}

func WithStatus(status string) OtherSetBucketVersioningOption {
	return func(o *SetBucketVersioningOptions) {
		o.Status = status
	}
}

func WithMFADelete(mfaDelete string) OtherSetBucketVersioningOption {
	return func(o *SetBucketVersioningOptions) {
		o.MFADelete = mfaDelete
	}
}

type OtherPutObjectOptions struct {
	ServerSideEncryption encrypt.ServerSide
	ObjectSize int64
}

type OtherPutObjectOption func(o *OtherPutObjectOptions)

func WithServerSideEncryption(ServerSideEncryption encrypt.ServerSide) OtherPutObjectOption {
	return func(o *OtherPutObjectOptions) {
		o.ServerSideEncryption = ServerSideEncryption
	}
}

func WithObjectSize(ObjectSize int64) OtherPutObjectOption {
	return func(o *OtherPutObjectOptions) {
		o.ObjectSize = ObjectSize
	}
}

type GetObjectOptions struct {
	ServerSideEncryption encrypt.ServerSide
}

type OtherGetObjectOption func(o *GetObjectOptions)

func WithGetServerSideEncryption(ServerSideEncryption encrypt.ServerSide) OtherGetObjectOption {
	return func(o *GetObjectOptions) {
		o.ServerSideEncryption = ServerSideEncryption
	}
}

type RemoveObjectOptions struct {
	GovernanceBypass bool
}

type OtherRemoveObjectOption func(o *RemoveObjectOptions)

func WithGovernaceBypass(GovernaceBypass bool) OtherRemoveObjectOption {
	return func(o *RemoveObjectOptions) {
		o.GovernanceBypass = GovernaceBypass
	}
}

type ListObjectsOptions struct {
	Prefix string
}

type OtherListObjectsOption func(o *ListObjectsOptions)

func WithPrefix(prefix string) OtherListObjectsOption {
	return func(o *ListObjectsOptions) {
		o.Prefix = prefix
	}
}
