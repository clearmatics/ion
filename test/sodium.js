
const Sodium = artifacts.require("Sodium");

contract('Sodium', (accounts) => {
    it('works', async () => {
        await Sodium.deployed();
    });
});
