package Mysql

type BuilderInterface interface {
	//地址
	SetHost(host string) BuilderInterface
	//获取地址
	GetHost() string
	//口端
	SetPort(port int) BuilderInterface

	//获取端口
	GetPort() int

	//用户名
	SetUser(user string) BuilderInterface
	//获取用户名
	GetUser() string

	//密码
	SetPwd(pwd string) BuilderInterface
	//获取密码
	GetPwd() string

	//数据库
	SetDbName(dbName string) BuilderInterface
	//获取数据库
	GetDbName() string

	//设置编码
	SetChartSet(chartSet string) BuilderInterface
	//获取编码
	GetChartSet() string

	//是否开启日志调试
	SetIsDebug(isDebug bool) BuilderInterface
	//获取是否开启日志调试
	GetIsDebug() bool

	//设置空闲连接池中的最大连接数(默认为10)
	SetMaxIdleConns(maxIdleConns int) BuilderInterface
	//获取空闲连接池中最大连接数
	GetMaxIdleConns() int

	//设置与数据库的最大打开连接数(默认为10)
	SetMaxOpenConns(maxOpenConns int) BuilderInterface
	//获取与数据中最大打开连接数
	GetMaxOpenConns() int

	//日志保存目录(当isDebug为true时开启)
	SetLogDir(logDir string) BuilderInterface
	//获取日志保存目录
	GetLogDir() string
}
