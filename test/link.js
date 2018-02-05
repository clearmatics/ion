
const IonLink = artifacts.require("IonLink");

contract('IonLink', (accounts) => {
    it('works', async () => {
        await IonLink.deployed();
    });
});
