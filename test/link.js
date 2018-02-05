
const IonLink = artifacts.require("IonLink");

contract('IonLink', (accounts) => {
	let obj;

    beforeEach(async function() {
		obj = await IonLink.new(0);
    });

    it('works', async () => {
        console.log("Obj address", obj.address);
    });
});
