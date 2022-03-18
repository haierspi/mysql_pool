package Test

import (
	"fmt"
	Mysql "gitee.com/tym_hmm/mysql_pool"
	"testing"
)

func init()  {
	mysqlBuilder:=Mysql.NewMysqlBuilder("test","192.168.1.169", 3306, "root", "fn123456", "test", "utf8", true)
	mysqlBuilder2:=Mysql.NewMysqlBuilder("test2","192.168.1.169", 3306, "root", "fn123456", "ad_task", "utf8", true)
	mysqlPool:=Mysql.NewMysqlPool()
	mysqlPool.SetBuilders(mysqlBuilder, mysqlBuilder2)
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

	t2Model:=new(T2Model)
	err=t2Model.FindOne()
	if err !=nil{
		fmt.Println(err)
	}
}