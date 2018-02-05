
const IonLock = artifacts.require("IonLock");

contract('IonLock', (accounts) => {
    it('works', async () => {
        await IonLock.deployed();
    });
});
