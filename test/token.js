const Token = artifacts.require("./Token.sol");

contract('Token', (accounts) => {
    it('works', async () => {
        console.log(await Token.deployed());

        /*
        let lock = await IonLock.deployed();
        console.log(lock);
        */
    });
});
