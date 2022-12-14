package Test

import (
	"fmt"
	Mysql "gitee.com/tym_hmm/mysql_pool"
	"gitee.com/tym_hmm/mysql_pool/Test/Field"
)

type TModel struct {
	Mysql.BaseManyModel
	Field.FieldTest
}

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
