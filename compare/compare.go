package compare

// Comapre对比HTTP 响应头部和内容
func Compare(status1 int, header1 map[string]string, body1 string, status2 int, header2 map[string]string, body2 string) (result map[string]Diff) {
	// 比较HTTP 内容
	httpInfo1 := ParseHttpInfo(status1, header1)
	httpInfo2 := ParseHttpInfo(status2, header2)
	httpInfo1.Compare(httpInfo2, result)

	// 比较文件内容
	info1 := ParseInfo(body1)
	info2 := ParseInfo(body2)
	info1.Compare(info2, result)
}
