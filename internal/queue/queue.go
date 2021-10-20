package queue

type Producer interface {
	SendToQueue(fileID, filename, sourceFormat, targetFormat, requestID string, ratio int) error
}
