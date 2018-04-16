module.exports = {
  networks: {
    development: {
      host: "localhost",
      port: 8545,
      network_id: "*" // Match any network id
    },
    ci: {
      host: "localhost",
      port: 8545,
      network_id: "*" // Match any network id
    },
    coverage: {
      host: "localhost",
      port: 8555,
      network_id: "*", // Match any network id
      gas: 0xFFFFFFF,
      gasprice: 0x1
    },
    testrpca: {
      host: "localhost",
      port: 8545,
      network_id: "*" // Match any network id
    },
    testrpcb: {
      host: "localhost",
      port: 8546,
      network_id: "*" // Match any network id
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
