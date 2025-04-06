package repository

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"servers_status/model"
)

func InsertServer(data model.Servers) (model model.Servers, err error) {
	result := DB.Create(&data).Find(&data)
	return data, result.Error
}

func GetServersChannels(channelId uuid.UUID) []model.Servers {
	var servers []model.Servers
	DB.Debug().Model(model.Servers{}).Where("channel_id = '" + fmt.Sprintf("%s", channelId) + "'").Order("server_order ASC").Find(&servers)
	return servers
}

func GetServer(query string) model.Servers {
	var server model.Servers
	DB.Debug().Model(model.Servers{}).Where(query).Order("created_at DESC").Limit(1).Find(&server)
	return server
}

func UpdateServer(data model.Servers) (model.Servers, error) {
	result := DB.Debug().Model(&data).Omit("uuid").Updates(data).Find(&data)
	if result.Error != nil {
		log.Print(result.Error)
		return data, result.Error
	}
	return data, nil
}

func DelServer(query string) error {
	result := DB.Debug().Where(query).Model(&model.Servers{}).Unscoped().Delete(&model.Servers{})
	if result.Error != nil {
		log.Print(result.Error)
		return result.Error
	}
	return nil
}
