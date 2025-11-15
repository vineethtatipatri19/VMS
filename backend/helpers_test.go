
package main

import "testing"

func TestGenerateLotNumber(t *testing.T) {
	n1 := generateLotNumber()
	n2 := generateLotNumber()
	if n1 == "" || n2 == "" { t.Fatal("empty lot") }
	if n1 == n2 { t.Fatal("lot numbers should be unique") }
}
