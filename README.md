docker run --name nats -it -p 4222:4222 nats --js

# Este comando ejecuta un servidor NATS con JetStream habilitado para persistencia. Se utiliza -it para ejecutar el contenedor en modo interactivo con un terminal asignado, si no se necesita para ver logs o ejecutar más comandos, se puede sustituir por -d para ejecutarlo en segundo plano.

# Creo un stream persistente para el canal

nats stream add CHAT --subject=chat.\* --storage=file --retention=limits --max-msgs=-1 --max-age=1h

# Se crea un nuevo stream de nombre CHAT

# Los subject introducidos son todos los que comienzan por "chat."

# El tipo de almacenamiento definido es en fichero, para que sea persistente incluso cuando el servidor se reinicia

# El tipo de tetención es "limits" que indica que se almacenan mientras no se excdan los límites definidos, por ejemplo el número de mensajes

# en max-msgs se ha puesto -1, que quiere decir que no hay número máximo de mensajes

# El tiempo máximo de vida de los mensajes en el stream es de 1 hora

go run main.go nats://localhost:4222 chat Felipe
