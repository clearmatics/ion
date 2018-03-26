const merkle = require('./merkle')
const Sodium = artifacts.require("Sodium");

contract.only('Sodium', (accounts) => {

  it.only('test JS Merkle', async () => {
    const merkleRoot = merkle.createMerkle(["1","2","3","4","5","6","7"])
    console.log(JSON.stringify(merkleRoot,2,2))
  })

  it('works', async () => {
    const obj = await Sodium.deployed();

    const root = "0x1a792cf089bfa56eae57ffe87e9b22f9c9bfe52c1ac300ea1f43f4ab53b4b794"
    const leafHash = "0x2584db4a68aa8b172f70bc04e2e74541617c003374de6eb4b295e823e5beab01"
    const path = [
      "0x1ab0c6948a275349ae45a06aad66a8bd65ac18074615d53676c09b67809099e0"
      ,"0x093fd25755220b8f497d65d2538c01ed279c131f63e42b2942867f2bd6622486"
      ,"0xb1d101d9a9d27c3a8ed9d1b6548626eacf3d19546306117eb8af547d1e97189e"
      ,"0xcb431dd627bc8dcfd858eae9304dc71a8d3f34a8de783c093188bb598eeafd04"
    ]
    const nextBlock = await obj.NextBlock()
    console.log(nextBlock.toString())

    const receiptUpdate = await obj.Update(nextBlock.toString(),[root])
    console.log(receiptUpdate)

    const valid = await obj.Verify(nextBlock.toString(),leafHash,path)
    console.log(valid)
  });
});
