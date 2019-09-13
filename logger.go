package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

//Multilogger mantém arquivo e metodos para log
type Multilogger struct {
	*log.Logger
	filename  string
	tag       string
	debugMode bool
}

// NewMultiloggerWithFile em arquivo e tag.
// Inicia arquivo logger.log ou adiciona ao existente, permite
// também especificar string tag padrão no arquivo
func NewMultiloggerWithFile(logfile, tag string) (*Multilogger, error) {

	arquivoLog, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("Nao foi possivel escrever no arquivo erro1:%s", err.Error())
	}

	mw := io.MultiWriter(os.Stdout, arquivoLog)

	logInstance := &Multilogger{
		Logger:    log.New(mw, tag, log.Ldate|log.Lmicroseconds),
		filename:  logfile,
		tag:       tag,
		debugMode: false,
	}
	return logInstance, nil
}

//NewMultilogger cria logger com tag
func NewMultilogger(tag string) (*Multilogger, error) {
	logger := &Multilogger{
		Logger:    log.New(os.Stdout, fmt.Sprintf("%s ", tag), log.Ldate|log.Lmicroseconds),
		tag:       tag,
		debugMode: false,
	}
	return logger, nil
}

//SetDebug configura modo debug
func (l *Multilogger) SetDebug(mode bool) {
	l.debugMode = mode
}

//Warning adiciona nova linha no arquivo de log com rotulo WARNING
//Warning log line "DATETIME TAG WARNING"
func (l *Multilogger) Warning(message string) {
	l.Println("WARNING ", message)
}

//Info adiciona nova linha no log com rotulo INFO
//
//Info log line "DATETIME TAG INFO"
func (l *Multilogger) Info(params ...interface{}) {
	l.Println("INFO	", params)
}

//Fatal log e finaliza
func (l *Multilogger) Fatal(params ...interface{}) {
	l.Logger.Fatal("FATAL ", params)
}

//Error adiciona nova linha no arquivo de log
//
//message é inserida no arquivo de log com rotulo ERROR
func (l *Multilogger) Error(params ...interface{}) {
	l.Println("ERROR ", params)
}

//Debug adiciona nova linha no arquivo de log
//
//TODO: adicionar uma configuração por variavel de ambiente
//que permite ligar/desligar
func (l *Multilogger) Debug(params ...interface{}) {
	if l.debugMode {
		l.Println("DEBUG ", params)
	}
}
