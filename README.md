基于 `github.com/pkg/errors` 包，增加对 `error code` 的支持，完全兼容 `github.com/pkg/errors`。

性能跟 `github.com/pkg/errors` 基本持平。


```go
// 包内注册公共的code码
// 1 /code 目录下添加你的code码
// 2 code.go 文件里init()方法注册你的code码
// 3 新版本发布，其他应用import对应版本的包，即可使用

// 使用者注册应用特有的code码
import (
"git.cai-inc.com/support/errors"
)
func TestRegister(t *testing.T) {
    errors.Register(100001,500,"abc")
}
// 使用
func TestWithCode(t *testing.T) {
	err := errors.WithCode(errors.code.ErrDecodingJSON,"abc")
	coder := errors.ParseCoder(err)
	fmt.Println(coder.String())
	fmt.Println(coder.Code())
	fmt.Println(coder.HTTPStatus())
	fmt.Println(coder.Reference())

	fmt.Println("---------------")
	// (# json), (+ detail),()
	fmt.Println(fmt.Sprintf("%+v",err))
	fmt.Println(fmt.Sprintf("%#v",err))
	fmt.Println(fmt.Sprintf("%#+v",err))
}
```
