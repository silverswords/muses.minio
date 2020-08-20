package bucket

import "os"

type client interface {
	PutObject(objectName string, object *os.File) error
	GetObject(objectName string) ([]byte, error)
}


