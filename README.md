# FUNCIONAMIENTO

Este programa crea un chat entre los diferentes usuarios (terminales) que se conecten a un canal a través de un servidor nats ejecutado con docker. Los mensajes enviados por los usuarios se capturan con un input en el terminal y son recibidos por todos los miembros del canal con el formato:

`[10/01/2025 11:57:29] usuario1: Hola, ¿cómo estás?`

Si un usuario quiere salir del chat notificando que se va, en el input se pone el texto: `exit`, de esta manera aparecerá la salida del canal al resto de usuarios de la siguiente forma:

`Felipe left the chat...`

Además aparecen notificaciones de uniones de los usuarios:

`Felipe joined the chat...`

Los usuarios que se incorporan a un chat reciben los mensajes enviados por este canal en la última hora a través de un stream de JetStream. (EXTRA)

Además, se propone en el stream habilitar varios canales, `chat` y `chat.*` para que haya varios chat diferentes y no tengan que estar todos en el mismo necesariamente. Sin embargo, cada terminal solo está desarrollado para conectarse a un canal a la vez.

# EJECUCIÓN

`docker run --name nats -it -p 4222:4222 nats --js`
Este comando ejecuta un servidor NATS con JetStream habilitado para persistencia. Se utiliza -it para ejecutar el contenedor en modo interactivo con un terminal asignado, si no se necesita para ver logs o ejecutar más comandos, se puede sustituir por -d para ejecutarlo en segundo plano.

Para ejecutar el programa de un usuario se utiliza el siguiente comando:
`go run main.go <url> <channel> <user>`
`go run main.go nats://localhost:4222 chat Felipe`

Se incluyen tres parámetros:

-   url: Dirección del servidor NATS
-   channel: Nombre del canal del chat
-   user: Nombre del usuario

    Se pueden abrir los terminales que se quieran con diferentes usuarios.
    Se ha seleccionado el canal "chat", sin embargo, también está habilitado de manera persistente cualquier chat por debajo, por ejemplo "char.informal"

# OBJETO

El objeto "ChatClient" tiene la inforamción del nombre del canal, del nombre de usuario, y la conexión nats. El objetivo de esta estructura es agrupar la información necesaria que se tiene que pasar entre funciones para poder separar la funcionalidad sin tener que pasar muchos parámetros por separado.

# EXPLICACIÓN DEL CÓDIGO

Inicialmente se valida que el número de argumentos es correcto, imprimiendo un error indicando cómo debe ejecutarse el programa en caso de no tener los argumentos necesarios (`go run main.go <url> <channel> <user>`).

Se realiza la conexión con el servidor NATS a partir de la dirección IP que se ha pasado como argumento (`nats.Connect(natsURL)`).

A continuación se intenta acceder al stream CHAT (`_, err = js.StreamInfo(streamName)`), si no existe se configura el JetStream:

-   Se crea un nuevo stream de nombre CHAT.
-   Los subjects introducidos son todos los que comienzan por "chat.", además de "chat".
-   El tipo de almacenamiento definido es en fichero, para que sea persistente incluso cuando el servidor se reinicia.
-   El tipo de retención es "limits" que indica que se almacenan mientras no se excedan los límites definidos, por ejemplo el número de mensajes.
-   En max-msgs se ha puesto -1, que quiere decir que no hay número máximo de mensajes.
-   El tiempo máximo de vida de los mensajes en el stream es de 1 hora.

streamConfig := &nats.StreamConfig{
Name: streamName,
Subjects: channels,
Retention: nats.LimitsPolicy,  
 MaxMsgs: -1,  
 MaxBytes: -1,  
 MaxAge: 1 \* time.Hour,  
 Storage: nats.FileStorage,  
 }

Se suscribe el cliente al channel, se reciben e imprimen los mensajes de la última hora y se envía un mensaje informativo de unión al chat.

sub, err := nc.js.Subscribe(client.Channel, func(msg \*nats.Msg) {
// Show received messages
fmt.Println(string(msg.Data))
}, subOpts...)

Para los envíos de mensajes se utiliza el método "sendMessage" que da formato al mensaje de la forma indicada en los ejemplos de arriba. Se llama al método publish de la conexión nats y este lo envía a través del canal a todos los usuarios que estén suscritos.

if err := client.Nc.Publish(client.Channel, []byte(message)); err != nil {
log.Printf("Error sending message: %v", err)
}

Para salir, si se escribe "exit", se produce un break que sale del bucle del scanner de los inputs de teclado, se llega al final de la ejecución del programa y se cierra la conexión del cliente, mandando un mensaje por pantalla. Se ejecutan entonces los defer que cierran la conexión al servidor NATS y terminan la suscripción del usuario al channel.

# DOCKER-COMPOSE

Se ha añadido la creación del contenedor que ejecuta el servidor NATS para evitar tener que lanzar el contenedor manualmente. Sin embargo, la ejecución de los usuarios que se conectan al chat no se incluyen en el docker-compose para que se vayan creando a demanda en diferentes terminales, en los canales deseados. La instrucción a ejecutar para crear un usuario está indicada en el apartado "EJECUCIÓN" del README.md
