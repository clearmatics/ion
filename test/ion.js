// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

/*
    Ion Mediator contract test

    Tests here are standalone unit tests for Ion functionality.
    Other contracts have been mocked to simulate basic behaviour.

    Tests the central mediator for block passing and validation registering.
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

const Ion = artifacts.require("Ion");
const MockValidation = artifacts.require("MockValidation");
const MockStorage = artifacts.require("MockStorage");

require('chai')
 .use(require('chai-as-promised'))
 .should();

const DEPLOYEDCHAINID = "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075"

const TESTCHAINID = "0x22b55e8a4f7c03e1689da845dd463b09299cb3a574e64c68eafc4e99077a7254"

/*
TESTRPC TEST DATA
*/

const TESTBLOCK = {
    difficulty: 2,
    extraData: '0xd68301080d846765746886676f312e3130856c696e7578000000000000000000583a78dd245604e57368cb2688e42816ebc86eff73ee219dd96b8a56ea6392f75507e703203bc2cc624ce6820987cf9e8324dd1f9f67575502fe6060d723d0e100',
    gasLimit: 7509409,
    gasUsed: 2883490,
    hash: '0x694752333dd1bd0f806cc6ef1063162f4f330c88f9dcd9e61174fcf5e4927eb7',
    logsBloom: '0x22440000020000090000000000000000041000080000008000088000080000000200000400000800000000000000400000000000000000000010000008020102000000000000080000000008800000000000022000000004000000010000000000080000000620400440100010200400082000000000000080040010000100020020000000000000080080000001000000000100000400480000000002000000002000080018000008108000100000000000000000020000050010001004000000000102000040004000000000000000000000004400000000000000000000000208000000000400008200020000004022400000000004000200848000000000',
    miner: '0x0000000000000000000000000000000000000000',
    mixHash: '0x0000000000000000000000000000000000000000000000000000000000000000',
    nonce: '0x0000000000000000',
    number: 2657422,
    parentHash: '0x3471555ab9a99528f02f9cdd8f0017fe2f56e01116acc4fe7f78aee900442f35',
    receiptsRoot: '0x907121bec78b40e8256fac47867d955c560b321e93fc9f046f919ffb5e3823ff',
    sha3Uncles: '0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347',
    size: 4848,
    stateRoot: '0xf526f481ffb6c3c56956d596f2b23e1f7ff17c810ba59efb579d9334a1765444',
    timestamp: 1531931421,
    totalDifficulty: 5023706,
    transactions:
     [ '0x7adbc5ee3712552a1e85962c3ea3d82394cfed7960d60c12d60ebafe67445450',
       '0x6be870e6dfb11894b64371560ec39e563cef91642afd193bfa67874f3508a282',
       '0x5ba6422455cb7127958df15c453bfe60d92921b647879864b531fd6589e36af4',
       '0xa2597e6fe6882626e12055b1378025aa64a85a03dd23f5dc66034f2ef3746810',
       '0x7ffb940740050ae3604f99a4eef07c83de5d75076cae42cb1561c370cba3a0a3',
       '0x4d6326a6d4cf606c7e44a4ae6710acd3876363bcaabd1b1b59d29fff4da223c5',
       '0x10b3360ef00cd7c4faf826365fddbd33938292c98c55a4cdb37194a142626f63',
       '0x655290cb44be2e64d3b1825a86d5647579015c5cffb03ede7f67eb34cea6b97f',
       '0x6b5e025ea558f4872112a39539ce9a819bfbb795b04eefcc45e1cf5ea947614c',
       '0xefd68b516babcf8a4ca74a358cfca925d9d2d5177ef7b859f3d9183ff522efe8',
       '0xa056eeeeb098fd5adb283e12e77a239797c96860c21712963f183937613d3391',
       '0xa5d1adf694e3442975a13685a9c7d9013c05a4fdcea5bc827566a331b2fead2b',
       '0x95a47360f89c48f0b1a484cbeee8816b6a0e2fc321bdb9db48082bd7272b4ebc',
       '0x896d29a87393c6607844fa545d38eb96056d5310a6b4e056dc00adde67c24be2',
       '0xef3ce2ad9259920094f7fd5ad00453b35888662696ae9b85a393e55cde3ec28d',
       '0x2de8af9b4e84b3ac93adfce81964cc69bafd0a2dbcac3a5f7628ee9e56fd1c8a',
       '0x2790cdb3377f556e8f5bc8eaaf9c6c0d36d0f242c2e4226af2aac0203f43019b',
       '0x98ae65246249785bd1ac8157900f7e1a2c69d5c3b3ffc97d55b9eacab3e212f0',
       '0x7d4f090c58880761eaaab1399864d4a52631db8f0b21bfb7051f9a214ad07993',
       '0xafc3ab60059ed38e71c7f6bea036822abe16b2c02fcf770a4f4b5fffcbfe6e7e',
       '0x2af8f6c49d1123077f1efd13764cb2a50ff922fbaf49327efc44c6048c38c968',
       '0x6d5e1753dc91dae7d528ab9b02350e726e006a5591a5d315a34a46e2a951b3fb',
       '0xdc864827159c7fde6bbd1672ed9a90ce5d69f5d0c81761bf689775d19a90387e',
       '0x22fb4d90a7125988b2857c50709e544483f898cb1e8036477f9ddd94b177bf93',
       '0x999c2e2ba342bed4ccedea01d638db3bbd1abd6d10784c317843880841db6dec',
       '0x11355abb5fe745ed458b2a78e116f4a8c2fe046a131eafe08f30d23bd9d10394' ],
    transactionsRoot: '0x07f36c7ad26564fa65daebda75a23dfa95d660199092510743f6c8527dd72586',
    uncles: []
}

