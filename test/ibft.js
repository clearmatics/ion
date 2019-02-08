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

const Ibft = artifacts.require("IBFT");
const MockIon = artifacts.require("MockIon");
const MockStorage = artifacts.require("MockStorage");

const web3 = new Web3();

web3.setProvider(new web3.providers.HttpProvider('http://localhost:8545'));

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
    "0x287d1058a7ae485552b9d18627570f8a46c4c924",
    "0x13ef33419a28f3d7fdc922b8a8696b4a5002050b",
    "0xf66aa7edb3b19cdf2486689039ad5af7bfae1471",
    "0x1e393c46d7cffc50c66a72067277bd9744a96c5c"
  ];

const GENESIS_HASH = "0xa4db1d14ac6d264cb0b30c6b3a641b634cb78b31747e7533403c3f54b0f78b43";
const BLOCK_HASH = "0x755a0a1145e70191c42bf9a8154e7932384e4a5e05f8eb4f9113dd9c6a0c7647";
const COMMIT_HASH = "0xc1b75da5f66996f3f3370142cba5b74fc0d1aeb3a2610ac7a7c1a1d3fc80983f";

const block = { 
    difficulty: '1',
    extraData: '0xd883010814846765746888676f312e31302e34856c696e757800000000000000f90164f85494287d1058a7ae485552b9d18627570f8a46c4c9249413ef33419a28f3d7fdc922b8a8696b4a5002050b94f66aa7edb3b19cdf2486689039ad5af7bfae1471941e393c46d7cffc50c66a72067277bd9744a96c5cb841e4db88ad7c924cd9c690affd756113bb82209a335eb7af9a3c88f372b36efa914e545de5bc755506f92b85dc4dafe0aace5d35453a1d29e738e7716f0367d2a900f8c9b841096810bb56276a8976367dd17e57eb31137bc1abcd5648b3112d8257298a8c25363eed24028904e076dd8cf89162fb5f0c34ac8a6cc70a6c2076f4335ff646eb01b84184bdaa8b4a389208324fa36d45270661df31548aec3fc339c1b43388b07ead65787a9550c6a571d9987c64b1611965f0bbff7c418926afb53a0956c3968657b801b841078a954e3a96ecdbd191a7e2796d51cb31b5d265dac6c5e9bf948fcdc286b9b17b4e780899855812df8f1989a76fc54687e20b9ab96b51b3254efc4559b5ec3700',
    gasLimit: 4704588,
    gasUsed: 0,
    hash: '0x755a0a1145e70191c42bf9a8154e7932384e4a5e05f8eb4f9113dd9c6a0c7647',
    logsBloom: '0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000',
    miner: '0x13Ef33419A28f3d7Fdc922B8a8696b4A5002050B',
    mixHash: '0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365',
    nonce: '0x0000000000000000',
    number: 1,
    parentHash: '0xa4db1d14ac6d264cb0b30c6b3a641b634cb78b31747e7533403c3f54b0f78b43',
    receiptsRoot: '0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421',
    sha3Uncles: '0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347',
    size: 901,
    stateRoot: '0x30e982d38e5e6ea77f130d8657649120d22fa240ab3ab2beea1212b534a9d5d6',
    timestamp: 1549556108,
    totalDifficulty: '2',
    transactions: [],
    transactionsRoot: '0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421',
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


contract.only('Ibft.js', (accounts) => {
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

  describe('Submit Block', () => {
      it.only('Authentic Submission Happy Path', async () => {
        await ibft.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH, storage.address);

        rlpHeader = encoder.encodeIbftHeader(block);

        // Submit block should succeed
        const validationReceipt = await ibft.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, storage.address);
        console.log(validationReceipt.logs)
      })

      it('Verify Validator', async () => {
        // await clique.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH, storage.address);

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

      it('Verify Seals', async () => {
        // await clique.RegisterChain(TESTCHAINID, VALIDATORS, GENESIS_HASH, storage.address);
        console.log(block)

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

        let sig = '0x' + rlpExtraData[2][0].toString('hex');
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
        rlpExtraData[2] = [];

        console.log("\nrlp extra data")
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
        console.log(Web3Utils.sha3(testBlockHeaderHash+"02"));
        console.log(block.hash);

        console.log("\nSignature Retrieved From Istanbul Extra:");
        console.log(sig);

        const blockHeaderHash = eth_util.sha3(testBlockHeaderHash+"02");
        console.log(blockHeaderHash)
        const {v, r, s} = eth_util.fromRpcSig(sig);

        const pubKey  = eth_util.ecrecover(blockHeaderHash, v, r, s);
        const addrBuf = eth_util.pubToAddress(pubKey);
        // assert.equal(VALIDATORS[0], '0x'+addrBuf.toString('hex'));
        console.log('0x'+addrBuf.toString('hex'));
 
      })

  })
});
