// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const Web3Utils = require('web3-utils');
const BN = require('bignumber.js')
const merkle = require('./helpers/merkle.js')
const encoder = require('./helpers/encoder.js')
const Trie = require('merkle-patricia-tree');
const rlp = require('rlp');
const async = require('async')
const EthereumTx = require('ethereumjs-tx');
const EthereumBlock = require('ethereumjs-block/from-rpc')
const levelup = require('levelup');
const sha3 = require('js-sha3').keccak_256
const EP = require('eth-proof');

const Web3 = require('web3');

const Ion = artifacts.require("Ion");
const PatriciaTrie = artifacts.require("PatriciaTrie");

require('chai')
 .use(require('chai-as-promised'))
 .should();

const DEPLOYEDCHAINID = "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075"

const TESTCHAINID = "0x22b55e8a4f7c03e1689da845dd463b09299cb3a574e64c68eafc4e99077a7254"
const TESTBLOCK = {
    difficulty: 2,
    extraData: '0xd88301080b846765746888676f312e31302e32856c696e757800000000000000dd2ba07230e2186ee83ef77d88298c068205167718d48ba5b6ba1de552d0c6ce156011a58b49ed91855de154346968a7eeaaf20914022e58e4f6c0e1e02567ec00',
    gasLimit: 5635559972940396,
    gasUsed: 273138,
    hash: '0x6f98a4b7bffb6c5b3dce3923be8a87eeef94ba22e3266cfcfd53407e70294fa4',
    logsBloom: '0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000',
    miner: '0x0000000000000000000000000000000000000000',
    mixHash: '0x0000000000000000000000000000000000000000000000000000000000000000',
    nonce: '0x0000000000000000',
    number: 5446,
    parentHash: '0xaa912ad61a8aa3e2d1144e4c76b746720e41682122a8b77eff890099a0ff6284',
    receiptsRoot: '0x1d000ef3f5ca9ebc62cc8aaa07e8fbd103583d1e3cbd28c13e62bc8eac5eb2f1',
    sha3Uncles: '0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347',
    size: 2027,
    stateRoot: '0xb347dd25d9a8a456448aed25e072c9db54f464be5e3ce1f505cc171cacf3a967',
    timestamp: 1531327572,
    totalDifficulty: 10893,
    transactions:
     [ '0x63eff998322fd9ec22bbe141ea74ab929197d2db65834e6f4db65743a214cea3',
       '0xa581c3669e5c927e624949d378a5a9df949d4e7f15e1e974c754929408e4b8a5',
       '0x51f1e414334270b7a338f4d81eb82a5560b406f992bf1b3a2371964425e7c0d8',
       '0xc199cd22b3285ea30d798204c3c2fdb8cebfb4648589aa9687aecd9296705ff6',
       '0x4da9368a70e4cfcee28f4c95d69d1256a7d649505f6971b0435bc90f963833f8',
       '0x3cd690a88f4eff005e85f12492afe84344355e9913ea391e52cc0c39debc19e1',
       '0x5dc2e7ea90a0b2630c8138d1357c78ec3d0f55ed23d2951f3c3754ccb9d47446',
       '0xc7f92719dd9f10e8e49ce31a1d271a268269f5c6103629b65869f595109d0462',
       '0x97ff99ad8a3ae45e933464d09b485b7e1adf2fae15ea88d4215cd676b9ca959e',
       '0x343b25b3c1140eb6bf24dbb7ef8595d62178e9ed686fb5d7e6431840c1194314',
       '0x15eb2874404febc7c5cf63bc8ee8100d3f66bf32b69c66805f2fd24732cee39d',
       '0xdfa64978248b67cd5941fe29fc4297ea311aca517ad0e43d71ca59b760fa9ede',
       '0x63f77993f0db424f3bfc202d6f2d3a4cc33979588ef156deff28987c352d44bc' ],
    transactionsRoot: '0xcb9ecdf5483a1435113250201f690124501cfb0c071b697fcfee88c9a368ef35',
    uncles: []
}

