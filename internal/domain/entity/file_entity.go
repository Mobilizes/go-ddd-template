package entity

import "time"

type File struct {
	ID string 
	UserID string
	Name string
	StorageKey string
	MimeType string
	SizeBytes int64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func NewFile(id string, userId string, name string, storageKey string, mimeType string, sizeBytes int64) *File {
	return &File{
		ID: id,
		UserID: userId,
		Name: name,
		StorageKey: storageKey,
		MimeType: mimeType,
		SizeBytes: sizeBytes,
	}
}
