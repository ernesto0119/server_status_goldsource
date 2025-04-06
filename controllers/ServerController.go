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

func UpdateServer(data model.Servers) (res string) {
	_, err := repository.UpdateServer(data)
	if err != nil {
		return "No fue posible completar un proceso importante, se procedera a intentar nuevamente."
	}
	return "Actualizado con exito"
}

func DeleteServer(channel uuid.UUID, ip string) (response string) {
	query := "channel_id = '" + fmt.Sprintf("%s", channel) + "' and server_ip = '" + ip + "'"
	res := repository.DelServer(query)
	if res != nil {
		fmt.Println(res)
		return "No fue posible eliminar el servidor, intente nuevamente."
	}
	return "Servidor Eliminado"
}

func DeleteServers(channel uuid.UUID) {
	query := "channel_id = '" + fmt.Sprintf("%s", channel) + "'"
	_ = repository.DelServer(query)
}
