package Mysql

import (
	"crypto/sha256"
	"encoding/hex"
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
func NewMysqlBuilder(builderName, host string, port int, user, pwd, dbName, chartSet string, debug bool) BuilderInterface {
	mysqlBuilder := &mysqlBuilder{
		builderName: builderName,
		host:        host,
		port:        port,
		user:        user,
		pwd:         pwd,
		dbName:      dbName,
		chartSet:    chartSet,
		isDebug:     debug,
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
	builderName  string
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

func (mb *mysqlBuilder) SetBuilderName(builderName string) BuilderInterface {
	mb.builderName = builderName
	return mb
}

func (mb *mysqlBuilder) GetBuildName() string {
	return mb.builderName
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

func (mb *mysqlBuilder) GetHash() (string, error) {
	//keyHash := Hash256(fmt.Sprintf("%s%d%s", mb.GetHost(), mb.GetPort(), mb.GetDbName()))
	if len(strings.TrimSpace(mb.builderName)) == 0 {
		return "", errors.New("builder name was not empty")
	}
	keyHash := hash256(mb.GetBuildName())
	return keyHash, nil
}

var mysqlOnce sync.Once
var mysqlInstance *mysqlPool

func NewMysqlPool() *mysqlPool {
	mysqlOnce.Do(func() {
		mysqlInstance = &mysqlPool{
			builders: []BuilderInterface{},
			conn:     make(map[string]*gorm.DB),
		}
	})
	return mysqlInstance
}

type mysqlPool struct {
	builders []BuilderInterface
	conn     map[string]*gorm.DB
}

func (mp *mysqlPool) SetBuilder(builder BuilderInterface) {
	mp.builders = append(mp.builders, builder)
}

func (mp *mysqlPool) SetBuilders(builders ...BuilderInterface) {
	mp.builders = builders
}

func (mp *mysqlPool) Init() error {
	if mp.builders != nil && len(mp.builders) > 0 {
		for _, v := range mp.builders {
			err := mp.initializationMysql(v)
			if err != nil {
				return err
			}
		}
	} else {
		return MYSQL_ERR_BUILDER
	}

	return nil
	//if mp.builder == nil {
	//	return MYSQL_ERR_BUILDER
	//}
	//mp.dns = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local", mp.builder.GetUser(), mp.builder.GetPwd(), mp.builder.GetHost(), mp.builder.GetPort(), mp.builder.GetDbName(), mp.builder.GetChartSet())
	//return mp.initializationMysql()
}

func (mp *mysqlPool) initializationMysql(builder BuilderInterface) error {
	//如果已存在则不加入
	bHansh, err := builder.GetHash()
	if err != nil {
		return err
	}
	if _, ok := mp.conn[bHansh]; ok {
		return nil
	}
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local", builder.GetUser(), builder.GetPwd(), builder.GetHost(), builder.GetPort(), builder.GetDbName(), builder.GetChartSet())
	db, err := gorm.Open("mysql", dns)
	if err != nil {
		return err
	}

	//设置空闲连接池中的最大连接数。
	maxIdleConns := 10
	//设置与数据库的最大打开连接数
	maxOpenConns := 10
	if builder.GetMaxIdleConns() > 0 {
		maxIdleConns = builder.GetMaxIdleConns()
	}
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)
	db.SingularTable(true)

	if builder.GetIsDebug() {
		logPathDir := DEFULA_LOG_DIR
		if len(strings.TrimSpace(builder.GetLogDir())) > 0 {
			logPathDir = strings.TrimSpace(builder.GetLogDir())
		}
		db.LogMode(true)
		db.SetLogger(log.New(logPath(logPathDir), "\r\n", log.Ldate|log.Ltime))
	}
	mp.conn[bHansh] = db
	return nil
}

func (mp *mysqlPool) getDb(builderName string) (*gorm.DB, error) {
	if len(strings.TrimSpace(builderName))==0{
		return nil, errors.New("builder name can not be empty")
	}
	builderNameHash:=hash256(builderName)

	if v, ok := mp.conn[builderNameHash]; ok {
		return v,nil
	}
	return nil, errors.New("conn cant not in build on builder name :"+ builderName)
}

func GetDb(buildName string) (*gorm.DB, error) {
	poll := NewMysqlPool()
	err := poll.Init()
	if err != nil {
		return nil, err
	}
	return poll.getDb(buildName)
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
生成hash256
*/
func hash256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	sum := h.Sum(nil)
	s := hex.EncodeToString(sum)
	return s
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
