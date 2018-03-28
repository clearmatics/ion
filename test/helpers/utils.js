const crypto = require('crypto')

const gasPrice = 100000000000 // truffle fixed gas price
const txGas = txReceipt => txReceipt.receipt.gasUsed * gasPrice
const txLoggedArgs = txReceipt => txReceipt.logs[0].args
const txContractId = txReceipt => txLoggedArgs(txReceipt).contractId
const oneFinney = web3.toWei(1, 'finney')

// Format required for sending bytes through eth client:
//  - hex string representation
//  - prefixed with 0x
const bufToStr = b => '0x' + b.toString('hex')

const sha256 = x =>
  crypto
    .createHash('sha256')
    .update(x)
    .digest()

const random32 = () => crypto.randomBytes(32)

const isSha256Hash = hashStr => /^0x[0-9a-f]{64}$/i.test(hashStr)

const newSecretHashPair = () => {
  const secret = random32()
  const hash = sha256(secret)
  return {
    secret: bufToStr(secret),
    hash: bufToStr(hash),
  }
}

const sleep = ms => {
  return new Promise(resolve => setTimeout(resolve, ms));
}

module.exports = {newSecretHashPair, oneFinney, sleep}
