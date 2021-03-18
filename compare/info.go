package compare

import (
	"net/http"
)

type Info struct {
	httpInfo  *HttpInfo
	fileInfo  *FileInfo
	imageInfo *ImageInfo
}

var DEFAUTL_INFO = &Info{}

func ParseInfo(status int, header map[string]string, body string) {
	info := &Info{
		httpInfo: ParseHttpInfo(status, header),
	}
	info.fileInfo = ParseFileInfo(body)
	if isImage(info.fileInfo.Format) {
		info.imageInfo = ParseImageInfo(body)
	}
	return info
}

func (this *Info) Compare(other *Info, result map[string]Diff) {
	if this == other {
		return
	}
	if this == nil {
		this = DEFAULT_INFO
	}
	if other == nil {
		other = DEFAULT_INFO
	}
	this.httpInfo.Compare(other.httpInfo, result)
	this.fileInfo.Compare(other.fileInfo, result)
	this.imageInfo.Compare(other.fileInfo, result)
}
