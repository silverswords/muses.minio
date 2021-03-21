# muses.minio

A private cloud storage solution based on Minio.

## Basic Features

The most basic features is to upload, download, and delete unstructured data objects. Users or clients can directly access these feactures through HTTP requests, and also support batch upload, download and delete.

### Middleware

The middleware mode designed here is for pluggable add new features. 
Middleware is based on the realization of basic features, adding some additional features, such as access control, batch request processing, access frequency limit, etc.
The reason for using the middleware mode is to realize that the new functions and basic functions do not interfere with each other, and each perform its own duties. Whether the new functions are used supports dynamic adjustment.

``` go
func Chain(j Judge, mw ...Middleware) Middleware {
	var arr = make([]bool, len(j))
	for _, value := range j {
		arr = append(arr, value)
	}
	return func(next Handler) Handler {
		for i := len(mw) - 1; i >= 0; i-- {
			if arr[i] {
				next = mw[i](next)
			} else {
				continue
			}
		}
		return next
	}
}
```

### New Features

- Current Limit: Record the total data size currently stored, and refuse to upload when the stored data size exceeds the threshold.

- MD5 Checkout: Check whether the uploaded object is complete. If the uploaded object is incomplete, request a re-upload. If it is incomplete for multiple times, it will return upload failure.

- Notification: Meet the needs of users who want to subscribe to the dynamic information of a certain object.

- Batch: Temporarily store user batch upload and download requests and process them in order.

- ACL: Determine whether the user has the permission to upload, download or delete according to the token that the user carries when sending the request.

- ObjectLock: Used to lock objects, which can only be used by the uploader for objects uploaded by themselves.

## Deployment

Deploy a minio cluster through k8s, and each cluster allows multiple tenants to be deployed. After deployment, we can expand the storage capacity by adding a new pool for each tenant.

## Architecture
![image](https://github.com/silverswords/muses.minio/blob/master/assets/bucket-architecture.png)

## Reference

- [Operator](https://docs.min.io/minio/k8s/reference/minio-operator-reference.html#minio-kubernetes-operator)

- [Expand tenant](https://docs.min.io/minio/k8s/reference/minio-kubectl-plugin.html#expand-a-minio-tenant)
