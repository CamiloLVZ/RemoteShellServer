package functions

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// FUNCION QUE RECIBE UNA CONTRASEÑA Y LA RETORNA EN HASH EN UN STRING
func HashPassword(password string) string {
	hashPassword := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hashPassword)
}

// FUNCION DE LOGIN RETORNA BOOL
func Login(userInput string, passwordInput string, intentosLogin int, socket *net.Conn, nombreArchivo string) bool {

	//Se obtiene la contraseña ingresada en hash
	hsPassword := HashPassword(passwordInput)
	//Leer el archivo de base de datos de usuarios
	archivo, err := os.ReadFile(nombreArchivo)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return false
	}
	//Convertir a string el contenido del archivo
	stringArchivo := string(archivo)
	//Dividir el contenido del archivo en lineas
	lineasArchivo := strings.Split(stringArchivo, "\n")
	var paramLogin, userAndPswd []string
	//Recorrer las lineas del archivo
	for _, str := range lineasArchivo {
		if str == "" {
			continue
		}
		//Dividir las lineas en usuario y contraseña (separados por :)
		userAndPswd = strings.Split(str, ":")
		//Guardar user y password consecutivamente en un array
		paramLogin = append(paramLogin, userAndPswd[0])
		paramLogin = append(paramLogin, userAndPswd[1])
	}
	//Crear un escritor para el socket
	env := bufio.NewWriter(*socket)
	var userExiste, pswdCorrecta bool
	//Ciclo que recorre el array de usuarios y contraseñas
	for i := 0; i < len(paramLogin); i++ {
		if paramLogin[i] == userInput {
			//Si el parametro actual es el user ingresado, el usuario existe
			userExiste = true
			//Si la siguiente posicion del arreglo es la contraseña ingresada, el login es exitoso
			if paramLogin[i+1] == hsPassword {
				//Enviar mensaje de exito al cliente
				env.WriteString("succes\n")
				env.Flush()
				pswdCorrecta = true
				break
			} else {
				//Si la contraseña ingresada no es correcta, se envia mensaje de fallo al cliente
				env.WriteString("failed\n")
				env.Flush()
				//For que pide contraseña las veces descritas en el archivo de configuracion
				for j := 0; j < intentosLogin; j++ {
					msg := "Password incorrecto, tiene " + strconv.Itoa(intentosLogin-j) + " intentos mas, Ingrese password:\n"
					//Enviar mensaje de intente de nuevo
					env.WriteString(msg)
					env.Flush()
					//Leer contraseña enviada por cliente
					password, err := bufio.NewReader(*socket).ReadString('\n')
					if err != nil {
						fmt.Println("Error al pedir contraseña: ", err)
					}
					//Se limpia la contraseña
					password = strings.TrimRight(password, "\n")
					//Hash a la contraseña
					hsPassword = HashPassword(password)
					//Se comprueba si la contraseña es correcta
					if paramLogin[i+1] == hsPassword {
						env.WriteString("succes\n")
						env.Flush()
						pswdCorrecta = true
						break
					}
				}
				env.WriteString("failed\n")
				env.Flush()
				break
			}
		}
	}
	if userExiste {
		if pswdCorrecta {
			//retorna un login exitoso
			fmt.Println("Sesion iniciada: ", userInput)
			return true
		}
	} else {
		//Envia señal al cliente de que el usuario no existe
		env.WriteString("notUser\n")
		env.Flush()
	}
	return false
}

func RegistrarUsuario(user string, password string, socket *net.Conn, nombreArchivo string, maxUsers int) {
	//Escritor en el socket
	env := bufio.NewWriter(*socket)
	//PROCEDIMIENTO PARA VERIFICAR SI EL USUARIO YA EXISTE
	//SIMILAR AL USADO EN LOGIN
	textoArchivo, err := os.ReadFile(nombreArchivo)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
	stringArchivo := string(textoArchivo)

	lineasArchivo := strings.Split(stringArchivo, "\n")
	if len(lineasArchivo) >= maxUsers {
		env.WriteString("Error al registrar: Maximo de Usuarios alcanzados\n")
		env.Flush()
		return
	}

	var paramLogin, userAndPswd []string
	for _, str := range lineasArchivo {
		if str == "" {
			continue
		}
		userAndPswd = strings.Split(str, ":")
		paramLogin = append(paramLogin, userAndPswd[0])
		paramLogin = append(paramLogin, userAndPswd[1])
	}
	var userExiste bool
	for i := 0; i < len(paramLogin); i++ {
		if paramLogin[i] == user {
			userExiste = true
			break
		}
	}
	if userExiste {
		//Si el usuario existe, se envia mensaje al cliente
		env.WriteString("Error al registrar: Usuario existente\n")
		env.Flush()
		return
	} //Si no existe, se procede a registrar el usuario

	// Hash de la contraseña
	hsPassword := HashPassword(password)
	// Crear el mensaje a guardar en el archivo
	mensaje := []byte(user + ":" + hsPassword + "\n")
	// Abrir el archivo en modo append, creándolo si no existe
	archivo, err := os.OpenFile(nombreArchivo, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer archivo.Close()

	// Escribir en el archivo
	_, err = archivo.Write(mensaje)
	if err != nil {
		fmt.Println("Error al escribir en el archivo")
		return
	}

	env.WriteString("Usuario Registrado correctamente\n")
	env.Flush()
}
