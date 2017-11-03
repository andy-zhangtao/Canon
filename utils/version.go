package utils

var _VERSION_ = "unknown"

func GetVersion() string {
	return _VERSION_
}

func SetVersion(v string) {
	_VERSION_ = v
}
