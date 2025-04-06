package model

type Channels struct {
	Base
	DiscordGuildId         string `gorm:"size:50" validate:"required,alpha" sql:"index" json:"discord_guild_id"`
	DiscordChannelId       string `gorm:"size:50" validate:"required,alpha" sql:"index" json:"discord_channel_id"`
	DiscordMessageId       string `gorm:"size:50" validate:"alpha" sql:"index" json:"discord_message_id"`        //ID mensaje original
	DiscordMessageIdDelete string `gorm:"size:50" validate:"alpha" sql:"index" json:"discord_message_id_delete"` //ID mensaje que debe ser borrado
	ContErr                int    `gorm:"default:0" validate:"required,numeric" sql:"default=0" json:"err"`
}
