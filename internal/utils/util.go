package utils

import "fmt"

var serverVersion string

func SetServerVersion(version, commit string) {
	serverVersion = fmt.Sprintf("%s.%s", version, commit)
}
func GetServerVersion() string {
	return serverVersion
}
