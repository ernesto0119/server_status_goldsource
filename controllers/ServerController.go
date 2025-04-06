package controllers

import (
	"fmt"
	"github.com/google/uuid"
	"servers_status/model"
	"servers_status/repository"
)

func GetServersChannel(channelId uuid.UUID) []model.Servers {
	return repository.GetServersChannels(channelId)
}

func CreateServer(channel model.Servers) (data model.Servers, res string) {
	response, err := repository.InsertServer(channel)
	if err != nil {
		fmt.Println(err)
		return model.Servers{}, "No se posible crear el servidor, intente de nuevo"
	}
	return response, "Servidor agregado a la lista"
}

func GetServer(channel string, ip string) model.Servers {
	query := "channel_id = '" + channel + "' and server_ip = '" + ip + "'"
	return repository.GetServer(query)
}
