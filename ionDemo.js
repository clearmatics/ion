const readline = require('readline');

require('events').EventEmitter.prototype._maxListeners = 15;

const Web3 = require('web3')

const Web3Utils = require('web3-utils');

const BN = require('bignumber.js')

const fs = require('fs');

const merkle = require('./test/helpers/merkle.js')

const color = {
  Reset: '\x1b[0m',
  Bright: '\x1b[1m',
  Dim: '\x1b[2m',
  Underscore: '\x1b[4m',
  Blink: '\x1b[5m',
  Reverse: '\x1b[7m',
  Hidden: '\x1b[8m',

  FgBlack: '\x1b[30m',
  FgRed: '\x1b[31m',
  FgGreen: '\x1b[32m',
  FgYellow: '\x1b[33m',
  FgBlue: '\x1b[34m',
  FgMagenta: '\x1b[35m',
  FgCyan: '\x1b[36m',
  FgWhite: '\x1b[37m',

  BgBlack: '\x1b[40m',
  BgRed: '\x1b[41m',
  BgGreen: '\x1b[42m',
  BgYellow: '\x1b[43m',
  BgBlue: '\x1b[44m',
  BgMagenta: '\x1b[45m',
  BgCyan: '\x1b[46m',
  BgWhite: '\x1b[47m',
}

const tokenJson = require('./build/contracts/Token.json')
const tokenAbi = tokenJson.abi
const lockJson = require('./build/contracts/IonLock.json')
const lockAbi = lockJson.abi
const linkJson = require('./build/contracts/IonLink.json')
const linkAbi = linkJson.abi

const deployContract = (web3, contractPath, ownerAcc, args) => {
  const contractData = require(contractPath)
  const abi = contractData.abi
  const bytecode = contractData.bytecode
  const contract = web3.eth.contract(abi)
  const txObj = {data: bytecode, from: ownerAcc, gas: '1000000'}
  //const contractInstance = contract.new(...(args || []), txObj)
  //return contractInstance
  const contractPromise =
    new Promise((resolve, reject) => contract.new(...(args || []), txObj, (err,contractInstance) => {
      if(err) reject(err)
      if(contractInstance.address) resolve(contractInstance)
    }))
  return contractPromise
}

const mintToken = async (token, owner, value) => {
  const mintTxObj = { from: owner, gas: '1000000' }
  const mintTxHash = await token.mint(value, mintTxObj)

  printBlock([
    ['Minted tokens to',owner],
    ['Value:',value],
    ['Mint TxHash:',mintTxHash],
  ])

  return mintTxHash
}

const transferToken = async (token, from, to, value, reference) => {
  const transfer = token.transfer['address,uint256,bytes']
  const transferTxObj = { from, gas: '1000000' }
  const transferTxHash = await transfer(to, value, reference, transferTxObj)

  printBlock([
    ['Transfered tokens from',from,'to',to],
    ['Value:', value],
    ['Reference:',reference],
    ['Transfer TxHash:',transferTxHash]
  ])

  return transferTxHash
}

const getIonLockEvent = async (web3, ionLock, reference) => {
  const refHash = web3.sha3(reference)
  const lockTransferEvent = ionLock.IonTransfer({ ref: refHash }) // filter by reference
  const promiseCallback = (resolve,reject) => {
    lockTransferEvent.get((err,result) => err ? reject(err) : resolve(result))
  }
  const lockEvent = await new Promise(promiseCallback)

  const printArr = lockEvent.reduce((prev,le) => prev.concat([
    ['IonLock event triggered'],
    ['Name:',le.event],
    ['TxHash:', le.transactionHash],
    ['Arguments of event:'],
    ['\tRecipient:',le.args._recipient],
    ['\tCurrency address:',le.args._currency],
    ['\tValue:',le.args.value.toString()],
    ['\tReference hash:',le.args.ref],
    ['\tData (hex of reference):',le.args.data],
  ]), [])
  printBlock(printArr)

  return lockEvent
}

const getWeb3 = (providerURL) => {
  return new Web3(new Web3.providers.HttpProvider(providerURL))
}

const waitForKeypress = async () => {
  // {}
  readline.emitKeypressEvents(process.stdin)
  process.stdin.setRawMode(true)
  const promiseKey = new Promise((resolve, reject) => process.stdin.on('keypress', (str, key) => {
    if (key.ctrl && key.name === 'c') {
      process.exit()
    } else {
      resolve(key)
    }
  }))
  console.log('\n=== Press any key to continue ===\n')
  return promiseKey
}

const printBlock = (args) => {
  console.log('===========================================================')
  args.forEach(a => console.log(...a))
  console.log('===========================================================')
}

