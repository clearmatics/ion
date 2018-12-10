// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const Web3Utils = require('web3-utils');
const rlp = require('rlp');
const utils = require('./utils.js');

const encoder = {};

// Encodes the block headers from clique returning the signed and unsigned instances
encoder.encodeBlockHeader = (block) => {
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
    
    return { 
      unsigned: encodedUnsignedHeader,
      signed: encodedSignedHeader,
      rawunsigned: unsignedHeader,
      rawsigned: signedHeader,
      extraDataSignature: extraDataSignature,
      extraDataShort: extraDataShort,
      extraBytesShort: extraBytesShort
    };
}

// Takes the extraData field from a clique genesis block and finds the validators
encoder.extractValidators = (extraData) => {
  genesisExtraData = utils.hexToBytes(extraData)

  // Remove dressin, 32 bytes pre validators, 65 bytes post validators, and extra byte for 0x
  extraDataValidators = genesisExtraData.splice(33, genesisExtraData.length-32-65-1)

  // Check that the validators length is factor of 20 
  assert.equal(extraDataValidators.length%20, 0);
  numValidators = extraDataValidators.length / 20;

  let validators = [];

  // Append each new validator to the array
  for (i = 0; i < numValidators; ++i) {
    validator = extraDataValidators.splice(0, 20);
    validators.push('0x' + utils.bytesToHex(validator));
  }

  return validators;
}

encoder.appendBlockHeaders = (signedHeaders, signedHeaderIndices, unsignedHeaders, unsignedHeaderIndices, rlpHeaders) => {
  // Start creating the long list of block headers
  signedHeaders.push(rlpHeaders.signed);
  unsignedHeaders.push(rlpHeaders.unsigned);

  // Need to append the cumulative length
  if (signedHeaderIndices.length==0) {
    signedHeaderIndices.push(utils.hexToBytes(rlpHeaders.signed).splice(1).length);
    unsignedHeaderIndices.push(utils.hexToBytes(rlpHeaders.unsigned).splice(1).length);
  } else {
    signedHeaderIndices.push(utils.hexToBytes(rlpHeaders.signed).splice(1).length + signedHeaderIndices[signedHeaderIndices.length - 1]);
    unsignedHeaderIndices.push(utils.hexToBytes(rlpHeaders.unsigned).splice(1).length + unsignedHeaderIndices[unsignedHeaderIndices.length - 1]);
  }

}

module.exports = encoder;