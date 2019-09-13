FROM golang:1.12.4-alpine

WORKDIR /app

#API GW ouve na porta 9000
#EXPOSE 9000

ADD gw /app/

ENTRYPOINT [ "./gw" ]