const depositIonLock = async (web3, token, ionLock, value, reference, ownerAccount ,senderAccount,senderName) => {
  console.log(`\nSTART - Deposit from ${senderName} to IonLock\n`)

  console.log('Mint tokens into the Owner account')
  const mintTxHash = await mintToken(token, ownerAccount, value)

  // await waitForKeypress()

  console.log(`Transfer tokens from Owner account to ${senderName}`)
  const transferTxHash = await transferToken(token, ownerAccount, senderAccount, value, reference)

  // await waitForKeypress()

  // setup filter to get ionlock event
  console.log(`Transfer tokens from ${senderName} to IonLock`)
  const lockTxHash = await transferToken(token, senderAccount, ionLock.address, value, reference)

  // await waitForKeypress()

  // Wait for IonLock Event
  console.log('IonLock event triggered')
  const lockEvent = await getIonLockEvent(web3, ionLock, reference)

  console.log(`\nEND - Deposit from ${senderName} to IonLock`)
  // await waitForKeypress()
}

const printTokenBalance = async (accountA, tokenA, nameA, accountB, tokenB, nameB, ionLockA, ionLockB) => {
  const balanceAA = await tokenA.balanceOf(accountA)
  const balanceBA = await tokenA.balanceOf(accountB)
  const balanceAB = await tokenB.balanceOf(accountA)
  const balanceBB = await tokenB.balanceOf(accountB)
  const balanceIonLockA = await tokenA.balanceOf(ionLockA.address)
  const balanceIonLockB = await tokenB.balanceOf(ionLockB.address)
  printBlock([
    ['Balance of Tokens'],
    [`Balance of ${nameA} on Token A ${tokenA.address}`],
    ['\tAccount:',accountA],
    ['\tBalance:',balanceAA.toString()],
    [`Balance of ${nameB} on Token A ${tokenA.address}`],
    ['\tAccount:',accountB],
    ['\tBalance:',balanceBA.toString()],
    [`Balance of ${nameA} on Token B ${tokenB.address}`],
    ['\tAccount:',accountA],
    ['\tBalance:',balanceAB.toString()],
    [`Balance of ${nameB} on Token B ${tokenB.address}`],
    ['\tAccount:',accountB],
    ['\tBalance:',balanceBB.toString()],
    [`Balance of IonLock on Token A ${tokenA.address}`],
    ['\tAccount:',ionLockA.address],
    ['\tBalance:',balanceIonLockA.toString()],
    [`Balance of IonLock on Token B ${tokenB.address}`],
    ['\tAccount:',ionLockB.address],
    ['\tBalance:',balanceIonLockB.toString()],
  ])
}

const printReferenceData = async (refA, refB, proofA, proofB, blockIdA, blockIdB) => {
  printBlock([
    ['Withdrawal Reference Data'],
    [`Alice on Chain B`],
    ['\tReference: ',refA],
    ['\tProof: ',proofA],
    ['\tBlockId: ',blockIdB.toString()],
    [`Bob on Chain A`],
    ['\tReference: ',refB],
    ['\tProof: ',proofB],
    ['\tBlockId: ',blockIdA.toString()],
  ])
}

const queryLithium = () => {
}

const joinIonLinkData = (receiverAddr,tokenAddr,ionLockAddr,value,reference) => {
  const valueHex = '0x'+Web3Utils.toBN(value).toString(16).padStart(64,'0') // make an hex that is good to sha3 in solidity uint256 -> 64 bytes
  const leaf = '0x' + [receiverAddr,tokenAddr,ionLockAddr,valueHex,reference].map(el=>el.slice(2)).join('') // joined args need to be added to the random leafs of the tree
  return leaf
}

//UNTESTED
const withdrawIonLock = async (ionLock, value, ref, blockId, proof, account) => {
  const withdrawTx = await ionLock.Withdraw(value, ref, blockId, proof, {from: account, gas: "0xFFFFFD"})
  return withdrawTx
}

const reader = (input) => {
  const output = fs.readFileSync(input, 'utf8');
  return output;
}

const arrayReader = (input) => {
  const output = fs.readFileSync(input, 'utf8').toString().split("\n");
  // Remove any empty elements
  var index = output.indexOf('');
  if (index > -1) {
      output.splice(index, 1);
  }
  return output;
}

const convertArrayBN = (input) => {
  var arrayLength = input.length;
  for (var i = 0; i < arrayLength; i++) {
      input[i] = new BN(input[i])
      // console.log(input[i]);
  }
  return input;
}

