FROM golang:latest

WORKDIR /home

COPY main.go /home/

RUN ["go", "get", "git.darknebu.la/GalaxySimulator/structs"]
RUN ["go", "get", "github.com/ajstarks/svgo"]
RUN ["go", "get", "github.com/gorilla/mux"]

ENTRYPOINT ["go", "run", "."]
