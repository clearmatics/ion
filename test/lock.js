const Web3Utils = require('web3-utils');
const web3Abi = require('web3-eth-abi');
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
const send2Lock = async (from,tokenAddr,ionLockAddr,value,rawRef) => {
  const overloadedTransferAbi = {
    "constant": false,
    "inputs": [
      { "name": "_to", "type": "address" },
      { "name": "_value", "type": "uint256" },
      { "name": "_data", "type": "bytes" }
    ],
    "name": "transfer",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "function"
  }
  const transferMethodTransactionData = web3Abi.encodeFunctionCall(
    overloadedTransferAbi,
    [ ionLockAddr, value, Web3Utils.toHex(rawRef) ]
  );
  const receiptTransfer1 = await web3.eth.sendTransaction({ from: from, to: tokenAddr, data: transferMethodTransactionData, value: 0 });
  return receiptTransfer1;
}
const waitLockEvent = async (lockContract,rawRef) => {
  const ionMintEventObj = lockContract.IonMint()
  const ionTransferEventObj = lockContract.IonTransfer()
  let ref
  try {
    const ionMintEvent = await watchEvent(ionMintEventObj)
    const ionTransferEvent = await watchEvent(ionTransferEventObj)
    assert.equal(ionTransferEvent.args.ref,Web3Utils.sha3(rawRef),'ref different than expected after transfer!')
    ref = ionTransferEvent.args.ref
  } catch (err) {
    console.log('event error:', err)
    assert.fail('event error:' + err)
  }
  ionMintEventObj.stopWatching()
  ionTransferEventObj.stopWatching()
  return ref
}
const joinIonLinkData = (receiverAddr,tokenAddr,ionLockAddr,value,reference) => {
  const valueHex = '0x'+Web3Utils.toBN(value).toString(16).padStart(64,'0') // make an hex that is good to sha3 in solidity uint256 -> 64 bytes
  const leaf = '0x' + [receiverAddr,ionLockAddr,tokenAddr,valueHex,reference].map(el=>el.slice(2)).join('') // joined args need to be added to the random leafs of the tree
  return leaf
}

const watchEvent = (eventObj) => new Promise((resolve,reject) => eventObj.watch((error,event) => error ? reject(error) : resolve(event)))

