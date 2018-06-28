module.exports = {
  networks: {
    development: {
      host: "localhost",
      port: 8545,
      network_id: "*",
      gas: 4600000
    },
    clique: {
      host: "localhost",
      port: 8501,
      network_id: "*",
      gas: 4600000
    }
  },
  mocha: {
    useColors: true,
    enableTimeouts: false
  },
  solc: {
    optimizer: {
      enabled: true,
        runs: 200
    }
  }
};
