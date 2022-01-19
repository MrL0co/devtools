# syntax=docker/dockerfile:1

FROM golang:1.17 AS build

WORKDIR /app

COPY go.mod ./
#COPY go.sum ./

RUN go mod download

COPY . ./

RUN make build


##
## Deploy
##
FROM ubuntu:latest

ARG USER_NAME="developer"
ARG USER_PASSWORD="test123"

ENV USER_NAME $USER_NAME
ENV USER_PASSWORD $USER_PASSWORD

WORKDIR /

RUN apt update && apt install -y ca-certificates wget sudo && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/devtools /usr/local/bin/devtools


#RUN devtools complete
#
## REMOVE DOCKER APT CLEANER TO EMULATE REAL PC (with apt cache files in order to check last apt update command)
#RUN rm /etc/apt/apt.conf.d/docker-clean
#ENV TERM xterm
#
### AUTOMATE THIgS IN CODE
#RUN apt update && apt install -y zsh git && `rm -rf /var/lib/apt/lists/*`
#
### SKIP
RUN adduser --quiet --disabled-password --shell /bin/bash --home /home/$USER_NAME --gecos "User" $USER_NAME \
    && echo "${USER_NAME}:${USER_PASSWORD}" | chpasswd && usermod -aG sudo $USER_NAME && adduser $USER_NAME sudo
#
USER ${USER_NAME}
### end SKIP
#
## terminal colors with xterm
ENV TERM xterm
## set the zsh theme
#ENV ZSH_THEME robbyrussell
#
#RUN echo $USER_PASSWORD | chsh -s $(which zsh) developer
#
#RUN wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh

EXPOSE 8080

WORKDIR "/home/$USER_NAME"
#CMD ["zsh"]
CMD ["bash"]
#ENTRYPOINT [ "/usr/local/bin/devtools" ]
