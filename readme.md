# README

## 测试

- 解析非标准DHCP格式的数据包程序是否崩溃 ✅
  - 错误处理 `https://xiaomi-info.github.io/2020/01/20/go-trample-panic-recover/`
- 配置文件 ✅
- 日志文件
  - DHCP MessageType 打印 ✅
  - 日志合并一条后打印
- 设置多个DNS服务器 ✅
- MAC地址大小写与数据库不符
- 广播或者单播数据包
- 广播包同一网卡响应
- 数据库
- 完善其他DHCP响应
- `func (m MessageType) OptBytes() []byte` 用指针不不用有什么区别？`func (e *argError) Error() string`

## 编译

GOOS=windows GOARCH=amd64 CGO_ENABLE=1 go install SHCP
