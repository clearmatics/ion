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
const rinkeby = new Web3();

web3.setProvider(new web3.providers.HttpProvider('http://localhost:8545'));
rinkeby.setProvider(new web3.providers.HttpProvider('https://rinkeby.infura.io'));

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

const VALIDATORS_START = ["0x42eb768f2244c8811c63729a21a3569731535f06", "0x7ffc57839b00206d1ad20c69a1981b489f772031", "0xb279182d99e65703f0076e4812653aab85fca0f0"];
const VALIDATORS_FINISH = ["0x42eb768f2244c8811c63729a21a3569731535f06", "0x6635f83421bf059cd8111f180f0727128685bae4", "0x7ffc57839b00206d1ad20c69a1981b489f772031", "0xb279182d99e65703f0076e4812653aab85fca0f0"];
const GENESIS_HASH = "0xf32b505a5ad95dfa88c2bd6904a1ba81a92a1db547dc17f4d7c0f64cf2cddbb1";
const ADD_VALIDATORS_GENESIS_HASH = "0xf32b505a5ad95dfa88c2bd6904a1ba81a92a1db547dc17f4d7c0f64cf2cddbb1";


contract('Clique.js', (accounts) => {
  const joinHex = arr => '0x' + arr.map(el => el.slice(2)).join('');

  const watchEvent = (eventObj) => new Promise((resolve,reject) => eventObj.watch((error,event) => error ? reject(error) : resolve(event)));

  // Fetch genesis from testrpc
  const genesisBlock = rinkeby.eth.getBlock(0);
  const VALIDATORS = encoder.extractValidators(genesisBlock.extraData);
  const GENESIS_HASH = genesisBlock.hash;
  
  let ion;
  let clique;
  let storage;
    
  beforeEach('setup contract for each test', async function () {
    ion = await MockIon.new(DEPLOYEDCHAINID);
    clique = await Clique.new(ion.address);
    storage = await MockStorage.new(ion.address);
  })

  it('Deploy Contract', async () => {
    let chainId = await ion.chainId();

    assert.equal(chainId, DEPLOYEDCHAINID);
  })

  describe('Register Chain', () => {
      it('Successful Register Chain', async () => {
        // Successfully add id of another chain
        let tx = await clique.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH, storage.address);
        console.log("\tGas used to register chain = " + tx.receipt.gasUsed.toString() + " gas");
        let chainExists = await clique.chains(TESTCHAINID);

        assert(chainExists);

        // Fail adding id of this chain
        await clique.RegisterChain(storage.address, DEPLOYEDCHAINID, VALIDATORS, GENESIS_HASH).should.be.rejected;

        // Fail adding id of chain already initialised
        await clique.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH, storage.address).should.be.rejected;
      })

      it('Check Validators', async () => {
        // Successfully add id of another chain
        await clique.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH, storage.address);

        let registeredValidators = await clique.getValidators.call(TESTCHAINID, GENESIS_HASH);

        for (let i = 0; i < VALIDATORS.length; i++) {
            let validatorExists = registeredValidators.some(v => { return v == VALIDATORS[i] });;
            assert(validatorExists);
        }
      })

      it('Check Genesis Hash', async () => {
        // Successfully add id of another chain
        await clique.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH, storage.address);

        let header = await clique.m_blockheaders(TESTCHAINID, GENESIS_HASH);
        let blockHeight = header[0];

        assert.equal(0, blockHeight);
      })
  })

  describe('Submit Block', () => {
      it('Authentic Submission Happy Path', async () => {
        await clique.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH, storage.address);

        // Fetch block 1 from rinkeby
        const block = rinkeby.eth.getBlock(1);

        const rlpHeaders = encoder.encodeBlockHeader(block);
        const signedHeaderHash = Web3Utils.sha3(rlpHeaders.signed);
        assert.equal(block.hash, signedHeaderHash);

        // Submit block should succeed
        const validationReceipt = await clique.SubmitBlock(TESTCHAINID, rlpHeaders.unsigned, rlpHeaders.signed, storage.address);
        let event = validationReceipt.receipt.logs.some(l => { return l.topics[0] == '0x' + sha3("AddedBlock(bytes32)") });
        assert.ok(event, "Stored event not emitted");

        const submittedEvent = validationReceipt.logs.find(l => { return l.event == 'BlockSubmitted' });
        assert.equal(signedHeaderHash, submittedEvent.args.blockHash);

        let blockHashExists = await clique.m_blockhashes(TESTCHAINID, block.hash);
        assert(blockHashExists);

        let header = await clique.m_blockheaders(TESTCHAINID, block.hash);

        // Separate fetched header info
        parentHash = header[2];

        // Assert that block was persisted correctly
        assert.equal(parentHash, block.parentHash);
      })

      // Here the block header is signed off chain but by a a non-whitelisted validator
      it('Fail Submit Block unkown validator - SubmitBlock()', async () => {
        // Successfully add id of another chain
        await clique.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH, storage.address);

        // Fetch block 1 from rinkeby
        const block = rinkeby.eth.getBlock(1);

        // Alter txHashin the unsigned header concatenation
        const rlpHeaders = encoder.encodeBlockHeader(block);
        const signedHeader = rlpHeaders.signed;
        const unsignedHeader = rlpHeaders.unsigned;

        // Remove last 65 Bytes of extraData
        const extraBytesShort = rlpHeaders.extraBytesShort;
        const extraDataSignature = rlpHeaders.extraDataSignature;
        const extraDataShort = rlpHeaders.extraDataShort;

        const signedHeaderHash = Web3Utils.sha3(signedHeader);
        const unsignedHeaderHash = Web3Utils.sha3(unsignedHeader);

        // Encode and sign the new header
        const encodedExtraData = '0x' + rlp.encode(extraDataShort).toString('hex');
        const newSignedHeaderHash = eth_util.sha3(unsignedHeader);

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

        await clique.SubmitBlock(TESTCHAINID, unsignedHeader, offchainSignedHeader, storage.address).should.be.rejected;

      })

      it('Fail Submit Block from unknown chain - SubmitBlock()', async () => {
        await clique.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH, storage.address);

        // Fetch block 1 from testrpc
        const block = web3.eth.getBlock(1);

        const rlpHeaders = encoder.encodeBlockHeader(block);

        // Submit block should fail
        await clique.SubmitBlock(TESTCHAINID.slice(0, -2) + "ff", rlpHeaders.unsigned, rlpHeaders.signed, storage.address).should.be.rejected;

      })

      it('Fail Submit Block with wrong unsigned header - SubmitBlock()', async () => {
        await clique.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH, storage.address);

        // Fetch block 1 from testrpc
        const block = rinkeby.eth.getBlock(1);

        const rlpHeaders = encoder.encodeBlockHeader(block);

        const signedHeaderHash = Web3Utils.sha3(rlpHeaders.signed);
        assert.equal(block.hash, signedHeaderHash);

        let unsignedHeader = rlpHeaders.rawunsigned;
        unsignedHeader[5] = unsignedHeader[5].slice(0, -2) + "fa";
        const encodedUnsignedHeader = '0x' + rlp.encode(unsignedHeader).toString('hex');
        const unsignedHeaderHash = Web3Utils.sha3(rlpHeaders.unsigned);

        // Submit block should fail
        await clique.SubmitBlock(TESTCHAINID, encodedUnsignedHeader, rlpHeaders.signed, storage.address).should.be.rejected;

      })


    // This test checks that new validators get added into the validator list as blocks are submitted to the contract.
    // Rinkeby adds its first non-genesis validator at block 873987 with the votes occuring at blocks 873983 and 873986
    // we will start following the chain from 873982 and then add blocks until the vote threshold, n/2 + 1, is passed.
    it('Add Validators Through Block Submission', async () => {
      await clique.RegisterChain(TESTCHAINID, VALIDATORS_START, ADD_VALIDATORS_GENESIS_HASH, storage.address);

      let registeredValidators = await clique.getValidators.call(TESTCHAINID, ADD_VALIDATORS_GENESIS_HASH);
      let voteThreshold = Math.floor((registeredValidators.length/2) + 1);
      assert.equal(voteThreshold, 2);

      let voteProposal = await clique.getProposal.call(TESTCHAINID, ADD_VALIDATORS_GENESIS_HASH, VALIDATORS_FINISH[1]);
      assert.equal(voteProposal, 0);

      // Fetch block 873982 from rinkeby
      let block = rinkeby.eth.getBlock(873982);
      let rlpHeaders = encoder.encodeBlockHeader(block);

      // Submit block should succeed
      let validationReceipt = await clique.SubmitBlock(TESTCHAINID, rlpHeaders.unsigned, rlpHeaders.signed, storage.address);
      console.log("\tGas used to submit block 873982 = " + validationReceipt.receipt.gasUsed.toString() + " gas");

      // Fetch block 873983 from rinkeby
      block = rinkeby.eth.getBlock(873983);
      rlpHeaders = encoder.encodeBlockHeader(block);

      // Submit block should succeed
      validationReceipt = await clique.SubmitBlock(TESTCHAINID, rlpHeaders.unsigned, rlpHeaders.signed, storage.address);
      console.log("\tGas used to submit block 873983 = " + validationReceipt.receipt.gasUsed.toString() + " gas");
      let submittedEvent = validationReceipt.logs.find(l => { return l.event == 'BlockSubmitted' });
      let blockHash = submittedEvent.args.blockHash;

      // Check proposal is added
      voteProposal = await clique.getProposal.call(TESTCHAINID, blockHash, VALIDATORS_FINISH[1]);
      assert.equal(voteProposal, 1);

      // Fetch block 873984 from rinkeby
      block = rinkeby.eth.getBlock(873984);
      rlpHeaders = encoder.encodeBlockHeader(block);

      // Submit block should succeed
      validationReceipt = await clique.SubmitBlock(TESTCHAINID, rlpHeaders.unsigned, rlpHeaders.signed, storage.address);
      console.log("\tGas used to submit block 873984 = " + validationReceipt.receipt.gasUsed.toString() + " gas");

      // Fetch block 873985 from rinkeby
      block = rinkeby.eth.getBlock(873985);
      rlpHeaders = encoder.encodeBlockHeader(block);

      // Submit block should succeed
      validationReceipt = await clique.SubmitBlock(TESTCHAINID, rlpHeaders.unsigned, rlpHeaders.signed, storage.address);
      console.log("\tGas used to submit block 873985 = " + validationReceipt.receipt.gasUsed.toString() + " gas");

      // Fetch block 873986 from rinkeby
      block = rinkeby.eth.getBlock(873986);
      rlpHeaders = encoder.encodeBlockHeader(block);

      // Submit block should succeed
      validationReceipt = await clique.SubmitBlock(TESTCHAINID, rlpHeaders.unsigned, rlpHeaders.signed, storage.address);
      console.log("\tGas used to submit block 873986 = " + validationReceipt.receipt.gasUsed.toString() + " gas");
      submittedEvent = validationReceipt.logs.find(l => { return l.event == 'BlockSubmitted' });
      blockHash = submittedEvent.args.blockHash;

      // Check proposal is added
      voteProposal = await clique.getProposal.call(TESTCHAINID, blockHash, VALIDATORS_FINISH[1]);
      assert.equal(voteProposal, 0);

      // Check all validators exist
      registeredValidators = await clique.getValidators.call(TESTCHAINID, blockHash);
      for (let i = 0; i < VALIDATORS_FINISH.length; i++) {
          let validatorExists = registeredValidators.some(v => { return v == VALIDATORS_FINISH[i] });;
          assert(validatorExists);
      }

      // Check that the vote threshold has increased with validator set size
      voteThreshold = Math.floor((registeredValidators.length/2) + 1);
      assert.equal(voteThreshold, 3);
    })
  })
});
