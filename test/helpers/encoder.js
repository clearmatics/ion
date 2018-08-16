// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const Web3Utils = require('web3-utils');
const rlp = require('rlp');
const utils = require('./utils.js');

const encoder = {};

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
      extraDataSignature: extraDataSignature,
      extraDataShort: extraDataShort,
      extraBytesShort: extraBytesShort
    };
}

module.exports = encoder;