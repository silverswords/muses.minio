package storage

type backend interface {
	Getter
	Saver
	Deleter
}
