package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type ServerInfo struct {
	Header     byte
	Address    string
	Name       string
	Map        string
	Folder     string
	Game       string
	Players    byte
	MaxPlayers byte
}

// Constantes para el protocolo A2S
const (
	a2sInfoRequest        = "\xFF\xFF\xFF\xFFTSource Engine Query\x00"
	a2sInfoResponseHeader = "\xFF\xFF\xFF\xFF\x54"
	s2cChallengeHeader    = "\xFF\xFF\xFF\xFF\x41"
)

func GetServerInfo(address string) (ServerInfo, error) {
	var vacio ServerInfo
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return vacio, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return vacio, err
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// Envía la consulta A2S_INFO inicial
	_, err = conn.Write([]byte(a2sInfoRequest))
	if err != nil {
		return vacio, err
	}

	// Establece un tiempo de espera para la respuesta
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return vacio, err
	}

	// Recibe la respuesta
	buffer := make([]byte, 1400) // Tamaño típico de los paquetes UDP
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return vacio, err
	}

	if n >= 5 && bytes.Equal(buffer[:5], []byte(s2cChallengeHeader)) {
		// Recibimos una respuesta de desafío
		if n >= 9 {
			challenge := buffer[5:9]
			fmt.Printf("Recibido desafío: %v\n", challenge)

			// Construye la segunda petición A2S_INFO con el desafío
			requestWithChallenge := append([]byte(a2sInfoRequest), challenge...)

			// Envía la consulta A2S_INFO con el desafío
			_, err = conn.Write(requestWithChallenge)
			if err != nil {
				return vacio, err
			}

			// Establece un tiempo de espera para la respuesta de información
			err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				return vacio, err
			}

			// Recibe la respuesta de información (segunda lectura)
			buffer = make([]byte, 1400)
			n, _, err = conn.ReadFromUDP(buffer)
			if err != nil {
				return vacio, err
			}

		} else {
			return vacio, fmt.Errorf("paquete de desafío incompleto: % x", buffer[:n])
		}
	}

	reader := bytes.NewReader(buffer[len(a2sInfoResponseHeader):n])

	//Primera lectura
	err = binary.Read(reader, binary.LittleEndian, &vacio.Header)
	if err != nil {
		return vacio, fmt.Errorf("error al leer el header: %w", err)
	}

	_, _ = readNullTerminatedString(reader)
	vacio.Address = address
	vacio.Name, _ = readNullTerminatedString(reader)
	vacio.Map, _ = readNullTerminatedString(reader)
	vacio.Folder, _ = readNullTerminatedString(reader)
	vacio.Game, _ = readNullTerminatedString(reader)

	err = binary.Read(reader, binary.LittleEndian, &vacio.Players)
	if err != nil {
		return vacio, fmt.Errorf("error al leer Players: %w", err)
	}
	err = binary.Read(reader, binary.LittleEndian, &vacio.MaxPlayers)
	if err != nil {
		return vacio, fmt.Errorf("error al leer Players: %w", err)
	}

	return vacio, nil
}

func readNullTerminatedString(r *bytes.Reader) (string, error) {
	var buffer bytes.Buffer
	for {
		var b byte
		err := binary.Read(r, binary.LittleEndian, &b)
		if err != nil {
			return "", err
		}
		if b == 0 {
			break
		}
		buffer.WriteByte(b)
	}
	return buffer.String(), nil
}
