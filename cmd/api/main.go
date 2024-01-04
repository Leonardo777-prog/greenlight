package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// declaro un string que contiene el numero de version. luego generaremos esto automaticamente
// en tiempo de contrucccion, pero por ahora solo nesesitamos almacenar la version

const version string = "1.0.0"

// definicion una structiura de configuracion para manterner todas las ajuestes de configuracion
// por ahora, las ajustes de configuracion solo seran el puerto de ret que quresmos usear y el nombre actual de entoro de la aplicacion
// (development, production,testin, etc.)
// leeremos estos valores por la liena de comando cuando se inicie la app

type config struct {
	port int
	env  string
}

// definimos una struct de applicacion para manterner las dependencias de nuestros controladores
// http elperls y middlewarePor el momento esto sólo contiene una copia de la estructura de configuración y un
// logger, pero crecerá para incluir muchos más a medida que avance nuestra compilación.

type application struct {
	config config
	logger *slog.Logger
}

func main() {

	// declaramos una instancia de la struct de configuaracion

	var configuration config

	// leemos los calores pasados por consola de coando
	// por defecto usamos el puerto 4000 y el entorno de development

	flag.IntVar(&configuration.port, "port", 4000, "api server port")
	flag.StringVar(&configuration.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Inicializa un nuevo registrador structura que escribe entradas de registro en la salida estándar
	// flujo, formateado como un objeto JSON.

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// declaramos una instancia de la struct application, contiene la configuracion y el logger

	app := &application{
		config: configuration,
		logger: logger,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", configuration.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", server.Addr, "env", configuration.env)

	err := server.ListenAndServe()

	logger.Error(err.Error())

	os.Exit(1)
}
