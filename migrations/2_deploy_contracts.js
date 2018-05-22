const Hydrogen = artifacts.require("Hydrogen");
const Token = artifacts.require("Token");

module.exports = async (deployer) => {
  try {
    deployer.deploy(Hydrogen)
    deployer.deploy(Token)

  } catch(err) {
    console.log('ERROR on deploy:',err);
  }

};
