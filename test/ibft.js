// Copyright (c) 2016-2019 Clearmatics Technologies Ltd
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
const sha3 = require('js-sha3').keccak_256
const config = require("./helpers/config.json")

const Ibft = artifacts.require("IBFT");
const MockIon = artifacts.require("MockIon");
const MockStorage = artifacts.require("MockStorage");

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


const VALIDATORS_BEFORE = [
    '0x4335d75841d8b85187cf651ed130774143927c79',
    '0x61d7d88dbc76259fcf1f26bc0b5763aebd67aead',
    '0x955425273ef777d6430d910f9a8b10adbe95fff6',
    '0xf00d3c728929e42000c8d92d1a7e6a666f12e6ed',
    '0xd42d697aa23f7b3e209259002b456c57af26edd6'
  ];

const VALIDATORS_AFTER = [
    '0x4335d75841d8b85187cf651ed130774143927c79',
    '0x61d7d88dbc76259fcf1f26bc0b5763aebd67aead',
    '0x955425273ef777d6430d910f9a8b10adbe95fff6',
    '0xf00d3c728929e42000c8d92d1a7e6a666f12e6ed'
  ];

const GENESIS_HASH = "0x6893c6fe9270461992e748db2f30aa1359babbd74d0392eb4c3476ef942eb5ec";

const block = {
    difficulty: 1,
    extraData: "0xdc83010000886175746f6e69747988676f312e31302e34856c696e7578000000f90164f854944335d75841d8b85187cf651ed130774143927c799461d7d88dbc76259fcf1f26bc0b5763aebd67aead94955425273ef777d6430d910f9a8b10adbe95fff694f00d3c728929e42000c8d92d1a7e6a666f12e6edb8410c11022a97fcb2248a2d757a845b4804755702125f8b7ec6c06503ae0277ad996dc22f81431e8036b6cf9ef7d3c1ff1b65a255c9cb70dd2f4925951503a6fdbf01f8c9b8412d3849c86c8ba3ed9a79cdd71b1684364c4c4efb1f01e83ca8cf663f3c95f7ac64b711cd297527d42fb3111b8f78d5227182f38ccc442be5ac4dcb52efede89a01b84135de3661d0191247c7f835c8eb6d7939052c0da8ae234baf8bd208c00225e706112df9bad5bf773120ba4bbc55f6d18e478de43712c0cd3de7a3e2bfd65abb7c01b841735f482a051e6ad7fb76a815907e68d903b73eff4e472006e56fdeca8155cb575f4c1d3e98cf3a4b013331c1bd171d0d500243ac0e073a5fd382294c4fe996f000",
    gasLimit: 4877543,
    gasUsed: 0,
    hash: "0xed607d816f792bff503fc01bf8903b50aae5bbc6d00293350e38bba92cde40ab",
    logsBloom: "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    miner: "0x955425273ef777d6430d910f9a8b10adbe95fff6",
    mixHash: "0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365",
    nonce: "0x0000000000000000",
    number: 38,
    parentHash: "0x6893c6fe9270461992e748db2f30aa1359babbd74d0392eb4c3476ef942eb5ec",
    receiptsRoot: "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    sha3Uncles: "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    size: 901,
    stateRoot: "0x4e64a3b5ab9c561f72836209e376d035a0aa23a1fc7251e5d21c3c8437fef58e",
    timestamp: 1549897775,
    totalDifficulty: 39,
    transactions: [],
    transactionsRoot: "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    uncles: []
  };

