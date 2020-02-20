// Copyright (c) 2020 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

/*
    Tendermint Validation contract test

    Tests here are standalone unit tests for tendermint module functionality.
    Other contracts have been mocked to simulate basic behaviour.

    Tests the tendermint scheme for block submission, validator signature verification and more.
*/


const MockIon = artifacts.require("MockIon");
const MockStorage = artifacts.require("MockStorage");
const Tendermint = artifacts.require("TendermintAutonity");
const testBlocks = require("./helpers/blockSamples").testBlocks.tendermint_autonity
const encoder = require('./helpers/encoder.js');
const { keccak256, bufferToHex } = require('ethereumjs-util');
const Web3EthAbi = require('web3-eth-abi');
const Web3Utils = require("web3-utils");
const Web3 = require('web3');
const web3 = new Web3();

web3.setProvider(new web3.providers.HttpProvider('http://localhost:8545'));

require('chai')
 .use(require('chai-as-promised'))
 .should();

contract("Tendermint Validation Module", (accounts) => {
    const DEPLOYEDCHAINID = "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075"
    const TESTCHAINID = "0x22b55e8a4f7c03e1689da845dd463b09299cb3a574e64c68eafc4e99077a7254"
    const INITIAL_VOTING_THRESHOLD = 0

    let txToBenchmark, duration, currentTestName

    beforeEach("Setup contract for each test", async () => {
        ion = await MockIon.new(DEPLOYEDCHAINID);
        tendermint = await Tendermint.new(ion.address);
        storage = await MockStorage.new(ion.address);

        //unset variables to check for benchmark after each test 
        txToBenchmark = undefined
        duration = 0

        //set current test name to use in afterEach hook
        // currentTestName = "tendermint-" + this.currentTest.title
    })

    afterEach("save to file tx hash and benchmark time", async () => {

        // if variables txToBenchmark has been set inside the current test
        if(txToBenchmark){
            // duration = duration ? duration : "Not estimated"
            // benchmark.saveStatsToFile(txToBenchmark.tx, currentTestName, txToBenchmark.receipt.gasUsed.toString(), duration)
            console.log("\tGas used to register chain = " + txToBenchmark.receipt.gasUsed.toString() + " gas");
        }
    
    })

    describe("Register Chain", () => {
        let block, validators 



        beforeEach(async () => {
            // get block and validators from samples
            block = testBlocks.validators_5.block
            validators = testBlocks.validators_5.validators 
            
            // register chain 
            tx = await tendermint.RegisterChain(TESTCHAINID, validators, INITIAL_VOTING_THRESHOLD, block.parentHash, storage.address);
        })

        it("Succesful Register Chain", async() => {
            // benchmark register chain call once 
            txToBenchmark = tx;

            // has registered the chain
            let chainExists = await tendermint.supportedChains(TESTCHAINID);
            assert(chainExists);

            // has correctly stored the genesis hash
            let chainHead = await tendermint.id_chainHeaders(TESTCHAINID, block.parentHash);
            assert.equal(chainHead.blockHash, block.parentHash);

            // has correctly calculated and stored the validators hash
            expectedHash = bufferToHex(keccak256(Web3EthAbi.encodeParameter("address[]", validators.map(x => x.toLowerCase()).sort())));
            assert.equal(await tendermint.getValidatorsRoot.call(TESTCHAINID, block.parentHash), expectedHash)

        })

        it("Fails Registering Chain Twice", async () => {
             // Fail adding id of this chain
            await tendermint.RegisterChain(DEPLOYEDCHAINID, validators, block.parentHash, storage.address).should.be.rejected;

            // Fail adding id of chain already initialised
            await tendermint.RegisterChain(TESTCHAINID, validators, block.parentHash, storage.address).should.be.rejected;
        })
    })

    describe("Submit Block", () => {

        afterEach("Perform all the checks", async () => {

            // check the events have been triggered
            let event = txToBenchmark.receipt.rawLogs.some(l => { return l.topics[0] == '0x' + keccak256("AddedBlock()") });
            assert.ok(event, "Stored event not emitted");

            const submittedEvent = txToBenchmark.logs.find(l => { return l.event == 'BlockSubmitted' });
            assert.equal(keccak256(rlpHeader.signed), submittedEvent.args.blockHash);

            // Check that block was persisted correctly
            let addedBlock = await tendermint.id_chainHeaders(TESTCHAINID, block.hash);

            assert.equal(addedBlock.hash, block.hash);
            assert.equal(addedBlock.parentHash, block.parentHash);
            assert.equal(addedBlock.validatorsHash, expectedHash)
            assert.equal(addedBlock.votingPower, 1)
            // TODO check voting power

        })

        it.only("Succesfully Submit Block - 5 validators", async () => {
            // get block and validators from samples
            block = testBlocks.validators_5.block
            validators = testBlocks.validators_5.validators 
            expectedHash = bufferToHex(keccak256(Web3EthAbi.encodeParameter("address[]", validators.map(x => x.toLowerCase()).sort())));

            // add genesis block
            await tendermint.RegisterChain(TESTCHAINID, validators, INITIAL_VOTING_THRESHOLD, block.parentHash, storage.address);

            // submit next block 
            rlpHeader = encoder.encodeIbftHeader(block);
            txToBenchmark = await tendermint.SubmitBlock(TESTCHAINID, rlpHeader.unsigned, rlpHeader.signed, rlpHeader.seal, validators, storage.address);
            assert(false)
        })
    })
})