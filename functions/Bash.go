package functions

import (
	"fmt"
	"os/exec"
	"strings"
)

func EjecutarComando(comando string) string {
	//Limpiar el string
	comando = strings.TrimRight(comando, "\r\n")
	comando = strings.TrimLeft(comando, "\b")
	//Dividir el comando en palabras en un array
	array_comando := strings.Fields(comando)

	// Ejecutar comando (Para windows, cambiar a linux)
	shell := exec.Command(array_comando[0], array_comando[1:]...)

	//Se recibe la salida del comando (array de bytes)
	salida, err := shell.Output()
	if err != nil {
		return "Comando No valido\n"
	}
	//Retornar la salida convertida a string
	return string(salida)
}

func GenerarReporte(comando string) string {
	// Ejecutar comando
	shell := exec.Command("bash", "-c", comando)

	// Se recibe la salida del comando (array de bytes)
	salida, err := shell.Output()
	if err != nil {
		return fmt.Sprintf("Error al ejecutar comando: %s\n", err)
	}

	// Retornar la salida convertida a string
	return string(salida)
}
