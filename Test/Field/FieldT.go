package Field

type FieldTest struct {
	Id     int    `gorm:"column:id"`
	Name   string `gorm:"column:name"`
	Number int    `gorm:"column:number"`
}

func (t FieldTest) TableName() string {
	return "test"
}

func (t *FieldTest) GetId() int      { return t.Id }
func (t *FieldTest) GetName() string { return t.Name }
func (t *FieldTest) GetNumber() int  { return t.Number }
