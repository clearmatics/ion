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
    channelId: "ch1",
    blocks: [{
        hash: "bTBvR22HJGYhxAiBb4HQapVmuQ2YhDfnvxP20T3rgA==",
        number: 1,
        prevHash: "fyBPC7a86s4JuL+4fBv9Px71ypvHzHxauc/RsBu8aDc=",
        dataHash: "BaC3tFqz0blHsA9vH9vFLvexSUehf+7ocWUNEw/ygVo=",
        timestampS: 1494408444,
        timestampN: 903695991,
        transactions: [{
            txId: "2a087008a13eb1341d5e921fded1e122b271be038443e65b5f40c6b34d6fa481",
            nsrw: [{
                namespace: "cc1",
                readsets: [{
                    key: "key1",
                    version: {
                        blockNumber: 0,
                        txNumber: 1
                    }
                }, {
                   key: "key2",
                   version: {
                       blockNumber: 0,
                       txNumber: 1
                   }
                }],
                writesets: [{
                    key: "key1",
                    value: "Hello"
                }, {
                    key: "key2",
                    value: "World"
                }]
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
                   TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].value
                ],[
                   TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].key,
                   TESTDATA[0].blocks[0].transactions[0].nsrw[0].writesets[1].value
                ]]
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

            console.log(formattedData);

//            await ion.storeBlock(storage.address, TESTCHAINID, TESTDATA[0].blocks[0].hash, TESTRLPENCODING);
        })
    })
})