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


        // Perform updates to sodium
        const expectEventA = helpers.joinHex([b_contract/*, topic*/,b_tradeId])
        const expectEventB = helpers.joinHex([a_contract/*, topic*/,a_tradeId])

        // First for blockchain A
        const testDataA = [expectEventA,"2","3","4","5","6","7"]
        const treeA = merkle.createMerkle(testDataA)
        const pathA = testDataA.map(value => merkle.pathMerkle(value,treeA[0]))

        const leafHashA = merkle.merkleHash(testDataA[0])
        const rootArgA = treeA[1]

        const nextBlockA = await a_sodium.NextBlock()
        const receiptUpdateA = await a_sodium.Update(nextBlockA,[rootArgA])
        const validA = await a_sodium.Verify(nextBlockA,leafHashA,pathA[0])
        assert(validA,'a_sodium.verify() failed!')

        // First for blockchain B
        const testDataB = [expectEventB,"2","3","4","5","6","7"]
        const treeB = merkle.createMerkle(testDataB)
        const pathB = testDataB.map(value => merkle.pathMerkle(value,treeB[0]))

        const leafHashB = merkle.merkleHash(testDataB[0])
        const rootBrgB = treeB[1]

        const nextBlockB = await b_sodium.NextBlock()
        const receiptUpdateB = await b_sodium.Update(nextBlockB,[rootBrgB])
        const validB = await b_sodium.Verify(nextBlockB,leafHashB,pathB[0])
        assert(validB,'b_sodium.verify() failed!')


    });

});
