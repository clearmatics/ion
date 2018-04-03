'use strict';

const BigNumber = web3.BigNumber;
const utils = require('web3-utils');
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
    });

    it.only("TestStart_OnAbyA(): Deposit from Alice on blockchain A", async function()
    {
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

        const txReceipt = await fluoride.VerifyTest(
          a_contract,
          a_expire,
          token_a,
          a_amount,
          a_sig,
          // b_contract,
          // b_state,
          // b_sig,
          {
            from: a_addr
          }
        )
        const logArgs = helpers.txLoggedArgs(txReceipt)
        // console.log(logArgs.meta)
        // console.log(a_addr)
        // console.log(logArgs.a_addr)
        assert.equal(logArgs.a_addr, a_addr)
        // assert.equal(logArgs.b_addr, a_addr)
    });

});
