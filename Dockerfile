FROM golang:latest

WORKDIR /home

COPY /src /home/

ENTRYPOINT ["go", "run", "."]
