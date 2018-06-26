// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const Util = require('ethereumjs-util');
const Web3 = require('web3');
const Web3Utils = require('web3-utils');
const Web3Abi = require('web3-eth-abi');
const Web3Accounts = require('web3-eth-accounts');
const rlp = require('rlp');

const Validation = artifacts.require("Validation");
const Recover = artifacts.require("Recover");

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

contract.only('test.js', (accounts) => {
  const blockNum = 10
  const joinHex = arr => '0x' + arr.map(el => el.slice(2)).join('');

  const watchEvent = (eventObj) => new Promise((resolve,reject) => eventObj.watch((error,event) => error ? reject(error) : resolve(event)));

  const ionAbi = [{"constant":false,"inputs":[{"name":"header","type":"bytes"}],"name":"ValidationTest","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"LatestBytes","outputs":[{"name":"_latestBytes","type":"bytes"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"m_blockheaders","outputs":[{"name":"prevBlockHash","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"LatestBlock","outputs":[{"name":"_latestBlock","type":"bytes32"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"header","type":"bytes"},{"name":"prefixHeader","type":"bytes"},{"name":"prefixExtraData","type":"bytes"}],"name":"ValidateBlock","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"GetValidators","outputs":[{"name":"_validators","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"_validators","type":"address[]"},{"name":"genHash","type":"bytes32"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"owner","type":"address"}],"name":"broadcastSig","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"header","type":"bytes"},{"indexed":false,"name":"parentHash","type":"bytes"},{"indexed":false,"name":"rootHash","type":"bytes"}],"name":"broadcastHashData","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"blockHash","type":"bytes32"}],"name":"broadcastHash","type":"event"}];
  const validators = ["0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d", "0x8671e5e08d74f338ee1c462340842346d797afd3"];
  const genHash = "0xc3bac257bbd04893316a76d41b6ff70de5f65c9f24db128864a6322d8e0e2f28";


  // Here the block header is signed off chain but by a whitelisted validator
  it('Test: Authentic Submission Off-Chain Signature - ValidateBlock()', async () => {
    const accounts = web3.eth.accounts;
    const signer = accounts[0];
    const ValidationContract = web3.eth.contract(ionAbi);

    // instantiate by address
    const validation = ValidationContract.at('0xb9fd43a71c076f02d1dbbf473c389f0eacec559f');

    // Get a single block
    const block = web3.eth.getBlock(blockNum);

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

    // Create new signed hash
    const extraBytes = hexToBytes(extraData);
    const extraBytesShort = extraBytes.splice(1, extraBytes.length-66);
    const extraDataSignature = '0x' + bytesToHex(extraBytes.splice(extraBytes.length-65));
    const extraDataShort = '0x' + bytesToHex(extraBytesShort);

    // Make some changes to the block
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
      extraDataShort,
      mixHash,
      nonce
    ];

    // Encode and sign the new header
    const encodedHeader = '0x' + rlp.encode(header).toString('hex');
    const headerHash = Util.sha3(encodedHeader);

    const privateKey = Buffer.from('e176c157b5ae6413726c23094bb82198eb283030409624965231606ec0fbe65b', 'hex')

    const sig = Util.ecsign(headerHash, privateKey)
    if (this._chainId > 0) {
      sig.v += this._chainId * 2 + 8
    }

    const pubKey  = Util.ecrecover(headerHash, sig.v, sig.r, sig.s);
    const addrBuf = Util.pubToAddress(pubKey);

    // const sigTest = Util.toRpcSig(sig.v, sig.r, sig.s);
    const newSigBytes = Buffer.concat([sig.r, sig.s]);
    let newSig = newSigBytes.toString('hex') + '00';
    if (sig.v==27) {
      newSig = newSigBytes.toString('hex') + '00';
    } else {
      newSig = newSigBytes.toString('hex') + '01';
    }

    // Append signature to the end of extraData
    const sigBytes = hexToBytes(newSig.toString('hex'));
    const newExtraDataBytes = extraBytesShort.concat(sigBytes);
    const newExtraData = '0x' + bytesToHex(newExtraDataBytes);

    const newBlockHeader = [
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
      newExtraData,
      mixHash,
      nonce
    ];

    console.log(newBlockHeader)
    const encodedBlockHeader = '0x' + rlp.encode(newBlockHeader).toString('hex');
    const blockHeaderHash = Web3Utils.sha3(encodedBlockHeader);
    assert.equal(block.hash, blockHeaderHash);

    // The new prefixes should be calculated off chain
    const prefixHeader = '0x0214';
    const prefixExtraData = '0xa0';

    console.log("Encoded Block Header:\n", encodedBlockHeader);
    const ecrecoveryReceipt = await validation.ValidationTest(encodedBlockHeader, {from: signer, gas: 1000000});
    const latestBlockReceipt = await validation.LatestBlock();
    const latestBytesReceipt = await validation.LatestBytes();
    console.log("Expected Block Hash:\n", block.hash);
    console.log("Latest Block Hash:\n", latestBlockReceipt);
    console.log("latest Bytes:\n", latestBytesReceipt);
    // const recoveredBlockHash = ecrecoveryReceipt.logs[0].args['blockHash'];
    // const recoveredSignature = ecrecoveryReceipt.logs[1].args['owner'];
    // console.log(block.hash, recoveredBlockHash)
    // assert.equal(block.hash, recoveredBlockHash)
    // assert.equal(recoveredSignature, signer);
  })

  // Here the block header is signed off chain but by a whitelisted validator
  it('Test: Check the latest data...', async () => {
    const accounts = web3.eth.accounts;
    const signer = accounts[0];
    const ValidationContract = web3.eth.contract(ionAbi);

    // instantiate by address
    const validation = ValidationContract.at('0xb9fd43a71c076f02d1dbbf473c389f0eacec559f');

    // Get a single block
    const block = web3.eth.getBlock(blockNum);

    const latestBlockReceipt = await validation.LatestBlock();
    const latestBytesReceipt = await validation.LatestBytes();
    console.log("Expected Block Hash:\n", block.hash);
    console.log("Latest Block Hash:\n", latestBlockReceipt);
    console.log("latest Bytes:\n", latestBytesReceipt);
  })

});