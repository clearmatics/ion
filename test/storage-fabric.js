// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

/*
    Fabric Storage contract test

    Tests here are standalone unit tests for Ion functionality.
    Other contracts have been mocked to simulate basic behaviour.

    Tests Fabric block structure decoding and verification of state transitions.
*/

const Web3Utils = require('web3-utils');
const utils = require('./helpers/utils.js');
const BN = require('bignumber.js')
const encoder = require('./helpers/encoder.js')
const rlp = require('rlp');
const async = require('async')
const levelup = require('levelup');
const sha3 = require('js-sha3').keccak_256

// Connect to the Test RPC running
const Web3 = require('web3');
const web3 = new Web3();
web3.setProvider(new web3.providers.HttpProvider('http://localhost:8545'));

const MockIon = artifacts.require("MockIon");
const FabricStore = artifacts.require("FabricStore");

require('chai')
 .use(require('chai-as-promised'))
 .should();

const DEPLOYEDCHAINID = "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075"
const TESTCHAINID = "0x22b55e8a4f7c03e1689da845dd463b09299cb3a574e64c68eafc4e99077a7254"

/*
TESTRPC TEST DATA
*/
const block = web3.eth.getBlock(1);

const TESTDATA = [{
    channelId: "orgchannel",
    blocks: [{
        hash: "ZWRtvD5Qw-qpV_Ss3TJIpjS-oc-Eh9vCzYRETHZLdIg",
        number: 4,
        prevHash: "Z2_xbMXvb6GmgwbVPjATH8-OmTExmro_qZgR7HR7ZwQ",
        dataHash: "honzEasuVR5cQx756QISDvADD2lgeov6k8WAx5WJ0iU",
        timestampS: 1547722778,
        timestampN: 111393784,
        transactions: [{
            txId: "b19cbdf267a5b41a6889cff3f3577aefb9da80ac597f7c25af482e47dc9d6eb0",
            nsrw: [{
                namespace: "ExampleCC",
                readsets: [{
                    key: "A",
                    version: {
                        blockNumber: 3,
                        txNumber: 0
                    }
                }, {
                   key: "B",
                   version: {
                       blockNumber: 3,
                       txNumber: 0
                   }
                }],
                writesets: [{
                    key: "A",
                    isDelete: "false",
                    value: "0"
                }, {
                    key: "B",
                    isDelete: "false",
                    value: "3"
                }]
            }, {
                namespace: "lscc",
                readsets: [{
                    key: "ExampleCC",
                    version: {
                        blockNumber: 3,
                        txNumber: 0
                    }
                }],
                writesets: []
            }]
        }]
    }]
}]

const formattedData = [[
    TESTDATA[0].channelId,
    [
        TESTDATA[0].blocks[0].hash,
        TESTDATA[0].blocks[0].number,
        TESTDATA[0].blocks[0].prevHash,
        TESTDATA[0].blocks[0].dataHash,
        TESTDATA[0].blocks[0].timestampS,
        TESTDATA[0].blocks[0].timestampN,
        [[
             TESTDATA[0].blocks[0].transactions[0].txId,
             [[
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].namespace,
                [[
                   TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].key,
                   [
                        TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].version.blockNumber,
                        TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].version.txNumber
                   ]
                ], [
                   TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[1].key,
                   [
                        TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[1].version.blockNumber,
                        TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[1].version.txNumber
                   ]
                ]],
                [[
                   TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].key,
                   TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].isDelete,
                   TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].value
                ],[
                   TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].key,
                   TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].isDelete,
                   TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].value
                ]]
             ], [
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].namespace,
                [[
                   TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].key,
                   [
                        TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].version.blockNumber,
                        TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].version.txNumber
                   ]
                ]],
                []
             ]]
        ]]
    ]
]];


contract('FabricStore.sol', (accounts) => {
    let ion;
    let storage;

    beforeEach('setup contract for each test', async function () {
        ion = await MockIon.new(DEPLOYEDCHAINID);
        storage = await FabricStore.new(ion.address);
    })

    describe('Register Chain', () => {
        it('Successful Register Chain', async () => {
            // Successfully add id of another chain
            await ion.addChain(storage.address, TESTCHAINID);

            let chainRegistered = storage.m_chains(TESTCHAINID);
            assert(chainRegistered);
        })

        it('Fail Register Current Chain', async () => {
            // Fail adding deployment chain id
            await ion.addChain(storage.address, DEPLOYEDCHAINID).should.be.rejected;
        })

        it('Fail Register Chain Twice', async () => {
            // Successfully add id of another chain
            await ion.addChain(storage.address, TESTCHAINID);

            let chainRegistered = storage.m_chains(TESTCHAINID);
            assert(chainRegistered);

            await ion.addChain(storage.address, TESTCHAINID).should.be.rejected;
        })
    })

    describe('Add Block', () => {
        it('Successful Add Block', async () => {
            // Successfully add id of another chain
            await ion.addChain(storage.address, TESTCHAINID);

            let rlpEncodedBlock = "0x" + rlp.encode(formattedData).toString('hex');

            await ion.storeBlock(storage.address, TESTCHAINID, "0x0", rlpEncodedBlock);

            let block = await storage.getBlock.call(TESTCHAINID, TESTDATA[0].channelId, TESTDATA[0].blocks[0].hash);
            console.log(block);

            let tx = await storage.getTransaction.call(TESTCHAINID, TESTDATA[0].channelId, TESTDATA[0].blocks[0].transactions[0].txId);
            console.log(tx);

//            assert.equal(chainId, TESTCHAINID);
        })
    })
})