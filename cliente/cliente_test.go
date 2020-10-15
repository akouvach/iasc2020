package main

import (
	"testing"
)

func TestConnections(t *testing.T) {
	t.Run(
		"conexiones concurrentes", func(t *testing.T) {
			rdo := ClienteAutomatico(5, 2, true)
			if rdo != "OK" {
				t.Fail()
			}

		})
}