const TESTRLPENCODING = "0xf9025fa0aa912ad61a8aa3e2d1144e4c76b746720e41682122a8b77eff890099a0ff6284a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347940000000000000000000000000000000000000000a0b347dd25d9a8a456448aed25e072c9db54f464be5e3ce1f505cc171cacf3a967a0cb9ecdf5483a1435113250201f690124501cfb0c071b697fcfee88c9a368ef35a01d000ef3f5ca9ebc62cc8aaa07e8fbd103583d1e3cbd28c13e62bc8eac5eb2f1b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002821546871405830e4c2a6c83042af2845b463454b861d88301080b846765746888676f312e31302e32856c696e757800000000000000dd2ba07230e2186ee83ef77d88298c068205167718d48ba5b6ba1de552d0c6ce156011a58b49ed91855de154346968a7eeaaf20914022e58e4f6c0e1e02567ec00a00000000000000000000000000000000000000000000000000000000000000000880000000000000000"
const TEST_NODE_VALUE = "0xf86982093f85174876e80083015f909407340652d03d131cd5737aac4a88623682e7e4c40180820bf9a070d26860a32ef4d08d6d91afa73c067af3211dd692a372770927dc9cbddd7869a05aac135e61c984c356509fc27d41b9f0c9c1f23c76d99571491bb0d15936608a"
const TEST_PATH = "0x80"
const TEST_PARENT_NODES = "0xf8c3f851a0448f4ee6a987bf17a91096e25247c3d7d78dbd08afddb5cfd4186d6a9f36bbc080808080808080a0c47289442eb85e0ca1f12c5ac6168f15513036935879931655dadfad3586dcb78080808080808080f86e30b86bf86982093f85174876e80083015f909407340652d03d131cd5737aac4a88623682e7e4c40180820bf9a070d26860a32ef4d08d6d91afa73c067af3211dd692a372770927dc9cbddd7869a05aac135e61c984c356509fc27d41b9f0c9c1f23c76d99571491bb0d15936608a"

