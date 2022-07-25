# go项目实战 blog_service

![image-20220716011358924](https://github.com/0x401000Nu1l/golang_blog_service/blob/main/md_img/image-20220716011358924.png)

原书中描述如上 

在GET方法中添加代码会返回提前设定的“错误处理”

但是尝试了几次之后都无法正确测试

![image-20220716011543679](https://github.com/0x401000Nu1l/golang_blog_service/blob/main/md_img/image-20220716011543679.png)

结果如上

在确定代码无误的情况下，我重新审计了代码

![image-20220716011627893](https://github.com/0x401000Nu1l/golang_blog_service/blob/main/md_img/image-20220716011627893.png)

发现了GET函数在不同参数会对应不同的操作

观察测试命令 

```powershell
curl -v http://127.0.0.1:8080/api/v1/articles/1
```

显然是后者，那么调用的即是

![image-20220716011745412](https://github.com/0x401000Nu1l/golang_blog_service/blob/main/md_img/image-20220716011745412.png)

List方法

所以在List方法下面加测试代码，测试之后果然成功

![image-20220716011825242](https://github.com/0x401000Nu1l/golang_blog_service/blob/main/md_img/image-20220716011825242.png)

继续尝试测试GET方法

将命令行修改为如下 


```powershell
curl -v http://127.0.0.1:8080/api/v1/articles
```

发现测试成功
![image-20220716011913790](https://github.com/0x401000Nu1l/golang_blog_service/blob/main/md_img/image-20220716011913790.png)

**记录在2.6章节遇到的bug以及解决方法**

**1.Update方法返回异常**

```powershell
curl -X PUT http://127.0.0.1:8000/api/v1/tags/{1} -F state=0 -F modified_by=eddycjy
```

结果返回 

![img](https://camo.githubusercontent.com/ed3407e45cdb3fe83d00df2d5181f9493d2c5be10e3748c6477fec7832a128c0/68747470733a2f2f692e626d702e6f76682f696d67732f323032322f30372f32352f636539643932616233306338626133612e706e67)

UpdateTagRequest结构体如下

```go
type UpdateTagRequest struct {
	ID         uint32 `form:"id" binding:"required,gte=1"`
	Name       string `from:"name" binding:"min=3,max=100"`
	State      uint8 `form:"state" binding:"required,oneof=0 1"`
	ModifiedBy string `form:"modified_by" binding:"required,min=3,max=100"`
}
```

我根据结构体标签的限定，将

```go
Name       string `from:"name" binding:"min=3,max=100"`
```

修改为

```go
Name       string `from:"name" binding:"min=0,max=100"`
```

并且在参考  [模块开发：标签管理 | Go 语言编程之旅 (eddycjy.com)](https://golang2.eddycjy.com/posts/ch2/06-api-tag-module/#2610-解决问题)

下面的issue后将state字段改为

```go
type UpdatedTagRequest struct {
	Id         uint32 `form:"id" binding:"required,gte=1"`
	Name       string `form:"name" binding:"min=0,max=100"`
	State      uint8 `form:"state" binding:"oneof=0 1"`
	ModifiedBy string `form:"modified_by" binding:"required,min=3,max=100"`
}
```

但是仍然会报错，参考之后采取另一种方法，直接将state字段类型改为指针

 ```go
 type UpdatedTagRequest struct {
 	Id         uint32 `form:"id" binding:"required,gte=1"`
 	Name       string `form:"name" binding:"min=3,max=100"`
 	State      *uint8 `form:"state" binding:"required,oneof=0 1"`
 	ModifiedBy string `form:"modified_by" binding:"required,min=3,max=100"`
 }
 ```

在这里给State指针赋值

```go
func (t Tag) Update(c *gin.Context) {
	newstate := convert.StrTo(c.Param("state")).MustUInt8()
	param := service.UpdateTagRequest{
	ID: convert.StrTo(c.Param("id")).MustUInt32(),
	State: &newstate,
}
//internal\routers\api\v1\tag.go
```

```go
func (svc *Service) UpdateTag(param *UpdateTagRequest) error {
	return svc.dao.UpdateTag(param.ID, param.Name, *param.State, param.ModifiedBy)
}
//internal\service\tag.go
```

![image](https://i.bmp.ovh/imgs/2022/07/26/527bc2ceb85597b3.png)

测试成功

**2.未知panic**

>runtime error: invalid memory address or nil pointer dereference
>...
>***/go-playground/universal-translator@v0.17.0/translator.go:335 (0x8b619a)
>(*translator).C: b = append(b, trans.text[:trans.indexes[0]]...)
>***/go-playground/validator/v10@v10.6.1/translations/zh/zh.go:184 (0xc2a315)
>RegisterDefaultTranslations.func4: c, err = ut.C("min-string-character", f64, digits, ut.FmtNumber(f64, digits))
>***/go-playground/validator/v10@v10.6.1/errors.go:274 (0x8efe82)
>(*fieldError).Translate: return fn(ut, fe)
>***/go-playground/validator/v10@v10.6.1/errors.go:76 (0x8ef975)
>ValidationErrors.Translate: trans[fe.ns] = fe.Translate(ut)
>***/go-programming-tour-book/blog-service/pkg/app/form.go:44 (0xbe8b55)
>BindAndValid: for key, value := range verrs.Translate(trans) {
>***/go-programming-tour-book/blog-service/internal/routers/api/v1/tag.go:104 (0xc33d84)
>Tag.Update: valid, errs := app.BindAndValid(c, &param)
>...

运行之后发生如下报错，定位报错位置

```go
for key, value := range verrs.Translate(trans) { // 这里报错
	errs = append(errs, &ValidError{
		Key:     key,
		Message: value,
	})
}
```

参考 [模块开发：标签管理 | Go 语言编程之旅 (eddycjy.com)](https://golang2.eddycjy.com/posts/ch2/06-api-tag-module/#2610-解决问题)

>locale := c.GetHeader("locale")
>trans, _ := uni.GetTranslator(locale)
>
>这一行如果 header 中没有 locale 就拿不到 translator，我改成
>
>if locale == "" {
>locale = "zh"
>}
>
>就不会报错了，总之要考虑拿不到 translator 的情况

