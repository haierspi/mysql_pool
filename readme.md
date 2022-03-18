#### mysql 连接池封装

> v1.0+ 单库连接<br>
> v2.0+ 支持多库链接, 使用多库时注意创建builderName 否则无法找到对应的数据库

### gorm 在线文档
[gorm仓库=>](https://github.com/jinzhu/gorm)
[gorm文档=>](https://learnku.com/docs/gorm/v1``)

### 依赖组件
```
github.com/go-sql-driver/mysql v1.6.0 
github.com/jinzhu/gorm v1.9.16 
```

### 使用方式
1.获取组件

```
go get -u gitee.com/tym_hmm/mysql_pool
```

# v2.0+ 支持多连连接

2.初始化
```
//第一个库
mysqlBuilder:=Mysql.NewMysqlBuilder("test","192.168.1.169", 3306, "root", "fn123456", "test", "utf8", true)
//第二个库
mysqlBuilder2:=Mysql.NewMysqlBuilder("test2","192.168.1.169", 3306, "root", "fn123456", "ad_task", "utf8", true)
mysqlPool:=Mysql.NewMysqlPool()
mysqlPool.SetBuilders(mysqlBuilder, mysqlBuilder2)
err:=mysqlPool.Init()
if err !=nil{
  fmt.Println("err",err)
}
```
3.构建gorm字段与数库字段对应
> 注）必须实现 `TableName string()`方法 返回数据库表名 接口为 BaseFieldInterface
```
type FieldTest struct {
  Id     int    `gorm:"column:id"`
  Name   string `gorm:"column:name"`
  Number int    `gorm:"column:number"`
}

func (t FieldTest) TableName() string {
  return "test"
}

func (t *FieldTest) GetId() int      { return t.Id }
func (t *FieldTest) GetName() string { return t.Name }
func (t *FieldTest) GetNumber() int  { return t.Number }

```

3.构建model
> 构建第二个库中的model

```
type T2Model struct {
	Mysql.BaseManyModel
	Field.Field2Test
}

//设置先择使用哪个builder 必须指定
func (t *T2Model) Build() *T2Model {
	t.SetBuilderName("test2")
	return t
}

func (t *T2Model) FindOne() error {
	db, err := t.Build().GetDb()
	if err != nil {
		return err
	}
	//db.Raw("select content from " + m.FieldWpCrawlerComment.TableName() + " order by rand() LIMIT 1").Scan(&crawlerComment)
	testFiled := []*Field.Field2Test{}
	db.Table(t.TableName()).Scan(&testFiled)
	for k, v := range testFiled {
		fmt.Printf("森test2  k:%d, v:%+v\n", k, v)
	}
	return nil
}
```

> 构建第一个库的model

```
type TModel struct {
	Mysql.BaseManyModel
	Field.FieldTest
}

//设置先择使用哪个builder 必须指定
func (t *TModel) Build() *TModel {
	t.SetBuilderName("test")
	return t
}

func (t *TModel) FindOne() error {
	db, err := t.Build().GetDb()
	if err != nil {
		return err
	}
	//db.Raw("select content from " + m.FieldWpCrawlerComment.TableName() + " order by rand() LIMIT 1").Scan(&crawlerComment)
	testFiled := []*Field.FieldTest{}
	db.Table(t.TableName()).Scan(&testFiled)
	for k, v := range testFiled {
		fmt.Printf("森纟 k:%d, v:%+v\n", k, v)
	}
	return nil
}
```

5. 使用
```
//查询第一个库
tModel:=new(TModel)
err :=tModel.FindOne()
if err !=nil{
  fmt.Println(err)
}

//查询第二个库
t2Model:=new(T2Model)
err=t2Model.FindOne()
if err !=nil{
  fmt.Println(err)
}
```

--------------------------------------------------

# v1.0.+ 单库链接

2.初始化

```
mysqlBuilder:=Mysql.NewMysqlBuilder("192.168.111.128", 3307, "root", "123456", "miao", "utf8", true)
mysqlPool:=Mysql.NewMysqlPool()
mysqlPool.SetBuilder(mysqlBuilder)
err:=mysqlPool.Init()
if err !=nil{
  fmt.Println("err",err)
}
```

3.构建gorm字段与数库字段对应
> 注）必须实现 `TableName string()`方法 返回数据库表名 接口为 BaseFieldInterface
```
type FieldTest struct {
  Id     int    `gorm:"column:id"`
  Name   string `gorm:"column:name"`
  Number int    `gorm:"column:number"`
}

func (t FieldTest) TableName() string {
  return "test"
}

func (t *FieldTest) GetId() int      { return t.Id }
func (t *FieldTest) GetName() string { return t.Name }
func (t *FieldTest) GetNumber() int  { return t.Number }

```

4.构建model

```
type TModel struct {
  Mysql.BaseModel
  Field.FieldTest
}
func (t *TModel)FindOne() (error) {
  db, err:=t.GetDb()
  if err !=nil{
    return err
  }
  testFiled := []*Field.FieldTest{}
  db.Table(t.TableName()).Scan(&testFiled)

  for k,v :=range testFiled{
    fmt.Printf("k:%d, v:%+v\n", k, v)
  }
  return nil
}
```

5.使用案例(可参见demo)

```	
//初始化mysql
mysqlBuilder:=Mysql.NewMysqlBuilder("192.168.111.128", 3307, "root", "123456", "miao", "utf8", true)
mysqlPool:=Mysql.NewMysqlPool()
mysqlPool.SetBuilder(mysqlBuilder)
err:=mysqlPool.Init()
if err !=nil{
	fmt.Println("err",err)
}
//查询数据
tModel:=new(TModel)
err :=tModel.FindOne()
if err !=nil{
  fmt.Println(err)
}

```


### builder 接口说明

```
//设置构建build名称 必须构建否则无法找到对应的库连接
SetBuilderName(builderName string) BuilderInterface
GetBuildName() string

//地址
SetHost(host string) BaseFieldInterface
//获取地址
GetHost() string
//口端
SetPort(port int) BaseFieldInterface

//获取端口
GetPort() int

//用户名
SetUser(user string) BaseFieldInterface
//获取用户名
GetUser() string

//密码
SetPwd(pwd string) BaseFieldInterface
//获取密码
GetPwd() string

//数据库
SetDbName(dbName string) BaseFieldInterface
//获取数据库
GetDbName() string

//设置编码
SetChartSet(chartSet string) BaseFieldInterface
//获取编码
GetChartSet() string

//是否开启日志调试
SetIsDebug(isDebug bool) BaseFieldInterface
//获取是否开启日志调试
GetIsDebug() bool

//设置空闲连接池中的最大连接数(默认为10)
SetMaxIdleConns(maxIdleConns int) BaseFieldInterface
//获取空闲连接池中最大连接数
GetMaxIdleConns() int

//设置与数据库的最大打开连接数(默认为10)
SetMaxOpenConns(maxOpenConns int) BaseFieldInterface
//获取与数据中最大打开连接数
GetMaxOpenConns() int

//日志保存目录(当isDebug为true时开启)
SetLogDir(logDir string) BaseFieldInterface
//获取日志保存目录
GetLogDir() string
```

