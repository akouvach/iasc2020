package main

import (
	"testing"
)

func TestConnections(t *testing.T) {
	t.Run(
		"conexiones concurrentes", func(t *testing.T) {
			rdo := ClienteAutomatico(10, 2, false)
			if rdo != "OK" {
				t.Fail()
			}

		})
}
