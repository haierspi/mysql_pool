package Test

import (
	"fmt"
	Mysql "gitee.com/tym_hmm/mysql_pool"
	"gitee.com/tym_hmm/mysql_pool/Test/Field"
)

type T2Model struct {
	Mysql.BaseManyModel
	Field.Field2Test
}

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
