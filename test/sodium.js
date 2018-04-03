const merkle = require('./merkle')
const Sodium = artifacts.require("Sodium");
const crypto = require('crypto');

const randomHex = () => crypto.randomBytes(32).toString('hex');
const randomArr = () => {
  const result = []
  const size =5// (Math.floor(Math.random() * 10) + 1);
  for(let i = size; 0 < i; i-- )
    result.push(randomHex())
  return result
}


contract.only('Sodium', (accounts) => {

  it('data test', async () => {
    const sodium = await Sodium.deployed();

    const testData = ["1","2","3","4","5","6","7"]
    const tree = merkle.createMerkle(testData)
    const path = testData.map(value => merkle.pathMerkle(value,tree[0]))

    testData.forEach((leaf,idx) => assert(merkle.proofMerkle(leaf,path[idx],tree[1])))

    const leafHash = merkle.merkleHash(testData[0])
    const rootArg = tree[1]

    const nextBlock = await sodium.NextBlock()
    const receiptUpdate = await sodium.Update(nextBlock,[rootArg])
    const valid = await sodium.Verify(nextBlock,leafHash,path[0])
    assert(valid,'Sodium.verify() failed!')
  })

  it('multiple root update', async () => {
    const sodium = await Sodium.deployed();

    const groupSize = await sodium.GroupSize()

    const testData1 = randomArr()
    const tree1 = merkle.createMerkle(testData1)
    const testData2 = randomArr()
    const tree2 = merkle.createMerkle(testData2)
    const testData3 = randomArr()
    const tree3 = merkle.createMerkle(testData3)

    const rootArr1 = [tree1[1],tree2[1],tree3[1]]

    const nextBlock1 = await sodium.NextBlock()
    const receiptUpdate1 = await sodium.Update(nextBlock1,rootArr1)
    const nextBlock2 = await sodium.NextBlock()
    const blocksSubmited = (nextBlock2.toString(10) - nextBlock1.toString(10))/groupSize
    assert.equal(blocksSubmited,rootArr1.length,'blocks submitted number wrong')


    const blockNumber = nextBlock2 - (2 * groupSize)
    const leafHash = merkle.merkleHash(testData2[0])
    const path = merkle.pathMerkle(testData2[0],tree2[0])
    const valid = await sodium.Verify(blockNumber,leafHash,path)
    assert(valid,'Sodium.verify() failed!')

    // separated multiple updates
    const testData4 = randomArr()
    const tree4 = merkle.createMerkle(testData4)
    const testData5 = randomArr()
    const tree5 = merkle.createMerkle(testData5)

    const rootArr2 = [tree4[1],tree5[1]]

    const nextBlock3 = await sodium.NextBlock()
    const receiptUpdate2 = await sodium.Update(nextBlock3,rootArr2)
    const nextBlock4 = await sodium.NextBlock()
    const blocksSubmited2 = (nextBlock4.toString(10) - nextBlock3.toString(10))/groupSize
    assert.equal(blocksSubmited2,rootArr2.length,'blocks submitted number wrong')
  })

  /* TEST WITHOUT CHANGES TO THE Merkle.sol
  it('harry data test', async () => {
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
    //console.log('0x'+nextBlock.toString(16))

    const receiptUpdate = await obj.Update(nextBlock.toString(),[root])
    //console.log(receiptUpdate)

    const valid = await obj.Verify(nextBlock.toString(),leafHash,path)
    //console.log(valid)
    assert(valid)
  });
  */
});
