package Test

import (
	"fmt"
	Mysql "gitee.com/tym_hmm/mysql_pool"
	"testing"
)

func init()  {
	mysqlBuilder:=Mysql.NewMysqlBuilder("192.168.111.128", 3307, "root", "123456", "miao", "utf8", true)
	mysqlPool:=Mysql.NewMysqlPool()
	mysqlPool.SetBuilder(mysqlBuilder)
	err:=mysqlPool.Init()
	if err !=nil{
		fmt.Println("err",err)
	}
}

func TestMysql(t *testing.T)  {

	tModel:=new(TModel)
	err :=tModel.FindOne()
	if err !=nil{
		fmt.Println(err)
	}
}