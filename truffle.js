module.exports = {
  networks: {
    development: {
      host: "localhost",
      port: 8545,
      gas: 0xFFFFFFFFFFF,
      network_id: "*"
    },
    clique: {
      host: "localhost",
      port: 8501,
      network_id: "*"
    },
    coverage: {
      host: "localhost",
      port: 8555,
      network_id: "*", // Match any network id
      gas: 0xFFFFFFF,
      gasprice: 0x1
    },
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
