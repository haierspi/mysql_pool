package Field

type Field2Test struct {
	Id   int    `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (t Field2Test) TableName() string {
	return "task_config"
}

func (t *Field2Test) GetId() int      { return t.Id }
func (t *Field2Test) GetName() string { return t.Name }
