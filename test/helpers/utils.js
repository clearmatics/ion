// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const crypto = require('crypto')

const utils = {};
  // Format required for sending bytes through eth client:
  //  - hex string representation
  //  - prefixed with 0x
utils.bufToStr = b => '0x' + b.toString('hex')

utils.gasPrice = 100000000000 // truffle fixed gas price
utils.joinHex = arr => '0x' + arr.map(el => el.slice(2)).join('')
utils.oneFinney = web3.toWei(1, 'finney')

utils.hexToBytes = (hex) => {
  for (var bytes = [], c = 0; c < hex.length; c += 2)
  bytes.push(parseInt(hex.substr(c, 2), 16));
  return bytes;
}

utils.bytesToHex = (bytes) => {
  for (var hex = [], i = 0; i < bytes.length; i++) {
      hex.push((bytes[i] >>> 4).toString(16));
      hex.push((bytes[i] & 0xF).toString(16));
  }
  return hex.join("");
}

utils.sha256 = x =>
  crypto
    .createHash('sha256')
    .update(x)
    .digest()

utils.random32 = () => crypto.randomBytes(32)

utils.randomHex = () => crypto.randomBytes(32).toString('hex');

utils.randomArr = () => {
  const result = []
  const size =(Math.floor(Math.random() * 10) + 1);
  for(let i = size; 0 < i; i-- )
    result.push(randomHex())
  return result
}

utils.isSha256Hash = hashStr => /^0x[0-9a-f]{64}$/i.test(hashStr)

const newSecretHashPair = () => {
  const secret = random32()
  const hash = sha256(secret)
  return {
    secret: bufToStr(secret),
    hash: bufToStr(hash),
  }
}

utils.sleep = ms => {
  return new Promise(resolve => setTimeout(resolve, ms));
}

utils.txGas = txReceipt => txReceipt.receipt.gasUsed * gasPrice
utils.txLoggedArgs = txReceipt => txReceipt.logs[0].args
utils.txContractId = txReceipt => txLoggedArgs(txReceipt).contractId

module.exports = utils;