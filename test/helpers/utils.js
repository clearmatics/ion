// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const crypto = require('crypto')

// Format required for sending bytes through eth client:
//  - hex string representation
//  - prefixed with 0x
const bufToStr = b => '0x' + b.toString('hex')

const gasPrice = 100000000000 // truffle fixed gas price
const joinHex = arr => '0x' + arr.map(el => el.slice(2)).join('')
const oneFinney = web3.toWei(1, 'finney')


const sha256 = x =>
  crypto
    .createHash('sha256')
    .update(x)
    .digest()

const random32 = () => crypto.randomBytes(32)

const randomHex = () => crypto.randomBytes(32).toString('hex');

const randomArr = () => {
  const result = []
  const size =(Math.floor(Math.random() * 10) + 1);
  for(let i = size; 0 < i; i-- )
    result.push(randomHex())
  return result
}

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

const txGas = txReceipt => txReceipt.receipt.gasUsed * gasPrice
const txLoggedArgs = txReceipt => txReceipt.logs[0].args
const txContractId = txReceipt => txLoggedArgs(txReceipt).contractId

module.exports = {bufToStr, joinHex, newSecretHashPair, oneFinney, random32, randomArr, randomHex, sha256, sleep, txGas, txLoggedArgs}
