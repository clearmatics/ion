// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const Web3Utils = require('web3-utils');
const HDWalletProvider = require('truffle-hdwallet-provider');
const mnemonic = "select they inform invite result believe equal daughter front arrest wagon miss same menu twenty";
const addr = "0x279884e133f9346f2fad9cc158222068221b613e";

const rlp = require('rlp');
const Web3 = require('web3');
const EP = require('eth-proof');
const deployedTrig = "0x61621bcf02914668f8404c1f860e92fc1893f74c";
const trigAbi = [{"anonymous":false,"inputs":[{"indexed":false,"name":"caller","type":"address"}],"name":"Triggered","type":"event"},{"constant":false,"inputs":[],"name":"fire","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}];

const TRIG_FIRED_RINKEBY_TXHASH = "0xafc3ab60059ed38e71c7f6bea036822abe16b2c02fcf770a4f4b5fffcbfe6e7e"
const TRIG_FIRED_RINKEBY_BLOCKNO = 2657422

function fireTrigger() {
    web3 = new Web3(new HDWalletProvider(mnemonic, 'https://rinkeby.infura.io/v3/973b00227ca84ced8266b2ab6d7592cb'));
    trig = new web3.eth.Contract(trigAbi, deployedTrig);

    trig.methods.fire().send({from: addr}).on('transactionHash', function(hash) { console.log(hash) });
}

async function getTxProof() {
    var eP = new EP(new HDWalletProvider(mnemonic, 'https://rinkeby.infura.io/v3/973b00227ca84ced8266b2ab6d7592cb'));

    await eP.getTransactionTrieRoot(TRIG_FIRED_RINKEBY_TXHASH).then( (root) => {
        console.log("EP TX Root hash = 0x" + root.toString('hex'))
    })

    var txValue;
    var txPath;
    var txParentNodes;
    await eP.getTransactionProof(TRIG_FIRED_RINKEBY_TXHASH).then( (proof) => {
        verified = EP.transaction(proof.path, proof.value, proof.parentNodes, proof.header, proof.blockHash);

        txValue = rlp.encode(proof.value);
        txPath = proof.path;
        txParentNodes = rlp.encode(proof.parentNodes);
    })
    console.log("EP TX VALUE = 0x" + txValue.toString('hex'));
    console.log("EP TX PATH = 0x" + txPath.toString('hex'));
    console.log("EP TX PARENT NODES = 0x" + txParentNodes.toString('hex'));

}

async function getReceiptProof() {
    var eP = new EP(new HDWalletProvider(mnemonic, 'https://rinkeby.infura.io/v3/973b00227ca84ced8266b2ab6d7592cb'));
    web3 = new Web3(new HDWalletProvider(mnemonic, 'https://rinkeby.infura.io/v3/973b00227ca84ced8266b2ab6d7592cb'));
//    receipt = await web3.eth.getTransactionReceipt(TRIG_FIRED_RINKEBY_TXHASH);
//    console.log(receipt);

    await eP.getReceiptTrieRoot(TRIG_FIRED_RINKEBY_TXHASH).then( (root) => {
        console.log("EP RECEIPT Root hash = 0x" + root.toString('hex'))
    })

    var txValue;
    var txPath;
    var txParentNodes;
    await eP.getReceiptProof(TRIG_FIRED_RINKEBY_TXHASH).then( (proof) => {
        verified = EP.receipt(proof.path, proof.value, proof.parentNodes, proof.header, proof.blockHash);
        console.log("VERIFIED: " + verified)

        txValue = rlp.encode(proof.value);
        txPath = proof.path;
        txParentNodes = rlp.encode(proof.parentNodes);
    })
    console.log("EP RECEIPT VALUE = 0x" + txValue.toString('hex'));
    console.log("EP RECEIPT PATH = 0x" + txPath.toString('hex'));
    console.log("EP RECEIPT PARENT NODES = 0x" + txParentNodes.toString('hex'));

}

function rlpEncodeReceipt(receipt) {
    let fields = [receipt.]
}


receipt = { blockHash: '0x694752333dd1bd0f806cc6ef1063162f4f330c88f9dcd9e61174fcf5e4927eb7',

            blockNumber: 2657422,

            contractAddress: null,

            cumulativeGasUsed: 2435175,

            from: '0x279884e133f9346f2fad9cc158222068221b613e',

            gasUsed: 22467,

            logs:

             [ { address: '0x61621BCf02914668F8404C1f860E92fC1893F74c',

                 topics: [Array],

                 data: '0x000000000000000000000000279884e133f9346f2fad9cc158222068221b613e',

                 blockNumber: 2657422,

                 transactionHash: '0xafc3ab60059ed38e71c7f6bea036822abe16b2c02fcf770a4f4b5fffcbfe6e7e',

                 transactionIndex: 19,

                 blockHash: '0x694752333dd1bd0f806cc6ef1063162f4f330c88f9dcd9e61174fcf5e4927eb7',

                 logIndex: 25,

                 removed: false,

                 id: 'log_ad29b82c' } ],

            logsBloom: '0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000010000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000',

            status: true,

            to: '0x61621bcf02914668f8404c1f860e92fc1893f74c',

            transactionHash: '0xafc3ab60059ed38e71c7f6bea036822abe16b2c02fcf770a4f4b5fffcbfe6e7e',

            transactionIndex: 19 }

console.log(rlpEncodeReceipt(receipt));
//getTxProof();
//getReceiptProof();

//console.log(rlp.decode("0xf902780182ab33b9010000000000020000080000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000008000000000000000000000004000000000000000000000000000000000000000000000000000000000000000000000010000000020020000000000000080000000000000000000000000000000000000000000000000000000000000008000000000000000000000000020000010010000000000000000002000000000000000000000000000000000000000000000000000000000000000000000400000000000000000022000000000000000000000000000000f9016ef89b949edcb9a9c4d34b5d6a082c86cb4f117a1394f831f863a0ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa00000000000000000000000001b31d19b6a9a942bbf3c197ca1e5efede3ff8ff2a0000000000000000000000000e05a6d421f1375864bc6e28567993e815eefec23a00000000000000000000000000000000000000000000000004563918244f40000f8cf941b31d19b6a9a942bbf3c197ca1e5efede3ff8ff2e1a0970a6f99f3b845055cfa2283651f03abd1202ee9ececa9ca4f034161dd02457bb896000000000000000000000000e05a6d421f1375864bc6e28567993e815eefec230000000000000000000000000000000000000000000000004563918244f400000000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000001652657761726420666f72206461696c79206c6f67696e"));
