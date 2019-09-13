package main

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) { //basePath string) {

	log, err := NewMultilogger("APIGW TESTING")
	if err != nil {
		t.Error(fmt.Sprintf("Falha ao testar info: %s", err.Error()))
	}
	apiGw := &apiGateway{log, []Service{}}
	apiGw.carregaMicroservicos("basePath")
	// if len(apiGw.Services) < 2 {
	// 	t.Error("Falha ao testar info", apiGw.Services)
	// }
}
