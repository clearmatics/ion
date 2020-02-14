// Copyright (c) 2016-2019 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

/*
    Ibft Validation contract test

    Tests here are standalone unit tests for ibft module functionality.
    Other contracts have been mocked to simulate basic behaviour.

    Tests the ibft scheme for block submission, validator signature verification and more.
*/

const benchmark= require("solidity-benchmark")
const encoder = require('./helpers/encoder.js');
const Web3 = require('web3');
const Web3Utils = require('web3-utils');
const rlp = require('rlp');
const sha3 = require('js-sha3').keccak_256
const Ibft = artifacts.require("IBFT");
const MockIon = artifacts.require("MockIon");
const MockStorage = artifacts.require("MockStorage");
const { MerkleTree } = require('./helpers/merkleTree.js');
const { keccak256, bufferToHex } = require('ethereumjs-util');
const Web3EthAbi = require('web3-eth-abi');
const testBlocks = require("./helpers/blockSamples").testBlocks.ibft //require ibft testBlocks

const web3 = new Web3();

web3.setProvider(new web3.providers.HttpProvider('http://localhost:8545'));

require('chai')
 .use(require('chai-as-promised'))
 .should();

function pad(n, width, z) {
  z = z || '0';
  n = n + '';
  return n.length >= width ? n : new Array(width - n.length + 1).join(z) + n;
}

const DEPLOYEDCHAINID = "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075"
const TESTCHAINID = "0x22b55e8a4f7c03e1689da845dd463b09299cb3a574e64c68eafc4e99077a7254"
const GENESIS_HASH = "0x6893c6fe9270461992e748db2f30aa1359babbd74d0392eb4c3476ef942eb5ec";

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


