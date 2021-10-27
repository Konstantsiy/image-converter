//+build ignore

package photoconverter

//go:generate mockgen -source=internal/service/service.go -destination=internal/service/mock/mock.go
