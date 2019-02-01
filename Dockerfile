FROM golang:1.8

RUN apt-get update && apt-get install -y \
        vim \
        curl \
        sudo \
        wget


# Install a recent version of nodejs
RUN curl -sL https://deb.nodesource.com/setup_10.x | sudo bash - && sudo apt-get install -y nodejs=10.15.1-1nodesource1

COPY . /go/src/github.com/clearmatics/ion

# Install the current compatible solc version
RUN wget https://github.com/ethereum/solidity/releases/download/v0.4.25/solc-static-linux -O solc
RUN chmod +x ./solc
RUN cp ./solc /go/src/github.com/clearmatics/ion
ENV PATH $PATH:/go/src/github.com/clearmatics/ion

WORKDIR /go/src/github.com/clearmatics/ion

CMD ["/bin/bash"]
