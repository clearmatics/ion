const Sodium = artifacts.require("./Sodium.sol");
const IonLink = artifacts.require("./IonLink.sol");
const IonLock = artifacts.require("./IonLock.sol");
const Token = artifacts.require("./Token.sol");

module.exports = async (deployer) => {
    await deployer.deploy(Sodium);

    await deployer.deploy(IonLink, 0);
    var link_deployed = await IonLink.deployed();

    await deployer.deploy(Token);
    var token_deployed = await Token.deployed();

    var lock_deployed = await deployer.deploy(IonLock, token_deployed.address, link_deployed.address);
};
