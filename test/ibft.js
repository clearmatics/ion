// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

/*
    Clique Validation contract test

    Tests here are standalone unit tests for clique module functionality.
    Other contracts have been mocked to simulate basic behaviour.

    Tests the clique scheme for block submission, validator signature verification and more.
*/

const eth_util = require('ethereumjs-util');
const utils = require('./helpers/utils.js');
const encoder = require('./helpers/encoder.js');
const Web3 = require('web3');
const Web3Utils = require('web3-utils');
const rlp = require('rlp');
const truffleAssert = require('truffle-assertions');
const sha3 = require('js-sha3').keccak_256

const Clique = artifacts.require("Clique");
const MockIon = artifacts.require("MockIon");
const MockStorage = artifacts.require("MockStorage");

const web3 = new Web3();
const autonity = new Web3();

web3.setProvider(new web3.providers.HttpProvider('http://localhost:8545'));
autonity.setProvider(new web3.providers.HttpProvider('http://localhost:9501'));
// autonity.setProvider(new web3.providers.HttpProvider('https://rinkeby.infura.io'));
// autonity.setProvider(new web3.providers.HttpProvider('http://34.243.204.94:30001'));

require('chai')
 .use(require('chai-as-promised'))
 .should();

// Takes a header and private key returning the signed data
// Needs extraData just to be sure of the final byte
signHeader = (headerHash, privateKey, extraData) => {
  const sig = eth_util.ecsign(headerHash, privateKey)
  if (this._chainId > 0) {
    sig.v += this._chainId * 2 + 8
  }
  
  const pubKey  = eth_util.ecrecover(headerHash, sig.v, sig.r, sig.s);
  const addrBuf = eth_util.pubToAddress(pubKey);
  
  const newSigBytes = Buffer.concat([sig.r, sig.s]);
  let newSig;
  
  const bytes = utils.hexToBytes(extraData)
  const finalByte = bytes.splice(bytes.length-1)
  if (finalByte.toString('hex')=="0") {
    newSig = newSigBytes.toString('hex') + '00';
  }
  if (finalByte.toString('hex')=="1") {
    newSig = newSigBytes.toString('hex') + '01';
  }

  return newSig;
}

function pad(n, width, z) {
  z = z || '0';
  n = n + '';
  return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
}

const DEPLOYEDCHAINID = "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075"
const TESTCHAINID = "0x22b55e8a4f7c03e1689da845dd463b09299cb3a574e64c68eafc4e99077a7254"


const VALIDATORS = [
  "0x4bf2776241283242b13f6501454cf27345280f13",
  "0x0d066e7626cadefa5465299d46b17d552c9da5bc",
  "0x7f4f1ee3c39b7be3e0ca54e6a3504004a873f934",
  "0x6c6fe934c860270fd57a541dd6f92e5ba456113c",
  "0x70b39355f8b1cd5106d1fce9cd0b00d77585c406",
  "0xa030cd716d381d165cbbe893087e679ac23d60aa",
  "0x0ce0cb6eb8df3075c1c4b9de2a59b24ce06b1a3b"
];

const VALIDATORS_START = ["0x42eb768f2244c8811c63729a21a3569731535f06", "0x7ffc57839b00206d1ad20c69a1981b489f772031", "0xb279182d99e65703f0076e4812653aab85fca0f0"];
const VALIDATORS_FINISH = ["0x42eb768f2244c8811c63729a21a3569731535f06", "0x6635f83421bf059cd8111f180f0727128685bae4", "0x7ffc57839b00206d1ad20c69a1981b489f772031", "0xb279182d99e65703f0076e4812653aab85fca0f0"];
const GENESIS_HASH = "0xf32b505a5ad95dfa88c2bd6904a1ba81a92a1db547dc17f4d7c0f64cf2cddbb1";
const ADD_VALIDATORS_GENESIS_HASH = "0xf32b505a5ad95dfa88c2bd6904a1ba81a92a1db547dc17f4d7c0f64cf2cddbb1";

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