const main = async () => {
  console.log('NodeA at http://localhost:8545')
  console.log('NodeB at http://localhost:8546')
  console.log()
  const web3A = getWeb3('http://localhost:8545')
  const web3B = getWeb3('http://localhost:8546')

  const owner = web3A.eth.accounts[0]
  const accountA = web3A.eth.accounts[1]
  const accountB = web3A.eth.accounts[2]

  // console.log(tokenAbi)
  // instantiate by address

  // // deploy contracts
  const TokenA = web3A.eth.contract(tokenAbi);
  const tokenA = TokenA.at('0x9561c133dd8580860b6b7e504bc5aa500f0f06a7');
  const LockA = web3A.eth.contract(lockAbi);
  const ionLockA = LockA.at('0xe982e462b094850f12af94d21d470e21be9d0e9c');
  const LinkA = web3A.eth.contract(linkAbi);
  const ionLinkA = LinkA.at('0xc89ce4735882c9f0f0fe26686c53074e09b0d550');

  // deploy contracts
  const TokenB = web3B.eth.contract(tokenAbi);
  const tokenB = TokenB.at('0x9561c133dd8580860b6b7e504bc5aa500f0f06a7');
  const LockB = web3B.eth.contract(lockAbi);
  const ionLockB = LockB.at('0xe982e462b094850f12af94d21d470e21be9d0e9c');
  const LinkB = web3B.eth.contract(linkAbi);
  const ionLinkB = LinkB.at('0xc89ce4735882c9f0f0fe26686c53074e09b0d550');

  printBlock([
    ['Deployed Contracts NodeA'],
    ['\tToken:',tokenA.address],
    ['\tIonLock:',ionLockA.address],
    ['\tIonLink:',ionLinkA.address],
    ['Deployed Contracts NodeB'],
    ['\tToken:',tokenB.address],
    ['\tIonLock:',ionLockB.address],
    ['\tIonLink:',ionLinkB.address]])

  await printTokenBalance(
    accountA, tokenA, 'Alice',
    accountB, tokenB, 'Bob',
    ionLockA, ionLockB)

  const value = 1000
  const date = Math.floor(Date.now() / 1000);
  const reference = Web3Utils.sha3('Reference from deposit on chain B')

  console.log('\n\n\n')
  console.log('|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||')
  console.log('VVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVV')
  console.log('\n')
  console.log('================= Deposit to IonLock on chain A =================')
  await waitForKeypress()
  await depositIonLock(
    web3A, tokenA, ionLockA,
    value, reference,
    owner, accountA, 'Alice')

  await printTokenBalance(
    accountA, tokenA, 'Alice',
    accountB, tokenB, 'Bob',
    ionLockA, ionLockB)

  console.log('\n\n\n')
  console.log('|||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||')
  console.log('VVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVV')
  console.log('\n')
  console.log('================= Deposit to IonLock on chain B =================')
  await waitForKeypress()
  await depositIonLock(
    web3B, tokenB, ionLockB,
    value, reference,
    owner, accountB, 'Bob')

  await printTokenBalance(
    accountA, tokenA, 'Alice',
    accountB, tokenB, 'Bob',
    ionLockA, ionLockB)

  // TODO: WAIT FOR UPDATE IN IONLINK
  console.log('\n\n\n')
  console.log('================= Lithium Withdraw =================')
  console.log('\n')

  // TODO: get IonLink block id (for the deposit) for chain A
  var refA = await reader('./data/reference8545.txt')
  var proofA = await arrayReader('./data/merklePath8545.txt')
  var blockIdB = await reader('./data/latestBlock8546.txt')

  // TODO: get IonLink proof (for the deposit) for chain B
  var refB = await reader('./data/reference8546.txt')
  var proofB = await arrayReader('./data/merklePath8546.txt')
  var blockIdA = await reader('./data/latestBlock8545.txt')

  // await printReferenceData(refA, refB, proofA, proofB, blockIdB, blockIdA)

  // Convert the strings into javascript stuff
  blockIdA = new BN(blockIdA)
  blockIdB = new BN(blockIdB)

  proofA = convertArrayBN(proofA)
  proofB = convertArrayBN(proofB)

  await waitForKeypress()

  // TODO: WITHDRAW
  var leafB = joinIonLinkData(accountA,tokenB.address,ionLockB.address,value,refA)
  leafhashB = merkle.merkleHash(leafB)
  var leafA = joinIonLinkData(accountB,tokenA.address,ionLockA.address,value,refB)
  leafhashA = merkle.merkleHash(leafA)

  const validB = await ionLinkB.Verify(blockIdB, leafhashB, proofA, {from: accountA, gas: "0xFFFD"})
  const valueB = await tokenB.balanceOf(ionLockB.address) // Alice on chain B
  console.log('\n\n\n')
  console.log('================= Check IonLink =================')
  console.log('\n')
  console.log("IonLink Verify on chain B: ", validB)
  console.log("Balance of IonLock on chain B: ", valueB)
  await waitForKeypress()

  const validA = await ionLinkA.Verify(blockIdA, leafhashA, proofB, {from: accountB, gas: "0xFFFD"})
  const valueA = await tokenA.balanceOf(ionLockA.address) // Bob on chain A
  console.log("IonLink Verify on chain A: ", validA)
  console.log("Balance of IonLock on chain A: ", valueA)
  await waitForKeypress()

  console.log('\n\n\n')
  console.log('================= Alice Withdraw on Chain B =================')
  console.log('\n')
  await waitForKeypress()

  const withdrawTxA = await withdrawIonLock(ionLockB, value, refA, blockIdB, proofA, accountA)
  // console.log(withdrawTxA)

  console.log('\n\n\n')
  console.log('================= Bob Withdraw on Chain A =================')
  console.log('\n')
  await waitForKeypress()
  const withdrawTxB = await withdrawIonLock(ionLockA, value, refB, blockIdA, proofB, accountB)
  // console.log(withdrawTxB)
  await printTokenBalance(
    accountA, tokenA, 'Alice',
    accountB, tokenB, 'Bob',
    ionLockA, ionLockB)

  await waitForKeypress()

  process.exit()
}

try {
main()
} catch (err) {
  console.log('error in main()',err)
}
