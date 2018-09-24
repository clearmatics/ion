// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const eth_util = require('ethereumjs-util');
const utils = require('./helpers/utils.js');
const encoder = require('./helpers/encoder.js');
const Web3 = require('web3');
const Web3Utils = require('web3-utils');
const rlp = require('rlp');

const Validation = artifacts.require("Validation");
const Ion = artifacts.require("Ion");

const web3 = new Web3();

web3.setProvider(new web3.providers.HttpProvider('http://localhost:8501'));

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

const DEPLOYEDCHAINID = "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075"
const TESTCHAINID = "0x22b55e8a4f7c03e1689da845dd463b09299cb3a574e64c68eafc4e99077a7254"

contract('Validation.js', (accounts) => {
  const joinHex = arr => '0x' + arr.map(el => el.slice(2)).join('');

  const watchEvent = (eventObj) => new Promise((resolve,reject) => eventObj.watch((error,event) => error ? reject(error) : resolve(event)));

  // Fetch genesis from testrpc
  const genesisBlock = web3.eth.getBlock(0);
  const VALIDATORS = encoder.extractValidators(genesisBlock.extraData);
  const GENESIS_HASH = genesisBlock.hash;

  it('Deploy Contract', async () => {
    const ion = await Ion.new(DEPLOYEDCHAINID);
    const validation = await Validation.new(DEPLOYEDCHAINID, ion.address);
    let chainId = await validation.chainId();

    assert.equal(chainId, DEPLOYEDCHAINID);
  })

  it('Register Chain', async () => {
    const ion = await Ion.new(DEPLOYEDCHAINID);
    const validation = await Validation.new(DEPLOYEDCHAINID, ion.address);

    // Successfully add id of another chain
    let tx = await validation.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH);
    console.log("\tGas used to register chain = " + tx.receipt.gasUsed.toString() + " gas");
    let chain = await validation.chains(TESTCHAINID);

    assert.equal(chain, true);

    // Fail adding id of this chain
    await validation.RegisterChain(DEPLOYEDCHAINID, VALIDATORS, GENESIS_HASH).should.be.rejected;

    // Fail adding id of chain already initialised
    await validation.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH).should.be.rejected;
  })

  it('Register Chain - Check Validators', async () => {
    const ion = await Ion.new(DEPLOYEDCHAINID);
    const validation = await Validation.new(DEPLOYEDCHAINID, ion.address);

    // Successfully add id of another chain
    await validation.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH);

    let validators = await validation.m_validators.call(TESTCHAINID, VALIDATORS[0]);
    assert.equal(validators, true);
  })


  it('Register Chain - Check Genesis Hash', async () => {
    const ion = await Ion.new(DEPLOYEDCHAINID);
    const validation = await Validation.new(DEPLOYEDCHAINID, ion.address);

    // Successfully add id of another chain
    await validation.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH);

    let header = await validation.m_blockheaders.call(TESTCHAINID, GENESIS_HASH);
    let blockHeight = header[0];

    assert.equal(0, blockHeight);
  })
  
  it('Authentic Submission Happy Path - SubmitBlock()', async () => {
    const ion = await Ion.new(DEPLOYEDCHAINID);
    const validation = await Validation.new(DEPLOYEDCHAINID, ion.address);
   
    await validation.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH);

    // Fetch block 1 from testrpc
    const block = web3.eth.getBlock(1);

    const rlpHeaders = encoder.encodeBlockHeader(block);
    const signedHeaderHash = Web3Utils.sha3(rlpHeaders.signed);
    assert.equal(block.hash, signedHeaderHash);

    // Submit block should succeed
    const validationReceipt = await validation.SubmitBlock(TESTCHAINID, rlpHeaders.unsigned, rlpHeaders.signed);
    const recoveredBlockHash = await validationReceipt.logs[0].args['blockHash'];
    assert.equal(signedHeaderHash, recoveredBlockHash);
    console.log("\tGas used to submit block = " + validationReceipt.receipt.gasUsed.toString() + " gas");

    let blockHash = await validation.m_blockhashes(TESTCHAINID, block.hash);
    assert.equal(blockHash, true);

    let header = await validation.m_blockheaders(TESTCHAINID, block.hash);

    // Separate fetched header info
    parentHash = header[2];

    // Assert that block was persisted correctly
    assert.equal(parentHash, block.parentHash);
  })

  // Here the block header is signed off chain but by a a non-whitelisted validator
  it('Fail Submit Block unkown validator - SubmitBlock()', async () => {
    const ion = await Ion.new(DEPLOYEDCHAINID);
    const validation = await Validation.new(DEPLOYEDCHAINID, ion.address);

    await validation.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH);

    // Fetch block 1 from testrpc
    const block = web3.eth.getBlock(1);

    // Alter txHashin the unsigned header concatenation
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

    const encodedUnsignedHeader = '0x' + rlp.encode(unsignedHeader).toString('hex');
    const unsignedHeaderHash = Web3Utils.sha3(encodedUnsignedHeader);

    // Encode and sign the new header
    const encodedExtraData = '0x' + rlp.encode(extraDataShort).toString('hex');
    const newSignedHeaderHash = eth_util.sha3(encodedUnsignedHeader);

    const privateKey = Buffer.from('4f35bad50b8b07fff875ec9d4dec6034b1cb0f7d283db4ce7df8fcfaa2030308', 'hex')

    let signature = await signHeader(newSignedHeaderHash, privateKey, block.extraData);

    // Append signature to the end of extraData
    const sigBytes = utils.hexToBytes(signature.toString('hex'));
    const newExtraDataBytes = extraBytesShort.concat(sigBytes);
    const newExtraData = '0x' + utils.bytesToHex(newExtraDataBytes);

    const newSignedHeader = [
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
      newExtraData, // Off-chain signed block
      block.mixHash,
      block.nonce
    ];

    // Encode the offchain signed header
    const offchainSignedHeader = '0x' + rlp.encode(newSignedHeader).toString('hex');
    const offchainHeaderHash = Web3Utils.sha3(offchainSignedHeader);

    await validation.SubmitBlock(TESTCHAINID, encodedUnsignedHeader, offchainSignedHeader).should.be.rejected;

  })

  it('Submission of two sequential blocks - SubmitBlocks()', async () => {
    const ion = await Ion.new(DEPLOYEDCHAINID);
    const validation = await Validation.new(DEPLOYEDCHAINID, ion.address);
   
    await validation.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH);

    // Generate the arrays which will store the block header and indices
    let signedHeaders = [];
    let signedHeaderIndices = [];
    let unsignedHeaders = [];
    let unsignedHeaderIndices = [];

    // Fetch block 1 from testrpc
    let block = web3.eth.getBlock(1); 

    let rlpHeaders = encoder.encodeBlockHeader(block);
    let signedHeaderHash = Web3Utils.sha3(rlpHeaders.signed);
    assert.equal(block.hash, signedHeaderHash);

    // Start creating the long list of block headers
    signedHeaders.push(rlpHeaders.signed);
    signedHeaderIndices.push(utils.hexToBytes(rlpHeaders.signed).splice(1).length);
    unsignedHeaders.push(rlpHeaders.unsigned);
    unsignedHeaderIndices.push(utils.hexToBytes(rlpHeaders.unsigned).splice(1).length);

    // Fetch block 1 from testrpc
    block = web3.eth.getBlock(2);

    rlpHeaders = encoder.encodeBlockHeader(block);
    signedHeaderHash = Web3Utils.sha3(rlpHeaders.signed);
    signedHeaders.push(rlpHeaders.signed);
    assert.equal(block.hash, signedHeaderHash);

    // Start creating the long list of block headers
    signedHeaders.push(rlpHeaders.signed);
    signedHeaderIndices.push(utils.hexToBytes(rlpHeaders.signed).splice(1).length + signedHeaderIndices[signedHeaderIndices.length - 1]);
    unsignedHeaders.push(rlpHeaders.unsigned);
    unsignedHeaderIndices.push(utils.hexToBytes(rlpHeaders.unsigned).splice(1).length + unsignedHeaderIndices[unsignedHeaderIndices.length - 1]);

    // Submit block should succeed
    let joinedSignedBlocks = utils.joinHex(signedHeaders);
    let joinedUnsignedBlocks = utils.joinHex(unsignedHeaders);


    // Submit two blocks at the same time
    const validationReceipt = await validation.SubmitBlocks(
      TESTCHAINID,
      joinedUnsignedBlocks,
      joinedSignedBlocks,
      unsignedHeaderIndices,
      signedHeaderIndices,
    );
 
    const recoveredBlockHash = await validationReceipt.logs[2].args['blockHash'];
    assert.equal(signedHeaderHash, recoveredBlockHash);
    console.log("\tGas used to submit block = " + validationReceipt.receipt.gasUsed.toString() + " gas");

    let blockHash = await validation.m_blockhashes(TESTCHAINID, block.hash);
    assert.equal(blockHash, true);

    let header = await validation.m_blockheaders(TESTCHAINID, block.hash);

    // Separate fetched header info
    parentHash = header[2];

    // Assert that block was persisted correctly
    assert.equal(parentHash, block.parentHash);
  })

  it('Fail Submit Block from unknown chain - SubmitBlock()', async () => {
    const ion = await Ion.new(DEPLOYEDCHAINID);
    const validation = await Validation.new(DEPLOYEDCHAINID, ion.address);

    await validation.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH);

    // Fetch block 1 from testrpc
    const block = web3.eth.getBlock(1);

    const rlpHeaders = encoder.encodeBlockHeader(block);

    // Submit block should succeed
    await validation.SubmitBlock(TESTCHAINID.slice(0, -2) + "ff", rlpHeaders.unsigned, rlpHeaders.signed).should.be.rejected;
    
  })

  it('Fail Submit Block with wrong unsigned header - SubmitBlock()', async () => {
    const ion = await Ion.new(DEPLOYEDCHAINID);
    const validation = await Validation.new(DEPLOYEDCHAINID, ion.address);

    await validation.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH);

    // Fetch block 1 from testrpc
    const block = web3.eth.getBlock(1);

    const signedHeader = [
        block.parentHash,
// web3.setProvider(new web3.providers.HttpProvider('https://rinkeby.infura.io'));
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
        block.receiptsRoot.slice(0, -2) + "fa",
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
    await validation.SubmitBlock(TESTCHAINID, encodedUnsignedHeader, encodedSignedHeader).should.be.rejected;
    
  })
});
