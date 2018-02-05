const Token = artifacts.require("Token");
const IonLink = artifacts.require("IonLink");
const IonLock = artifacts.require("IonLock");

contract('IonLock', (accounts) => {
	let link_obj;
	let token_obj;
	let lock_obj;

    beforeEach(async function() {
    	token_obj = await Token.new();
    	link_obj = await IonLink.new(0);
		lock_obj = await IonLock.new(token_obj.address, link_obj.address);
    });

    it('works', async () => {
        console.log("Obj address", lock_obj.address);
    });
});
