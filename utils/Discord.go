package utils

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"servers_status/controllers"
	"servers_status/model"
	"strings"
	"time"
)

type DiscordHTTPError struct {
	StatusCode int
	Body       map[string]interface{} // O una estructura más específica si la conoces
}

var (
	MessageID   string
	UpdateDelay = 15 * time.Second // Frecuencia de actualización del mensaje
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Responde con Pong!",
	},
	{
		Name:        "crear",
		Description: "Crea un nuevo servidor a mostrar, recuerda colocar IP ORDEN",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "ip",
				Description: "IP",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "orden",
				Description: "Orden",
				Required:    true,
			},
		},
	},
	{
		Name:        "editar",
		Description: "Editar un servidor registrado, recuerda colocar IP ORDEN",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "ip",
				Description: "IP Actual",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "ip_nueva",
				Description: "IP Nueva",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "orden",
				Description: "Orden",
				Required:    false,
			},
		},
	},
	{
		Name:        "eliminar",
		Description: "Elima un servidor registrado, recuerda colocar IP",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "ip",
				Description: "IP",
				Required:    true,
			},
		},
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pong!",
			},
		})
		if err != nil {
			log.Println("Error responding to ping:", err)
		}
	},
	"crear": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		guilid := i.GuildID
		channel := i.ChannelID
		message_id := i.Interaction.ID
		ip := ""
		orden := ""
		options := i.ApplicationCommandData().Options

		for _, option := range options {
			switch option.Name {
			case "ip":
				ip = option.StringValue()
			case "orden":
				orden = option.StringValue()
			}
		}
		//Buscar si la ip mandada existe
		servidor, err := GetServerInfo(ip)
		if err != nil {
			SendMsg(s, i, err.Error())
			return
		}

		//Buscar si el canal ya se encuentra registrado
		channel_info := controllers.GetChannel(channel)
		msg := ""
		if channel_info.DiscordChannelId == "" {
			//Si no encuentra registrado, se procede a crear
			fmt.Println("Se ha registrado el channel id: " + channel)
			channel_info, msg = controllers.CreateChannel(guilid, channel, message_id)
			if channel_info.DiscordChannelId == "" {
				SendMsg(s, i, msg)
				return
			}
		}
		//Actualizar la informacion del canal con el nuevo id del mensaje a eliminar
		channel_info.DiscordMessageIdDelete = message_id
		controllers.UpdateChannel(channel_info)

		//Buscar si el servidor se encuentra registrado en el canal
		server_found := controllers.GetServer(fmt.Sprintf("%s", channel_info.Uuid), ip)
		if server_found.ServerIp != "" {
			SendMsg(s, i, "El servidor que se intenta ingresar ya se encuentra registrado para este canal")
			return
		}
		//Se procede a registrar los servidores al canal
		var server model.Servers
		server.ChannelId = channel_info.Uuid
		server.ServerIp = ip
		server.ServerOrder = orden
		server.ServerName = servidor.Name
		_, res := controllers.CreateServer(server)
		if res != "" {
			SendMsg(s, i, res)
			return
		}
		//Si se logro, se procede a enviar un mensaje indicando que el registro fue exitoso
		SendMsg(s, i, "Servidor registrado")
	},
	"editar": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		channel := i.ChannelID
		ip := ""
		ip_nueva := ""
		orden := ""

		options := i.ApplicationCommandData().Options
		for _, option := range options {
			switch option.Name {
			case "ip":
				ip = option.StringValue()
			case "ip_nueva":
				ip_nueva = option.StringValue()
			case "orden":
				orden = option.StringValue()
			}
		}

		//Buscar si la ip mandada existe
		servidor, err := GetServerInfo(ip)
		if err != nil {
			SendMsg(s, i, err.Error())
			return
		}

		//Buscar el canal
		channel_info := controllers.GetChannel(channel)

		if channel_info.DiscordChannelId == "" {
			//Si el canal buscado no fue encontrado, se devuelve error
			SendMsg(s, i, "No hemos encontrado los registros en nuestra base de datos")
		}

		//Buscar si el servidor se encuentra registrado en el canal
		server_found := controllers.GetServer(fmt.Sprintf("%s", channel_info.Uuid), ip)
		if server_found.ServerIp == "" {
			SendMsg(s, i, "El servidor que se intenta editar no se encuentra registrado en este canal")
			return
		}

		//Si fue encontrada la ip, se procede a editar con la nueva ip en caso de que se haya mandando
		if ip_nueva != "" {
			server_found.ServerIp = ip_nueva
		}
		//Si orden fue mandado se actualiza
		if orden != "" {
			server_found.ServerOrder = orden
		}

		//Se actualiza el nombre del servidor
		server_found.ServerName = servidor.Name
		msg := controllers.UpdateServer(server_found)

		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		if err != nil {
			log.Println("Error al responder un mensaje:", err)
		}
	},
	"eliminar": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		channel := i.ChannelID
		ip := ""
		options := i.ApplicationCommandData().Options
		for _, option := range options {
			switch option.Name {
			case "ip":
				ip = option.StringValue()
			}
		}

		//Buscar el canal
		channel_info := controllers.GetChannel(channel)

		if channel_info.DiscordChannelId == "" {
			//Si el canal buscado no fue encontrado, se devuelve error
			SendMsg(s, i, "No hemos encontrado los registros en nuestra base de datos")
		}

		//Buscar si el servidor se encuentra registrado en el canal
		server_found := controllers.GetServer(fmt.Sprintf("%s", channel_info.Uuid), ip)
		if server_found.ServerIp == "" {
			SendMsg(s, i, "El servidor que se intenta eliminar no se encuentra registrado en este canal")
			return
		}

		msg := controllers.DeleteServer(server_found.ChannelId, server_found.ServerIp)

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
		if err != nil {
			log.Println("Error al eliminar un servidor:", err)
		}
	},
}