const block_add = {
    difficulty: 1,
    extraData: "0xdc83010000886175746f6e69747988676f312e31302e34856c696e7578000000f90179f869944335d75841d8b85187cf651ed130774143927c799461d7d88dbc76259fcf1f26bc0b5763aebd67aead94955425273ef777d6430d910f9a8b10adbe95fff694f00d3c728929e42000c8d92d1a7e6a666f12e6ed94d42d697aa23f7b3e209259002b456c57af26edd6b841a01291465dfa2b138d48f0f819c31ae9e707a2ee2f3bb93d1341371ab315c9473a4b93b6ccb2b9b29462da66c1a95b27e9254cdf9fcac731e84c7183772f091200f8c9b841ce258c674a9b7ec8bacd5386313c976cbf3dd3f63dd704f93b5e71155c3ce11f124bcf430e1c285e0bce060172930a2c8c15054a14b5629b5dcec069c87e570400b841640736f30ef4ee4baf68448d87020366da4ce6ad2d3872027bbcba8cbbad58e01f2e4e057075dad411f958753615e4141bce861f2780e0499a485741154c707601b841490aa29598b1a7ee0830799bc781b47bfb22c884e2ed2aedd6e9c7ca648e1b547cb469e92e5f375bc1bc3abc191cb180abc93bf3cb67009c75d397a1ab4717d901",
    gasLimit: 4882305,
    gasUsed: 52254,
    hash: "0xd9944319153421ebe524ad3648fbb733f8d8b4aaa75bca8e406fc3b8c171e568",
    logsBloom: "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    miner: "0xf00d3c728929e42000c8d92d1a7e6a666f12e6ed",
    mixHash: "0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365",
    nonce: "0x0000000000000000",
    number: 39,
    parentHash: "0xed607d816f792bff503fc01bf8903b50aae5bbc6d00293350e38bba92cde40ab",
    receiptsRoot: "0x5340517c0dcd60ef9d9735035fcd4a55607eff320684f48796ff57b0a28c8933",
    sha3Uncles: "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    size: 1066,
    stateRoot: "0x68ebd003e05d477e02be898089958e509ca2bff03fe4a9ca1bef2b24aefda03d",
    timestamp: 1549897776,
    totalDifficulty: 40,
    transactions: ["0x8c0faa1990b8b4e0ec8129cd8e2ccf5578be92ee9540361efad993b51179594c"],
    transactionsRoot: "0xd21dcc8688b5b3ab638474485516cda326615f0e8a9853e97589d198b01916b9",
    uncles: []
  };



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
  
  let ion;
  let ibft;
  let storage;

  beforeEach('setup contract for each test', async function () {
    ion = await MockIon.new(DEPLOYEDCHAINID);
    ibft = await Ibft.new(ion.address);
    storage = await MockStorage.new(ion.address);

  })

  describe('Register Chain', () => {
    it('Successful Register Chain', async () => {
      // Successfully add id of another chain
      let tx = await ibft.RegisterChain(TESTCHAINID, VALIDATORS_BEFORE, GENESIS_HASH, storage.address);
      console.log("\tGas used to register chain = " + tx.receipt.gasUsed.toString() + " gas");
      utils.saveGas(config.BENCHMARK_FILEPATH, tx.tx, "ibft-registerChain", tx.receipt.gasUsed.toString())

      let chainExists = await ibft.chains(TESTCHAINID);

      assert(chainExists);

      let chainHead = await ibft.m_chainHeads(TESTCHAINID);
      assert.equal(chainHead, GENESIS_HASH);
    })

    it('Fail Register Chain Twice', async () => {
      // Successfully add id of another chain
      let tx = await ibft.RegisterChain(TESTCHAINID, VALIDATORS_BEFORE, GENESIS_HASH, storage.address);
      console.log("\tGas used to register chain = " + tx.receipt.gasUsed.toString() + " gas");
      let chainExists = await ibft.chains(TESTCHAINID);

      assert(chainExists);

      let chainHead = await ibft.m_chainHeads(TESTCHAINID);
      assert.equal(chainHead, GENESIS_HASH);

      // Fail adding id of this chain
      await ibft.RegisterChain(DEPLOYEDCHAINID, VALIDATORS_BEFORE, GENESIS_HASH, storage.address).should.be.rejected;

      // Fail adding id of chain already initialised
      await ibft.RegisterChain(TESTCHAINID, VALIDATORS_BEFORE, GENESIS_HASH, storage.address).should.be.rejected;
    })

    it('Check Validators', async () => {
      // Successfully add id of another chain
      await ibft.RegisterChain(TESTCHAINID, VALIDATORS_BEFORE, GENESIS_HASH, storage.address);

      let registeredValidators = await ibft.getValidators.call(TESTCHAINID);

      for (let i = 0; i < VALIDATORS_BEFORE.length; i++) {
          let validatorExists = registeredValidators.map(v => v.toLowerCase()).some(v => { return v == VALIDATORS_BEFORE[i] });;
          assert(validatorExists);
      }
    })

    it('Check Genesis Hash', async () => {
      // Successfully add id of another chain
      await ibft.RegisterChain(TESTCHAINID, VALIDATORS_BEFORE, GENESIS_HASH, storage.address);

      let chainHead = await ibft.m_chainHeads(TESTCHAINID);
      assert.equal(chainHead, GENESIS_HASH);
    })
  })

  describe('Submit Block', () => {
      it('Successful Submit block', async () => {
        await ibft.RegisterChain(TESTCHAINID, VALIDATORS_BEFORE, GENESIS_HASH, storage.address);

        let chainHead = await ibft.m_chainHeads(TESTCHAINID);
        assert.equal(chainHead, GENESIS_HASH);

        rlpHeader = encoder.encodeIbftHeader(block);

        // Submit block should succeed
        const validationReceipt = await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, storage.address);
        console.log("\tGas used to submit block = " + validationReceipt.receipt.gasUsed.toString() + " gas");
        utils.saveGas(config.BENCHMARK_FILEPATH, validationReceipt.tx, "ibft-submitBlock-1", validationReceipt.receipt.gasUsed.toString())

        let event = validationReceipt.receipt.rawLogs.some(l => { return l.topics[0] == '0x' + sha3("AddedBlock()") });
        assert.ok(event, "Stored event not emitted");

        const submittedEvent = validationReceipt.logs.find(l => { return l.event == 'BlockSubmitted' });
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
      })

      it('Submit Sequential Blocks with Additional Validator', async () => {
        await ibft.RegisterChain(TESTCHAINID, VALIDATORS_BEFORE, GENESIS_HASH, storage.address);

        rlpHeader = encoder.encodeIbftHeader(block);

        // Submit block should succeed
        let validationReceipt = await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, storage.address);
        console.log("\tGas used to submit block = " + validationReceipt.receipt.gasUsed.toString() + " gas");
        utils.saveGas(config.BENCHMARK_FILEPATH, validationReceipt.tx, "ibft-submitBlock-2", validationReceipt.receipt.gasUsed.toString())

        let event = validationReceipt.receipt.rawLogs.some(l => { return l.topics[0] == '0x' + sha3("AddedBlock()") });
        assert.ok(event, "Stored event not emitted");

        rlpHeader = encoder.encodeIbftHeader(block_add);

        validationReceipt = await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, storage.address);
        event = validationReceipt.receipt.rawLogs.some(l => { return l.topics[0] == '0x' + sha3("AddedBlock()") });
        assert.ok(event, "Stored event not emitted");

        const submittedEvent = validationReceipt.logs.find(l => { return l.event == 'BlockSubmitted' });
        assert.equal(Web3Utils.sha3(rlpHeader.signed), submittedEvent.args.blockHash);

        let addedBlockHash = await ibft.m_chainHeads.call(TESTCHAINID);
        assert.equal(addedBlockHash, block_add.hash);

        let header = await ibft.m_blockheaders(TESTCHAINID, block_add.hash);

        // Separate fetched header info
        parentHash = header[2];

        // Assert that block was persisted correctly
        assert.equal(parentHash, block_add.parentHash);

        // Check new validators
        let registeredValidators = await ibft.getValidators.call(TESTCHAINID);
        for (let i = 0; i < VALIDATORS_AFTER.length; i++) {
          let validatorExists = registeredValidators.map(v => v.toLowerCase()).some(v => { return v == VALIDATORS_AFTER[i] });;
          assert(validatorExists);
        }
      })

      it('Fail Submit Block with Unknown Validator', async () => {
        await ibft.RegisterChain(TESTCHAINID, VALIDATORS_BEFORE, GENESIS_HASH, storage.address);

        block.extraData = "0xdc83010000886175746f6e69747988676f312e31302e34856c696e7578000000f90164f854941cb62855cd70774634c85c9acb7c3070ce692936946b2f468af3d0ba2f3a09712faea4d379c2e891a194a667ea98809a69724c6672018bd7db799cd3fefc94c2054df3acfdbe5b221866b25e09026734ca5572b841012edd2e5936deaf4c0ee17698dc0fda832bb51a81d929ae3156d73e5475123c19d162cf1e434637c16811d63d1d3b587906933d75e25cedf7bef59e8fa8375d01f8c9b841719c5bc521721e71ff7fafff09fdff4037e678a77a816b08d45b89d55f35edc94b5c51cc3eeba79d3de291c3c46fbf04faec4952e7d0836be9ad5d855f525c9301b841a7c9eed0337f92a5d4caf6f57b3b59ba10a14ea615c6264fc82fcf5b2e4b626f701fd3596cd1f8639b37a41cb4f3a7582bb530790441de73e6e3449284127b4d00b841210db6ef89906ef1c77538426d29b8440a1c987d508e396776e63515df2a345767c195dc540cfabdf86d696c73b4a24632445565d322d8e45fa2668ec5e6c0e000";

        rlpHeader = encoder.encodeIbftHeader(block);

        // Submit block should not succeed
        await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, storage.address).should.be.rejected;
        
      })

      it('Fail Submit Block with Insufficient Seals', async () => {
        await ibft.RegisterChain(TESTCHAINID, VALIDATORS_BEFORE, GENESIS_HASH, storage.address);

        let badExtraData = "0xf90164f854944335d75841d8b85187cf651ed130774143927c799461d7d88dbc76259fcf1f26bc0b5763aebd67aead94955425273ef777d6430d910f9a8b10adbe95fff694f00d3c728929e42000c8d92d1a7e6a666f12e6edb8410c11022a97fcb2248a2d757a845b4804755702125f8b7ec6c06503ae0277ad996dc22f81431e8036b6cf9ef7d3c1ff1b65a255c9cb70dd2f4925951503a6fdbf01f8c9b8412d3849c86c8ba3ed9a79cdd71b1684364c4c4efb1f01e83ca8cf663f3c95f7ac64b711cd297527d42fb3111b8f78d5227182f38ccc442be5ac4dcb52efede89a01b84135de3661d0191247c7f835c8eb6d7939052c0da8ae234baf8bd208c00225e706112df9bad5bf773120ba4bbc55f6d18e478de43712c0cd3de7a3e2bfd65abb7c01b841735f482a051e6ad7fb76a815907e68d903b73eff4e472006e56fdeca8155cb575f4c1d3e98cf3a4b013331c1bd171d0d500243ac0e073a5fd382294c4fe996f000";

        // Remove seal from extradata
        const decodedExtraData = rlp.decode(badExtraData);
        decodedExtraData[2].pop()

        // Reapply the rlp encoded istanbul extra minus single seal
        encodedExtraData = rlp.encode(decodedExtraData).toString('hex');
        block.extraData = "0xdc83010000886175746f6e69747988676f312e31302e34856c696e7578000000" + encodedExtraData;
        rlpHeader = encoder.encodeIbftHeader(block);

        // Submit block should not succeed
        await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, storage.address).should.be.rejected;
        
      })

  })

});
