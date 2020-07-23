package storage

func (b *Bucket) AccessBaseOnWeight() {
	// var minioClient *minio.Client
	configs := getConfig()
	for i, v := range configs.minioClients {
		if v.weight == 1 {
			b = &Bucket{minioClient: getMinioClients(configs)[i]}
			b.MakeBucket()
		}
		if v.weight > 0 && v.weight < 1 {

		}
	}
}
