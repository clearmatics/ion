const Ion = artifacts.require("Ion");
const Clique = artifacts.require("Clique");
const EthereumStore = artifacts.require("EthereumStore");
const EventFunction = artifacts.require("Function");
const EventVerifier = artifacts.require("TriggerEventVerifier");

module.exports = async (deployer) => {
  try {
      deployer.deploy(Ion, "0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177")
      .then(() => Ion.deployed)
      .then(() => deployer.deploy(EthereumStore, Ion.address))
      .then(() => EthereumStore.deployed)
      .then(() => deployer.deploy(Clique, Ion.address))
      .then(() => Clique.deployed)
      .then(() => deployer.deploy(EventVerifier))
      .then(() => EventVerifier.deployed)
      .then(() => deployer.deploy(EventFunction, Ion.address, EventVerifier.address))
      .then(() => EventFunction.deployed)
  } catch(err) {
    console.log('ERROR on deploy:',err);
  }

};