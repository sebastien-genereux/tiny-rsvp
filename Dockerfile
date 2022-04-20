## 
# ARGS
#

#used to ensure container non-root account has write access to mounted host folder
ARG UID
ARG GID

ARG PORT

##
# BUILD
##

FROM golang:1.18-stretch AS build 

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./

RUN go build -ldflags "-s -w" -o /tiny-rsvp

##
# DEPLOY
##

FROM gcr.io/distroless/base-debian11

WORKDIR /

# copy over the app and required files
COPY --from=build /tiny-rsvp /tiny-rsvp
COPY configs/ ./configs
COPY --chown=1000:1000 /databases /databases
COPY web/ ./web

EXPOSE $PORT

USER $UID:$GID

ENTRYPOINT ["/tiny-rsvp"]