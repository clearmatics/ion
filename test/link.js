const crypto = require('crypto');

const merkle = require('./merkle')

const IonLink = artifacts.require("IonLink");
const utils = require('./helpers/utils.js')

const randomHex = () => crypto.randomBytes(32).toString('hex');
const randomArr = () => {
  const result = []
  const size =(Math.floor(Math.random() * 10) + 1);
  for(let i = size; 0 < i; i-- )
    result.push(randomHex())
  return result
}

contract('IonLink', (accounts) => {
  it('GetRoot', async () => {
    //const ionLink = await IonLink.new(10);
    const ionLink = await IonLink.deployed();

    const testData1 = randomArr()
    const tree1 = merkle.createMerkle(testData1)
    const testData2 = randomArr()
    const tree2 = merkle.createMerkle(testData2)
    const testData3 = randomArr()
    const tree3 = merkle.createMerkle(testData3)

    const rootArr1 = [tree1[1],tree2[1],tree3[1]]
    const receipt1 = await ionLink.Update(rootArr1)

    const testData4 = randomArr()
    const tree4 = merkle.createMerkle(testData4)
    const testData5 = randomArr()
    const tree5 = merkle.createMerkle(testData5)
    const testData6 = randomArr()
    const tree6 = merkle.createMerkle(testData6)

    const rootArr2 = [tree4[1],tree5[1],tree6[1]]
    const receipt2 = await ionLink.Update(rootArr2)

    const latestBlock = await ionLink.GetLatestBlock()
    const previousBlock = await ionLink.GetPrevious(latestBlock)
    const latestRoot = await ionLink.GetRoot(latestBlock)
    const previousRoot = await ionLink.GetRoot(previousBlock)
    assert.equal(tree6[1].toString(16),latestRoot.toString(16),'latest root is wrong')
    assert.equal(tree5[1].toString(16),previousRoot.toString(16),'previous root is wrong')
  })

  it('Update', async () => {
    //const ionLink = await IonLink.new(10);
    const ionLink = await IonLink.deployed();

    const testData1 = randomArr()
    const tree1 = merkle.createMerkle(testData1)
    const testData2 = randomArr()
    const tree2 = merkle.createMerkle(testData2)

    const leaf = testData2[0]
    const leafHash = merkle.merkleHash(leaf)
    const path = merkle.pathMerkle(leaf,tree2[0])
    const rootArg = [tree1[1],tree2[1]]

    const receiptUpdate = await ionLink.Update(rootArg)
    const latestBlock = await ionLink.GetLatestBlock()
    const valid = await ionLink.Verify(latestBlock,leafHash,path)
    assert(valid,'IonLink.verify() failed!')

    assert( receiptUpdate.logs.length > 0)
    assert.equal( receiptUpdate.logs[0].event, 'IonLinkUpdated', 'IonLinkUpdated event not found in logs' )
  })

  it('duplicate root', async () => {
    //const ionLink = await IonLink.new(10);
    const ionLink = await IonLink.deployed();

    const testData = randomArr()
    const tree = merkle.createMerkle(testData)

    const receiptUpdate = await ionLink.Update([tree[1],tree[1]])
    const latestBlock = await ionLink.GetLatestBlock()
    const previousBlock = await ionLink.GetPrevious(latestBlock)
    assert.notEqual(latestBlock.toString(16),previousBlock.toString(16),'submitted smae root 2x should have different hashes!')
  })
});
