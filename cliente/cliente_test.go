package main

import (
	"testing"
)

func TestConnections(t *testing.T) {
	t.Run(
		"Go can add", func(t *testing.T) {
			cant, err := ClienteAutomatico(5)
			if err != nil {
				t.Fail()
			}
			if cant != 5 {
				t.Fail()
			}

		})
}
