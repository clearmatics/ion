// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

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

const Ion = artifacts.require("Ion");
const ShareSettle = artifacts.require("ShareSettle");
const Token = artifacts.require("Token");
const FabricStore = artifacts.require("FabricStore");
const BaseValidation = artifacts.require("Base");

require('chai')
 .use(require('chai-as-promised'))
 .should();

const DEPLOYEDCHAINID = "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075"
const TESTCHAINID = "0x22b55e8a4f7c03e1689da845dd463b09299cb3a574e64c68eafc4e99077a7254"

let TRANSFER_DATA = [{
    channelId: "sharechannel",
    blocks: [{
        hash: "he4OZ961NAoayXkDhmPOFOt_1loz3hpXXDuH0IyXTXw",
        number: 4,
        prevHash: "CGly4yQSmLeFiJObft1-3e7LMXa7ka96sAE14Df0mRQ",
        dataHash: "83Rehn_T0vc47w3IXKjKkmaWU0f_qsy_LIdqzgC1LdM",
        timestampS: 1548674588,
        timestampN: 186274873,
        transactions: [{
            txId: "4a1b253d40d4e330ac44a7a7cfd38b70a07312ad5564d7e6d54a6abb8dc214c7",
            nsrw: [{
                namespace: "Shares",
                readsets: [{
                    key: "shares", 
                    version: {
                        blockNumber: 3,
                        txNumber: 0
                    }
                }],
                writesets: [{
                    key: "shares",
                    isDelete: "false",
                    value: "0xf84df84bf84989756e69717565526566f83d864e65774f726794e7cf944311eabff15b1b091422a2ecada1dd053d949ecd4a8ca1560c4bb92ca9ebfa5eab448048db93056489756e69717565526566"
                }]
            }, {
                namespace: "lscc",
                readsets: [{
                    key: "Shares",
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
   
// Create a formatted block from the example block
createData = (DATA) => {
    const formattedData = [[
        DATA[0].channelId,
        [
            DATA[0].blocks[0].hash,
            DATA[0].blocks[0].number,
            DATA[0].blocks[0].prevHash,
            DATA[0].blocks[0].dataHash,
            DATA[0].blocks[0].timestampS,
            DATA[0].blocks[0].timestampN,
            [[
                DATA[0].blocks[0].transactions[0].txId,
                [[
                    DATA[0].blocks[0].transactions[0].nsrw[0].namespace,
                    [[
                        DATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].key,
                        [
                            DATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].version.blockNumber,
                            DATA[0].blocks[0].transactions[0].nsrw[0].readsets[0].version.txNumber
                        ]
                    ]],
                    [[
                        DATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].key,
                        DATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].isDelete,
                        DATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].value  
                    ]]
                ], [
                    DATA[0].blocks[0].transactions[0].nsrw[1].namespace,
                    [[
                       DATA[0].blocks[0].transactions[0].nsrw[1].readsets[0].key,
                       [
                            DATA[0].blocks[0].transactions[0].nsrw[1].readsets[0].version.blockNumber,
                            DATA[0].blocks[0].transactions[0].nsrw[1].readsets[0].version.txNumber
                       ]
                    ]],
                    []
                ]]
            ]]
        ]
    ]];
    return formattedData;
  }

contract('Base-Fabric Integration', (accounts) => {
    let ion;
    let validation;
    let storage;

    // Update transfer data with addresses of accounts
    let updateValue = TRANSFER_DATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].value;
    updateValue = updateValue.replace("e7cf944311eabff15b1b091422a2ecada1dd053d", accounts[0].toString().slice(2));
    updateValue = updateValue.replace("9ecd4a8ca1560c4bb92ca9ebfa5eab448048db93", accounts[1].toString().slice(2));
    TRANSFER_DATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].value = updateValue;
    const formattedData = createData(TRANSFER_DATA);

    let rlpEncodedBlock = "0x" + rlp.encode(formattedData).toString('hex');

    beforeEach('setup contract for each test', async function () {
        ion = await Ion.new(DEPLOYEDCHAINID);
        validation = await BaseValidation.new(ion.address);
        storage = await FabricStore.new(ion.address);
    })


    describe('Chaincode usage Contract', () => {
        const Barclays = accounts[0]
        const LBBW = accounts[1]


        it('Submit Block, retrieve state and execute', async () => {
            await validation.register();
            const token = await Token.new()
            const shareSettle = await ShareSettle.new(token.address, storage.address)

            const value = 5;
            const price = 100;
            const reference = Web3Utils.sha3('uniqueRef');

            // Mint ERC223 tokens, funding LBBW
            await token.mint(value*price, {from: LBBW});

            // LBBW initiates a trade agreement which it will settle later
            let tx = await shareSettle.initiateTrade(
                "NewOrg",
                Barclays,
                value,
                price,
                reference,
                {
                    from: LBBW
                },
            );
            console.log("\tGas used to initiate trade: " + tx.receipt.gasUsed.toString());

            tx = await validation.RegisterChain(TESTCHAINID, storage.address);

            let receipt = await validation.SubmitBlock(TESTCHAINID, rlpEncodedBlock, storage.address);
            console.log("\tGas used to store fabric block: %d", receipt.receipt.gasUsed);

            /// Escrow tokens in IonLock contract under specific reference
            let receiptTransfer = await token.metadataTransfer(
                shareSettle.address,
                value*price,
                reference,
                {
                    from: LBBW
                },
            );
            console.log("\tGas used to escrow tokens: " + receiptTransfer.receipt.gasUsed.toString());

            let balance = await token.balanceOf(shareSettle.address);
            assert.equal(value*price, balance)
            
            tx = await shareSettle.retrieveAndExecute(TESTCHAINID, TRANSFER_DATA[0].channelId, TRANSFER_DATA[0].blocks[0].transactions[0].nsrw[0].writesets[0].key);
            console.log("\tGas used to retrieve and execute: " + tx.receipt.gasUsed.toString());
            
            assert.equal(500, await token.balanceOf(Barclays));
            assert.equal(0, await token.balanceOf(LBBW));
            assert.equal(0, await token.balanceOf(shareSettle.address));
            
        })

    })

})