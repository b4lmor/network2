package cli

import "flag"

func Parse() int {
	port := flag.Int("port", 8080, "Порт для прокси")

	flag.Parse()

	return *port
}
