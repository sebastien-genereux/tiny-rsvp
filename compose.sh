#!/bin/sh

if [ $1 = "build" ]; then
    docker build -t tiny-rsvp --build-arg PORT=8080 --build-arg UID=$(id -u) --build-arg GID=$(id -g) .
elif [ $1 = "up" ]; then
    docker run -d -v /$PWD/databases/:/databases -p 8080:8080 --name tiny-rsvp tiny-rsvp
fi