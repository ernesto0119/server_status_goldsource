package controllers

import (
	"fmt"
	"servers_status/model"
	"servers_status/repository"
)

func CreateChannel(channel_id string, message_id string) (data model.Channels, res string) {
	var channel model.Channels
	channel.DiscordChannelId = channel_id
	channel.DiscordMessageIdDelete = message_id
	response, err := repository.InsertChannel(channel)
	if err != nil {
		fmt.Println(err)
		return model.Channels{}, "No se posible crear el servidor, intente de nuevo"
	}
	return response, "Servidor agregado a la lista"
}

func GetChannels() []model.Channels {
	return repository.GetAllChannels()
}

func GetChannel(channel_id string) model.Channels {
	query := "discord_channel_id = '" + channel_id + "'"
	return repository.GetChannel(query)
}

func UpdateChannel(data model.Channels) (ok bool, res string) {
	_, err := repository.UpdateChannel(data)
	if err != nil {
		return false, "No fue posible completar un proceso importante, se procedera a intentar nuevamente."
	}
	return true, "Actualizado con exito"
}
