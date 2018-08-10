FROM ubuntu:16.04

LABEL version="1.0"
LABEL maintainer="mgb@clearmatics.com"

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install --yes software-properties-common
RUN add-apt-repository ppa:ethereum/ethereum
RUN apt-get update && apt-get install --yes geth

RUN adduser --disabled-login --gecos "" eth_user

COPY docker_build /home/eth_user/docker_build
RUN chown -R eth_user:eth_user /home/eth_user/docker_build

USER eth_user

WORKDIR /home/eth_user

RUN geth --datadir docker_build/account/ init docker_build/clique.json

EXPOSE 8545

ENTRYPOINT bash

