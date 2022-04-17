package Mysql

import (
	"github.com/jinzhu/gorm"
	"math"
	"reflect"
)

/**
多表
*/
type BaseManyModel struct {
	buildName string
}

func (b *BaseManyModel) GetDb() (*gorm.DB, error) {
	return GetDb(b.buildName)
}
func (b *BaseManyModel) SetBuilderName(buildName string) {
	b.buildName = buildName
}

/**
获取字段下的gorm
@param structs BaseFieldInterface  字段结构体
@param key string 结构体属性字段
@param hasTable bool 返回时是否需字段前带表名
@param aliasTable string 别名
*/
func (b *BaseManyModel) getFieldTag(structs BaseFieldInterface, key string, hasTable bool, aliasTable string) []string {
	object := reflect.ValueOf(structs)
	myref := object
	typeOfType := myref.Type()
	var dbField []string
	for i := 0; i < myref.NumField(); i++ {
		itemField := typeOfType.Field(i)
		if key != "" && itemField.Name == key {
			tagColumn, err := getColumn(itemField.Tag.Get("gorm"))
			if err == nil {
				if hasTable {
					if aliasTable != "" {
						dbField = []string{aliasTable + ".`" + tagColumn + "`"}
					} else {
						dbField = []string{"`" + structs.TableName() + "`.`" + tagColumn + "`"}
					}

				} else {
					dbField = []string{tagColumn}
				}
			}
			break
		} else {
			tagColumn, err := getColumn(itemField.Tag.Get("gorm"))
			if err == nil {
				if tagColumn != "" && tagColumn != "-" {
					if hasTable {
						if aliasTable != "" {
							dbField = []string{aliasTable + ".`" + tagColumn + "`"}
						} else {
							dbField = append(dbField, "`"+structs.TableName()+"`.`"+tagColumn+"`")
						}

					} else {
						dbField = append(dbField, tagColumn)
					}
				}
			}
		}

	}
	return dbField
}

/**
获取所有 tag
*/
func (b *BaseManyModel) GetAllFieldTag(structs BaseFieldInterface) []string {
	return b.getFieldTag(structs, "", true, "")
}

/**
获取所有 tag 不加table前缀
*/
func (b *BaseManyModel) GetAllFieldTagNotTable(structs BaseFieldInterface) []string {
	return b.getFieldTag(structs, "", true, "")
}

/**
获取指定tag
*/
func (b *BaseManyModel) GetItemTag(structs BaseFieldInterface, key string) string {
	item := b.getFieldTag(structs, key, true, "")
	return item[0]
}

/**
获取指定标签 并转为别名
*/
func (b *BaseManyModel) GetItemTagAliasTable(structs BaseFieldInterface, key string, AliasTable string) string {
	item := b.getFieldTag(structs, key, true, AliasTable)
	return item[0]
}

func (b *BaseManyModel) GetItemTagNotTable(structs BaseFieldInterface, key string) string {
	item := b.getFieldTag(structs, key, false, "")
	return item[0]
}

/**
@title 总页数计算
@param page int 当前页
@param pageSize int 每页显示
@param totalNum int 总记录数
*/
func (b *BaseManyModel) Paginator(page int, pageSize int, totalNum int) int {
	totalPages := int(math.Ceil(float64(totalNum) / float64(pageSize))) //page总数
	return totalPages
}
