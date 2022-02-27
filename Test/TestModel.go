package Test

import (
	"fmt"
	Mysql "gitee.com/tym_hmm/mysql_pool"
	"gitee.com/tym_hmm/mysql_pool/Test/Field"
)

type TModel struct {
	Mysql.BaseModel
	Field.FieldTest
}
func (t *TModel)FindOne() (error) {
	db, err:=t.GetDb()
	if err !=nil{
		return err
	}
//db.Raw("select content from " + m.FieldWpCrawlerComment.TableName() + " order by rand() LIMIT 1").Scan(&crawlerComment)
	testFiled := []*Field.FieldTest{}
	db.Table(t.TableName()).Scan(&testFiled)


	for k,v :=range testFiled{
		fmt.Printf("k:%d, v:%+v\n", k, v)
	}
	return nil
}