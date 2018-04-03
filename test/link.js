// TODO: merkle tree library for javascript
//       must be compatible with the Python implementation

const IonLink = artifacts.require("IonLink");
const utils = require('./helpers/utils.js')

contract('IonLink', (accounts) => {
  	let link;

    beforeEach(async function() {
		link = await IonLink.new(0);
    });

    const sender = accounts[0]

    it('works', async () => {
        console.log("Obj address", link.address);
    });

		it("Update(): determine whether a new block can be added", async function()
		{

      const test = 1234
			const txReceipt = await link.Update(
				test,
        {
          from: sender
        }
			)


      const receipt = utils.txLoggedArgs(txReceipt)
      console.log(txReceipt.logs[0].args)
      const out = txReceipt.logs[1].args
      console.log(out)

      // const root = await link.GetRoot(
      //   1,
      //   {
      //     from: sender
      //   }
      // )
      // console.log(root)


      // latestBlock = await link.LatestBlock.call()

		});


    // TODO: verify Update() works with multiple items in-sequence

    // TODO: verify that same roots applied twice result in different hashes
});
