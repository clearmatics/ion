const Ion = artifacts.require("Ion");
const Clique = artifacts.require("Clique");
const EthereumStore = artifacts.require("EthereumStore");
const EventFunction = artifacts.require("Function");
const EventVerifier = artifacts.require("TriggerEventVerifier");

const benchmarkHelpers = require("../benchmark/helpers")
const config = require("../benchmark/config.json")

module.exports = async (deployer, network) => {
  
  try {
      deployer.deploy(Ion, "0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177")
      .then((res) => writeGasToFile(res.transactionHash, "Deploy Ion"))
      .then(() => deployer.deploy(EthereumStore, Ion.address))
      .then((res) => writeGasToFile(res.transactionHash, "Deploy Ethereum Store"))
      .then(() => deployer.deploy(Clique, Ion.address))
      .then((res) => writeGasToFile(res.transactionHash, "Deploy Clique Validation"))
      .then(() => deployer.deploy(EventVerifier))
      .then((res) => writeGasToFile(res.transactionHash, "Deploy Event verifier"))
      .then(() => deployer.deploy(EventFunction, Ion.address, EventVerifier.address))
      .then((res) => writeGasToFile(res.transactionHash, "Deploy Event Function"))
  } catch(err) {
    console.log('ERROR on deploy:',err);
  }

};

writeGasToFile = async (txHash, contractName) => {
  let duration = "Not estimated"
  receipt = await web3.eth.getTransactionReceipt(txHash)
  benchmarkHelpers.saveStatsToFile(config.BENCHMARK_DEPLOYMENT_FILEPATH, txHash, contractName, receipt.cumulativeGasUsed, duration)
}