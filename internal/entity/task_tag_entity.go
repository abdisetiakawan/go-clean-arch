package entity

type TaskTag struct {
	ID     uint `gorm:"column:id;primaryKey;autoIncrement"`
	TaskId uint `gorm:"column:task_id;not null;index"`
	TagId  uint `gorm:"column:tag_id;not null;index"`
	Task   Task `gorm:"foreignKey:task_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tag    Tag  `gorm:"foreignKey:tag_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (TaskTag) TableName() string {
	return "task_tags"
}