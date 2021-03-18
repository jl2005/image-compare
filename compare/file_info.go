package compare

import (
	"crypto/md5"

	"github.com/h2non/filetype"
)

const (
	SIZE   = "size"
	FORMAT = "format"
	MD5    = "md5"
)

type FileInfo struct {
	Size   int
	Format string
	MD5    string
}

var DEFAULT_FILE_INFO = &FileInfo{}

func ParseFileInfo(data []byte) *FileInfo {
	info := &FileInfo{
		Size: len(data),
	}
	if kind, err := filetype.Match(data); err != nil {
		info.Format = kind.Extension
	}
	info.MD5 = string(md5.Sum(data))
	return info
}

func (this *FileInfo) Compare(other *FileInfo, result map[string]Diff) {
	if this == other {
		return
	}
	if this == nil {
		this = DEFAULT_FILE_INFO
	}
	if other == nil {
		other = DEFAULT_FILE_INFO
	}
	if this.Size != other.Size {
		result[SIZE] = &IntDiff{this.Size, other.Size}
	}
	if this.Format != other.Format {
		result[FORMAT] = &StringDiff{this.Format, other.Format}
	}
	if this.MD5 != other.MD5 {
		result[MD5] = &StringDiff{this.MD5, other.MD5}
	}
}
