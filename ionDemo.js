const readline = require('readline');

const Web3 = require('web3')

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
    ['\tCurrency hash:',le.args._currency],
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

  await waitForKeypress()

  console.log(`Transfer tokens from Owner account to ${senderName}`)
  const transferTxHash = await transferToken(token, ownerAccount, senderAccount, value, reference)

  await waitForKeypress()

  // setup filter to get ionlock event
  console.log(`Transfer tokens ${senderName} to IonLock`)
  const lockTxHash = await transferToken(token, senderAccount, ionLock.address, value, reference)

  await waitForKeypress()

  // Wait for IonLock Event
  console.log('IonLock event triggered')
  const lockEvent = await getIonLockEvent(web3, ionLock, reference)

  console.log(`\nEND - Deposit from ${senderName} to IonLock`)
  await waitForKeypress()
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

const queryLithium = () = {
}

//UNTESTED
const withdrawIonLock = async (ionLock, value, ref, blockId, proof) => {
  const withdrawTx = await ionLock.Withdraw(value, ref, blockId, proof)
  return withdrawTx
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

  // deploy contracts
  const tokenA = await deployContract(web3A, './build/contracts/Token.json', owner)
  const ionLinkA = await deployContract(web3A, './build/contracts/IonLink.json', owner)
  const ionLockA = await deployContract(web3A, './build/contracts/IonLock.json', owner, [tokenA.address, ionLinkA.address])
  // deploy contracts
  const tokenB = await deployContract(web3B, './build/contracts/Token.json', owner)
  const ionLinkB = await deployContract(web3B, './build/contracts/IonLink.json', owner)
  const ionLockB = await deployContract(web3B, './build/contracts/IonLock.json', owner, [tokenB.address, ionLinkB.address])

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
  const reference = 'Simple Ion Example'

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
  console.log('================= Wait for Lithium to continue =================')
  console.log('\n')
  // TODO: get IonLink block id (for the deposit) for chain A and B
  // TODO: get IonLink proof (for the deposit) for chain A and B
  await waitForKeypress()

  // TODO: WITHDRAW
  //const withdrawTxA = await withdrawIonLock(ionLockA, valueA, refA, blockIdA, proofA)
  //const withdrawTxB = await withdrawIonLock(ionLockB, valueB, refB, blockIdB, proofB)

  await printTokenBalance(
    accountA, tokenA, 'Alice',
    accountB, tokenB, 'Bob',
    ionLockA, ionLockB)


  process.exit()
}

try {
main()
} catch (err) {
  console.log('error in main()',err)
}
