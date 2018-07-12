// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const Web3 = require('web3');
const Web3Utils = require('web3-utils');
const Web3Abi = require('web3-eth-abi');
const rlp = require('rlp');

const web3 = new Web3();

web3.setProvider(new web3.providers.HttpProvider('http://localhost:8501'));

function hexToBytes(hex) {
    for (var bytes = [], c = 0; c < hex.length; c += 2)
    bytes.push(parseInt(hex.substr(c, 2), 16));
    return bytes;
}

function bytesToHex(bytes) {
    for (var hex = [], i = 0; i < bytes.length; i++) {
        hex.push((bytes[i] >>> 4).toString(16));
        hex.push((bytes[i] & 0xF).toString(16));
    }
    return hex.join("");
}


const block = web3.eth.getBlock(7);
// console.log("\n\n\n\nBlock = " + JSON.stringify(block))
// Decompose the values in the block to hash
const parentHash = block.parentHash;
const sha3Uncles = block.sha3Uncles;
const coinbase = block.miner;
const root = block.stateRoot;
const txHash = block.transactionsRoot;
const receiptHash = block.receiptsRoot;
const logsBloom = block.logsBloom;
const difficulty = Web3Utils.toBN(block.difficulty);
const number = Web3Utils.toBN(block.number);
const gasLimit = block.gasLimit;
const gasUsed = block.gasUsed;
const timestamp = Web3Utils.toBN(block.timestamp);
const extraData = block.extraData;
const mixHash = block.mixHash;
const nonce = block.nonce;

// Remove last 65 Bytes of extraData
const extraBytes = hexToBytes(extraData);
const extraBytesShort = extraBytes.splice(1, extraBytes.length-66);
const extraDataSignature = '0x' + bytesToHex(extraBytes.splice(extraBytes.length-65));
const extraDataShort = '0x' + bytesToHex(extraBytesShort);

const header = [
  parentHash,
  sha3Uncles,
  coinbase,
  root,
  txHash,
  receiptHash,
  logsBloom,
  difficulty,
  number,
  gasLimit,
  gasUsed,
  timestamp,
  extraData,
  mixHash,
  nonce
];

const encodedHeader = rlp.encode(header);

const headerHash = Web3Utils.sha3(encodedHeader);

// console.log("\n\n\nRLP-encoded header: " +bytesToHex(encodedHeader));
// console.log("\n\n\nHashed header: " +headerHash)
// console.log("\n\n\nExpected hash: " +block.hash)
// console.log("\n\n\n" + ((rlp.encode(parentHash))))