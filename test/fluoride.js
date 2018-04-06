'use strict';

const BigNumber = web3.BigNumber;
const utils = require('web3-utils');
const web3Abi = require('web3-eth-abi');
const helpers = require('./helpers/utils.js')
const merkle = require('./merkle')

const should = require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bignumber')(BigNumber))
    .should();

const assert = require('assert');

// We need to have some sodium stuff so depploy this also
const Sodium = artifacts.require("Sodium");
const Token = artifacts.require("Token");
const Fluoride = artifacts.require("Fluoride");

contract.only('Fluoride', (accounts) => {
  	let a_sodium;
  	let a_token;
    let a_fluoride;
  	let b_sodium;
  	let b_token;
    let b_fluoride;

    // Annoyingly two copies of each need to be deployed so they can interact...
    beforeEach(async function() {
  		a_sodium = await Sodium.new();
      a_token = await Token.new();
      a_fluoride = await Fluoride.new(a_sodium.address);
  		b_sodium = await Sodium.new();
      b_token = await Token.new();
      b_fluoride = await Fluoride.new(b_sodium.address);
    });

    const a_Alice = accounts[1]
    const a_Bob = accounts[2]
    const b_Alice = accounts[3]
    const b_Bob = accounts[4]
    const date = Math.floor(Date.now() / 1000);


    it('Check that contracts have deployed correctly...', async () => {
        console.log("  Sodium address chain a:", a_sodium.address);
        console.log("  Fluoride address chain a:", a_fluoride.address);
        console.log("  Token address chain a:", a_token.address);
        console.log("  Sodium address chain b:", b_sodium.address);
        console.log("  Fluoride address chain b:", b_fluoride.address);
        console.log("  Token address chain b:", b_token.address);
    });

    it("Perform Start_OnAbyA then Start_OnBbyB", async function()
    {
        // Setup the contract for Alice on A
        const a_contract = a_fluoride.address
        const a_expire = date + 60;
        const token_a = a_token.address;
        const a_amount = 100;

        const a_state = utils.soliditySha3(a_expire, token_a, a_amount)
        const a_hash = utils.soliditySha3(a_contract, a_state)

        const a_sig = web3.eth.sign(a_Alice, a_hash)

        // Setup the contract for Bob on B
        const b_contract = b_fluoride.address
        const b_expire = date + 60;
        const token_b = b_token.address;
        const b_amount = 100;

        const b_state = utils.soliditySha3(b_expire, token_b, b_amount)
        const b_hash = utils.soliditySha3(b_contract, b_state)

        const b_sig = web3.eth.sign(b_Bob, b_hash)

        const ab_hash = utils.soliditySha3(a_hash, a_Alice, b_hash)
        const ab_sig = web3.eth.sign(b_Bob, ab_hash)

        const abc_hash = utils.soliditySha3(ab_hash, b_Bob)
        const ac_sig = web3.eth.sign(a_Alice, abc_hash)
        const bc_sig = web3.eth.sign(a_Alice, abc_hash)

        const txReceiptA = await a_fluoride.Start_OnAbyA(
          a_contract,
          a_expire,
          token_a,
          a_amount,
          a_sig,
          b_contract,
          b_state,
          ab_sig,
          ac_sig,
          {
            from: a_Alice
          }
        )
        let logArgs = helpers.txLoggedArgs(txReceiptA)
        assert.equal(logArgs.a_addr, a_Alice)
        assert.equal(logArgs.b_addr, b_Bob)

        const txReceiptB = await b_fluoride.Start_OnBbyB(
          a_contract,
          a_state,
          a_sig,
          b_contract,
          b_expire,
          token_b,
          b_amount,
          ab_sig,
          bc_sig,
          {
            from: b_Bob
          }
        )
        logArgs = helpers.txLoggedArgs(txReceiptB)
        assert.equal(logArgs.a_addr, a_Alice)
        assert.equal(logArgs.b_addr, b_Bob)

    });


    it.only("Start_OnAbyA(): Deposit from Alice on blockchain A, then deposit",
      async function()
    {
        // Mint some a_token and give some CCY to Alice
        await a_token.mint(500)
        await a_token.transfer(a_Alice, 250)
        await b_token.mint(500)
        await b_token.transfer(b_Bob, 250)
        const initialBalance = await a_token.balanceOf(a_Alice)
        const bobInitBalance = await b_token.balanceOf(b_Bob)

        // Setup the contract for Alice on A
        const a_contract = a_fluoride.address
        const a_expire = date + 60;
        const token_a = a_token.address;
        const a_amount = 100;

        const a_state = utils.soliditySha3(a_expire, token_a, a_amount)
        const a_hash = utils.soliditySha3(a_contract, a_state)

        const a_sig = web3.eth.sign(a_Alice, a_hash)

        // Setup the contract for Bob on B
        const b_contract = b_fluoride.address
        const b_expire = date + 60;
        const token_b = b_token.address;
        const b_amount = 100;

        const b_state = utils.soliditySha3(b_expire, token_b, b_amount)
        const b_hash = utils.soliditySha3(b_contract, b_state)

        const b_sig = web3.eth.sign(b_Bob, b_hash)

        const ab_hash = utils.soliditySha3(a_hash, a_Alice, b_hash)
        const ab_sig = web3.eth.sign(b_Bob, ab_hash)

        const abc_hash = utils.soliditySha3(ab_hash, b_Bob)
        const ac_sig = web3.eth.sign(a_Alice, abc_hash)
        const bc_sig = web3.eth.sign(a_Alice, abc_hash)

        const txReceiptA = await a_fluoride.Start_OnAbyA(
          a_contract,
          a_expire,
          token_a,
          a_amount,
          a_sig,
          b_contract,
          b_state,
          ab_sig,
          ac_sig,
          {
            from: a_Alice
          }
        )
        let logArgs = helpers.txLoggedArgs(txReceiptA)
        const a_tradeId = logArgs.trade_id
        assert.equal(logArgs.a_addr, a_Alice)
        assert.equal(logArgs.b_addr, b_Bob)

        const txReceiptB = await b_fluoride.Start_OnBbyB(
          a_contract,
          a_state,
          a_sig,
          b_contract,
          b_expire,
          token_b,
          b_amount,
          ab_sig,
          bc_sig,
          {
            from: b_Bob
          }
        )
        logArgs = helpers.txLoggedArgs(txReceiptB)
        const b_tradeId = logArgs.trade_id
        assert.equal(logArgs.a_addr, a_Alice)
        assert.equal(logArgs.b_addr, b_Bob)


        // This allows us to use the overloaded transfer from Token.sol
        const overloadedTransferAbi = {
          "constant": false,
          "inputs": [
            { "name": "_to", "type": "address" },
            { "name": "_value", "type": "uint256" },
            { "name": "_data", "type": "bytes" }
          ],
          "name": "transfer",
          "outputs": [],
          "payable": false,
          "stateMutability": "nonpayable",
          "type": "function"
        }

        // Transfer funds to appropriate escrow on chain A
        const tokenTxReceipt = web3Abi.encodeFunctionCall(
          overloadedTransferAbi,
          [
            a_fluoride.address,
            a_amount,
            a_tradeId
          ]
        );

        const transferReceipt = await web3.eth.sendTransaction(
           {
             from: a_Alice,
             to: a_token.address,
             data: tokenTxReceipt,
             value: 0
          }
        );
        const depositBalance = await a_token.balanceOf(a_Alice)
        assert.equal(initialBalance - a_amount, depositBalance)

        // Transfer funds to appropriate escrow on chain B
        const b_tokenTxReceipt = web3Abi.encodeFunctionCall(
          overloadedTransferAbi,
          [
            b_fluoride.address,
            b_amount,
            b_tradeId
          ]
        );

        const b_transferReceipt = await web3.eth.sendTransaction(
          {
            from: b_Bob,
            to: b_token.address,
            data: b_tokenTxReceipt,
            value: 0
          }
        );
        const bobDepBalance = await b_token.balanceOf(b_Bob)
        assert.equal(initialBalance - b_amount, bobDepBalance)

        console.log("================================================================================")
        console.log("First prove sodium is functioning...")
        console.log("================================================================================")
        // Create a random block
        const testData1 = helpers.randomArr()
        const tree1 = merkle.createMerkle(testData1)
        const testData2 = helpers.randomArr()
        const tree2 = merkle.createMerkle(testData2)
        const testData3 = helpers.randomArr()
        const tree3 = merkle.createMerkle(testData3)
        const rootArr1 = [tree1[1],tree2[1],tree3[1]]
        console.log("Random merkle tree:\n", rootArr1)
        console.log("Test tree:\n", tree2[1])
        console.log("Tree data:\n", testData2)

        const groupSize = await b_sodium.GroupSize()

        const nextBlock1 = await b_sodium.NextBlock()
        console.log("Block number 1:\n", nextBlock1)

        const receiptUpdate1 = await b_sodium.Update(nextBlock1,rootArr1)
        console.log("Update Blockchain:")

        const nextBlock2 = await b_sodium.NextBlock()
        console.log("Block number 2:\n", nextBlock2)

        const blocksSubmitted = (nextBlock2.toString(10) - nextBlock1.toString(10))/groupSize
        assert.equal(blocksSubmitted,rootArr1.length,'blocks submitted number wrong')
        console.log("Blocks submitted:\n", blocksSubmitted)

        console.log("Verify that rootArr1 is in the merkle tree:\n")

        const blockNumber = nextBlock2 - (2 * groupSize)
        console.log("Block number:\n", blockNumber)

        const leafHash = merkle.merkleHash(testData2[0])
        console.log("Verify testData2 is in merkletree:\n", testData2[0])
        console.log("leaf hash:\n", leafHash)
        const path = merkle.pathMerkle(testData2[0], tree2[0])
        console.log("Merkle path:\n", path)
        console.log("Comprising of:\n")
        console.log("testData2[0]:\n", testData2[0])
        console.log("tree2[0]:\n", tree2[0])

        const valid = await b_sodium.Verify(blockNumber,leafHash,path)
        console.log("Verification result:\n", valid)
        console.log("Which was passed:")
        console.log("BlockNumber:\n", blockNumber)
        console.log("leafHash:\n", leafHash)
        console.log("path:\n", path)
        assert(valid,'Sodium.verify() failed!')


        console.log("================================================================================")
        console.log("Verfiy that transaction is in the block...")
        console.log("================================================================================")

        // Hash details to be used in sodium
        const reference = utils.sha3(b_tradeId)
        const valueHex = '0x' + utils.toBN(b_amount).toString(16).padStart(64,'0') // make an hex that is good to sha3 in solidity uint256 -> 64 bytes
        const lockAddr = b_sodium.address
        const tokenAddr = b_token.address
        const withdrawReceiver = b_Bob

        //concat the arguments of sha3 in solidity (in the same way solidity does)
        const joinedArgs = '0x' + [a_contract, b_tradeId].map(el=>el.slice(2)).join('')

        const hashData = utils.sha3(joinedArgs)
        console.log("Transaction data to be verified:\n", hashData)

        // submit hashdata to IonLink
        const leaf = joinedArgs // joined args need to be added to the random leafs of the tree
        const testData6 = helpers.randomArr()
        testData6[0] = leaf
        console.log("Transaction leaf:\n", leaf)

        console.log("Create more a new root:")
        const tree6 = merkle.createMerkle(testData6)
        const testData4 = helpers.randomArr()
        const tree4 = merkle.createMerkle(testData4) // IonLink needs 2 roots min to update
        const testData5 = helpers.randomArr()
        const tree5 = merkle.createMerkle(testData5) // IonLink needs 2 roots min to update
        const rootArr2 = [tree4[1],tree5[1],tree6[1]]
        console.log("Random merkle tree:\n", rootArr2)
        console.log("Test tree:\n", tree6[1])
        console.log("Tree data:\n", testData6)

        const nextBlock3 = await b_sodium.NextBlock()
        console.log("Block number 3:\n", nextBlock3)

        const receiptUpdate2 = await b_sodium.Update(nextBlock2,rootArr2)
        console.log("Update Blockchain:")

        const nextBlock4 = await b_sodium.NextBlock()
        console.log("Block number 4:\n", nextBlock4)

        console.log("Verify that rootArr2 is in the merkle tree:\n")

        const blockNumber2 = nextBlock4 - (2 * groupSize)
        console.log("Block number2:\n", blockNumber2)

        const leafHash2 = merkle.merkleHash(leaf)
        console.log("Verify testData6 is in merkletree:\n", testData6[0])
        console.log("leaf hash2:\n", leafHash2)
        const path2 = merkle.pathMerkle(testData6[0],tree6[0])
        console.log("Merkle path:\n", path)
        console.log("Comprising of:\n")
        console.log("testData6[0]:\n", testData6[0])
        console.log("tree6[0]:\n", tree6[0])

        const valid2 = await b_sodium.Verify(blockNumber2,leafHash2,path2)
        console.log("Verification result:\n", valid2)
        console.log("Which was passed:")
        console.log("BlockNumber:\n", blockNumber2)
        console.log("leafHash2:\n", leafHash2)
        console.log("path:\n", path2)
        assert(valid2,'Sodium.verify() failed!')


    });

});
