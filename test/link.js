// TODO: merkle tree library for javascript
//       must be compatible with the Python implementation

const IonLink = artifacts.require("IonLink");

const crypto = require('crypto')

const gasPrice = 100000000000 // truffle fixed gas price
const txGas = txReceipt => txReceipt.receipt.gasUsed * gasPrice
const txLoggedArgs = txReceipt => txReceipt.logs[0].args
const txContractId = txReceipt => txLoggedArgs(txReceipt).contractId
const oneFinney = web3.toWei(1, 'finney')

// Format required for sending bytes through eth client:
//  - hex string representation
//  - prefixed with 0x
const bufToStr = b => '0x' + b.toString('hex')

const sha256 = x =>
  crypto
    .createHash('sha256')
    .update(x)
    .digest()

const random32 = () => crypto.randomBytes(32)

const isSha256Hash = hashStr => /^0x[0-9a-f]{64}$/i.test(hashStr)

const newSecretHashPair = () => {
  const secret = random32()
  const hash = sha256(secret)
  return {
    secret: bufToStr(secret),
    hash: bufToStr(hash),
  }
}

contract.only('IonLink', (accounts) => {
  	let link;

    beforeEach(async function() {
		link = await IonLink.new(0);
    });

    const sender = accounts[0]

    it('works', async () => {
        console.log("Obj address", link.address);
    });

		it.only("Update(): determine whether a new block can be added", async function()
		{

      let latestBlock = await link.LatestBlock.call()
      console.log(latestBlock.root)

			const test = 0
			const txReceipt = await link.Update(
				test,
        {
          from: sender
        }
			)

      latestBlock = await link.LatestBlock.call()
      console.log(latestBlock)

		});


    // TODO: verify Update() works with multiple items in-sequence

    // TODO: verify that same roots applied twice result in different hashes
});
