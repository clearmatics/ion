const IonLink = artifacts.require("IonLink");
const IonLock = artifacts.require("IonLock");
const Token = artifacts.require("Token");

module.exports = async (deployer) => {
  try {
    deployer.deploy(IonLink, 0)
      .then(() => IonLink.deployed)
      .then(() => deployer.deploy(Token))
      .then(() => Token.deployed)
      .then(() => deployer.deploy(IonLock, Token.address, IonLink.address))

  } catch(err) {
    console.log('ERROR on deploy:',err);
  }

};
