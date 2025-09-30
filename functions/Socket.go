package functions

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func CrearSocket(puerto string) net.Conn {

	//Crear la direccion TCP para el socket en el puerto recibido
	tcpApdress, err := net.ResolveTCPAddr("tcp4", puerto)
	if err != nil {
		fmt.Println("Error al resolver la dirección:", err)
		return nil
	}
	//Crear un socket
	socket, err := net.ListenTCP("tcp", tcpApdress)
	if err != nil {
		fmt.Println("Error al escuchar en el puerto:", err)
		return nil
	}

	fmt.Println("Esperando conexion")
	//El socket empieza a esperar una conexion
	socketServer, err := socket.Accept()
	if err != nil {
		fmt.Println("Error al aceptar conexion")
	}
	//Aviso de conexion de cliente
	fmt.Println("Se ha conectado [", socketServer.RemoteAddr(), "]")

	return socketServer
}

// FUNCION QUE ENVIA REPORTES DE RECURSOS PERIODICAMENTE
func EnviarReporte(socket *net.Conn, tiempoReporte time.Duration) {

	for {
		//Espera el tiempo indicado por el usuario
		time.Sleep(tiempoReporte)
		//Crear un nuevo escritor en el socket
		env := bufio.NewWriter(*socket)
		//Se declara el mensaje de reporte (arbitrario)
		comandos := []string{
			`free -m | awk 'NR==2{printf "%.2f\n", $3/$2*100}'`,
			`top -bn1 | grep "Cpu(s)" | awk '{print $2+$4}'`,
			`df -h / | awk 'NR==2 {print $5}' | cut -d'%' -f1`,
		}
		usoMemoria := strings.TrimRight(GenerarReporte(comandos[0]), "\n")
		usoCPU := strings.TrimRight(GenerarReporte(comandos[1]), "\n")
		usoDisco := strings.TrimRight(GenerarReporte(comandos[2]), "\n")

		reporte := "USO DE MEMORIA: " + usoMemoria + "%\n" +
			"USO DE CPU: " + usoCPU + "%\n" +
			"USO DE DISCO: " + usoDisco + "%\n"

		//Se convierte a un array de bytes
		reporteBytes := []byte(reporte)
		//Se obtiene la longitud del reporte
		sizeReporte := uint16(len(reporteBytes))
		//Se carga el reporte en el socket
		err := binary.Write(env, binary.LittleEndian, sizeReporte)
		_, err = env.Write(reporteBytes)
		if err != nil {
			fmt.Println("Error, Cliente desconectado\n", err)
			return
		}
		//Se envia lo cargado en el socket
		env.Flush()
	}
}

// FUNCION QUE RECIBE LOS COMANDO ENVIADOS POR EL USUARIO
func RecibeMensaje(socket *net.Conn) {
	//Lector del socket
	lector := bufio.NewReader(*socket)
	for {
		// Leer la longitud del comando recibido
		var sizeComando uint16
		err := binary.Read(*socket, binary.LittleEndian, &sizeComando)
		if err != nil {
			fmt.Println("Error, Cliente desconectado\n", err)
			return
		}
		// Leer los datos del comando
		comandoBytes := make([]byte, sizeComando)
		_, err = lector.Read(comandoBytes)
		if err != nil {
			fmt.Println("Error, Cliente desconectado\n", err)
			return
		}
		//Convertir el comando a un string
		comando := string(comandoBytes)
		fmt.Println("comando recibido: $", comando)

		//Si el comando es "bye" termina el programa
		if comando == "bye\r\n" {
			os.Exit(0)
			return
		}

		//Ejecutar el comando en la terminal y guardar la salida en un string
		salida := EjecutarComando(comando)
		// Crear un escritor para el socket
		env := bufio.NewWriter(*socket)
		if err != nil {
			fmt.Println("Error, Cliente desconectado")
			return
		}

		// Convertir la salida a bytes
		salidaBytes := []byte(salida)
		// Obtener la longitud de la salida
		sizeSalida := uint16(len(salidaBytes))
		// Escribir la longitud de la salida en el socket
		err = binary.Write(env, binary.LittleEndian, sizeSalida)
		// Escribir los datos de la salida en el socket
		_, err = env.Write(salidaBytes)
		if err != nil {
			fmt.Println("Error, Cliente desconectado\n", err)
			return
		}
		//Enviar lo cargado en el socket
		env.Flush()
		//}
	}
}
