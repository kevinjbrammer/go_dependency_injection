# Dependency Injection in Go

In Go there are four different methods of [dependency injection](https://en.wikipedia.org/wiki/Dependency_injection):

* Constructor Injection
* Method Injection
* Config Injection
* Just-in-Time Injection

[Hands-On Dependency Injection in Go](https://www.amazon.com/Hands-Dependency-Injection-Go-maintain-ebook/dp/B07FSKBWVR/ref=sr_1_3?keywords=dependency+injection+in+go&qid=1571877477&sr=8-3) by Corey Scott provides an excellent survey of each of these methods.

For our purposes, we are going to focus primarily on Constructor Injection, as this is the primary method utilized by most dependency injection frameworks. 

## The Spaceship Application

Inorder to explore dependency injection in Go I have developed a semi-realistic REST web service which retrieves spaceships.  This application has the following endpoints:

| REST API        | Decription |
|-----------------|------------|
| GET /spaceships | Retrieves all spaceships from a sqlite3 database |
| GET /spaceships/:id | Retrieves a spaceship by its ID from a sqlite3 database |

The first iteration of this application serves as a baseline, which uses no dependency injection (code can be found in the *no_di* folder).

The second iteration of this application demonstrates how to implement Constructor Injection manually, i.e. without any dependency injection framework (code can be found in the *manual_di* folder). 

The third iteration of this application builds on top of the second iteration and explores the use of [Uber Dig](https://github.com/uber-go/dig).  The code for this application can be found in the *dig_di* folder. 

The fourth iteration of this application is also an extension of the second iteration but explores the use of [Google Wire](https://github.com/google/wire).  The code for this application can be found in the *wire_di* folder. 

## Constructor Injection


*manual_di/service/spacehip.go*
```
// SpaceshipService ...
type SpaceshipService struct {
	config     config.IConfig
	repository repository.ISpaceshipRepository
}

// NewSpaceshipService constructs a SpaceshipService
func NewSpaceshipService(cfg config.IConfig, rep repository.ISpaceshipRepository) *SpaceshipService {
	return &SpaceshipService{config: cfg, repository: rep}
}
```

As can be seen above the SpaceshipService is having its fields initialized by the parameters passed into its constructor function.  To contrast this let's look at the non dependency injection version of the above code.

*no_di/service/spacehip.go*
```
// SpaceshipService ...
type SpaceshipService struct {
	config     *config.Config
	repository *repository.SpaceshipRepository
}

// NewSpaceshipService constructs a SpaceshipService
func NewSpaceshipService() *SpaceshipService {
	cfg := config.NewConfig()
	rep := repository.NewSpaceshipRepository()
	return &SpaceshipService{config: cfg, repository: rep}
}
```

In the above code the constructor function is creating the SpaceshipService's dependencies and then initializing the SpaceshipService.

It should be apparent from the examples above that the SpaceshipService utiling constructor injection is more testable than its non-dependency injection counterpart.   

Unfortunately Constructor Injection can complicate the creation of objects.  Let's compare the creation of the non-dependency injection spaceship server with the spaceship server utilizing dependency injection.

*no_di/main.go*
```
func main() {
	server := server.NewServer()
	server.Run()
}
```

In the non-dependency injection version the creation of the spaceship server is simple.  In the dependency injection version shown below all of the server's dependencies have to be created and wired together before the server can be constructed. 

*manual_di/main.go*
```
func main() {
	var cfg config.IConfig = config.NewConfig("./spaceships.db", "8000", true)

	db, err := openDatabaseConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	repository := repository.NewSpaceshipRepository(db)
	service := service.NewSpaceshipService(cfg, repository)
	server := server.NewServer(cfg, service)
	server.Run()
}
```

This management of dependencies can become complicated quickly in big projects.  However, there are third party dependency injection frameworks which we can use to alleviate this complexity.

## Uber Dig

[Uber Dig](https://github.com/uber-go/dig) is a reflection based dependency injection framework for Go.  It was not difficult to convert the manual dependency injection Spaceship application to use Dig.  The first step is to build a container and to provide it with provider functions.

*dig_di/main.go*
```
func buildContainer(databasePath, serverPort string, serviceEnabled bool) *dig.Container {
	container := dig.New()

	//Note: dig has the dig.As provider option that is slated for
	//next release which will make it so that you do not have
	//to wrap your constructors with functions that return the
	//appropriate interface.

	container.Provide(func() config.IConfig {
		return config.NewConfig(databasePath, serverPort, serviceEnabled)
	})

	container.Provide(openDatabaseConnection)
	container.Provide(func(db *sql.DB) repository.ISpaceshipRepository {
		return repository.NewSpaceshipRepository(db)
	})
	container.Provide(func(cfg config.IConfig, rep repository.ISpaceshipRepository) service.ISpaceshipService {
		return service.NewSpaceshipService(cfg, rep)
	})
	container.Provide(server.NewServer)

	return container
}
```

A provider function's parameters represent the returned object's dependencies.  Dig uses these provider functions to inorder to construct a specifed object via the container's Invoke method, as can be seen below.

*dig_di/main.go*
```
func main() {
	container := buildContainer("./spaceships.db", "8000", true)

	err := container.Invoke(func(server *server.Server) {
		server.Run()
	})

	if err != nil {
		log.Fatal(err)
	}
}
```

The container's Invoke takes a function whose parameter specifies the object for dig to construct.  In this case it is the *server.Server object.

### Links
* [Uber Dig](https://github.com/uber-go/dig)
* [Dig GoDoc](https://godoc.org/go.uber.org/dig)
* [Dependency Injection in Go](https://blog.drewolson.org/dependency-injection-in-go)

## Google Wire
[Google Wire](https://github.com/google/wire) is a compile-time dependency injection framework.  It uses code generation tool, aptly called wire, for wiring together objects.  With Wire it was fairly simple to convert the manual dependency injection Spaceship application.  The first step is create a new file and to build a provider set and to prov.  This is new file uses build tag (+build wireinject) so that it is not included in any of the go build.

*wire_di/container.go*
```
//+build wireinject

package main

import (
	"spaceships/config"
	"spaceships/repository"
	"spaceships/server"
	"spaceships/service"

	"github.com/google/wire"
)

func buildServer(config config.IConfig) (*server.Server, error) {
	wire.Build(
		openDatabaseConnection,

		repository.NewSpaceshipRepository,
		wire.Bind(new(repository.ISpaceshipRepository), new(*repository.SpaceshipRepository)),

		service.NewSpaceshipService,
		wire.Bind(new(service.ISpaceshipService), new(*service.SpaceshipService)),

		server.NewServer)

	return &server.Server{}, nil
}
```

In the buildServer function we construct a provider set using wire.Build().  A provider like in dig uses the function paramaters to specify the returned object's dependencies.  

One of the things that I like about Wire over Dig is that Wire provides the ability to bind an interface with a concrete type.  This is done by first specifying the provider function for a concrete type and then following it with a interface binding via the wire.Bind() method.  For example:

```
wire.Bind(new(service.ISpaceshipService), new(*service.SpaceshipService))
```

The above call binds the service.ISpaceshipService interface with the *service.SpaceshipService concreate type.

Once the wire file (container.go in our case) has been setup, you then run the wire code generator in the same directory.  This will produce a wire_gen.go file.

*wire_di/wire_gen.go*
```
// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"spaceships/config"
	"spaceships/repository"
	"spaceships/server"
	"spaceships/service"
)

import (
	_ "github.com/mattn/go-sqlite3"
)

// Injectors from container.go:

func buildServer(config2 config.IConfig) (*server.Server, error) {
	db, err := openDatabaseConnection(config2)
	if err != nil {
		return nil, err
	}
	spaceshipRepository := repository.NewSpaceshipRepository(db)
	spaceshipService := service.NewSpaceshipService(config2, spaceshipRepository)
	serverServer := server.NewServer(config2, spaceshipService)
	return serverServer, nil
}
```

The output of Wire for our case is an injector function which we can use in our main to create a spaceship server.

*wire_di/main.go*
```
func main() {
	cfg := config.NewConfig("./spaceships.db", "8000", true)
	server, err := buildServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server.Run()
}
```

### Links
* [Google Wire](https://github.com/google/wire)
* [Compile-time Dependency Injection With Go Cloud's Wire](https://blog.golang.org/wire)
* [Tutorial](https://github.com/google/wire/blob/master/_tutorial/README.md)
* [User Guide](https://github.com/google/wire/blob/master/docs/guide.md)
* [Best Practices](https://github.com/google/wire/blob/master/docs/best-practices.md)
* [FAQ](https://github.com/google/wire/blob/master/docs/faq.md)
* [Go Dependency Injection with Wire](https://blog.drewolson.org/go-dependency-injection-with-wire)