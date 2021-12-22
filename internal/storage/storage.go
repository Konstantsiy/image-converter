package storage

import (
	io "io"
)

// Storage represents images storage.
type Storage interface {
	UploadFile(file io.ReadSeeker, fileID string) error
	DownloadFile(fileID string) (io.ReadSeeker, error)
	GetDownloadURL(fileID string) (string, error)
}