contract('Ibft.js', (accounts) => {
  const joinHex = arr => '0x' + arr.map(el => el.slice(2)).join('');

  const watchEvent = (eventObj) => new Promise((resolve,reject) => eventObj.watch((error,event) => error ? reject(error) : resolve(event)));

  const expectedHashValidatorsBefore = bufferToHex(keccak256(Web3EthAbi.encodeParameter("address[]", testBlocks.validators_5.validators.map(x => x.toLowerCase()).sort())));
  const expectedHashValidatorsAfter= bufferToHex(keccak256(Web3EthAbi.encodeParameter("address[]", testBlocks.validators_4.validators.map(x => x.toLowerCase()).sort())));

  let ion;
  let ibft;
  let storage;
  let txToBenchmark, duration, currentTestName

  beforeEach('setup contract for each test', async function () {
    ion = await MockIon.new(DEPLOYEDCHAINID);
    ibft = await Ibft.new(ion.address);
    storage = await MockStorage.new(ion.address);

    //unset variables to check for benchmark after each test 
    txToBenchmark = undefined
    duration = 0
    
    //set current test name to use in afterEach hook
    currentTestName = "ibft-" + this.currentTest.title

  })

  afterEach("save to file tx hash and benchmark time", async () => {

    // if variables txToBenchmark has been set inside the current test
      if(txToBenchmark){
        duration = duration ? duration : "Not estimated"
        benchmark.saveStatsToFile(txToBenchmark.tx, currentTestName, txToBenchmark.receipt.gasUsed.toString(), duration)
      }

  })

  after("Trace the transactions benchmarked in this test suite", async () => {
    await benchmark.trace()
  })

  describe('Register Chain', () => {
    it('Successful Register Chain', async () => {
      // get block and validators from samples
      block = testBlocks.validators_5.block
      validators = testBlocks.validators_5.validators

      // Successfully add id of another chain
      let start = Date.now()
      txToBenchmark = await ibft.RegisterChain(TESTCHAINID, validators, block.parentHash, storage.address);
      duration = Number( (Date.now() - start) / 1000 ).toFixed(3)

      console.log("\tGas used to register chain = " + txToBenchmark.receipt.gasUsed.toString() + " gas");

      let chainExists = await ibft.chains(TESTCHAINID);

      assert(chainExists);

      let chainHead = await ibft.m_chainHeads(TESTCHAINID);
      assert.equal(chainHead, block.parentHash);
    })

    it('Fail Register Chain Twice', async () => {
      // get block and validators from samples
      block = testBlocks.validators_5.block
      validators = testBlocks.validators_5.validators

      // Successfully add id of another chain
      await ibft.RegisterChain(TESTCHAINID, validators, block.parentHash, storage.address);

      let chainExists = await ibft.chains(TESTCHAINID);

      assert(chainExists);

      let chainHead = await ibft.m_chainHeads(TESTCHAINID);
      assert.equal(chainHead, block.parentHash);

      // Fail adding id of this chain
      await ibft.RegisterChain(DEPLOYEDCHAINID, validators, block.parentHash, storage.address).should.be.rejected;

      // Fail adding id of chain already initialised
      await ibft.RegisterChain(TESTCHAINID, validators, block.parentHash, storage.address).should.be.rejected;
    })

    it('Check Validators', async () => {
       // get block and validators from samples
       block = testBlocks.validators_5.block
       validators = testBlocks.validators_5.validators

      // Successfully add id of another chain
      await ibft.RegisterChain(TESTCHAINID, validators, block.parentHash, storage.address);

      let chainValidatorsRoot = await ibft.getValidatorsRoot.call(TESTCHAINID);
      assert.equal(chainValidatorsRoot, expectedHashValidatorsBefore)

    })

    it('Check Genesis Hash', async () => {
       // get block and validators from samples
       block = testBlocks.validators_5.block
       validators = testBlocks.validators_5.validators
 
      // Successfully add id of another chain
      await ibft.RegisterChain(TESTCHAINID, validators, block.parentHash, storage.address);

      let chainHead = await ibft.m_chainHeads(TESTCHAINID);
      assert.equal(chainHead, block.parentHash);
    })
  })

  describe('Submit Block', () => {
      it('Submit Sequential Blocks 4 validators', async () => {
        // get blocks and validators from samples
        firstBlock = testBlocks.validators_5.block
        firstValidators = testBlocks.validators_5.validators

        secondBlock = testBlocks.validators_4.block
        secondValidators = testBlocks.validators_4.validators

        await ibft.RegisterChain(TESTCHAINID, firstValidators, firstBlock.parentHash, storage.address);

        rlpHeader = encoder.encodeIbftHeader(firstBlock);

        // Submit block should succeed   
        tx = await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, storage.address, firstValidators);  

        // Check validators have changed
        let chainValidatorsRoot = await ibft.getValidatorsRoot.call(TESTCHAINID);
        assert.equal(chainValidatorsRoot, expectedHashValidatorsAfter)

        let event = tx.receipt.rawLogs.some(l => { return l.topics[0] == '0x' + sha3("AddedBlock()") });
        assert.ok(event, "Stored event not emitted");

        // submit another block
        rlpHeader = encoder.encodeIbftHeader(secondBlock);

        let start = Date.now()
        txToBenchmark = await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, storage.address, secondValidators);
        duration = Number( (Date.now() - start) / 1000 ).toFixed(3)

        console.log("\tGas used to submit block with 4 validators= " + txToBenchmark.receipt.gasUsed.toString() + " gas");

        event = txToBenchmark.receipt.rawLogs.some(l => { return l.topics[0] == '0x' + sha3("AddedBlock()") });
        assert.ok(event, "Stored event not emitted");

        const submittedEvent = txToBenchmark.logs.find(l => { return l.event == 'BlockSubmitted' });
        assert.equal(Web3Utils.sha3(rlpHeader.signed), submittedEvent.args.blockHash);

        let addedBlockHash = await ibft.m_chainHeads.call(TESTCHAINID);
        assert.equal(addedBlockHash, secondBlock.hash);

        let header = await ibft.m_blockheaders(TESTCHAINID, secondBlock.hash);

        // Separate fetched header info
        parentHash = header[2];

        // Assert that block was persisted correctly
        assert.equal(parentHash, secondBlock.parentHash);

        // Check new validators
        chainValidatorsRoot = await ibft.getValidatorsRoot.call(TESTCHAINID);
        assert.equal(chainValidatorsRoot, expectedHashValidatorsBefore)
      })

      it('Successful Submit block - 5 validators', async () => {
        // get block and validators from samples
        block = testBlocks.validators_5.block
        validators = testBlocks.validators_5.validators

        // start
        await ibft.RegisterChain(TESTCHAINID, validators, block.parentHash, storage.address);

        let chainHead = await ibft.m_chainHeads(TESTCHAINID);
        assert.equal(chainHead, block.parentHash);

        rlpHeader = encoder.encodeIbftHeader(block);

        // Submit block should succeed
        let start = Date.now()
        txToBenchmark = await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, storage.address, validators);
        duration = Number( (Date.now() - start) / 1000 ).toFixed(3)
        
        console.log("\tGas used to submit block with 5 validators= " + txToBenchmark.receipt.gasUsed.toString() + " gas");

        let event = txToBenchmark.receipt.rawLogs.some(l => { return l.topics[0] == '0x' + sha3("AddedBlock()") });
        assert.ok(event, "Stored event not emitted");

        const submittedEvent = txToBenchmark.logs.find(l => { return l.event == 'BlockSubmitted' });
        assert.equal(Web3Utils.sha3(rlpHeader.signed), submittedEvent.args.blockHash);

        let addedBlockHash = await ibft.m_chainHeads.call(TESTCHAINID);
        assert.equal(addedBlockHash, block.hash);

        let header = await ibft.m_blockheaders(TESTCHAINID, block.hash);

        // Separate fetched header info
        parentHash = header[2];

        // Assert that block was persisted correctly
        assert.equal(parentHash, block.parentHash);

        chainHead = await ibft.m_chainHeads(TESTCHAINID);
        assert.equal(chainHead, block.hash);

        // Check new validators that are changed in this block
        let chainValidatorsRoot = await ibft.getValidatorsRoot.call(TESTCHAINID);
        assert.equal(chainValidatorsRoot, expectedHashValidatorsAfter)

      })

      it('Successful Submit block - 8 validators', async () => {
        // get block and validators from samples
        block = testBlocks.validators_8.block
        validators = testBlocks.validators_8.validators
        expectedHash = bufferToHex(keccak256(Web3EthAbi.encodeParameter("address[]", validators.map(x => x.toLowerCase()).sort())));

        // start
        await ibft.RegisterChain(TESTCHAINID, validators, block.parentHash, storage.address);

        let chainHead = await ibft.m_chainHeads(TESTCHAINID);
        assert.equal(chainHead, block.parentHash);

        rlpHeader = encoder.encodeIbftHeader(block);

        // Submit block should succeed
        let start = Date.now()
        txToBenchmark = await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, storage.address, validators);
        duration = Number( (Date.now() - start) / 1000 ).toFixed(3)
        
        console.log("\tGas used to submit block with 8 validators= " + txToBenchmark.receipt.gasUsed.toString() + " gas");

        let event = txToBenchmark.receipt.rawLogs.some(l => { return l.topics[0] == '0x' + sha3("AddedBlock()") });
        assert.ok(event, "Stored event not emitted");

        const submittedEvent = txToBenchmark.logs.find(l => { return l.event == 'BlockSubmitted' });
        assert.equal(Web3Utils.sha3(rlpHeader.signed), submittedEvent.args.blockHash);

        let addedBlockHash = await ibft.m_chainHeads.call(TESTCHAINID);
        assert.equal(addedBlockHash, block.hash);

        let header = await ibft.m_blockheaders(TESTCHAINID, block.hash);

        // Separate fetched header info
        parentHash = header[2];

        // Assert that block was persisted correctly
        assert.equal(parentHash, block.parentHash);

        chainHead = await ibft.m_chainHeads(TESTCHAINID);
        assert.equal(chainHead, block.hash);

        // Check validators hash is correct
        let chainValidatorsRoot = await ibft.getValidatorsRoot.call(TESTCHAINID);
        assert.equal(chainValidatorsRoot, expectedHash)
      })

      it('Successful Submit block - 16 validators', async () => {
        // get block and validators from samples
        block = testBlocks.validators_16.block
        validators = testBlocks.validators_16.validators
        expectedHash = bufferToHex(keccak256(Web3EthAbi.encodeParameter("address[]", validators.map(x => x.toLowerCase()).sort())));

        // start
        await ibft.RegisterChain(TESTCHAINID, validators, block.parentHash, storage.address);

        let chainHead = await ibft.m_chainHeads(TESTCHAINID);
        assert.equal(chainHead, block.parentHash);

        rlpHeader = encoder.encodeIbftHeader(block);

        // Submit block should succeed
        let start = Date.now()
        txToBenchmark = await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, storage.address, validators);
        duration = Number( (Date.now() - start) / 1000 ).toFixed(3)
        
        console.log("\tGas used to submit block with 16 validators= " + txToBenchmark.receipt.gasUsed.toString() + " gas");

        let event = txToBenchmark.receipt.rawLogs.some(l => { return l.topics[0] == '0x' + sha3("AddedBlock()") });
        assert.ok(event, "Stored event not emitted");

        const submittedEvent = txToBenchmark.logs.find(l => { return l.event == 'BlockSubmitted' });
        assert.equal(Web3Utils.sha3(rlpHeader.signed), submittedEvent.args.blockHash);

        let addedBlockHash = await ibft.m_chainHeads.call(TESTCHAINID);
        assert.equal(addedBlockHash, block.hash);

        let header = await ibft.m_blockheaders(TESTCHAINID, block.hash);

        // Separate fetched header info
        parentHash = header[2];

        // Assert that block was persisted correctly
        assert.equal(parentHash, block.parentHash);

        chainHead = await ibft.m_chainHeads(TESTCHAINID);
        assert.equal(chainHead, block.hash);

        // Check validators hash is correct
        let chainValidatorsRoot = await ibft.getValidatorsRoot.call(TESTCHAINID);
        assert.equal(chainValidatorsRoot, expectedHash)
      })
   
      it('Fail Submit Block with Unknown Validator', async () => {
        block = testBlocks.validators_5.block 
        validators = testBlocks.validators_5.validators

        await ibft.RegisterChain(TESTCHAINID, validators, block.parentHash, storage.address);

        block.extraData = "0xdc83010000886175746f6e69747988676f312e31302e34856c696e7578000000f90164f854941cb62855cd70774634c85c9acb7c3070ce692936946b2f468af3d0ba2f3a09712faea4d379c2e891a194a667ea98809a69724c6672018bd7db799cd3fefc94c2054df3acfdbe5b221866b25e09026734ca5572b841012edd2e5936deaf4c0ee17698dc0fda832bb51a81d929ae3156d73e5475123c19d162cf1e434637c16811d63d1d3b587906933d75e25cedf7bef59e8fa8375d01f8c9b841719c5bc521721e71ff7fafff09fdff4037e678a77a816b08d45b89d55f35edc94b5c51cc3eeba79d3de291c3c46fbf04faec4952e7d0836be9ad5d855f525c9301b841a7c9eed0337f92a5d4caf6f57b3b59ba10a14ea615c6264fc82fcf5b2e4b626f701fd3596cd1f8639b37a41cb4f3a7582bb530790441de73e6e3449284127b4d00b841210db6ef89906ef1c77538426d29b8440a1c987d508e396776e63515df2a345767c195dc540cfabdf86d696c73b4a24632445565d322d8e45fa2668ec5e6c0e000";

        rlpHeader = encoder.encodeIbftHeader(block);

        // Submit block should not succeed
        await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, storage.address, validators).should.be.rejected;
        
      })

      it('Fail Submit Block with Insufficient Seals', async () => {
        block = testBlocks.validators_5.block 
        validators = testBlocks.validators_5.validators

        await ibft.RegisterChain(TESTCHAINID, validators, block.parentHash, storage.address);

        let badExtraData = "0xf90164f854944335d75841d8b85187cf651ed130774143927c799461d7d88dbc76259fcf1f26bc0b5763aebd67aead94955425273ef777d6430d910f9a8b10adbe95fff694f00d3c728929e42000c8d92d1a7e6a666f12e6edb8410c11022a97fcb2248a2d757a845b4804755702125f8b7ec6c06503ae0277ad996dc22f81431e8036b6cf9ef7d3c1ff1b65a255c9cb70dd2f4925951503a6fdbf01f8c9b8412d3849c86c8ba3ed9a79cdd71b1684364c4c4efb1f01e83ca8cf663f3c95f7ac64b711cd297527d42fb3111b8f78d5227182f38ccc442be5ac4dcb52efede89a01b84135de3661d0191247c7f835c8eb6d7939052c0da8ae234baf8bd208c00225e706112df9bad5bf773120ba4bbc55f6d18e478de43712c0cd3de7a3e2bfd65abb7c01b841735f482a051e6ad7fb76a815907e68d903b73eff4e472006e56fdeca8155cb575f4c1d3e98cf3a4b013331c1bd171d0d500243ac0e073a5fd382294c4fe996f000";

        // Remove seal from extradata
        const decodedExtraData = rlp.decode(badExtraData);
        decodedExtraData[2].pop()

        // Reapply the rlp encoded istanbul extra minus single seal
        encodedExtraData = rlp.encode(decodedExtraData).toString('hex');
        testBlocks.validators_5.block.extraData = "0xdc83010000886175746f6e69747988676f312e31302e34856c696e7578000000" + encodedExtraData;
        rlpHeader = encoder.encodeIbftHeader(block);

        // Submit testBlocks.validators_5.block should not succeed
        await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, storage.address, validators).should.be.rejected;
        
      })
      
      it("Fails when the provided set of validators is not the one in the previous block", async () => {
        block = testBlocks.validators_5.block 
        validators = testBlocks.validators_5.validators

        await ibft.RegisterChain(TESTCHAINID, validators, block.parentHash, storage.address);

        let chainHead = await ibft.m_chainHeads(TESTCHAINID);
        assert.equal(chainHead, block.parentHash);

        rlpHeader = encoder.encodeIbftHeader(block);

        // Submit block should fail cause i provide a set of validators that are different from the one stored in the previous block
        await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, storage.address, testBlocks.validators_4.validators).should.be.rejected;
      })
  })

});
