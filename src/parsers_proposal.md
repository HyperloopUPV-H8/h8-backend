# Propuesta de como estructurar el parser

Por una parte podemos tener una sola función / gorutina que reciba todos los raw y luego hacemos de mapa a parser y mandamos por un canal de tipo especifico el mensaje descodificado.

El problema con este approach es que estamos obligando a un sitio a conocer todos los tipos de mensajes que llegan del vehiculo a pesar de que no tienen nada que ver unos con otros.

Por otra parte, debido a que diferenciamos los mensajes según el id, tiene sentido que la discriminacion de los menajes suceda en un lugar y a partir de ahí se va al respectivo sitio donde se procesa.

Esta interfaz tiene un método que devuelve las ids que están asociadas a ese parser. Nosotros los vamos a estructurar de manera que solo hay un parser por tipo de id.

```go
type Parser interface {
	Ids() []uint16
}
```

Vehicle tiene un canal de output por cada tipo de datos:

-   Datos placas
-   Ordenes placas
-   Mensajes placas
-   Ordenes GUI

A una go rutina del Vehicle llegan todos los tipos de mensajes posibles.
La función acepta de parametros el raw y la timestamp. Los dos primeros bytes
del raw son siempre.

Si dentro del mapa tenemos los parsers, seria llamar una funcion, osea, seria blocking.
Si dentro del mapa tenemos un canal que lo envia a una gorutina que siempre esta activa,
seria non-blocking pero seria más complejo al tener que gestionar el inicial estas go routinas.

## Primera opción: almacenamos parser

Si almacenamos el parser como un valor en el mapa, no podemos tener metodos que muten el parser. Si lo almacenamos como un puntero, todos usan el mismo y ademas podemos tener metodo que muten el parser.

## Segunda opción: almacenamos un canal

Para que esto funcionase, los parsers tendrian que tener una rutina que fuese `Listen()`. Enviar al canal es bloqueante con lo que en esos términos estaríamos igual que almacenando el parser, tenemos que esperar a que el parseo anterior acabe para poder recibir el valor. Como detalle, se le puede meter un buffer a los canales, lo que suavizaría este.

Inicialmente, lo haría almacenando un puntero al parser ya que seguir la ejecución del código va a ser más sencillo. Si nos encontramos con problemas de rendimiento, podemos usar gorutinas.

La función Parse de cada Parser podría ejecutarse como una gorutina. El acceso a las propiedades del parser debe de ser solo de lectura si no se quieren usar mutexes. De esta manera, desacoplamos el procesamiento del mensaje.

Si todas las gorutinas de una parser envian el resultado por el mismo canal de output, pueden surgir problemas de mantener el orden de los mensajes procesados porque algunas gorutinas pueden acabar antes que otra de procesar el mensaje.

Ejemplo de esperar a una gorutina:
Empieza la primera gorutina, la segunda y la tercera. La segunda y la tercera acaban antes que la primera. Se tienen que esperar a la primera para enviar el mensaje. La diferencia es el tiempo de bloqueo. Esto introduce un pequeño retraso pero no debería colapsar el parser, simplemente es un pequeño retraso que se propaga.

Puede surgir un problema si enviar el resultado bloquea.

En el caso de los mensajes que nos llegan de las placas, si que queremos mantener el orden. Una manera de mantener el orden es pasarle a cada goroutina un wait group. El wait group tendra una cuenta de 1 y este sera decrementado por la gorutina anterior. Dicho de otra manera, bloqueamos justo al final de la gorutina de parseo hasta que la anterior haga `wait.Done()`.

# Tipos

El tipo de llegada de los mensajes siempre es el mismo: un timestamp y el raw que es `[]byte` (estos son los parámetros de la función que recibe todo).

## ¿Cómo devuelven el mensaje sin perder el tipo (sin usar any)?

Vamos a ver como sería añadir un parser al `Vehicle` para ver donde le podemos dar un canal de output al parser:

```go
type Vehicle struct {
    ...
    updateChan chan models.Update
    messageChan models.Message

    parsers map[uint16]Parser/parseFunction
}

func NewVehicle(...) Vehicle {
    parsers := make(Parsers) ó Parsers{}

    updateParser, updateChannel := NewUpdateParser()
    parsers.AddParser(updateParser)

    messageParser, messageChan := NewMessageParser()
    parsers.AddParser(messageParser)


    myVehicle := Vehicle{
        parsers: parsers
        updateChan: updateChan
        messageChan: messageChan
    }
}

type Parsers map[uint16]Parser

func(parsers *Parsers) AddParser(parser *Parser) {
    ids := parser.Ids()
    for _, id := range ids {
        parsers[id] = parser
    }
}
```

Esta es la parte del `Vehicle` que hace que todos los mensajes lleguen por el mismo canal:

```go
type Data struct {
    Raw       []byte
    From      string
	To        string
    Timestamp time.Time
    Seq       uint32
}

func NewVehicle() Vehicle {
    ... // Crear los parsers

    generalChan := make(chan Data)

    sniffer := createSniffer(generalChan) // Hacer sniffer.Listen() dentro del sniffer
                                          // Igual que pasa con las pipes, que se connectan
                                          // al crearse
    pipes := createPipes(generalChan)
}

// De esta manera es bloqueante porque tengo que esperar a que acabe parser.Parse()

func (vehicle *Vehicle) Listen() {
    for data := range vehicle.generalChan {
        vehicle.parsers[data.Id].Parse(data)
    }
}

```

Ahora vamos a hablar de como procesa cada modulo la información obtenida de data.

## Data -> PacketUpdate

PacketUpdate contiene en `update.Values` un `map[string]any`, donde any es `float64`, `string`, `bool`

Luego se le tiene que añadir el average si la update es numerica. De eso se encarga el update_factory

## Data -> BoardOrder

Usa el mismo parser que `PacketUpdate` pero debe de salir por otro canal porque tiene que procesarse de otra manera. Otra manera de hacerlo es mandandolo por el mismo sitio y luego discriminar fuera el tipo que es el función de la id. Esto tiene sentido porque son identicos en todos los aspectos estructurales, solo que decidimos tratarlos de manera diferente luego.

## Data -> Message

`Message` se genera a partir del string al completo, no hay que incluir nada de información externa.
