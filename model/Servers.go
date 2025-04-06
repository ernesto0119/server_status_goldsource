package model

import "github.com/google/uuid"

type Servers struct {
	Base
	ChannelId   uuid.UUID `sql:"index" json:"channel_id" validate:"required"`
	Channel     Channels  `gorm:"foreignKey:ChannelId" json:"channel" validate:"-"`
	ServerIp    string    `gorm:"size:50" validate:"required,alpha" sql:"index" json:"server_ip"`
	ServerName  string    `gorm:"size:50" validate:"alpha" sql:"index" json:"server_name"`
	ServerOrder string    `gorm:"size:50" validate:"alpha" sql:"index" json:"server_order"`
	ContErr     int       `gorm:"default:0" validate:"required,numeric" sql:"default=0" json:"err"`
}
