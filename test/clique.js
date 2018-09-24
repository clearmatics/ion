// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const eth_util = require('ethereumjs-util');
const utils = require('./helpers/utils.js');
const encoder = require('./helpers/encoder.js');
const Web3 = require('web3');

const Validation = artifacts.require("Validation");
const Ion = artifacts.require("Ion");

const web3 = new Web3();
const rinkeby = new Web3();

web3.setProvider(new web3.providers.HttpProvider('http://localhost:8501'));
rinkeby.setProvider(new web3.providers.HttpProvider('https://rinkeby.infura.io'));

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
const VALIDATORS_START = ["0x42eb768f2244c8811c63729a21a3569731535f06", "0x7ffc57839b00206d1ad20c69a1981b489f772031", "0xb279182d99e65703f0076e4812653aab85fca0f0"];
const VALIDATORS_FINISH = ["0x42eb768f2244c8811c63729a21a3569731535f06", "0x6635f83421bf059cd8111f180f0727128685bae4", "0x7ffc57839b00206d1ad20c69a1981b489f772031", "0xb279182d99e65703f0076e4812653aab85fca0f0"];
const GENESIS_HASH = "0xf32b505a5ad95dfa88c2bd6904a1ba81a92a1db547dc17f4d7c0f64cf2cddbb1";

require('chai')
 .use(require('chai-as-promised'))
 .should();

