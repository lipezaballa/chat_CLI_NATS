# FUNCIONAMIENTO

Este programa crea un chat entre los diferentes usuarios (terminales) que se conecten a un canal a través de un servidor nats ejecutado en docker. Los mensajes enviados por los usuarios se capturan con un input en el terminal y son recibidos por todos los miembros del canal con el formato:

[10/01/2025 11:57:29] usuario1: Hola, ¿cómo estás?

Si un usuario quiere salir del chat notificando que se va, en el input se pone el texto: `exit`, de esta manera aparecerá la salida del canal al resot de usuarios de la siguiente forma:

Felipe left the chat...

Además aparecen notificaciones de uniones de los usuarios:

Felipe joined the chat...

Los usuarios que se incorporan a un chat reciben los mensajes enviados por este canal en la última hora a través de un stream de JetStream.

Además, se porpone en el stream habilitar varios canales, `chat` y `chat.*` para que haya varios chat diferentes y no tengan que estar todos en el mismo necesariamente. Sin embargo, cada terminal solo está desarrollado para conectarse a un canal a la vez.

# EJECUCIÓN

`docker run --name nats -it -p 4222:4222 nats --js`
Este comando ejecuta un servidor NATS con JetStream habilitado para persistencia. Se utiliza -it para ejecutar el contenedor en modo interactivo con un terminal asignado, si no se necesita para ver logs o ejecutar más comandos, se puede sustituir por -d para ejecutarlo en segundo plano.

Creo un stream persistente para el canal

`nats stream add CHAT --subjects=chat,chat.* --storage=file --retention=limits --max-msgs=-1 --max-age=1h`

Se crea un nuevo stream de nombre CHAT.
Los subject introducidos son todos los que comienzan por "chat." y "chat".
El tipo de almacenamiento definido es en fichero, para que sea persistente incluso cuando el servidor se reinicia.
El tipo de retención es "limits" que indica que se almacenan mientras no se excdan los límites definidos, por ejemplo el número de mensajes.
En max-msgs se ha puesto -1, que quiere decir que no hay número máximo de mensajes.
El tiempo máximo de vida de los mensajes en el stream es de 1 hora.
Al ejecutarlo sale una serie de opciones que puedes configurar (algunas de ellas se modifican con los parámetros), con las elecciones por defecto funciona correctamente.

Para ejecutar el programa se utiliza el siguiente comando:
`go run main.go nats://localhost:4222 chat Felipe`
Se pueden abrir los terminales que se quieran con diferentes usuarios.
Se ha seleccionado el canal "chat", sin embargo, también está habilitado de manera persistente cualquier chat por debajo, por ejemplo "char.informal"

# OBJETO

El objeto "ChatClient" tiene la inforamción sel nombre del canal, del nombre de usuario, y la conexión nats. El objetivo de este objeto es agrupar la información necesaria que se tiene que pasar entre funciones para poder separar la funcionalidad sin tener que pasar muchos parámetros por separado.
