package Mysql

type BuilderInterface interface {
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
}
