package jwtUtil

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetToken(t *testing.T) {
	Convey("测试获取token", t, func() {
		account := "zipper"
		expire := 100
		key := "admin@cloudminds.com"
		tokenStr, err := GetToken(account, expire, key)
		t.Log(tokenStr)
		So(err, ShouldBeNil)
		tokenInfo, err := ParseToken(tokenStr, key)

		t.Log(tokenInfo)
		So(err, ShouldBeNil)

	})
}
