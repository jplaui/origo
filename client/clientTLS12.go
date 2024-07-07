package main

import (
	r "client/request"

	"flag"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main2() {

	// logging settings
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// checks logging flag if program is called as ./main.go -debug
	debug := flag.Bool("debug", false, "sets log level to debug.")

	// checks for handshake only -hsonly flag
	hsonly := flag.Bool("hsonly", false, "establishes tcp tls session with server and returns.")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// activated check
	log.Debug().Msg("Debugging activated.")

	req := r.NewRequest()
	_, err := req.Call2(*hsonly)
	if err != nil {
		log.Error().Msg("req.Call()")
	}

}
