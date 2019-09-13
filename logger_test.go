package main

import "testing"

func TestNew(t *testing.T) {
	logger, err := NewMultilogger("teste unitatio")
	if err != nil {
		t.Error("Falha ao criar logger")
	}
	logger.Info("Linha de info")
	logger.Warning("Linha de Warning")
}
