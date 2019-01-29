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
const util = require('util');

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

const TESTDATA = [{
    channelId: "orgchannel",
    blocks: [{
        hash: "vBmkcC8xbLMAUK-wkLMYGDz9qFdu1n8SbsHp62Of_-o",
        number: 4,
        prevHash: "hnw1EQE3SXA_LCUsRWXAj5nZ_JjLPm6DGiRn-g7g4Pc",
        dataHash: "_xuKFW3Po3gNBXjXac11M39a-o-_92_PC6DWBUWnk4I",
        timestampS: 1548756504,
        timestampN: 121452526,
        transactions: [{
            txId: "d4a03d5b71ac3fab92b90bae047c9c5e6ccf0b4396be6807c1724fc0139f999b",
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
                TESTDATA[0].blocks[0].transactions[0].nsrw[1].namespace,
                [[
                   TESTDATA[0].blocks[0].transactions[0].nsrw[1].readsets[0].key,
                   [
                        TESTDATA[0].blocks[0].transactions[0].nsrw[1].readsets[0].version.blockNumber,
                        TESTDATA[0].blocks[0].transactions[0].nsrw[1].readsets[0].version.txNumber
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
    let rlpEncodedBlock_Gen = "0x" + rlp.encode(formattedData).toString('hex');
    // Generated by CLI https://github.com/Shirikatsu/fabric-examples/blob/fc3ff7243f282a4edf21c6667e71ab02e759c3c5/fabric-cli/cmd/fabric-cli/printer/encoder.go#L40
    let rlpEncodedBlock = "0xf90127f901248a6f72676368616e6e656cf90116ab76426d6b63433878624c4d41554b2d776b4c4d5947447a3971466475316e38536273487036324f665f2d6f04ab686e7731455145335358415f4c435573525758416a356e5a5f4a6a4c506d36444769526e2d673767345063ab5f78754b465733506f33674e42586a58616331314d3339612d6f2d5f39325f50433644574255576e6b3449845c50261884073d37eef885f883b84064346130336435623731616333666162393262393062616530343763396335653663636630623433393662653638303763313732346663303133396639393962f83fe8894578616d706c654343cac441c20380c442c20380d2c8418566616c736530c8428566616c736533d5846c736363cecd894578616d706c654343c20380c0";

    beforeEach('setup contract for each test', async function () {
        ion = await MockIon.new(DEPLOYEDCHAINID);
        storage = await FabricStore.new(ion.address);
    })

    describe('Block Encode', () => {
        it('Correct Encoding', async () => {
            assert.equal(rlpEncodedBlock_Gen, rlpEncodedBlock);
        })
    })

    describe('Register Chain', () => {
        it('Successful Register Chain', async () => {
            // Successfully add id of another chain
            await ion.addChain(storage.address, TESTCHAINID);

            let chainRegistered = await storage.m_chains(TESTCHAINID);
            assert(chainRegistered);

            let chainId = await storage.m_networks.call(TESTCHAINID);
            assert.equal(TESTCHAINID, chainId);
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

            let receipt = await ion.storeBlock(storage.address, TESTCHAINID, rlpEncodedBlock);
            console.log("Gas used to store fabric block: %d", receipt.receipt.gasUsed);

            let block = await storage.getBlock.call(TESTCHAINID, TESTDATA[0].channelId, TESTDATA[0].blocks[0].hash);

            assert.equal(block[0], TESTDATA[0].blocks[0].number);
            assert.equal(block[1], TESTDATA[0].blocks[0].hash);
            assert.equal(block[2], TESTDATA[0].blocks[0].prevHash);
            assert.equal(block[3], TESTDATA[0].blocks[0].dataHash);
            assert.equal(block[4], TESTDATA[0].blocks[0].timestampS);
            assert.equal(block[5], TESTDATA[0].blocks[0].timestampN);
            assert.equal(block[6], TESTDATA[0].blocks[0].transactions[0].txId);

            let tx = await storage.getTransaction.call(TESTCHAINID, TESTDATA[0].channelId, TESTDATA[0].blocks[0].transactions[0].txId);

            assert.equal(tx[0], TESTDATA[0].blocks[0].hash);
            assert.equal(tx[1], TESTDATA[0].blocks[0].transactions[0].nsrw[0].namespace + "," + TESTDATA[0].blocks[0].transactions[0].nsrw[1].namespace);

            let txExists = await storage.isTransactionExists.call(TESTCHAINID, TESTDATA[0].channelId, TESTDATA[0].blocks[0].transactions[0].txId);

            assert(txExists);

            let nsrw = await storage.getNSRW.call(TESTCHAINID, TESTDATA[0].channelId, TESTDATA[0].blocks[0].transactions[0].txId, TESTDATA[0].blocks[0].transactions[0].nsrw[0].namespace);

            let expectedReadset = util.format("{ key: %s, version: { blockNo: %d, txNo: %d } } { key: %s, version: { blockNo: %d, txNo: %d } } ",
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].key,
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].version.blockNumber,
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].version.txNumber,
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[1].key,
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[1].version.blockNumber,
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].readsets[1].version.txNumber
            )

            assert.equal(expectedReadset, nsrw[0]);

            let expectedWriteset = util.format("{ key: %s, isDelete: %s, value: %s } { key: %s, isDelete: %s, value: %s } ",
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].key,
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].isDelete,
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].value,
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].key,
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].isDelete,
                TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].value
            )

            assert.equal(expectedWriteset, nsrw[1]);

            nsrw = await storage.getNSRW.call(TESTCHAINID, TESTDATA[0].channelId, TESTDATA[0].blocks[0].transactions[0].txId, TESTDATA[0].blocks[0].transactions[0].nsrw[1].namespace);

            expectedReadset = util.format("{ key: %s, version: { blockNo: %d, txNo: %d } } ",
                TESTDATA[0].blocks[0].transactions[0].nsrw[1].readsets[0].key,
                TESTDATA[0].blocks[0].transactions[0].nsrw[1].readsets[0].version.blockNumber,
                TESTDATA[0].blocks[0].transactions[0].nsrw[1].readsets[0].version.txNumber
            )

            assert.equal(expectedReadset, nsrw[0]);

            expectedWriteset = "";
            assert.equal(expectedWriteset, nsrw[1]);

            let state = await storage.getState.call(TESTCHAINID, TESTDATA[0].channelId, TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].key);

            assert.equal(state[0], TESTDATA[0].blocks[0].number);
            assert.equal(state[1], 0);
            assert.equal(state[2], TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].value);

            state = await storage.getState.call(TESTCHAINID, TESTDATA[0].channelId, TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].key);

            assert.equal(state[0], TESTDATA[0].blocks[0].number);
            assert.equal(state[1], 0);
            assert.equal(state[2], TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].value);
        })

        it('Fail Add Block from unregistered chain', async () => {
            await ion.storeBlock(storage.address, TESTCHAINID, rlpEncodedBlock).should.be.rejected;
        })

        it('Fail Add Block from non-ion', async () => {
            await ion.addChain(storage.address, TESTCHAINID);

            await storage.addBlock(TESTCHAINID, rlpEncodedBlock).should.be.rejected;
        })

        it('Fail Add Block with malformed data', async () => {
            // Successfully add id of another chain
            await ion.addChain(storage.address, TESTCHAINID);

            await ion.storeBlock(storage.address, TESTCHAINID, "0xf86707843b9aca008257c39461621bcf02914668f8404c1f860e92fc1893f74c8084457094cc1ba07e2ebe15f4ece2fd8ffc9a49d7e9e4e71a30534023ca6b24ab4000567709ad53a013a61e910eb7145aa93e865664c54846f26e09a74bd577eaf66b5dd00d334288").should.be.rejected;
        })

        it('Fail Add Same Block Twice', async () => {
            // Successfully add id of another chain
            await ion.addChain(storage.address, TESTCHAINID);

            await ion.storeBlock(storage.address, TESTCHAINID, rlpEncodedBlock);

            await ion.storeBlock(storage.address, TESTCHAINID, rlpEncodedBlock).should.be.rejected;
        })
    })
})