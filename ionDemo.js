const Web3 = require('web3')

const deployContract = (contractPath, ownerAcc, args) => {
  const contractData = require(contractPath)
  const abi = contractData.abi
  const bytecode = contractData.bytecode
  const contract = web3.eth.contract(abi)
  const txObj = {data: bytecode, from: ownerAcc, gas: '1000000'}
  //const contractInstance = contract.new(...(args || []), txObj)
  //return contractInstance
  const contractPromise = new Promise((resolve, reject) => contract.new(...(args || []), txObj, (err,contractInstance) => {
    if(err) reject(err)
    if(contractInstance.address) resolve(contractInstance)
  }))
  return contractPromise
}

const loadContract = (address,abiPath) => {
  const abi = require(abiPath).abi
  const contract = web3.eth.contract(abi);
  const instance = contract.at(address)
  return instance
}

const waitForKeypress = () => {
  console.log('\n=== Press any key to continue ===\n')
  // TODO: HALT UNTIL KEYPRESS
}

const printBlock = (args) => {
  console.log('===========================================================')
  args.forEach(a => console.log(...a))
  console.log('===========================================================')
}

const main = async () => {
  if (typeof web3 !== 'undefined') {
    web3 = new Web3(web3.currentProvider);
  } else {
    // set the provider you want from Web3.providers
    const providerURL = 'http://localhost:8545'
    web3 = new Web3(new Web3.providers.HttpProvider(providerURL));
  }

  const owner = web3.eth.accounts[0]
  const accountA = web3.eth.accounts[1]
  const accountB = web3.eth.accounts[2]

  // expects the contracts to have been deployed previously
  //const token = loadContract('0x9561c133dd8580860b6b7e504bc5aa500f0f06a7','./build/contracts/Token.json')
  //const ionLink = loadContract('0xc89ce4735882c9f0f0fe26686c53074e09b0d550','./build/contracts/IonLink.json')
  //const ionLock = loadContract('0xe982e462b094850f12af94d21d470e21be9d0e9c','./build/contracts/IonLock.json')
  const token = await deployContract('./build/contracts/Token.json', owner)
  const ionLink = await deployContract('./build/contracts/IonLink.json', owner)
  const ionLock = await deployContract('./build/contracts/IonLock.json', owner, [token.address, ionLink.address])

  printBlock([
    ['Deployed Contracts'],
    ['Token:',token.address],
    ['IonLock:',ionLock.address],
    ['IonLink:',ionLink.address]])

  waitForKeypress()

  const value = 1000
  const reference = 'Example A'
  const mintTxHash = await token.mint(value,{ from: owner, gas: '1000000' })
  const transferTxHash = await token.transfer['address,uint256,bytes'](accountA, value, reference, { from: owner, gas: '1000000' })

  printBlock([
    ['Minted tokens to Owner and transfered from Owner to A'],
    ['Owner Account:', owner],
    ['Account A:', accountA],
    ['Value:',value],
    ['Mint TxHash:',mintTxHash],
    ['Transfer TxHash:',transferTxHash]
  ])

  waitForKeypress()

  // setup filter to get ionlock event
  const lockTxHash = await token.transfer['address,uint256,bytes'](ionLock.address, value, reference, { from: accountA, gas: '1000000' })

  printBlock([
    ['Transfered tokens from A to IonLock'],
    ['Account A:', accountA],
    ['IonLock:', ionLock.address],
    ['Value:',value],
    ['Reference:',reference],
    ['Transfer to IonLock TxHash:',lockTxHash],
  ])

  waitForKeypress()

  // Wait for IonLock Event
  const lockTransferEvent = ionLock.IonTransfer({ ref: web3.sha3(reference) }) // filter by reference
  const lockEvent = await new Promise((resolve,reject) => lockTransferEvent.get((err,result) => err ? reject(err) : resolve(result)))

  printBlock([
    ['IonLock event triggered'],
    ['TxHash:', lockEvent[0].transactionHash],
    ['Arguments of event:'],
    ['\tRecipient:',lockEvent[0].args._recipient],
    ['\tCurrency hash:',lockEvent[0].args._currency],
    ['\tValue:',lockEvent[0].args.value.toString()],
    ['\tReference hash:',lockEvent[0].args.ref],
    ['\tData (hex of reference):',lockEvent[0].args.data],
  ])
}

try {
main()
} catch (err) {
  console.log('error in main()',err)
}
