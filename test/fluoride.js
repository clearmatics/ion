'use strict';

const BigNumber = web3.BigNumber;
const utils = require('web3-utils');
const web3Abi = require('web3-eth-abi');
const helpers = require('./helpers/utils.js')

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
  	let sodium;
  	let token;
    let fluoride;

    // Annoyingly two copies of each need to be deployed so they can interact...
    beforeEach(async function() {
  		sodium = await Sodium.new();
      token = await Token.new();
      fluoride = await Fluoride.new(sodium.address);
    });

    const a_addr = accounts[1]
    const b_addr = accounts[2]
    const date = Math.floor(Date.now() / 1000);


    it('Check that contracts have deployed correctly...', async () => {
        console.log("  Sodium address:", sodium.address);
        console.log("  Fluoride address:", fluoride.address);
        console.log("  Token address:", token.address);
    });


    it.only("Start_OnAbyA(): Deposit from Alice on blockchain A", async function()
    {
        // Mint some token and give some CCY to Alice
        await token.mint(500)
        await token.transfer(a_addr, 250)
        const initialBalance = await token.balanceOf(a_addr)

        // Fix all the variables to pass in...
        const a_contract = fluoride.address
        const a_expire = date + 60;
        const token_a = token.address;
        const a_amount = 100;

        const meta_hash = utils.soliditySha3(a_expire, token_a, a_amount)
        const a_hash = utils.soliditySha3(a_contract, meta_hash)

        const a_sig = web3.eth.sign(a_addr, a_hash)

        const b_state = utils.soliditySha3("fake state of b")
        const b_contract = token.address
        const b_hash = utils.soliditySha3(b_contract, b_state)
        const ab_hash = utils.soliditySha3(a_hash, a_addr, b_hash)

        const b_sig = web3.eth.sign(b_addr, ab_hash)

        const abc_hash = utils.soliditySha3(ab_hash, b_addr)
        const c_sig = web3.eth.sign(a_addr, abc_hash)

        const txReceipt = await fluoride.Start_OnAbyA(
          a_contract,
          a_expire,
          token_a,
          a_amount,
          a_sig,
          b_contract,
          b_state,
          b_sig,
          c_sig,
          {
            from: a_addr
          }
        )
        let logArgs = helpers.txLoggedArgs(txReceipt)

        assert.equal(logArgs.a_addr, a_addr)
        assert.equal(logArgs.b_addr, b_addr)

        const input = logArgs.trade_id

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

        const tokenTxReceipt = web3Abi.encodeFunctionCall(
          overloadedTransferAbi,
          [
            fluoride.address,
            a_amount,
            input
          ]
        );

         const transferReceipt = await web3.eth.sendTransaction(
           {
             from: a_addr,
             to: token.address,
             data: tokenTxReceipt,
             value: 0
           }
         );
         const depositBalance = await token.balanceOf(a_addr)
         assert.equal(initialBalance - a_amount, depositBalance)

    });

    it("Start_OnBbyB(): Deposit from Bob on blockchain B", async function()
    {
        // Fix all the variables to pass in...
        const a_contract = token.address
        const a_state = utils.soliditySha3("fake state of a")
        const a_hash = utils.soliditySha3(a_contract, a_state)
        const a_sig = web3.eth.sign(a_addr, a_hash)


        const b_contract = fluoride.address
        const b_expire = date + 60;
        const token_b = token.address;
        const b_amount = 100;

        const meta_hash = utils.soliditySha3(b_expire, token_b, b_amount)
        const b_hash = utils.soliditySha3(b_contract, meta_hash)

        const ab_hash = utils.soliditySha3(a_hash, a_addr, b_hash)
        const b_sig = web3.eth.sign(b_addr, ab_hash)

        const abc_hash = utils.soliditySha3(ab_hash, b_addr)
        const c_sig = web3.eth.sign(a_addr, abc_hash)

        const txReceipt = await fluoride.Start_OnBbyB(
          a_contract,
          a_state,
          a_sig,
          b_contract,
          b_expire,
          token_b,
          b_amount,
          b_sig,
          c_sig,
          {
            from: b_addr
          }
        )
        const logArgs = helpers.txLoggedArgs(txReceipt)

        const input = logArgs.trade_id
        console.log(input)

        assert.equal(logArgs.a_addr, a_addr)
        assert.equal(logArgs.b_addr, b_addr)
    });

});
