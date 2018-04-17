const readline = require('readline');

const Web3 = require('web3')

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

const main = async () => {
  const web3A = getWeb3('http://localhost:8545')
  //const web3B = getWeb3('http://localhost:8546')

  const owner = web3A.eth.accounts[0]
  const accountA = web3A.eth.accounts[1]
  //const accountB = web3A.eth.accounts[2]

  // deploy contracts
  const tokenA = await deployContract(web3A, './build/contracts/Token.json', owner)
  const ionLinkA = await deployContract(web3A, './build/contracts/IonLink.json', owner)
  const ionLockA = await deployContract(web3A, './build/contracts/IonLock.json', owner, [tokenA.address, ionLinkA.address])

  printBlock([
    ['Deployed Contracts'],
    ['Token:',tokenA.address],
    ['IonLock:',ionLockA.address],
    ['IonLink:',ionLinkA.address]])

  //TODO: PRINT BALANCES WHEN RELEVANT

  await waitForKeypress()

  const value = 1000
  const reference = 'Example A'

  console.log('Mint tokens into the Owner account')
  const mintTxHash = await mintToken(tokenA, owner, value)

  console.log('Transfer tokens from Owner account to Account A')
  const transferTxHash = await transferToken(tokenA, owner, accountA, value, reference)

  await waitForKeypress()

  // setup filter to get ionlock event
  console.log('Transfer tokens Account A to IonLock')
  const lockTxHash = await transferToken(tokenA, accountA, ionLockA.address, value, reference)

  await waitForKeypress()

  // Wait for IonLock Event
  console.log('IonLock event triggered')
  const lockEvent = await getIonLockEvent(web3A, ionLockA, reference)


  // TODO: WAIT FOR UPDATE IN IONLINK
  // TODO: WITHDRAW
  process.exit()
}

try {
main()
} catch (err) {
  console.log('error in main()',err)
}
