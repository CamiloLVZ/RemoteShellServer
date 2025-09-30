package main

import (
	"ServerOper/functions"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	//Ruta del Archivo que contiene parametros de configuracion
	archivoConfig := "/etc/RemoteShellServer.conf"
	var archivoUser, puerto string
	var intentosLogin, maxUsers int

	//Leer contenido del archivo conf
	textoArchivo, err := os.ReadFile(archivoConfig)
	if err != nil {
		fmt.Println("Error al leer archivo de configuracion:", err)
	}
	//Guardar el contenido del archivo en un string
	stringArchivo := string(textoArchivo)
	//Dividir el archivo en lineas en un array de strings
	lineasArchivo := strings.Split(stringArchivo, "\n")
	var parametros, configs []string
	//Ciclo que recorre el array de lineas del archivo
	for _, str := range lineasArchivo {
		if str == "" {
			continue
		}
		//Dividir linea del archivo en antes y despues del igual
		parametros = strings.Split(str, "=")
		//Guardar consecutivamente parametro y valor en el array configs
		configs = append(configs, parametros[0])
		configs = append(configs, parametros[1])
	}
	//ciclo que recorre el array configs
	for i := 0; i < len(configs); i++ {
		//Si el parametro actual es "usersDatabase", el siguiente elemento en el array será el valor
		if configs[i] == "usersDatabase" {
			archivoUser = strings.TrimRight(configs[i+1], "\r")
			continue
		} else if configs[i] == "intentosLogin" {
			//Si el parametro actual es "intentosLogin", el siguiente elemento en el array será el valor
			intentosLogin, _ = strconv.Atoi(strings.TrimRight(configs[i+1], "\r"))
			continue
		} else if configs[i] == "puerto" {
			//Si el parametro actual es "puerto", el siguiente elemento en el array será el valor
			puerto = ":" + strings.TrimRight(configs[i+1], "\r")
			continue
		} else if configs[i] == "maxUsers" {
			//Si el parametro actual es "maxUsers", el siguiente elemento en el array será el valor
			maxUsers, _ = strconv.Atoi(strings.TrimRight(configs[i+1], "\r"))
			continue
		}
	}

	fmt.Println("============================================")
	fmt.Println("||       Servidor Operativos 2023.2       ||")
	fmt.Println("============================================")
	//Se crea el socket en el puerto obtenido
	socket := functions.CrearSocket(puerto)

	//Recibir tiempo reporte del cliente como un string
	tiempoReporte, err := bufio.NewReader(socket).ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer", err)
		return
	}
	//Se limpia el string
	tiempoReporte = strings.TrimRight(tiempoReporte, "\r\n") + "s" //La "s" indica segundos
	//Se convierte a tipo de dato de tiempo
	tiempoReporte1, err := time.ParseDuration(tiempoReporte)
	if err != nil {
		fmt.Println("Error al convertir la duración:", err)
		return
	}

	for {
		//Recibe string usuario del cliente
		user, err := bufio.NewReader(socket).ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer", err)
			return
		}
		//Se limpia el string usuario
		user = strings.TrimRight(user, "\r\n")

		//Recibe string password del cliente
		password, err := bufio.NewReader(socket).ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer", err)
			return
		}
		//Se limpia el string password
		password = strings.TrimRight(password, "\n")

		//Recibe la opcion que eligió el cliente en string
		opcion, err := bufio.NewReader(socket).ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer", err)
			return
		}
		//Se limpia el string
		opcion = strings.TrimRight(opcion, "\n")

		switch opcion {
		case "login":
			//Se llama la funcion login exitoso que retorna un bool
			loginExitoso := functions.Login(user, password, intentosLogin, &socket, archivoUser)
			if loginExitoso == true {
				//Se ejecuta en segundo plano EnviarReporte()
				go functions.EnviarReporte(&socket, tiempoReporte1)
				//Se ejecuta en el hilo principal RecibeMensaje
				functions.RecibeMensaje(&socket)
			}
		case "registrar":
			functions.RegistrarUsuario(user, password, &socket, archivoUser, maxUsers)
		}
	}
	socket.Close()
}
