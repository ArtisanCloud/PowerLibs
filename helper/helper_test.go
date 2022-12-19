package helper

import (
	"github.com/ArtisanCloud/PowerLibs/v2/fmt"
	"testing"
)

func Test_CheckPassword(t *testing.T) {

	content := "111111"

	// sha256 明文
	encodedPassword, _ := EncodePassword(content)
	fmt.Dump(encodedPassword)

	// Hash编码
	hashedPassword, _ := HashPassword(content)
	fmt.Dump(hashedPassword)

	// 目标保存在数据库的密码
	targetPassword, _ := EncodePlainPassword(content)
	//fmt.Dump(targetPassword)

	result, _ := CheckPassword(targetPassword, encodedPassword)
	fmt.Dump(result)

}
