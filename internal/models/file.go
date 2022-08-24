package models

import "mime/multipart"

type FileInfo struct {
	ObjectName  string
	FileBuffer  multipart.File
	ContentType string
	FileSize    int64
}
