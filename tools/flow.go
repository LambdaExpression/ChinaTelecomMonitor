package tools

import (
	"errors"
	"strconv"
	"strings"
)

func ToInt64(s string) (int64, error) {
	if strings.Contains(s, "GB") {
		sn := strings.ReplaceAll(s, "GB", "")
		n, err := strconv.ParseFloat(sn, 64)
		if err != nil {
			return 0, err
		}
		return int64(n * 1024 * 1024), nil
	} else if strings.Contains(s, "MB") {
		sn := strings.ReplaceAll(s, "MB", "")
		n, err := strconv.ParseFloat(sn, 64)
		if err != nil {
			return 0, err
		}
		return int64(n * 1024), nil
	} else if strings.Contains(s, "KB") {
		sn := strings.ReplaceAll(s, "KB", "")
		n, err := strconv.ParseFloat(sn, 64)
		if err != nil {
			return 0, err
		}
		return int64(n), nil
	} else if strings.Contains(s, "B") {
		return 0, nil
	}
	return 0, errors.New("文本格式不符合规范")
}