const signedHeader = [
    TESTBLOCK.parentHash,
    TESTBLOCK.sha3Uncles,
    TESTBLOCK.miner,
    TESTBLOCK.stateRoot,
    TESTBLOCK.transactionsRoot,
    TESTBLOCK.receiptsRoot,
    TESTBLOCK.logsBloom,
    Web3Utils.toBN(TESTBLOCK.difficulty),
    Web3Utils.toBN(TESTBLOCK.number),
    TESTBLOCK.gasLimit,
    TESTBLOCK.gasUsed,
    Web3Utils.toBN(TESTBLOCK.timestamp),
    TESTBLOCK.extraData,
    TESTBLOCK.mixHash,
    TESTBLOCK.nonce
    ];

const TEST_SIGNED_HEADER = '0x' + rlp.encode(signedHeader).toString('hex');

contract('Ion.js', (accounts) => {
    let ion;
    let validation;
    let storage;

    beforeEach('setup contract for each test', async function () {
        ion = await Ion.new(DEPLOYEDCHAINID);
        validation = await MockValidation.new(ion.address);
        storage = await MockStorage.new(ion.address);
    })

    it('Deploy Ion', async () => {
        let chainId = await ion.chainId();

        assert.equal(chainId, DEPLOYEDCHAINID);
    })

    describe('Register Validation', () => {
        it('Successful registration', async () => {
            // Successfully add id of another chain
            let registered = await validation.register.call();
            await validation.register();

            assert(registered);
        })

        it('Fail second registration', async () => {
            // Successfully add id of another chain
            let registered = await validation.register.call();
            await validation.register();

            assert(registered);

            // Fail second attempt to register validation
            validation.register.call().should.be.rejected;
        })

        it('Fail registration by non-contract', async () => {
            ion.registerValidationModule().should.be.rejected;
        })
    })

    describe('Store Block', () => {
        it('Successful Store Block', async () => {
            await validation.register();

            const tx = await validation.SubmitBlock(storage.address, TESTCHAINID, TEST_SIGNED_HEADER);
            let event = tx.receipt.logs.some(l => { return l.topics[0] == '0x' + sha3("AddedBlock()") });
            assert.ok(event, "Block not stored");
        })

        it('Fail Store Block by unregistered validation', async () => {
            validation.SubmitBlock(storage.address, TESTCHAINID, TEST_SIGNED_HEADER).should.be.rejected;
        })

        it('Fail Store Block by non-contract', async () => {
            ion.storeBlock(storage.address, TESTCHAINID, TEST_SIGNED_HEADER).should.be.rejected;
        })

        it('Fail Store Block with non contract storage address', async () => {
            ion.storeBlock(accounts[0], TESTCHAINID, TEST_SIGNED_HEADER).should.be.rejected;
        })
    })
})