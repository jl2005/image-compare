package compare

type Info struct {
	httpInfo  *HttpInfo
	fileInfo  *FileInfo
	imageInfo *ImageInfo

	Data []byte
}

var DEFAULT_INFO = &Info{}

func ParseInfo(status int, header map[string][]string, data []byte) *Info {
	info := &Info{
		httpInfo: ParseHttpInfo(status, header),
		Data:     data,
	}
	info.fileInfo = ParseFileInfo(data)
	if isImage(info.fileInfo.Format) {
		info.imageInfo = ParseImageInfo(data)
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
	this.imageInfo.Compare(other.imageInfo, result)
}
