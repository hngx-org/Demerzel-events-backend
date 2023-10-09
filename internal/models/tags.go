package models

type Tag struct {
	Id    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title string `json:"title"`

	Groups []Group `json:"groups,omitempty" gorm:"many2many:group_tags"`
}
