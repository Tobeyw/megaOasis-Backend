package nullstring

import "strings"

func IsNull(str string) bool {
	// 首尾去掉空格之后检查长度
	str = strings.TrimSpace(str)
	if len(str) == 0 {
		return false
	} else {

		return true
	}
}
func IsZero(num int64) bool {

	if num > 0 {
		return true
	} else {

		return false
	}
}
