package Mysql

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

/**参见 https://github.com/jinzhu/gorm**/

var (
	MYSQL_ERR_BUILDER = errors.New("mysql builder error")
	DEFULA_LOG_DIR    = "/Log"
	UTFALL_DATE       = "2006-01-02"
	cstZone           = time.FixedZone("CST", 8*3600)
)

/**
带参初始化
*/
func NewMysqlBuilder(host string, port int, user, pwd, dbName, chartSet string, debug bool) BuilderInterface {
	mysqlBuilder := &mysqlBuilder{
		host:     host,
		port:     port,
		user:     user,
		pwd:      pwd,
		dbName:   dbName,
		chartSet: chartSet,
		isDebug:  debug,
	}
	return mysqlBuilder
}

/**
无参数初始化
*/
func NewMysqlBuilderEm() BuilderInterface {
	return &mysqlBuilder{}
}

type mysqlBuilder struct {
	host         string
	port         int
	user         string
	pwd          string
	dbName       string
	chartSet     string
	isDebug      bool
	maxIdleConns int //设置空闲连接池中的最大连接数
	maxOpenConns int //设置与数据库的最大打开连接数

	logDir string //日志保存目录(当isDebug为true时开启)
}

func (mb *mysqlBuilder) SetHost(host string) BuilderInterface {
	mb.host = host
	return mb
}

func (mb *mysqlBuilder) GetHost() string {
	return mb.host
}

func (mb *mysqlBuilder) SetPort(port int) BuilderInterface {
	mb.port = port
	return mb
}
func (mb *mysqlBuilder) GetPort() int {
	return mb.port
}

func (mb *mysqlBuilder) SetUser(user string) BuilderInterface {
	mb.user = user
	return mb
}

func (mb *mysqlBuilder) GetUser() string {
	return mb.user
}

func (mb *mysqlBuilder) SetPwd(pwd string) BuilderInterface {
	mb.pwd = pwd
	return mb
}
func (mb *mysqlBuilder) GetPwd() string {
	return mb.pwd
}

func (mb *mysqlBuilder) SetDbName(dbName string) BuilderInterface {
	mb.dbName = dbName
	return mb
}

func (mb *mysqlBuilder) GetDbName() string {
	return mb.dbName
}

func (mb *mysqlBuilder) SetChartSet(chartSet string) BuilderInterface {
	mb.chartSet = chartSet
	return mb
}
func (mb *mysqlBuilder) GetChartSet() string {
	return mb.chartSet
}
func (mb *mysqlBuilder) SetIsDebug(isDebug bool) BuilderInterface {
	mb.isDebug = isDebug
	return mb
}

func (mb *mysqlBuilder) GetIsDebug() bool {
	return mb.isDebug
}

func (mb *mysqlBuilder) SetMaxIdleConns(maxIdleConns int) BuilderInterface {
	mb.maxIdleConns = maxIdleConns
	return mb
}

func (mb *mysqlBuilder) GetMaxIdleConns() int {
	return mb.maxIdleConns
}

func (mb *mysqlBuilder) SetMaxOpenConns(maxOpenConns int) BuilderInterface {
	mb.maxOpenConns = maxOpenConns
	return mb
}

func (mb *mysqlBuilder) GetMaxOpenConns() int {
	return mb.maxOpenConns
}

func (mb *mysqlBuilder) SetLogDir(logDir string) BuilderInterface {
	mb.logDir = logDir
	return mb
}

func (mb *mysqlBuilder) GetLogDir() string {
	return mb.logDir
}

var mysqlOnce sync.Once
var mysqlInstance *mysqlPool

func NewMysqlPool() *mysqlPool {
	mysqlOnce.Do(func() {
		mysqlInstance = &mysqlPool{}
	})
	return mysqlInstance
}

type mysqlPool struct {
	builder *mysqlBuilder

	dns string
	db  *gorm.DB
}

func (mp *mysqlPool) SetBuilder(builder *mysqlBuilder) {
	mp.builder = builder
}

func (mp *mysqlPool) Init() error {
	if mp.builder == nil {
		return MYSQL_ERR_BUILDER
	}
	mp.dns = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local", mp.builder.GetUser(), mp.builder.GetPwd(), mp.builder.GetHost(), mp.builder.GetPort(), mp.builder.GetDbName(), mp.builder.GetChartSet())
	return mp.initializationMysql()
}

func (mp *mysqlPool) getDb() *gorm.DB {
	return mp.db
}

func (mp *mysqlPool) initializationMysql() error {
	db, err := gorm.Open("mysql", mp.dns)
	if err != nil {
		return err
	}
	mp.db = db
	//设置空闲连接池中的最大连接数。
	maxIdleConns := 10
	//设置与数据库的最大打开连接数
	maxOpenConns := 10
	if mp.builder.maxIdleConns > 0 {
		maxIdleConns = mp.builder.maxIdleConns
	}

	if mp.builder.maxOpenConns > 0 {
		maxOpenConns = mp.builder.maxOpenConns
	}
	mp.db.DB().SetMaxIdleConns(maxIdleConns)
	mp.db.DB().SetMaxOpenConns(maxOpenConns)
	mp.db.SingularTable(true)

	if mp.builder.isDebug {
		logPathDir := DEFULA_LOG_DIR
		if len(strings.TrimSpace(mp.builder.GetLogDir())) > 0 {
			logPathDir = strings.TrimSpace(mp.builder.GetLogDir())
		}
		mp.db.LogMode(true)
		mp.db.SetLogger(log.New(logPath(logPathDir), "\r\n", log.Ldate|log.Ltime))
	}
	return nil
}

func GetDb() (*gorm.DB, error) {
	poll := NewMysqlPool()
	err := poll.Init()
	if err != nil {
		return nil, err
	}
	return poll.getDb(), nil
}

func logPath(logPathDir string) *os.File {
	file, err := os.OpenFile(getLogPath(logPathDir, "mysql", true), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln("fail to create mysql.log file!")
	}
	return file
}

func getLogPath(logPathDir, logFileName string, isDataLogFormat bool) string {
	//	createLock.Lock()
	formatString := ""
	if isDataLogFormat {
		timeStamp := time.Now().In(cstZone)
		formatString = timeStamp.Format(UTFALL_DATE)
		formatString = "_" + formatString
	}
	logDir := getCurrentPath() + string(os.PathSeparator) + logPathDir + string(os.PathSeparator) + logFileName
	logPath := logDir + string(os.PathSeparator) + logFileName + formatString + ".log"
	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		errs := os.MkdirAll(logDir, 777)
		if errs != nil {
			panic(fmt.Sprintf("%s 创建失败", errs))
			return ""
		}
	}
	return logPath
}

/**
aqt获b取当前项目录
*/
func getCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
