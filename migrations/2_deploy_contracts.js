const Sodium = artifacts.require("Sodium");
const IonLink = artifacts.require("IonLink");
const IonLock = artifacts.require("IonLock");
const Token = artifacts.require("Token");

module.exports = async (deployer) => {
    deployer.deploy(Sodium);

    deployer.deploy(Token).then( function () {
        deployer.deploy(IonLink, 0).then( async function () {
            await deployer.deploy(IonLock, Token.address, IonLink.address);
        } );
    } );
};
