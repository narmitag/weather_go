FROM golang:1.19-alpine

# RUN apk update && apk upgrade
# RUN apk add git 

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY data/. ./data

RUN go mod download

COPY *.go ./

RUN go build -o /weather

RUN ls -la ./

EXPOSE 8081

CMD [ "/weather", "--http=true", "--datapath=data/ISTGBUCH2/hourly" ]
