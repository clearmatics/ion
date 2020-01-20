const Ion = artifacts.require("Ion");
const Clique = artifacts.require("Clique");
const EthereumStore = artifacts.require("EthereumStore");
const EventFunction = artifacts.require("Function");
const EventVerifier = artifacts.require("TriggerEventVerifier");

const saveGasInfo = require("../test/helpers/utils").saveGas
const config = require("../test/helpers/config.json")

module.exports = async (deployer, network) => {
  try {
      deployer.deploy(Ion, "0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177")
      .then((res) => writeGasToFile(res.transactionHash, "Ion"))
      .then(() => deployer.deploy(EthereumStore, Ion.address))
      .then((res) => writeGasToFile(res.transactionHash, "Ethereum Store"))
      .then(() => deployer.deploy(Clique, Ion.address))
      .then((res) => writeGasToFile(res.transactionHash, "Clique Validation"))
      .then(() => deployer.deploy(EventVerifier))
      .then((res) => writeGasToFile(res.transactionHash, "Event verifier"))
      .then(() => deployer.deploy(EventFunction, Ion.address, EventVerifier.address))
      .then((res) => writeGasToFile(res.transactionHash, "Event Function"))
  } catch(err) {
    console.log('ERROR on deploy:',err);
  }

};

writeGasToFile = async (txHash, contractName) => {
  receipt = await web3.eth.getTransactionReceipt(txHash)
  saveGasInfo(config.BENCHMARK_DEPLOYMENT_FILEPATH, txHash, contractName, receipt.cumulativeGasUsed)
}