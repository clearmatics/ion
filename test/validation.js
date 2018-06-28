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

contract.only('validation.js', (accounts) => {
  const joinHex = arr => '0x' + arr.map(el => el.slice(2)).join('');

  const watchEvent = (eventObj) => new Promise((resolve,reject) => eventObj.watch((error,event) => error ? reject(error) : resolve(event)));

  const blockNum = 1;

  const validators = ["0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d", "0x8671e5e08d74f338ee1c462340842346d797afd3"];
  const genHash = "0xc3bac257bbd04893316a76d41b6ff70de5f65c9f24db128864a6322d8e0e2f28";

  it('Test: GetValidators()', async () => {
    const validation = await Validation.new(validators, genHash);
    const accounts = web3.eth.accounts;
    const signer = validators[1];

    const validatorsReceipt = await validation.GetValidators();
    assert.equal(validators[0], validatorsReceipt[0])
  })

  // This test takes a block and makes no changes to the block and submits it to the contract
  it('Test: Authentic Submission Happy Path - ValidateBlock()', async () => {
    const validation = await Validation.new(validators, genHash);
    const accounts = web3.eth.accounts;
    const signer = validators[1];

    // Get a single block
    const block = web3.eth.getBlock(1);

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

    const blockHeader = [
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

    const encodedBlockHeader = '0x' + rlp.encode(blockHeader).toString('hex');
    const blockHeaderHash = Web3Utils.sha3(encodedBlockHeader);
    assert.equal(block.hash, blockHeaderHash);

    const encodedHeader = '0x' + rlp.encode(header).toString('hex');
    const headerHash = Web3Utils.sha3(encodedHeader);

    // The new prefixes should be calculated off chain
    const prefixHeader = '0x0214';
    const prefixExtraData = '0xa0';

    const ecrecoveryReceipt = await validation.ValidateBlock(encodedBlockHeader, prefixHeader, prefixExtraData);
    const recoveredBlockHash = ecrecoveryReceipt.logs[0].args['blockHash'];
    const recoveredSignature = ecrecoveryReceipt.logs[1].args['owner'];
    assert.equal(block.hash, recoveredBlockHash)
    assert.notEqual(validators.indexOf(recoveredSignature), -1);

  })

  // Here the block header is signed off chain but by a whitelisted validator
  it('Test: Authentic Submission Off-Chain Signature - ValidateBlock()', async () => {
    const validation = await Validation.new(validators, genHash);
    const accounts = web3.eth.accounts;
    const signer = validators[1];

    // Get a single block
    const block = web3.eth.getBlock(1);

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

    const privateKey = Buffer.from('d18bc3878eb28192238d92ae085cdb9438527e36faa92484dea2e3baa047b083', 'hex')

    const sig = Util.ecsign(headerHash, privateKey)
    if (this._chainId > 0) {
      sig.v += this._chainId * 2 + 8
    }

    const pubKey  = Util.ecrecover(headerHash, sig.v, sig.r, sig.s);
    const addrBuf = Util.pubToAddress(pubKey);

    const newSigBytes = Buffer.concat([sig.r, sig.s]);
    let newSig;
    if (sig.v==28)
      newSig = newSigBytes.toString('hex') + '00';
    else (sig.v==27)
      newSig = newSigBytes.toString('hex') + '01';

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

    const encodedBlockHeader = '0x' + rlp.encode(newBlockHeader).toString('hex');
    const blockHeaderHash = Web3Utils.sha3(encodedBlockHeader);
    assert.equal(block.hash, blockHeaderHash);

    // The new prefixes should be calculated off chain
    const prefixHeader = '0x0214';
    const prefixExtraData = '0xa0';

    const ecrecoveryReceipt = await validation.ValidateBlock(encodedBlockHeader, prefixHeader, prefixExtraData);
    const recoveredBlockHash = ecrecoveryReceipt.logs[0].args['blockHash'];
    const recoveredSignature = ecrecoveryReceipt.logs[1].args['owner'];
    assert.equal(block.hash, recoveredBlockHash)
    assert.equal(recoveredSignature, signer);
  })

  it('Test: Inauthentic Block Submission - ValidateBlock()', async () => {
    const validation = await Validation.new(validators, genHash);
    const accounts = web3.eth.accounts;
    const signer = validators[1];

    // Get a single block
    const block = web3.eth.getBlock(1);

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
    const newTxHash = Web3Utils.sha3("Test Data");
    const header = [
      parentHash,
      sha3Uncles,
      coinbase,
      root,
      newTxHash,
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

    const privateKey = Buffer.from('d18bc3878eb28192238d92ae085cdb9438527e36faa92484dea2e3baa047b083', 'hex')

    const sig = Util.ecsign(headerHash, privateKey)
    if (this._chainId > 0) {
      sig.v += this._chainId * 2 + 8
    }

    const pubKey  = Util.ecrecover(headerHash, sig.v, sig.r, sig.s);
    const addrBuf = Util.pubToAddress(pubKey);

    const vPseudo = new Buffer(0);
    const newSigBytes = Buffer.concat([sig.r, sig.s]);
    let newSig;
    if (sig.v==28)
      newSig = newSigBytes.toString('hex') + '00';
    else (sig.v==27)
      newSig = newSigBytes.toString('hex') + '01';

    // Append signature to the end of extraData
    const sigBytes = hexToBytes(newSig.toString('hex'));
    const newExtraDataBytes = extraBytesShort.concat(sigBytes);
    const newExtraData = '0x' + bytesToHex(newExtraDataBytes);

    const newBlockHeader = [
      parentHash,
      sha3Uncles,
      coinbase,
      root,
      newTxHash,
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

    const encodedBlockHeader = '0x' + rlp.encode(newBlockHeader).toString('hex');
    const blockHeaderHash = Web3Utils.sha3(encodedBlockHeader);

    // The new prefixes should be calculated off chain
    const prefixHeader = '0x0214';
    const prefixExtraData = '0xa0';

    const ecrecoveryReceipt = await validation.ValidateBlock(encodedBlockHeader, prefixHeader, prefixExtraData);
    const recoveredBlockHash = ecrecoveryReceipt.logs[0].args['blockHash'];
    const recoveredSignature = ecrecoveryReceipt.logs[1].args['owner'];
    assert.equal(recoveredSignature, signer);
  })

  it('Test: Authentic Block Unkown Validator Submission - ValidateBlock()', async () => {
    const validation = await Validation.new(validators, genHash);
    const accounts = web3.eth.accounts;
    const signer = validators[1];

    // Get a single block
    const block = web3.eth.getBlock(1);

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
    const newTxHash = Web3Utils.sha3("Test Data");
    const header = [
      parentHash,
      sha3Uncles,
      coinbase,
      root,
      newTxHash,
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

    const privateKey = Buffer.from('4f35bad50b8b07fff875ec9d4dec6034b1cb0f7d283db4ce7df8fcfaa2030308', 'hex')

    const sig = Util.ecsign(headerHash, privateKey)
    if (this._chainId > 0) {
      sig.v += this._chainId * 2 + 8
    }

    const pubKey  = Util.ecrecover(headerHash, sig.v, sig.r, sig.s);
    const addrBuf = Util.pubToAddress(pubKey);
    const addr    = Util.bufferToHex(addrBuf);

    const vPseudo = new Buffer(0);
    const newSigBytes = Buffer.concat([sig.r, sig.s]);
    const newSig = newSigBytes.toString('hex') + '00';

    // Append signature to the end of extraData
    const sigBytes = hexToBytes(newSig.toString('hex'));
    const newExtraDataBytes = extraBytesShort.concat(sigBytes);
    const newExtraData = '0x' + bytesToHex(newExtraDataBytes);

    const newBlockHeader = [
      parentHash,
      sha3Uncles,
      coinbase,
      root,
      newTxHash,
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

    const encodedBlockHeader = '0x' + rlp.encode(newBlockHeader).toString('hex');
    const blockHeaderHash = Web3Utils.sha3(encodedBlockHeader);

    // The new prefixes should be calculated off chain
    const prefixHeader = '0x0214';
    const prefixExtraData = '0xa0';

    try {
          const ecrecoveryReceipt = await validation.ValidateBlock(encodedBlockHeader, prefixHeader, prefixExtraData);
        } catch (err) {
          assert.isDefined(err, "transaction should have thrown");
    }
  })

});
