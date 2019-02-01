FROM golang:1.8

RUN apt-get update && apt-get install -y \
        vim \
        curl \
        sudo

# Install a recent version of nodejs
RUN curl -sL https://deb.nodesource.com/setup_10.x | sudo bash - && sudo apt-get install -y nodejs
RUN npm install -g truffle ganache-cli

COPY . /go/src/github.com/Clearmatics/ion

WORKDIR /go/src/github.com/Clearmatics/ion

CMD ["/bin/bash"]
