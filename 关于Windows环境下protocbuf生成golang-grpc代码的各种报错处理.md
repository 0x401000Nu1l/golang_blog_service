**关于Windows环境下protocbuf生成golang-grpc代码的各种报错处理** 

![image-20220714123736319](C:\Users\Ikhan\AppData\Roaming\Typora\typora-user-images\image-20220714123736319.png)

有如下目录结构，需要通过helloworld.proto文件生成golang代码 或者golang-grpc代码

根据官方示例命令行 

**前提是已经安装完成 protoc 与生成golang代码插件 protoc-go-gen**

```powershell
protoc --go_out=plugins=grpc:. proto/helloworld.proto
# returns
# plugins=grpc/: No such file or directory
protoc --go_out=./ proto/helloworld.proto
#生成普通pb.go代码
protoc --go-grpc_out=./proto/ proto/helloworld.proto
#生成grpc代码

```

