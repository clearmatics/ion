const Sodium = artifacts.require("Sodium");
const IonLink = artifacts.require("IonLink");
const IonLock = artifacts.require("IonLock");
const Token = artifacts.require("Token");
const HTLC = artifacts.require("HTLC");

module.exports = async (deployer) => {
    deployer.deploy(HTLC);

    deployer.deploy(Sodium);

    await deployer.deploy(Token);

    await deployer.deploy(IonLink, 0);

    let ionlink_deployed = await IonLink.deployed();
    console.log("IonLink address", ionlink_deployed.address);

    await deployer.deploy(IonLock, Token.address, ionlink_deployed.address);
    let ionlock_deployed = await IonLock.deployed();
    console.log("IonLock address", ionlock_deployed.address);
};
