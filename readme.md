#### mysql 连接池封装

### gorm 在线文档
[gorm文档=>](https://learnku.com/docs/gorm/v2)

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

3.构建gorm字段 与数据据对应
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