func Parametros(s *discordgo.Session) {
	//Se asignan los Parametros
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(os.Getenv("DISCORD_APLICATION_ID"), "", v)
		if err != nil {
			log.Fatalf("No se pudo crear el comando '%s': %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	fmt.Println("Comandos registrados exitosamente.")
}

// ready se ejecuta cuando el bot se conecta a Discord.
func Ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Printf("Bot %s está listo.\n", s.State.User.Username)
	Parametros(s)
	ticker := time.NewTicker(UpdateDelay)
	defer ticker.Stop()
	for range ticker.C {
		//Buscar todos los canales
		channels := controllers.GetChannels()
		for _, channel := range channels {
			//Si el servidor tiene pendiente eliminar mensajes
			if channel.DiscordMessageIdDelete != "" {
				err := deleteAllMessagesExcept(s, channel.DiscordChannelId, channel.DiscordMessageId)
				if err != nil {
					log.Printf("Error al eliminar los mensajes (2): %v", err)
					VerfiicarRespuesta(err.Error(), channel.DiscordChannelId)
					continue
				}
				//Actualizar la informacion del canal con el nuevo id del mensaje a eliminar
				channel.DiscordMessageIdDelete = ""
				controllers.UpdateChannel(channel)
			}

			//Se procede a validar si no tiene un mensaje para editar
			if channel.DiscordMessageId == "" {
				//Se manda un mensaje y se captura el ID
				initialMessage, err := s.ChannelMessageSend(channel.DiscordChannelId, "Mensaje inicial...")
				if err != nil {
					log.Printf("Error al enviar el mensaje inicial: %v", err)
					VerfiicarRespuesta(err.Error(), channel.DiscordChannelId)
					continue
				}
				MessageID = initialMessage.ID
				fmt.Printf("Mensaje inicial creado con ID: %s en el canal: %s\n", MessageID, channel.DiscordChannelId)
				//Se actualiza el registro con el nuevo mensaje id
				channel.DiscordMessageId = MessageID
				res, errores := controllers.UpdateChannel(channel)
				if !res {
					//Si ocurrio un error se envia un mensaje
					_, err = s.ChannelMessageSend(channel.DiscordChannelId, errores)
					if err != nil {
						log.Printf("Error al enviar el mensaje inicial: %v", err)
						VerfiicarRespuesta(err.Error(), channel.DiscordChannelId)
					}
					continue
				}
			}

			//Se procede a buscar los servidores asociados al canal
			servers := controllers.GetServersChannel(channel.Uuid)
			if len(servers) > 0 {
				go UpdateMessageLoop(s, channel, servers)
			} else {
				maxIntentos := 20
				if channel.ContErr <= maxIntentos {
					s.ChannelMessageSend(channel.DiscordChannelId, fmt.Sprintf("No hemos encontrado un servidor a enlistar, se llevan %d intentos fallidos, luego de %d intentos fallidos, dejare de enviar mensajes", channel.ContErr, maxIntentos))
					channel.ContErr = channel.ContErr + 1
					controllers.UpdateChannel(channel)
				} else {
					s.ChannelMessageSend(channel.DiscordChannelId, "Se ha eliminado de la lista de seguimiento")
					controllers.DeleteChannels(channel.DiscordGuildId)
				}
			}
		}
	}
}

// updateMessageLoop es el bucle principal que actualiza el mensaje de Discord periódicamente.
func UpdateMessageLoop(s *discordgo.Session, channel model.Channels, servers []model.Servers) {
	var info_servers []string
	for _, server := range servers {
		serverAddress := server.ServerIp

		info, err := GetServerInfo(serverAddress)
		if err != nil {
			fmt.Println("Error al consultar el servidor:", serverAddress, err)
			info.Name = server.ServerName
			info.Address = serverAddress
			info.Players = 0
			info.MaxPlayers = 0
			info.Map = "N/A"
		}
		server1 := fmt.Sprintf("%s :arrow_right: %s (%d/%d) Mapa: %s", info.Name, info.Address, info.Players, info.MaxPlayers, info.Map)

		info_servers = append(info_servers, server1)
	}
	finalMessage := strings.Join(info_servers, "\n")

	// Genera el nuevo contenido del mensaje
	newMessage := fmt.Sprintf("%s", finalMessage)

	// Edita el mensaje existente
	_, err := s.ChannelMessageEdit(channel.DiscordChannelId, channel.DiscordMessageId, newMessage)
	if err != nil {
		//Si hubo algun error al tratar de editar el mensaje, posiblemente fue elimiando, se procede a dejar limpio el valor en channel para que lo vuelva a crear
		controllers.SetNullMessages(channel)
		VerfiicarRespuesta(err.Error(), channel.DiscordChannelId)
		log.Printf("Error al editar el mensaje: %v", err)
	} else {
		log.Println("Mensaje actualizado exitosamente.")
	}
}

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	}
}

