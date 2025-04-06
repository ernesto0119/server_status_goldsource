package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"servers_status/repository"
	"servers_status/utils"
	"syscall"
)

// Variables de configuración
var (
	BotToken string
)

func main() {
	loadEnv()
	repository.Init()

	// Lee las variables de entorno
	BotToken = os.Getenv("DISCORD_BOT_TOKEN")

	if BotToken == "" {
		log.Fatalf("Debes configurar las variables de entorno: DISCORD_BOT_TOKEN, DISCORD_CHANNEL_ID, DISCORD_MESSAGE_ID")
		return
	}

	// Crea una nueva sesión de Discord
	dg, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Error al crear la sesión de Discord: %v", err)
		return
	}

	//Registra el handler para el evento de "listo" (bot conectado)
	dg.AddHandler(utils.Ready)
	dg.AddHandler(utils.InteractionCreate)
	dg.AddHandler(utils.GuildDelete)

	// Inicia la conexión con Discord
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error al abrir la conexión con Discord: %v", err)
		return
	}

	// Espera hasta que se interrumpa el bot (Ctrl+C o señal de terminación)
	fmt.Println("Bot iniciado. Presiona Ctrl+C para salir.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cierra la conexión con Discord de forma limpia
	err = dg.Close()
	if err != nil {
		log.Println("Error al cerrar la conexión de Discord:", err)
	}
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}
