package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/airtonGit/version"
	_ "github.com/joho/godotenv/autoload"
)

//ServiceEndpointInfo configuração minima para cada serviço
//Formato do JSON para informar na VARIAVEL de AMBIENTE lista dos serviços configurados
type ServiceEndpointInfo struct {
	Name string `json:"service"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Path string `json:"path"`
}

//Service agrupa informação do endpoint e http handler
type Service struct {
	ServiceEndpointInfo
	handlerFunc http.HandlerFunc
}

// //GatewayConfig configuracao total do GW, lista de serviços, etc...
// type GatewayConfig struct {
// 	Services []Service
// }

type apiGateway struct {
	log      *Multilogger
	Services []Service
}

//Microservices
// TODO: Seguir modelo do reverse-proxy e carregar lista path/service de ENV
func (a *apiGateway) carregaMicroservicos(basePath string) error {
	// TODO CARREGAR microservicos do ENV VARS
	// Config string [{'service':'auth', host: 'auth', port: 9001, path: '/token'}]
	var services []ServiceEndpointInfo

	if err := json.Unmarshal([]byte(os.Getenv("SERVICES")), &services); err != nil {
		erroMsg := fmt.Sprintf("Erro ao decode json SERVICES: %s, info:%s", os.Getenv("SERVICES"), err.Error())
		return fmt.Errorf(erroMsg)
	}

	for _, service := range services {
		//Adicionar serviço
		itemService := Service{
			ServiceEndpointInfo: service,
			handlerFunc: func(res http.ResponseWriter, req *http.Request) {
				a.log.Debug("GW encaminhando para serviço ", service.Name, service.Host, service.Port, service.Path, req.URL)
				a.forwardMicroservice(fmt.Sprintf("http://%s:%d/%s", service.Host, service.Port, basePath), res, req)
			},
		}
		a.Services = append(a.Services, itemService)
		a.log.Debug(fmt.Sprintf("API GW micro servico adicionado"))
	}
	return nil
}

func (a *apiGateway) forwardMicroservice(endpointHostPort string, res http.ResponseWriter, req *http.Request) {
	//parse the url
	url, err := url.Parse(endpointHostPort)
	if err != nil {
		a.log.Error("forwardMicroservice url.Parse:", err)
	}

	//create de reverse proxy
	a.log.Debug("forwaredMicroservice: URLs", endpointHostPort, url)
	proxy := httputil.NewSingleHostReverseProxy(url)
	a.log.Debug("forward req", req.URL.Host, req.URL.Scheme, req.Host)

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func main() {
	var listenPort string

	listenPort = os.Getenv("LISTEN_PORT")
	version.ParseAll("0.4")

	//Informar o caminho do arquivo de log ex "./logs/api-gw.log" e tag no arquivo
	log, err := NewMultilogger("APIGW")
	if err != nil {
		panic(fmt.Sprintf("APIGW - Não foi possível iniciar log info:%s", err.Error()))
	}

	basePath := os.Getenv("BASE_PATH")
	apiGw := &apiGateway{log: log, Services: []Service{}}
	err = apiGw.carregaMicroservicos(basePath)
	if err != nil {
		log.Error("Carrega serviços info:", err)
		os.Exit(1)
	}

	//Preparando um modelo de carregar rotas de configuracao //http.HandleFunc("/carga-auto/desativar", desativarHandler())
	for _, item := range apiGw.Services {
		http.HandleFunc(item.Path, item.handlerFunc)
	}

	// Aguarda por requisiçoes na porta 9000
	hostPorta := fmt.Sprintf(":%s", listenPort)
	log.Info("Waiting requests", hostPorta)
	if err := http.ListenAndServe(hostPorta, nil); err != nil {
		log.Error("ListenAndServe info:", err)
	}
}
