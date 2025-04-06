package repository

import (
	"log"
	"servers_status/model"
)

func InsertChannel(data model.Channels) (model model.Channels, err error) {
	result := DB.Create(&data).Find(&data)
	return data, result.Error
}

func GetAllChannels() []model.Channels {
	var channels []model.Channels
	DB.Debug().Model(model.Channels{}).Order("created_at DESC").Find(&channels)
	return channels
}

func GetChannel(query string) model.Channels {
	var channel model.Channels
	DB.Debug().Model(model.Channels{}).Where(query).Order("created_at DESC").Limit(1).Find(&channel)
	return channel
}

func UpdateChannel(data model.Channels) (model.Channels, error) {
	result := DB.Debug().Model(&data).Omit("uuid").Updates(data).Find(&data)
	if result.Error != nil {
		log.Print(result.Error)
		return data, result.Error
	}
	return data, nil
}

func SetNullMessage(model model.Channels) {
	DB.Debug().Unscoped().Model(&model).Updates(map[string]interface{}{"discord_message_id": nil})
}

func DeleteChannel(query string) {
	DB.Debug().Where(query).Model(&model.Channels{}).Unscoped().Limit(1).Delete(&model.Channels{})
}
