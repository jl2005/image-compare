package compare

const (
	STATUS = "status"
	HEADER = "header"
)

var defaultHeader = make(map[string][]string)

var compareHeaders = []string{
	"Content-Length",
	"Content-Type",
}

// TODO 增加对header进行修改的接口

type HttpInfo struct {
	Status int
	Header map[string][]string
}

func ParseHttpInfo(status int, header map[string][]string) *HttpInfo {
	info := &HttpInfo{
		Status: status,
		Header: header,
	}
	if info.Header == nil {
		info.Header = defaultHeader
	}
	return info
}

func (this *HttpInfo) Compare(other *HttpInfo, result map[string]Diff) {
	// 比较status
	if this.Status != other.Status {
		result[STATUS] = &IntDiff{this.Status, other.Status}
	}

	// 只比较指定header就可以
	for _, key := range compareHeaders {
		var val1, val2 []string
		val1, _ = this.Header[key]
		val2, _ = other.Header[key]
		// TODO 将Header转换为统一小写
		if len(val1) == len(val2) {
			isSame := true
			for i := range val1 {
				if val1[i] != val2[i] {
					isSame = false
					break
				}
			}
			if isSame {
				continue
			}
		}
		result[key] = &StringArrayDiff{val1, val2}
	}
}
