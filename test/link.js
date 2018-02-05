// TODO: merkle tree library for javascript
//       must be compatible with the Python implementation

const IonLink = artifacts.require("IonLink");

contract('IonLink', (accounts) => {
	let obj;

    beforeEach(async function() {
		obj = await IonLink.new(0);
    });

    it('works', async () => {
        console.log("Obj address", obj.address);
    });

    // TODO: verify GetRoot() works

    // TODO: verify Update() works with multiple items in-sequence

    // TODO: verify that same roots applied twice result in different hashes
});
