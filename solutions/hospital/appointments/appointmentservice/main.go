package main

import (
	"log"
)

func main() {
	s := newServer(newAppointmentStore())
	err := s.GRPCListenBlocking(":60001")
	if err != nil {
		log.Fatalf("Error starting rest-notification server: %s", err)
	}
}