contract('IonLock', (accounts) => {

  it('tokenFallback is called by Token.transfer', async () => {
    const token = await Token.new();
    const ionLink = await IonLink.new(0);
		const ionLock = await IonLock.new(token.address, ionLink.address);

    const owner = accounts[0]
    const totalSupply = 1000
    const value = 10 // value transferred

    const receiptMint = await token.mint(totalSupply)

    const receiptTransfer1 = await token.rawTransfer(ionLock.address,value)

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

  })

  it('withdraw', async () => {
    const token = await Token.new();
    const ionLink = await IonLink.new(0);
		const ionLock = await IonLock.new(token.address, ionLink.address);

    const owner = accounts[0]
    const totalSupply = 1000
    const value = 10 // value transferred

    const receiptMint = await token.mint(totalSupply)

    const receiptTransfer1 = await token.transfer(ionLock.address,value)

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

    assert(ref,'reference is empty!') //check that ref has something

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


  // using the overloaded function is a problem for truffle so it is better to have a different test for that
  it('withdraw with reference', async () => {
    const token = await Token.new();
    const ionLink = await IonLink.new(0);
		const ionLock = await IonLock.new(token.address, ionLink.address);

    const owner = accounts[0]
    const totalSupply = 1000
    const value = 10 // value transferred
    const rawRef = 'Hello world!'

    const receiptMint = await token.mint(totalSupply)

    //const receiptTransfer1 = await token.transfer(ionLock.address,value)
    const overloadedTransferAbi = {
      "constant": false,
      "inputs": [
        { "name": "_to", "type": "address" },
        { "name": "_value", "type": "uint256" },
        { "name": "_data", "type": "bytes" }
      ],
      "name": "transfer",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    }
    const transferMethodTransactionData = web3Abi.encodeFunctionCall(
      overloadedTransferAbi,
      [
        ionLock.address,
        value,
        Web3Utils.toHex(rawRef)
      ]
    );
    const receiptTransfer1 = await web3.eth.sendTransaction(
      {
        from: owner,
        to: token.address,
        data: transferMethodTransactionData,
        value: 0
      }
    );

    // get reference from events
    const ionMintEventObj = ionLock.IonMint()
    const ionTransferEventObj = ionLock.IonTransfer()
    let ref
    try {
      const ionMintEvent = await watchEvent(ionMintEventObj)
      const ionTransferEvent = await watchEvent(ionTransferEventObj)
      assert.equal(ionTransferEvent.args.ref,Web3Utils.sha3(rawRef),'ref different than expected after transfer!')

    } catch (err) {
      console.log('event error:', err)
      assert.fail('event error:' + err)
    }
    ionMintEventObj.stopWatching()
    ionTransferEventObj.stopWatching()

    // hash details to be added to IonLink
    const reference = Web3Utils.sha3(rawRef) //hash our reference
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

  it('withdraw different chains with reference', async () => {
    const token = await Token.new();
    const ionLink = await IonLink.new(0);
		const ionLock = await IonLock.new(token.address, ionLink.address);
    const tokenB = await Token.new();
    const ionLinkB = await IonLink.new(0);
		const ionLockB = await IonLock.new(tokenB.address, ionLinkB.address);

    const owner = accounts[0]
    const totalSupply = 1000
    const value = 10 // value transferred
    const rawRef = Web3Utils.sha3('Reference from deposit on chain A')
    const totalSupplyB = 1000
    const valueB = 10 // value transferred
    const rawRefB = Web3Utils.sha3('Reference from deposit on chain B')

    const sender = accounts[3]
    const withdrawReceiver = accounts[5]
    const senderB = accounts[4]
    const withdrawReceiverB = accounts[6]

    const receiptMint = await token.mint(totalSupply)
    const receiptTransfer = await token.transfer(sender,totalSupply)
    const receiptMintB = await tokenB.mint(totalSupplyB)
    const receiptTransferB = await tokenB.transfer(senderB,totalSupplyB)

    // wait lock events blocks the rest of the test from running if no event is triggered
    // A -> LOCK_A // get REFERENCE_A
    const receiptSend2Lock = await send2Lock(sender,token.address,ionLock.address,value,rawRef)
    const ref = await waitLockEvent(ionLock,rawRef)

    // B -> LOCK_B // get REFERENCE_B
    const receiptSend2LockB = await send2Lock(senderB,tokenB.address,ionLockB.address,value,rawRefB)
    const refB = await waitLockEvent(ionLockB,rawRefB)

    // hash details to be added to IonLink
    // MERKLE_ROOT(REFERENCE_B) -> LINK_A
    // this marks B as the recipient of the tokens
    const leaf = joinIonLinkData(withdrawReceiverB,token.address,ionLock.address,value,refB)

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

    // MERKLE_ROOT(REFERENCE_A) -> LINK_B
    // this marks A as the recipient of the tokens
    const leafB = joinIonLinkData(withdrawReceiver,tokenB.address,ionLockB.address,valueB,ref)

    const testDataB = randomArr()
    testDataB[0] = leafB
    const treeB = merkle.createMerkle(testDataB)
    const treeExtraB = merkle.createMerkle(randomArr()) // IonLink needs 2 roots min to update

    const leafHashB = merkle.merkleHash(leafB)
    const pathB = merkle.pathMerkle(leafB,treeB[0])
    const rootArgB = [treeExtraB[1],treeB[1]]

    const receiptUpdateB = await ionLinkB.Update(rootArgB)
    const latestBlockB = await ionLinkB.GetLatestBlock()
    const validB = await ionLinkB.Verify(latestBlockB,leafHashB,pathB)
    assert(validB,'leaf not found in tree')

    // withdraw from ionlock
    // LOCK_A -> B
    const receiptWithdraw = await ionLock.Withdraw(value,refB,latestBlock,path,{ from: withdrawReceiverB })

    const balanceSender = await token.balanceOf(sender)
    const balanceReceiver = await token.balanceOf(withdrawReceiverB)

    assert.equal(balanceSender,totalSupply - value, 'sender balance wrong!')
    assert.equal(balanceReceiver,value, 'receiver balance wrong!')

    // LOCK_B -> A
    const receiptWithdrawB = await ionLockB.Withdraw(valueB,ref,latestBlockB,pathB,{ from: withdrawReceiver })

    const balanceSenderB = await tokenB.balanceOf(senderB)
    const balanceReceiverB = await tokenB.balanceOf(withdrawReceiver)

    assert.equal(balanceSenderB,totalSupplyB - valueB, 'sender balance wrong!')
    assert.equal(balanceReceiverB,valueB, 'receiver balance wrong!')
  })
});
