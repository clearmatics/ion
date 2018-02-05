/*

*/
const IonLink = artifacts.require("./IonLink.sol");
const IonLock = artifacts.require("./IonLock.sol");

contract('IonLink', (accounts) => {
    it('works', async () => {
        console.log(await IonLink.deployed());
    });
});
