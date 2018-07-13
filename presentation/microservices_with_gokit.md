# Building Microservices in Go Using Go kit

## Rob Warner

* GitHub: @hoop33
* Twitter: @hoop33
* Email: hoop33@gmail.com
* Blog: grailbox.com

![right,fit,filtered](gokit-logo-header.png)

---

## Agenda

* Just Enough Go
* Introduce Go kit
* Build a CLI CRUD app
* Add HTTP
* Add gRPC
* Build a Node client

---

## Just Enough Go

* "C 2.0"
* Cross-platform
* Compiled
* Garbage collected
* Search `golang` instead of `go`
* https://golang.org

![right,filtered](Go-Logo_Black.png)

---

```go
type Person struct {
  Name string
}

func NewPerson(name string) (*Person, error) {
  if name == "" {
    return nil, errors.New("Name required")
  }
  return &Person{ Name: name }, nil
}

func (p *Person) Greet() {
  fmt.Println("Hello", p.Name)
}
```

---

## Interfaces

```go
type Greetable interface{
  Greet()
}
```

### Everything conforms to an empty interface

```go
interface{}
```

---

## Go kit

> "A toolkit for microservices"

* Created by Peter Bourgon (https://peter.bourgon.org/about/)
* Collection of tools
* Separation of concerns
* https://gokit.io

---

## Step 0: The Data

```sql
CREATE TABLE IF NOT EXISTS players (
  id SERIAL PRIMARY KEY,
  name TEXT,
  number TEXT,
  position TEXT,
  height TEXT,
  weight TEXT,
  age TEXT,
  experience INTEGER,
  college TEXT
)
```

https://www.jaguars.com/team/players-roster/

![right,fit](jacksonville-jaguars-logo-transparent.png)

---

## Step 1: A Service

```go
type Service interface {
	ListPlayers(context.Context, string) ([]models.Player, error)
	GetPlayer(context.Context, int) (*models.Player, error)
	SavePlayer(context.Context, *models.Player) (*models.Player, bool, error)
	DeletePlayer(context.Context, int) error
}
```

---

## Step 2: Add Logger

--- 

## Step 3: Add Endpoints

```go
type Endpoint func(ctx context.Context, 
  request interface{}) (response interface{}, err error)
```

---

## Step 4: Add HTTP Transport

---

## Step 5: Protocol Buffers

* https://developers.google.com/protocol-buffers/
* Install `protoc` compiler
* Install language-specific plugin
* Create `.proto` file(s)

---

## Step 6: Add gRPC Transport

* https://grpc.io
* https://github.com/ktr0731/evans

---

## Next Steps

* Security
* Other transports?
* **Go Programming Blueprints, Second Edition** by Mat Ryer
* https://gokit.io

---

## Questions?

* GitHub: @hoop33
* Twitter: @hoop33
* Email: hoop33@gmail.com
* Blog: grailbox.com