contract('Clique.js', (accounts) => {
  const joinHex = arr => '0x' + arr.map(el => el.slice(2)).join('');

  const watchEvent = (eventObj) => new Promise((resolve,reject) => eventObj.watch((error,event) => error ? reject(error) : resolve(event)));

  // This test checks that new validators get added into the validator list as blocks are submitted to the contract.
  // Rinkeby adds its first non-genesis validator at block 873987 with the votes occuring at blocks 873983 and 873986
  // we will start following the chain from 873982 and then add blocks until the vote threshold, n/2 + 1, is passed.
  it('Add Validators Through Individual Block Submission', async () => {
    const ion = await Ion.new(DEPLOYEDCHAINID);
    const validation = await Validation.new(DEPLOYEDCHAINID, ion.address);

    await validation.RegisterChain(TESTCHAINID, VALIDATORS_START, GENESIS_HASH);

    let voteThreshold = await validation.m_threshold(TESTCHAINID);
    assert.equal(voteThreshold, 2);

    let voteProposal = await validation.m_proposals(TESTCHAINID, VALIDATORS_FINISH[1]);
    assert.equal(voteProposal, 0);

    // Fetch block 873982 from rinkeby
    let block = rinkeby.eth.getBlock(873982);
    let rlpHeaders = encoder.encodeBlockHeader(block);

    // Submit block should succeed
    let validationReceipt = await validation.SubmitBlock(TESTCHAINID, rlpHeaders.unsigned, rlpHeaders.signed);
    console.log("\tGas used to submit block 873982 = " + validationReceipt.receipt.gasUsed.toString() + " gas");

    // Fetch block 873983 from rinkeby
    block = rinkeby.eth.getBlock(873983);
    rlpHeaders = encoder.encodeBlockHeader(block);

    // Submit block should succeed
    validationReceipt = await validation.SubmitBlock(TESTCHAINID, rlpHeaders.unsigned, rlpHeaders.signed);
    console.log("\tGas used to submit block 873983 = " + validationReceipt.receipt.gasUsed.toString() + " gas");

    // Check proposal is added
    voteProposal = await validation.m_proposals(TESTCHAINID, VALIDATORS_FINISH[1]);
    assert.equal(voteProposal, 1);


    // Fetch block 873984 from rinkeby
    block = rinkeby.eth.getBlock(873984);
    rlpHeaders = encoder.encodeBlockHeader(block);

    // Submit block should succeed
    validationReceipt = await validation.SubmitBlock(TESTCHAINID, rlpHeaders.unsigned, rlpHeaders.signed);
    console.log("\tGas used to submit block 873984 = " + validationReceipt.receipt.gasUsed.toString() + " gas");

    // Fetch block 873985 from rinkeby
    block = rinkeby.eth.getBlock(873985);
    rlpHeaders = encoder.encodeBlockHeader(block);

    // Submit block should succeed
    validationReceipt = await validation.SubmitBlock(TESTCHAINID, rlpHeaders.unsigned, rlpHeaders.signed);
    console.log("\tGas used to submit block 873985 = " + validationReceipt.receipt.gasUsed.toString() + " gas");

    // Fetch block 873986 from rinkeby
    block = rinkeby.eth.getBlock(873986);
    rlpHeaders = encoder.encodeBlockHeader(block);

    // Submit block should succeed
    validationReceipt = await validation.SubmitBlock(TESTCHAINID, rlpHeaders.unsigned, rlpHeaders.signed);
    console.log("\tGas used to submit block 873986 = " + validationReceipt.receipt.gasUsed.toString() + " gas");
    
    // Check proposal is added
    voteProposal = await validation.m_proposals(TESTCHAINID, VALIDATORS_FINISH[1]);
    assert.equal(voteProposal, 0);

    // Check all new validators are added
    let validators = await validation.m_validators(TESTCHAINID, VALIDATORS_FINISH[0]);
    assert.equal(validators, true);
    validators = await validation.m_validators(TESTCHAINID, VALIDATORS_FINISH[1]);
    assert.equal(validators, true);
    validators = await validation.m_validators(TESTCHAINID, VALIDATORS_FINISH[2]);
    assert.equal(validators, true);
    validators = await validation.m_validators(TESTCHAINID, VALIDATORS_FINISH[3]);
    assert.equal(validators, true);
  })

  // This test checks that new validators get added into the validator list as blocks are submitted to the contract.
  // Rinkeby adds its first non-genesis validator at block 873987 with the votes occuring at blocks 873983 and 873986
  // we will start following the chain from 873982 and then add blocks until the vote threshold, n/2 + 1, is passed.
  it('Add Validators Through Simultaneous Block Submission', async () => {
    const ion = await Ion.new(DEPLOYEDCHAINID);
    const validation = await Validation.new(DEPLOYEDCHAINID, ion.address);

    // Generate the arrays which will store the block header and indices
    let signedHeaders = [];
    let signedHeaderIndices = [];
    let unsignedHeaders = [];
    let unsignedHeaderIndices = [];

    await validation.RegisterChain(TESTCHAINID, VALIDATORS_START, GENESIS_HASH);

    let voteThreshold = await validation.m_threshold(TESTCHAINID);
    assert.equal(voteThreshold, 2);

    let voteProposal = await validation.m_proposals(TESTCHAINID, VALIDATORS_FINISH[1]);
    assert.equal(voteProposal, 0);

    // Fetch block 873982 from rinkeby
    let block = rinkeby.eth.getBlock(873982);
    let rlpHeaders = encoder.encodeBlockHeader(block);
    encoder.appendBlockHeaders(
      signedHeaders,
      signedHeaderIndices,
      unsignedHeaders,
      unsignedHeaderIndices,
      rlpHeaders,
    );

    // Fetch block 873983 from rinkeby
    block = rinkeby.eth.getBlock(873983);
    rlpHeaders = encoder.encodeBlockHeader(block);
    encoder.appendBlockHeaders(
      signedHeaders,
      signedHeaderIndices,
      unsignedHeaders,
      unsignedHeaderIndices,
      rlpHeaders,
    );

    // Fetch block 873984 from rinkeby
    block = rinkeby.eth.getBlock(873984);
    rlpHeaders = encoder.encodeBlockHeader(block);
    encoder.appendBlockHeaders(
      signedHeaders,
      signedHeaderIndices,
      unsignedHeaders,
      unsignedHeaderIndices,
      rlpHeaders,
    );

    // Fetch block 873985 from rinkeby
    block = rinkeby.eth.getBlock(873985);
    rlpHeaders = encoder.encodeBlockHeader(block);
    encoder.appendBlockHeaders(
      signedHeaders,
      signedHeaderIndices,
      unsignedHeaders,
      unsignedHeaderIndices,
      rlpHeaders,
    );

    // Fetch block 873986 from rinkeby
    block = rinkeby.eth.getBlock(873986);
    rlpHeaders = encoder.encodeBlockHeader(block);
    encoder.appendBlockHeaders(
      signedHeaders,
      signedHeaderIndices,
      unsignedHeaders,
      unsignedHeaderIndices,
      rlpHeaders,
    );

    // Submit block should succeed
    let joinedSignedBlocks = utils.joinHex(signedHeaders);
    let joinedUnsignedBlocks = utils.joinHex(unsignedHeaders);

    // Submit multiple blocks at the same time
    const validationReceipt = await validation.SubmitBlocks(
      TESTCHAINID,
      joinedUnsignedBlocks,
      joinedSignedBlocks,
      unsignedHeaderIndices,
      signedHeaderIndices,
    );
    console.log("\tGas used to simultaneously submit blocks (873982-873986) = " + validationReceipt.receipt.gasUsed.toString() + " gas");

    // Check all new validators are added
    let validators = await validation.m_validators(TESTCHAINID, VALIDATORS_FINISH[0]);
    assert.equal(validators, true);
    validators = await validation.m_validators(TESTCHAINID, VALIDATORS_FINISH[1]);
    assert.equal(validators, true);
    validators = await validation.m_validators(TESTCHAINID, VALIDATORS_FINISH[2]);
    assert.equal(validators, true);
    validators = await validation.m_validators(TESTCHAINID, VALIDATORS_FINISH[3]);
    assert.equal(validators, true);
  })

});
