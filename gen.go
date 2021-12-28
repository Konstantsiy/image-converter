package photoconverter

//go:generate mockgen -source=internal/service/service.go -destination=internal/service/mock/mock.go

//go:generate mockgen -source=internal/storage/storage.go -destination=internal/storage/mock/mock.go

//go:generate mockgen -source=internal/queue/queue.go -destination=internal/queue/mock/mock.go
