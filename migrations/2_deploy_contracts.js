const Sodium = artifacts.require("Sodium");
const IonLink = artifacts.require("IonLink");
const IonLock = artifacts.require("IonLock");
const Token = artifacts.require("Token");
const HTLC = artifacts.require("HTLC");

module.exports = async (deployer) => {
  try {
    deployer.deploy(HTLC);

    deployer.deploy(Sodium);

    deployer.deploy(IonLink, 10);

    deployer.deploy(Token)

    const ionlink_deployed = await IonLink.deployed();
    const token_deployed = await Token.deployed();
    await deployer.deploy(IonLock, token_deployed.address, ionlink_deployed.address);

  } catch(err) {
    console.log('ERROR on deploy:',err);
  }

};
