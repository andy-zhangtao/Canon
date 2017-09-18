package util

const (
	CCTV = "cctv"
)

// MakeChanMap 填充视频列表
func MakeChanMap() map[string]string {
	chanMap := make(map[string]string)

	// CCTV
	chanMap["1001"] = CCTV
	chanMap["1002"] = CCTV
	chanMap["1003"] = CCTV
	chanMap["1004"] = CCTV
	chanMap["1005"] = CCTV

	return chanMap
}