func SendMsg(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		log.Println("Error al enviar el mensaje al crear:", err)
	}
}

func deleteAllMessagesExcept(s *discordgo.Session, channelID string, excludedMessageID string) error {
	const fetchLimit = 100
	var beforeID string

	for {
		messages, err := s.ChannelMessages(channelID, fetchLimit, beforeID, "", "")
		if err != nil {
			return fmt.Errorf("Error al obtener los mensajes: %w", err)
		}
		if len(messages) == 0 {
			break // No more messages in the channel
		}

		messagesToDelete := make([]string, 0)
		for _, msg := range messages {
			if msg.ID != excludedMessageID {
				messagesToDelete = append(messagesToDelete, msg.ID)
			}
			beforeID = msg.ID // For fetching the next batch of older messages
		}

		if len(messagesToDelete) > 0 {
			if len(messagesToDelete) == 1 {
				err = s.ChannelMessageDelete(channelID, messagesToDelete[0])
				if err != nil {
					log.Printf("Error al eliminar un mensaje (1) %s: %v", messagesToDelete[0], err)
					return err
				}
			} else if len(messagesToDelete) > 1 {
				err = s.ChannelMessagesBulkDelete(channelID, messagesToDelete)
				if err != nil {
					log.Printf("Error al eliminar los mensajes (1): %v", err)
					// If bulk delete fails, try individual deletion as a fallback
					for _, msgID := range messagesToDelete {
						err := s.ChannelMessageDelete(channelID, msgID)
						if err != nil {
							log.Printf("Error al eliminar un mensaje (2) %s: %v", msgID, err)
							return err
						}
					}
				}
			}
		}

		// If we fetched less than the limit, we've likely reached the end
		if len(messages) < fetchLimit {
			break
		}
	}

	return nil
}

func GuildDelete(s *discordgo.Session, g *discordgo.GuildDelete) {
	guilid := g.ID
	channel := controllers.GetChannelGuildId(guilid)
	//Se procede a eliminar los servidores
	controllers.DeleteServers(channel.Uuid)
	//Se procede a eliminar el canal
	controllers.DeleteChannels(channel.DiscordGuildId)
}

func VerfiicarRespuesta(response string, channelId string) {
	// 1. Encontrar el inicio del JSON
	startIndex := strings.Index(response, "{")
	if startIndex == -1 {
		fmt.Println("No se encontró el inicio del JSON")
		return
	}

	// 2. Extraer la parte JSON de la cadena
	jsonString := response[startIndex:]

	// 3. Definir una estructura para el JSON
	var errorResponse struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	// 4. Deserializar la cadena JSON en la estructura
	err := json.Unmarshal([]byte(jsonString), &errorResponse)
	if err != nil {
		fmt.Println("Error al deserializar JSON:", err)
		return
	}

	// 5. Acceder al valor del mensaje
	message := errorResponse.Message
	fmt.Println("Mensaje del error:", message)
	switch message {
	case "Missing Access":
	case "Missing Permissions":
		//Si se perdieron los permisos, se procede a eliminar cualquier registro
		//Buscar el canal
		fmt.Println("Se intento eliminar el canal id", channelId)
		channel := controllers.GetChannel(channelId)
		if channel.DiscordChannelId != "" {
			fmt.Println("Se intento borrar el mensaje")
			//Si se encontro, se procede a eliminar los servidores
			controllers.DeleteServers(channel.Uuid)
			//Se elimina el canal
			controllers.DeleteChannels(channel.DiscordGuildId)
		}
		break
	}

}
