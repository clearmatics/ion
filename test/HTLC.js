'use strict';

const BigNumber = web3.BigNumber;

const utils = require('./helpers/utils.js')

const should = require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bignumber')(BigNumber))
    .should();

const assert = require('assert');
const HTLC = artifacts.require("HTLC");


contract('HTLC', (accounts) => {
  	let htlc;
    let htlc_owner;

    beforeEach(async function() {
		htlc = await HTLC.new();
    });

    const sender = accounts[1]
    const receiver = accounts[2]
    const date = Math.floor(Date.now() / 1000);

    it("Deposit(): Testing the depositing of funds", async function()
    {
        const hashPair = utils.newSecretHashPair()
        const timeout = date + 10;
        const txReceipt = await htlc.Deposit(
          timeout,
          hashPair.hash,
          receiver,
          {
            value: utils.oneFinney,
            from: sender
          }
        )
        const logArgs = utils.txLoggedArgs(txReceipt)

        assert.equal(logArgs.lock_id.toNumber(), 0);
        assert.equal(logArgs.receiver, receiver);

    });

    it("Claim(): Should send receiver funds when they give the correct preimage", async function()
    {
        const hashPair = utils.newSecretHashPair()
        const timeout = date + 10;
        const newContractTx = await htlc.Deposit(
          timeout,
          hashPair.hash,
          receiver,
          {
            value: utils.oneFinney,
            from: sender
          }
        )

        const receiverBalBefore = web3.eth.getBalance(receiver)
        const contractArgs = utils.txLoggedArgs(newContractTx)

        // Posit a claim
        var sig = web3.eth.sign(receiver, hashPair.secret).slice(2)
        let r = '0x' + sig.substr(0, 64)
        let s = '0x' + sig.substr(64, 64)
        let v = web3.toDecimal(sig.slice(128, 130)) + 27

        const claimTx = await htlc.Claim(
          contractArgs.lock_id,
          hashPair.secret,
          v,
          r,
          s,
          {
            from: receiver
          }
        )
        const expectedBal = receiverBalBefore
          .plus(utils.oneFinney)
          .minus(utils.txGas(claimTx))

        const logArgs = utils.txLoggedArgs(claimTx)
        const currentBal = web3.eth.getBalance(receiver)
        assert.equal(logArgs.verified, receiver)
        assert.equal(currentBal.c[0], expectedBal.c[0]) // check gas cost
        assert.equal(currentBal.c[1], expectedBal.c[1]) // check balances

    });

    it("Refund(): Should return sender funds when they give the correct preimage", async function()
    {
        const hashPair = utils.newSecretHashPair()
        const timeout = date + 2;
        const newContractTx = await htlc.Deposit(
          timeout,
          hashPair.hash,
          receiver,
          {
            value: utils.oneFinney,
            from: sender
          }
        )

        const senderBalBefore = web3.eth.getBalance(sender)
        const contractArgs = utils.txLoggedArgs(newContractTx)

        // Wait sometime to ensure timeout is passed
        await utils.sleep(3000)

        // Posit a claim
        var sig = web3.eth.sign(sender, hashPair.secret).slice(2)
        let r = '0x' + sig.substr(0, 64)
        let s = '0x' + sig.substr(64, 64)
        let v = web3.toDecimal(sig.slice(128, 130)) + 27

        const refundTx = await htlc.Refund(
          contractArgs.lock_id,
          hashPair.secret,
          v,
          r,
          s,
          {
            from: sender
          }
        )
        const expectedBal = senderBalBefore
          .plus(utils.oneFinney)
          .minus(utils.txGas(refundTx))

        const logArgs = utils.txLoggedArgs(refundTx)
        const currentBal = web3.eth.getBalance(sender)
        assert.equal(logArgs.verified, sender)
        assert.equal(currentBal.c[0], expectedBal.c[0]) // check gas cost
        assert.equal(currentBal.c[1], expectedBal.c[1]) // check balances

    });

    it('Verify(): ecrecover result matches address', async function() {
      var msg = 'hello'
      var h = web3.sha3(msg)
      var sig = web3.eth.sign(receiver, h).slice(2)
      var r = '0x' + sig.substr(0, 64)
      var s = '0x' + sig.substr(64, 64)
      var v = web3.toDecimal(sig.slice(128, 130)) + 27

      var result = await htlc.Verify(h, v, r, s)
      const logArgs = utils.txLoggedArgs(result)

      assert.equal(logArgs.verified, receiver)
    })

});
