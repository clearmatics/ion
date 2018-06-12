const IonLink = artifacts.require("IonLink");
const IonLock = artifacts.require("IonLock");
const Token = artifacts.require("Token");
const HTLC = artifacts.require("HTLC");

module.exports = async (deployer) => {
  try {
    deployer.deploy(IonLink, 0)
      .then(() => IonLink.deployed)
      .then(() => deployer.deploy(Token))
      .then(() => Token.deployed)
      .then(() => deployer.deploy(IonLock, Token.address, IonLink.address))
      .then(() => deployer.deploy(HTLC))
      .then(() => HTLC.deployed)

  } catch(err) {
    console.log('ERROR on deploy:',err);
  }

};
