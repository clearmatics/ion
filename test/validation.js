// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const eth_util = require('ethereumjs-util');
const utils = require('./helpers/utils.js');
const Web3 = require('web3');
const Web3Utils = require('web3-utils');
const Web3Abi = require('web3-eth-abi');
const Web3Accounts = require('web3-eth-accounts');
const rlp = require('rlp');

const Validation = artifacts.require("Validation");
const Recover = artifacts.require("Recover");

const web3 = new Web3();

web3.setProvider(new web3.providers.HttpProvider('http://localhost:8501'));

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

const DEPLOYEDCHAINID = "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075"
const TESTCHAINID = "0x22b55e8a4f7c03e1689da845dd463b09299cb3a574e64c68eafc4e99077a7254"
const VALIDATORS = ["0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d"];
const GENESIS_HASH = "0xaf0d377824ecc16cfdd5946ad0cd0da904cbcfff8c6cd31628c9c9e5bed2c95b";


contract('Validation.js', (accounts) => {
  const joinHex = arr => '0x' + arr.map(el => el.slice(2)).join('');

  const watchEvent = (eventObj) => new Promise((resolve,reject) => eventObj.watch((error,event) => error ? reject(error) : resolve(event)));

  const blockNum = 1;

  // Hash of the genesis block
  const genesisBlock = web3.eth.getBlock(0);
  // const genHash = genesisBlock.hash;
  const genHash = "0xaf0d377824ecc16cfdd5946ad0cd0da904cbcfff8c6cd31628c9c9e5bed2c95b";

  // Find the validator of block 1 as it is not known a priori
  const validators = ["0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d", "0x8671e5e08d74f338ee1c462340842346d797afd3"];

  it('Test: Deploy Contract', async () => {
    const validation = await Validation.new(DEPLOYEDCHAINID);
    let chainId = await validation.chainId();

    assert.equal(chainId, DEPLOYEDCHAINID);
  })

  it('Test: Initialise Chain', async () => {
    const validation = await Validation.new(DEPLOYEDCHAINID);

    // Successfully add id of another chain
    await validation.InitChain(TESTCHAINID, VALIDATORS, GENESIS_HASH);
    let chain = await validation.chains(TESTCHAINID);

    assert.equal(chain, true);

    // Fail adding id of this chain
    await validation.InitChain(DEPLOYEDCHAINID, VALIDATORS, GENESIS_HASH).should.be.rejected;

    // Fail adding id of chain already initialised
    await validation.InitChain(TESTCHAINID, VALIDATORS, GENESIS_HASH).should.be.rejected;
  })

  it('Validate Block', async () => {
    const validation = await Validation.new(DEPLOYEDCHAINID);

    await validation.InitChain(TESTCHAINID, VALIDATORS, GENESIS_HASH);

    // Fetch block 1 from testrpc
    const block = web3.eth.getBlock(1);

    const signedHeader = [
        block.parentHash,
        block.sha3Uncles,
        block.miner,
        block.stateRoot,
        block.transactionsRoot,
        block.receiptsRoot,
        block.logsBloom,
        Web3Utils.toBN(block.difficulty),
        Web3Utils.toBN(block.number),
        block.gasLimit,
        block.gasUsed,
        Web3Utils.toBN(block.timestamp),
        block.extraData,
        block.mixHash,
        block.nonce
      ];

    // Remove last 65 Bytes of extraData
    const extraBytes = utils.hexToBytes(block.extraData);
    const extraBytesShort = extraBytes.splice(1, extraBytes.length-66);
    const extraDataSignature = '0x' + utils.bytesToHex(extraBytes.splice(extraBytes.length-65));
    const extraDataShort = '0x' + utils.bytesToHex(extraBytesShort);

    const unsignedHeader = [
        block.parentHash,
        block.sha3Uncles,
        block.miner,
        block.stateRoot,
        block.transactionsRoot,
        block.receiptsRoot,
        block.logsBloom,
        Web3Utils.toBN(block.difficulty),
        Web3Utils.toBN(block.number),
        block.gasLimit,
        block.gasUsed,
        Web3Utils.toBN(block.timestamp),
        extraDataShort, // extraData minus the signature
        block.mixHash,
        block.nonce
      ];

    const encodedSignedHeader = '0x' + rlp.encode(signedHeader).toString('hex');
    const signedHeaderHash = Web3Utils.sha3(encodedSignedHeader);
    assert.equal(block.hash, signedHeaderHash);

    const encodedUnsignedHeader = '0x' + rlp.encode(unsignedHeader).toString('hex');
    const unsignedHeaderHash = Web3Utils.sha3(encodedUnsignedHeader);

    // Submit block should succeed
    const validationReceipt = await validation.ValidateBlock(TESTCHAINID, encodedUnsignedHeader, encodedSignedHeader);
    const recoveredBlockHash = await validationReceipt.logs[0].args['blockHash'];
    assert.equal(signedHeaderHash, recoveredBlockHash);

    let blockHash = await validation.m_blockhashes(TESTCHAINID, block.hash);
    assert.equal(blockHash, true);

    let header = await validation.m_blockheaders(TESTCHAINID);

    // Separate fetched header info
    parentHash = header[2];

    // Assert that block was persisted correctly
    assert.equal(parentHash, block.parentHash);
  })

  // This test takes a block and makes no changes to the block and submits it to the contract
//   it('Test: Authentic Submission Happy Path - ValidateBlock()', async () => {
//     const validation = await Validation.new(validators, genHash);
//     const accounts = web3.eth.accounts;
//     const signer = validators[0];

//     // Get a single block
//     const block = web3.eth.getBlock(blockNum);

//     // Decompose the values in the block to hash
//     const parentHash = block.parentHash;
//     const sha3Uncles = block.sha3Uncles;
//     const coinbase = block.miner;
//     const root = block.stateRoot;
//     const txHash = block.transactionsRoot;
//     const receiptHash = block.receiptsRoot;
//     const logsBloom = block.logsBloom;
//     const difficulty = Web3Utils.toBN(block.difficulty);
//     const number = Web3Utils.toBN(block.number);
//     const gasLimit = block.gasLimit;
//     const gasUsed = block.gasUsed;
//     const timestamp = Web3Utils.toBN(block.timestamp);
//     const extraData = block.extraData;
//     const mixHash = block.mixHash;
//     const nonce = block.nonce;

//     // Remove last 65 Bytes of extraData
//     const extraBytes = utils.hexToBytes(extraData);
//     const extraBytesShort = extraBytes.splice(1, extraBytes.length-66);
//     const extraDataSignature = '0x' + utils.bytesToHex(extraBytes.splice(extraBytes.length-65));
//     const extraDataShort = '0x' + utils.bytesToHex(extraBytesShort);

//     const blockHeader = [
//       parentHash,
//       sha3Uncles,
//       coinbase,
//       root,
//       txHash,
//       receiptHash,
//       logsBloom,
//       difficulty,
//       number,
//       gasLimit,
//       gasUsed,
//       timestamp,
//       extraData,
//       mixHash,
//       nonce
//     ];

//     const header = [
//       parentHash,
//       sha3Uncles,
//       coinbase,
//       root,
//       txHash,
//       receiptHash,
//       logsBloom,
//       difficulty,
//       number,
//       gasLimit,
//       gasUsed,
//       timestamp,
//       extraDataShort,
//       mixHash,
//       nonce
//     ];

//     const encodedBlockHeader = '0x' + rlp.encode(blockHeader).toString('hex');
//     const blockHeaderHash = Web3Utils.sha3(encodedBlockHeader);
//     assert.equal(block.hash, blockHeaderHash);
    
//     const encodedHeader = '0x' + rlp.encode(header).toString('hex');
//     const encodedExtraData = '0x' + rlp.encode(extraDataShort).toString('hex');
//     const headerHash = Web3Utils.sha3(encodedHeader);

//     // Get Prefixes
//     const prefixHeader = '0x' + encodedHeader.substring(4, 8);
//     const prefixExtraData = '0x' + encodedExtraData.substring(2,4);

//     const ecrecoveryReceipt = await validation.ValidateBlock(encodedBlockHeader, prefixHeader, prefixExtraData);
//     const recoveredBlockHash = ecrecoveryReceipt.logs[0].args['blockHash'];
//     const recoveredSignature = ecrecoveryReceipt.logs[1].args['owner'];
//     assert.equal(block.hash, recoveredBlockHash)
//     assert.notEqual(validators.indexOf(recoveredSignature), -1);

//   })

//   // Here the block header is signed off chain but by a whitelisted validator
//   it('Test: Authentic Submission Off-Chain Signature - ValidateBlock()', async () => {
//     const validation = await Validation.new(validators, genHash);
//     const accounts = web3.eth.accounts;
//     const signer = validators[0];

//     // Get a single block
//     const block = web3.eth.getBlock(blockNum);

//     // Decompose the values in the block to hash
//     const parentHash = block.parentHash;
//     const sha3Uncles = block.sha3Uncles;
//     const coinbase = block.miner;
//     const root = block.stateRoot;
//     const txHash = block.transactionsRoot;
//     const receiptHash = block.receiptsRoot;
//     const logsBloom = block.logsBloom;
//     const difficulty = Web3Utils.toBN(block.difficulty);
//     const number = Web3Utils.toBN(block.number);
//     const gasLimit = block.gasLimit;
//     const gasUsed = block.gasUsed;
//     const timestamp = Web3Utils.toBN(block.timestamp);
//     const extraData = block.extraData;
//     const mixHash = block.mixHash;
//     const nonce = block.nonce;

//     // Create new signed hash
//     const extraBytes = hexToBytes(extraData);
//     const extraBytesShort = extraBytes.splice(1, extraBytes.length-66);
//     const extraDataSignature = '0x' + utils.bytesToHex(extraBytes.splice(extraBytes.length-65));
//     const extraDataShort = '0x' + utils.bytesToHex(extraBytesShort);

//     // Make some changes to the block
//     const header = [
//       parentHash,
//       sha3Uncles,
//       coinbase,
//       root,
//       txHash,
//       receiptHash,
//       logsBloom,
//       difficulty,
//       number,
//       gasLimit,
//       gasUsed,
//       timestamp,
//       extraDataShort,
//       mixHash,
//       nonce
//     ];

//     // Encode and sign the new header
//     const encodedHeader = '0x' + rlp.encode(header).toString('hex');
//     const encodedExtraData = '0x' + rlp.encode(extraDataShort).toString('hex');
//     const headerHash = eth_util.sha3(encodedHeader);

//     const privateKey = Buffer.from('e176c157b5ae6413726c23094bb82198eb283030409624965231606ec0fbe65b', 'hex')

//     let signature = await signHeader(headerHash, privateKey, extraData);

//     // Append signature to the end of extraData
//     const sigBytes = hexToBytes(signature.toString('hex'));
//     const newExtraDataBytes = extraBytesShort.concat(sigBytes);
//     const newExtraData = '0x' + utils.bytesToHex(newExtraDataBytes);
//     assert.equal(extraDataSignature, '0x'+signature.toString('hex'))

//     const newBlockHeader = [
//       parentHash,
//       sha3Uncles,
//       coinbase,
//       root,
//       txHash,
//       receiptHash,
//       logsBloom,
//       difficulty,
//       number,
//       gasLimit,
//       gasUsed,
//       timestamp,
//       newExtraData,
//       mixHash,
//       nonce
//     ];

//     const encodedBlockHeader = '0x' + rlp.encode(newBlockHeader).toString('hex');
//     const blockHeaderHash = Web3Utils.sha3(encodedBlockHeader);
//     assert.equal(block.hash, blockHeaderHash);

//     // Get Prefixes
//     const prefixHeader = '0x' + encodedHeader.substring(4, 8);
//     const prefixExtraData = '0x' + encodedExtraData.substring(2,4);

//     const ecrecoveryReceipt = await validation.ValidateBlock(encodedBlockHeader, prefixHeader, prefixExtraData);
//     const recoveredBlockHash = ecrecoveryReceipt.logs[0].args['blockHash'];
//     const recoveredSignature = ecrecoveryReceipt.logs[1].args['owner'];
//     assert.equal(block.hash, recoveredBlockHash)
//     assert.equal(recoveredSignature, signer);
//   })

//   it('Test: Inauthentic Block Submission - ValidateBlock()', async () => {
//     const validation = await Validation.new(validators, genHash);
//     const accounts = web3.eth.accounts;
//     const signer = validators[0];

//     // Get a single block
//     const block = web3.eth.getBlock(blockNum);

//     // Decompose the values in the block to hash
//     const parentHash = block.parentHash;
//     const sha3Uncles = block.sha3Uncles;
//     const coinbase = block.miner;
//     const root = block.stateRoot;
//     const txHash = block.transactionsRoot;
//     const receiptHash = block.receiptsRoot;
//     const logsBloom = block.logsBloom;
//     const difficulty = Web3Utils.toBN(block.difficulty);
//     const number = Web3Utils.toBN(block.number);
//     const gasLimit = block.gasLimit;
//     const gasUsed = block.gasUsed;
//     const timestamp = Web3Utils.toBN(block.timestamp);
//     const extraData = block.extraData;
//     const mixHash = block.mixHash;
//     const nonce = block.nonce;

//     // Create new signed hash
//     const extraBytes = hexToBytes(extraData);
//     const extraBytesShort = extraBytes.splice(1, extraBytes.length-66);
//     const extraDataSignature = '0x' + utils.bytesToHex(extraBytes.splice(extraBytes.length-65));
//     const extraDataShort = '0x' + utils.bytesToHex(extraBytesShort);

//     // Make some changes to the block
//     const newTxHash = Web3Utils.sha3("Test Data");
//     // console.log(txHash, newTxHash);
//     const header = [
//       parentHash,
//       sha3Uncles,
//       coinbase,
//       root,
//       // txHash,
//       newTxHash,
//       receiptHash,
//       logsBloom,
//       difficulty,
//       number,
//       gasLimit,
//       gasUsed,
//       timestamp,
//       extraDataShort,
//       mixHash,
//       nonce
//     ];

//     // Encode and sign the new header
//     const encodedHeader = '0x' + rlp.encode(header).toString('hex');
//     const headerHash = eth_util.sha3(encodedHeader);

//     const privateKey = Buffer.from('e176c157b5ae6413726c23094bb82198eb283030409624965231606ec0fbe65b', 'hex')

//     const sig = eth_util.ecsign(headerHash, privateKey)
//     if (this._chainId > 0) {
//       sig.v += this._chainId * 2 + 8
//     }

//     const pubKey  = eth_util.ecrecover(headerHash, sig.v, sig.r, sig.s);
//     const addrBuf = eth_util.pubToAddress(pubKey);
//     assert.equal(signer, '0x'+addrBuf.toString('hex'));

//     const newSigBytes = Buffer.concat([sig.r, sig.s]);
//     let newSig;

//     // Need to understand why but signature requires different v than in others to recover correctly
//     if (sig.v=="27") {
//       newSig = newSigBytes.toString('hex') + '00';
//     }
//     if (sig.v=="28") {
//       newSig = newSigBytes.toString('hex') + '01';
//     }

//     // Append signature to the end of extraData
//     const sigBytes = hexToBytes(newSig.toString('hex'));
//     const newExtraDataBytes = extraBytesShort.concat(sigBytes);
//     const newExtraData = '0x' + utils.bytesToHex(newExtraDataBytes);

//     const newBlockHeader = [
//       parentHash,
//       sha3Uncles,
//       coinbase,
//       root,
//       // txHash,
//       newTxHash,
//       receiptHash,
//       logsBloom,
//       difficulty,
//       number,
//       gasLimit,
//       gasUsed,
//       timestamp,
//       newExtraData,
//       mixHash,
//       nonce
//     ];

//     const encodedBlockHeader = '0x' + rlp.encode(newBlockHeader).toString('hex');
//     const encodedExtraData = '0x' + rlp.encode(extraDataShort).toString('hex');
//     const blockHeaderHash = Web3Utils.sha3(encodedBlockHeader);

//     // Get Prefixes
//     const prefixHeader = '0x' + encodedHeader.substring(4, 8);
//     const prefixExtraData = '0x' + encodedExtraData.substring(2,4);

//     const ecrecoveryReceipt = await validation.ValidateBlock(encodedBlockHeader, prefixHeader, prefixExtraData);
//     const recoveredBlockHash = ecrecoveryReceipt.logs[0].args['blockHash'];
//     const recoveredSignature = ecrecoveryReceipt.logs[1].args['owner'];
//     assert.equal(recoveredBlockHash, blockHeaderHash);
//     assert.equal(recoveredSignature, signer);

//   })

//   it('Test: Authentic Block Unkown Validator Submission - ValidateBlock()', async () => {
//     const validation = await Validation.new(validators, genHash);
//     const accounts = web3.eth.accounts;
//     const signer = validators[1];

//     // Get a single block
//     const block = web3.eth.getBlock(blockNum);

//     // Decompose the values in the block to hash
//     const parentHash = block.parentHash;
//     const sha3Uncles = block.sha3Uncles;
//     const coinbase = block.miner;
//     const root = block.stateRoot;
//     const txHash = block.transactionsRoot;
//     const receiptHash = block.receiptsRoot;
//     const logsBloom = block.logsBloom;
//     const difficulty = Web3Utils.toBN(block.difficulty);
//     const number = Web3Utils.toBN(block.number);
//     const gasLimit = block.gasLimit;
//     const gasUsed = block.gasUsed;
//     const timestamp = Web3Utils.toBN(block.timestamp);
//     const extraData = block.extraData;
//     const mixHash = block.mixHash;
//     const nonce = block.nonce;

//     // Create new signed hash
//     const extraBytes = hexToBytes(extraData);
//     const extraBytesShort = extraBytes.splice(1, extraBytes.length-66);
//     const extraDataSignature = '0x' + utils.bytesToHex(extraBytes.splice(extraBytes.length-65));
//     const extraDataShort = '0x' + utils.bytesToHex(extraBytesShort);

//     // Make some changes to the block
//     const newTxHash = Web3Utils.sha3("Test Data");
//     const header = [
//       parentHash,
//       sha3Uncles,
//       coinbase,
//       root,
//       newTxHash,
//       receiptHash,
//       logsBloom,
//       difficulty,
//       number,
//       gasLimit,
//       gasUsed,
//       timestamp,
//       extraDataShort,
//       mixHash,
//       nonce
//     ];

//     // Encode and sign the new header
//     const encodedHeader = '0x' + rlp.encode(header).toString('hex');
//     const headerHash = eth_util.sha3(encodedHeader);

//     const privateKey = Buffer.from('4f35bad50b8b07fff875ec9d4dec6034b1cb0f7d283db4ce7df8fcfaa2030308', 'hex')

//     const sig = eth_util.ecsign(headerHash, privateKey)
//     if (this._chainId > 0) {
//       sig.v += this._chainId * 2 + 8
//     }

//     const pubKey  = eth_util.ecrecover(headerHash, sig.v, sig.r, sig.s);
//     const addrBuf = eth_util.pubToAddress(pubKey);
//     const addr    = eth_util.bufferToHex(addrBuf);

//     const newSigBytes = Buffer.concat([sig.r, sig.s]);
//     let newSig;

//     const bytes = hexToBytes(extraData)
//     const finalByte = bytes.splice(bytes.length-1)
//     if (finalByte.toString('hex')=="00")
//       newSig = newSigBytes.toString('hex') + '00';
//     else (finalByte.toString('hex')=="01")
//       newSig = newSigBytes.toString('hex') + '01';

//     // Append signature to the end of extraData
//     const sigBytes = hexToBytes(newSig.toString('hex'));
//     const newExtraDataBytes = extraBytesShort.concat(sigBytes);
//     const newExtraData = '0x' + utils.bytesToHex(newExtraDataBytes);

//     const newBlockHeader = [
//       parentHash,
//       sha3Uncles,
//       coinbase,
//       root,
//       newTxHash,
//       receiptHash,
//       logsBloom,
//       difficulty,
//       number,
//       gasLimit,
//       gasUsed,
//       timestamp,
//       newExtraData,
//       mixHash,
//       nonce
//     ];

//     const encodedBlockHeader = '0x' + rlp.encode(newBlockHeader).toString('hex');
//     const encodedExtraData = '0x' + rlp.encode(extraDataShort).toString('hex');
//     const blockHeaderHash = Web3Utils.sha3(encodedBlockHeader);

//     // Get Prefixes
//     const prefixHeader = '0x' + encodedHeader.substring(4, 8);
//     const prefixExtraData = '0x' + encodedExtraData.substring(2,4);

//     try {
//           const ecrecoveryReceipt = await validation.ValidateBlock(encodedBlockHeader, prefixHeader, prefixExtraData);
//         } catch (err) {
//           assert.isDefined(err, "transaction should have thrown");
//     }
//   })

});