contract.only('Ibft.js', (accounts) => {
  const joinHex = arr => '0x' + arr.map(el => el.slice(2)).join('');

  const watchEvent = (eventObj) => new Promise((resolve,reject) => eventObj.watch((error,event) => error ? reject(error) : resolve(event)));
  
  describe('Submit Block', () => {
      it.only('Authentic Submission Happy Path', async () => {
        // await clique.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH, storage.address);
        let block = await autonity.eth.getBlock(1);

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

        console.log("\nExtra Data:")
        console.log(extraData);
        let istExtraData = extraData.slice(66);
        
        console.log("\nIstanbul Extra Data:")
        console.log('0x' + istExtraData);
        let rlpExtraData = rlp.decode('0x' + istExtraData);

        let sig = '0x' + rlpExtraData[1].toString('hex');
        let validators = rlpExtraData[0];

        console.log("\nValidators:");
        validators.forEach( function(entry) {
            console.log(entry.toString('hex'))
          }
        );

        // Remove the committed seals
        committedSeals = rlpExtraData[2];
        rlpExtraData[2] = [];

        let rlpEncodedExtraDataSeal = rlp.encode(rlpExtraData);

        console.log('0x' + rlpEncodedExtraDataSeal.toString('hex'));

        // Remove last 65 Bytes of extraData
        let extraBytes = hexToBytes(extraData);
        let extraBytesShort = extraBytes.splice(1, 32);
        let extraDataShort = '0x' + bytesToHex(extraBytesShort) + rlpEncodedExtraDataSeal.toString('hex');


        console.log("\nSigned Commit Block\n")
        
        console.log("\nExtraData:")
        console.log(extraData);
        console.log("\nExtraData Short:")
        console.log(extraDataShort);

        let header = [
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
          extraDataShort,
          mixHash,
          nonce
        ];

        let testBlockHeader = '0x' + rlp.encode(header).toString('hex');
        let testBlockHeaderHash = Web3Utils.sha3(testBlockHeader);

        console.log("\nBlock hash:")
        console.log(testBlockHeaderHash);
        console.log(block.hash);

        // Create the rlp encoded extra data
        rlpExtraData[1] = new Buffer([]);
        rlpExtraData[2] = [];

        console.log("\rlp extra data")
        console.log(rlpExtraData)

        rlpEncodedExtraDataSeal = rlp.encode(rlpExtraData);

        console.log('0x' + rlpEncodedExtraDataSeal.toString('hex'));

        // Remove last 65 Bytes of extraData
        extraBytes = hexToBytes(extraData);
        extraBytesShort = extraBytes.splice(1, 32);
        extraDataShort = '0x' + bytesToHex(extraBytesShort) + rlpEncodedExtraDataSeal.toString('hex');

        console.log("\nSigned Seal Block\n")


        console.log("\nExtraData:")
        console.log(extraData);
        console.log("\nExtraData Short:")
        console.log(extraDataShort);

        header = [
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
          extraDataShort,
          mixHash,
          nonce
        ];

        testBlockHeader = '0x' + rlp.encode(header).toString('hex');
        testBlockHeaderHash = Web3Utils.sha3(testBlockHeader);

        console.log("\nBlock hash:")
        console.log(testBlockHeaderHash);
        console.log(Web3Utils.sha3(testBlockHeaderHash));
        console.log(block.hash);

        console.log("\nSignature Retrieved From Istanbul Extra:");
        console.log(sig);

        const blockHeaderHash = eth_util.sha3(testBlockHeaderHash);
        console.log(blockHeaderHash)
        const {v, r, s} = eth_util.fromRpcSig(sig);

        const pubKey  = eth_util.ecrecover(blockHeaderHash, v, r, s);
        const addrBuf = eth_util.pubToAddress(pubKey);
        // assert.equal(VALIDATORS[0], '0x'+addrBuf.toString('hex'));
        console.log('0x'+addrBuf.toString('hex'));
 
      })

  })
});
