const Web3Utils = require('web3-utils');
const crypto = require('crypto');

const merkle = require('./merkle')
const Token = artifacts.require("Token");
const IonLink = artifacts.require("IonLink");
const IonLock = artifacts.require("IonLock");

const randomHex = () => crypto.randomBytes(32).toString('hex');
const randomArr = () => {
  const result = []
  const size =(Math.floor(Math.random() * 10) + 1);
  for(let i = size; 0 < i; i-- )
    result.push(randomHex())
  return result
}

contract.only('IonLock', (accounts) => {
  const watchEvent = (eventObj) => new Promise((resolve,reject) => eventObj.watch((error,event) => error ? reject(error) : resolve(event)))

  it('tokenFallback is called by Token.transfer', async () => {
    const token = await Token.new();
    const ionLink = await IonLink.new(0);
		const ionLock = await IonLock.new(token.address, ionLink.address);
    /*
    const token = await Token.deployed();
    const ionLink = await IonLink.deployed();
    const ionLock = await IonLock.deployed();
    */
    /*
    const lockEvents = ionLock.allEvents()
    const linkEvents = ionLink.allEvents()
    const tokenEvents = token.allEvents()
    const printEvent = (error,event) => {
      if(error) {
        console.log('Event ERROR:',error)
        return
      }
      const address = '0x' + event.address.toString(16)
      const args = event.args
      const blockHash = event.blockHash
      const blockNumber = event.blockNumber
      const logIndex = event.logIndex
      const name = event.event
      const removed = event.removed /// success?
      const txIndex = event.transactionIndex
      const txHash = event.transactionHash
      const argsName = Object.keys(args)
      console.log(name,address)//,args);
      console.log(argsName.map(key => Web3Utils.isBN(args[key]) ? `${key}: ${args[key].toString(16)}` : `${key}: ${args[key]}`))
    }
    lockEvents.watch(printEvent)
    linkEvents.watch(printEvent)
    tokenEvents.watch(printEvent)
    */

    const owner = accounts[0]
    const totalSupply = 1000
    const value = 10 // value transferred

    const receiptMint = await token.mint(totalSupply)

    //const balance = await token.balanceOf(owner)

    //console.log('accounts',accounts)
    //console.log('token.address',token.address)
    //console.log('ionLink.address',ionLink.address)
    //console.log('ionLock.address',ionLock.address)
    //const receiptTansfer = await token.transfer(accounts[1],1)
    const receiptTansfer1 = await token.transfer(ionLock.address,value)

    const ionMintEventObj = ionLock.IonMint()
    const ionTransferEventObj = ionLock.IonTransfer()
    try {
      const ionMintEvent = await watchEvent(ionMintEventObj)
      const ionTransferEvent = await watchEvent(ionTransferEventObj)

      assert.equal(ionMintEvent.args.value.toString(10),''+value,'IonMint event unexpected value!')
      assert.equal(ionMintEvent.args.ref,ionTransferEvent.args.ref,'IonMint event ref not equal to IonTransfer event ref!')

    } catch (err) {
      console.log('event error:', err)
      assert.fail('event error:' + err)
    }
    ionMintEventObj.stopWatching()
    ionTransferEventObj.stopWatching()

    /*
    lockEvents.stopWatching()
    linkEvents.stopWatching()
    tokenEvents.stopWatching()
    */
  })

  it('withdraw', async () => {
    const token = await Token.new();
    const ionLink = await IonLink.new(0);
		const ionLock = await IonLock.new(token.address, ionLink.address);

    const owner = accounts[0]
    const totalSupply = 1000
    const value = 10 // value transferred

    const receiptMint = await token.mint(totalSupply)

    const receiptTansfer1 = await token.transfer(ionLock.address,value)

    // get reference from events
    const ionMintEventObj = ionLock.IonMint()
    const ionTransferEventObj = ionLock.IonTransfer()
    let ref
    try {
      const ionMintEvent = await watchEvent(ionMintEventObj)
      const ionTransferEvent = await watchEvent(ionTransferEventObj)
      ref = ionTransferEvent.args.ref

    } catch (err) {
      console.log('event error:', err)
      assert.fail('event error:' + err)
    }
    ionMintEventObj.stopWatching()
    ionTransferEventObj.stopWatching()

    assert(ref) //check that ref has something

    // hash details to be added to IonLink
    const reference = ref
    const valueHex = '0x'+Web3Utils.toBN(value).toString(16).padStart(64,'0') // make an hex that is good to sha3 in solidity uint256 -> 64 bytes
    const lockAddr = ionLock.address
    const tokenAddr = token.address
    const withdrawReceiver = accounts[1]

    //concat the arguments of sh3 in solidity (in the same way solidity does)
    const joinedArgs = '0x' + [withdrawReceiver,lockAddr,tokenAddr,valueHex,reference].map(el=>el.slice(2)).join('')

    //const hashData = Web3Utils.soliditySha3(withdrawReceiver,lockAddr,tokenAddr,value,reference) //it is the same
    const hashData = Web3Utils.sha3(joinedArgs)
    //console.log([withdrawReceiver,lockAddr,tokenAddr,value,reference],hashData)
    //console.log( joinedArgs,hashData2)

    // submit hashdata to IonLink
    const leaf = joinedArgs // joined args need to be added to the random leafs of the tree
    const testData = randomArr()
    testData[0] = leaf
    const tree = merkle.createMerkle(testData)
    const treeExtra = merkle.createMerkle(randomArr()) // IonLink needs 2 roots min to update

    const leafHash = merkle.merkleHash(leaf)
    const path = merkle.pathMerkle(leaf,tree[0])
    const rootArg = [treeExtra[1],tree[1]]

    const receiptUpdate = await ionLink.Update(rootArg)
    const latestBlock = await ionLink.GetLatestBlock()
    const valid = await ionLink.Verify(latestBlock,leafHash,path)
    assert(valid,'leaf not found in tree')

    // withdraw from ionlock
    const receiptWithdraw = await ionLock.Withdraw(value,reference,latestBlock,path,{ from: withdrawReceiver })

    const balanceOwner = await token.balanceOf(owner)
    const balanceReceiver = await token.balanceOf(withdrawReceiver)

    assert.equal(balanceOwner,totalSupply - value, 'sender balance wrong!')
    assert.equal(balanceReceiver,value, 'receiver balance wrong!')
  })
});
