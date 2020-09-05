package bucketStorage

import (
	"encoding/xml"
	"github.com/minio/minio-go/v7/pkg/encrypt"
)

type OtherBucketConfigOptions struct {
	useCache      bool
}

type OtherBucketConfigOption func(*OtherBucketConfigOptions)

func WithUseCache(useCache bool) OtherBucketConfigOption {
	return func(o *OtherBucketConfigOptions) {
		o.useCache = useCache
	}
}

type OtherMakeBucketOptions struct{
	Region string
	ObjectLocking bool
}

type OtherMakeBucketOption func(o *OtherMakeBucketOptions)

func WithRegion(region string) OtherMakeBucketOption {
	return func(o *OtherMakeBucketOptions) {
		o.Region = region
	}
}

func WithObjectLocking(objectLocking bool) OtherMakeBucketOption {
	return func(o *OtherMakeBucketOptions) {
		o.ObjectLocking = objectLocking
	}
}

type OtherSetBucketVersioningOptions struct {
	XMLName   xml.Name
	Status    string
	MFADelete string
}

type OtherSetBucketVersioningOption func(o *OtherSetBucketVersioningOptions)

func WithXMLName(xmlName xml.Name) OtherSetBucketVersioningOption {
	return func(o *OtherSetBucketVersioningOptions) {
		o.XMLName = xmlName
	}
}

func WithStatus(status string) OtherSetBucketVersioningOption {
	return func(o *OtherSetBucketVersioningOptions) {
		o.Status = status
	}
}

func WithMFADelete(mfaDelete string) OtherSetBucketVersioningOption {
	return func(o *OtherSetBucketVersioningOptions) {
		o.MFADelete = mfaDelete
	}
}

type OtherPutObjectOptions struct {
	ServerSideEncryption    encrypt.ServerSide
}

type OtherPutObjectOption func(o *OtherPutObjectOptions)

func WithServerSideEncryption(ServerSideEncryption encrypt.ServerSide) OtherPutObjectOption {
	return func(o *OtherPutObjectOptions) {
		o.ServerSideEncryption = ServerSideEncryption
	}
}

type OtherGetObjectOptions struct {
	ServerSideEncryption encrypt.ServerSide
}

type OtherGetObjectOption func(o *OtherGetObjectOptions)

func WithGetServerSideEncryption(ServerSideEncryption encrypt.ServerSide) OtherGetObjectOption {
	return func(o *OtherGetObjectOptions) {
		o.ServerSideEncryption = ServerSideEncryption
	}
}

type OtherRemoveObjectOptions struct {
	GovernanceBypass bool
}

type OtherRemoveObjectOption func(o *OtherRemoveObjectOptions)

func WithGovernaceBypass(GovernaceBypass bool) OtherRemoveObjectOption {
	return func(o *OtherRemoveObjectOptions) {
		o.GovernanceBypass = GovernaceBypass
	}
}

type OtherListObjectsOptions struct {
	prefix string
}

type OtherListObjectsOption func(o *OtherListObjectsOptions)

func WithPrefix(prefix string) OtherListObjectsOption {
	return func(o *OtherListObjectsOptions) {
		o.prefix = prefix
	}
}