contract('Ion.js', (accounts) => {
    it('Deploy Ion', async () => {
        const ion = await Ion.new(DEPLOYEDCHAINID);
        let chainId = await ion.chainId();

        assert.equal(chainId, DEPLOYEDCHAINID);
    })

    it('Register Chain', async () => {
        const ion = await Ion.new(DEPLOYEDCHAINID);

        // Successfully add id of another chain
        await ion.RegisterChain(TESTCHAINID);
        let chain = await ion.chains.call(0);

        assert.equal(chain, TESTCHAINID);

        // Fail adding id of this chain
        await ion.RegisterChain(DEPLOYEDCHAINID).should.be.rejected;

        // Fail adding id of chain already registered
        await ion.RegisterChain(TESTCHAINID).should.be.rejected;
    })

    it('Submit Block', async () => {
        const ion = await Ion.new(DEPLOYEDCHAINID);

        await ion.RegisterChain(TESTCHAINID);

        // Submit block should succeed
        await ion.SubmitBlock(TESTCHAINID, TESTBLOCK.hash, TESTRLPENCODING)

        let blockHash = await ion.m_blockhashes(TESTCHAINID, 0);
        let header = await ion.getBlockHeader.call(blockHash);

        // Separate fetched header info
        parentHash = header[0];
        txRootHash = header[1];
        receiptRootHash = header[2];

        // Assert that block was persisted correctly
        assert.equal(blockHash, TESTBLOCK.hash);
        assert.equal(parentHash, TESTBLOCK.parentHash);
        assert.equal(txRootHash, TESTBLOCK.transactionsRoot);
        assert.equal(receiptRootHash, TESTBLOCK.receiptsRoot);
    })

    it('Fail Submit Block from unknown chain', async () => {
        const ion = await Ion.new(DEPLOYEDCHAINID);

        await ion.RegisterChain(TESTCHAINID);

        await ion.SubmitBlock(TESTCHAINID.slice(0, -2) + "ff", "0xe40cd510f5e415980a2a18ab97b1983c7da43ee56b299cf931c35d9c9ce435f2", "0xf9025ea0f4d7435eff2fcff295eca2c97a1299eeb1d2ce479b4c6e0e799f4a7bed6e4f72a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347940000000000000000000000000000000000000000a019ac400db348a4975008c6e75c537bce261d116bcd74d8b75a9d6992e3b161eda087c9f55218d8784fa39a773791633e9d007a99bef43c12233ebf980810d47464a05ad439bb61e71db83d139847424ac55990546a1b55cc5dd12a57850fd47af845b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000281d2880d08334ef5308dff826928845b23c06eb861d88301080b846765746888676f312e31302e32856c696e757800000000000000461bc1df80fdafba4508e41ef01a570b7998fa0c64eaae65d62e57929afc232a0656a0a43e10387ffebcc8837d1c0d28ab801313e18775f574e73f119452b42e01a00000000000000000000000000000000000000000000000000000000000000000880000000000000000").should.be.rejected;
    })

    it('Fail Submit Block with wrong block hash', async () => {
        const ion = await Ion.new(DEPLOYEDCHAINID);

        await ion.RegisterChain(TESTCHAINID);

        await ion.SubmitBlock(TESTCHAINID, "0xe4" + "1" + "cd510f5e415980a2a18ab97b1983c7da43ee56b299cf931c35d9c9ce435f2", "0xf9025ea0f4d7435eff2fcff295eca2c97a1299eeb1d2ce479b4c6e0e799f4a7bed6e4f72a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347940000000000000000000000000000000000000000a019ac400db348a4975008c6e75c537bce261d116bcd74d8b75a9d6992e3b161eda087c9f55218d8784fa39a773791633e9d007a99bef43c12233ebf980810d47464a05ad439bb61e71db83d139847424ac55990546a1b55cc5dd12a57850fd47af845b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000281d2880d08334ef5308dff826928845b23c06eb861d88301080b846765746888676f312e31302e32856c696e757800000000000000461bc1df80fdafba4508e41ef01a570b7998fa0c64eaae65d62e57929afc232a0656a0a43e10387ffebcc8837d1c0d28ab801313e18775f574e73f119452b42e01a00000000000000000000000000000000000000000000000000000000000000000880000000000000000").should.be.rejected;
    })

    it('Check Tx Proof', async () => {
        const ion = await Ion.new(DEPLOYEDCHAINID);

        await ion.RegisterChain(TESTCHAINID);

        await ion.SubmitBlock(TESTCHAINID, TESTBLOCK.hash, TESTRLPENCODING);

        await ion.CheckTxProof(TESTCHAINID, TESTBLOCK.hash, TEST_NODE_VALUE, TEST_PARENT_NODES, TEST_PATH);
    })

//    it('Check EP Proofs', async () => {
//        const ion = await Ion.new(DEPLOYEDCHAINID);
//
//        // Building transactions in a block for better trie constructions
//        for (let i = 0; i < 10; i++) {
//            web3.eth.sendTransaction({from: web3.eth.accounts[0], to: ion.address, value:1})
//        }
//
//        block = web3.eth.getBlock(5446);
//        blockNumber = block.number;
//
//        block.difficulty = block.difficulty.toNumber();
//        block.totalDifficulty = block.totalDifficulty.toNumber();
//        console.log(block);
//
//        while (block.transactions < 1) {
//            console.log("No transactions found for block: " + blockNumber.toString());
//            console.log("Trying again...");
//            blockNumber -= 1;
//            block = web3.eth.getBlock(blockNumber);
//        }
//
//        txHash = block.transactions[0];
//        transaction = web3.eth.getTransaction(txHash);
//
//        var eP = new EP(web3.currentProvider);
//
//        await eP.getTransactionTrieRoot(txHash).then( (root) => {
//            console.log("EP Root hash = 0x" + root.toString('hex'))
//        })
//
//        var txValue;
//        var txPath;
//        var txParentNodes;
//        await eP.getTransactionProof(txHash).then( (proof) => {
//            verified = EP.transaction(proof.path, proof.value, proof.parentNodes, proof.header, proof.blockHash);
//            assert(verified);
//
//            txValue = rlp.encode(proof.value);
//            txPath = proof.path;
//            txParentNodes = rlp.encode(proof.parentNodes);
//        })
//
//        console.log("EP VALUE = 0x" + txValue.toString('hex'));
//        console.log("EP PATH = 0x" + txPath.toString('hex'));
//        console.log("EP PARENT NODES = 0x" + txParentNodes.toString('hex'));
//
//    })

//    it('Check Infura Proofs', async () => {
//
//        var eP = new EP(new Web3.providers.HttpProvider("https://gmainnet.infura.io"));
//
//        await eP.getTransactionProof("0xb53f752216120e8cbe18783f41c6d960254ad59fac16229d4eaec5f7591319de").then( (proof) => {
//            verified = EP.transaction(proof.path, proof.value, proof.parentNodes, proof.header, proof.blockHash);
//            assert(verified);
//
//            txValue = rlp.encode(proof.value);
//            txPath = rlp.encode(proof.path);
//            txParentNodes = rlp.encode(proof.parentNodes);
//        })
//
//
//    })

    it('Fail Tx Proof', async () => {
        const ion = await Ion.new(DEPLOYEDCHAINID);

        await ion.RegisterChain(TESTCHAINID);

        await ion.SubmitBlock(TESTCHAINID, TESTBLOCK.hash, TESTRLPENCODING);

        // Fail with wrong chain ID
        await ion.CheckTxProof(DEPLOYEDCHAINID, TESTBLOCK.hash, TEST_NODE_VALUE, TEST_PARENT_NODES, TEST_PATH).should.be.rejected;

        // Fail with wrong block hash
        await ion.CheckTxProof(TESTCHAINID, TESTBLOCK.hash.substring(0, 30) + "ff", TEST_NODE_VALUE, TEST_PARENT_NODES, TEST_PATH).should.be.rejected;

        // Fail with wrong path
        await ion.CheckTxProof(TESTCHAINID, TESTBLOCK.hash, TEST_NODE_VALUE, TEST_PARENT_NODES, "0xff").should.be.rejected;
    })

    it('Check Receipt Proof', () => {

    })

    it('Fail Receipt Proof', () => {

    })

    it('Check Roots Proof', () => {

    })

    it('Fail Roots Proof', () => {

    })
})
