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

contract.only('IonLock', (accounts) => {
  const watchEvent = (eventObj) => new Promise((resolve,reject) => eventObj.watch((error,event) => error ? reject(error) : resolve(event)))

  // using the overloaded function is a problem for truffle so it is better to have a different test for that
  it.only('Make a cross chain payment', async () => {
    const tokenA = await Token.new();
    const tokenB = await Token.new();
    const ionLinkA = await IonLink.new(0);
		const ionLockA = await IonLock.new(tokenA.address, ionLinkA.address);

    const sender = accounts[0]
    const receiver = accounts[1]
    const totalSupply = 1000
    const value = 10 // value transferred
    const rawRef = 'Hello world!'

    const receiptMint = await tokenA.mint(totalSupply)

    console.log("Perform minting and transferring of tokens on chain A:\n")
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
        ionLockA.address,
        value,
        Web3Utils.toHex(rawRef)
      ]
    );
    const receiptTransfer1 = await web3.eth.sendTransaction(
      {
        from: sender,
        to: tokenA.address,
        data: transferMethodTransactionData,
        value: 0
      }
    );

    // get reference from events
    console.log(await tokenA.balanceOf(sender))
    const ionMintEventObj = ionLockA.IonMint()
    const ionTransferEventObj = ionLockA.IonTransfer()
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

    console.log(await tokenA.balanceOf(sender))

    // // hash details to be added to IonLink
    // const reference = Web3Utils.sha3(rawRef) //hash our reference
    // const valueHex = '0x'+Web3Utils.toBN(value).toString(16).padStart(64,'0') // make an hex that is good to sha3 in solidity uint256 -> 64 bytes
    // const lockAddr = ionLockA.address
    // const tokenAAddr = tokenA.address
    // const withdrawReceiver = accounts[1]
    //
    // //concat the arguments of sh3 in solidity (in the same way solidity does)
    // const joinedArgs = '0x' + [withdrawReceiver,lockAddr,tokenAAddr,valueHex,reference].map(el=>el.slice(2)).join('')
    //
    // //const hashData = Web3Utils.soliditySha3(withdrawReceiver,lockAddr,tokenAAddr,value,reference) //it is the same
    // const hashData = Web3Utils.sha3(joinedArgs)
    // //console.log([withdrawReceiver,lockAddr,tokenAAddr,value,reference],hashData)
    // //console.log( joinedArgs,hashData2)
    //
    // // submit hashdata to IonLink
    // const leaf = joinedArgs // joined args need to be added to the random leafs of the tree
    // const testData = randomArr()
    // testData[0] = leaf
    // const tree = merkle.createMerkle(testData)
    // const treeExtra = merkle.createMerkle(randomArr()) // IonLink needs 2 roots min to update
    //
    // const leafHash = merkle.merkleHash(leaf)
    // const path = merkle.pathMerkle(leaf,tree[0])
    // const rootArg = [treeExtra[1],tree[1]]
    //
    // const receiptUpdate = await ionLinkA.Update(rootArg)
    // const latestBlock = await ionLinkA.GetLatestBlock()
    // const valid = await ionLinkA.Verify(latestBlock,leafHash,path)
    // assert(valid,'leaf not found in tree')
    //
    // // withdraw from ionlock
    // const receiptWithdraw = await ionLockA.Withdraw(value,reference,latestBlock,path,{ from: withdrawReceiver })
    //
    // const balanceOwner = await tokenA.balanceOf(sender)
    // const balanceReceiver = await tokenA.balanceOf(withdrawReceiver)
    //
    // assert.equal(balanceOwner,totalSupply - value, 'sender balance wrong!')
    // assert.equal(balanceReceiver,value, 'receiver balance wrong!')
  })
});
