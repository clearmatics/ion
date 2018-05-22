// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const crypto = require('crypto');
const web3 = require('web3');
const Web3Utils = require('web3-utils');
const merkle = require('./helpers/merkle.js')

const Hydrogen = artifacts.require("Hydrogen");
const Token = artifacts.require("Token");
const utils = require('./helpers/utils.js')

require('chai')
 .use(require('chai-as-promised'))
 .should();

const REVERT_MSG = 'VM Exception while processing transaction: revert';

contract('Hydrogen', (accounts) => {
    it('Initiate Trade', async () => {
        console.log("\n==== Test: Initiate Trade ====");
        const token = await Token.new();
        const hydrogen = await Hydrogen.new();

        sender = accounts[0];
        recipient = accounts[1];

        value = 1000;

        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeId;
        console.log("Initiating trade...");
        // 1. Initiate Trade Agreement
        await hydrogen.InitiateTradeAgreement(token.address, recipient, value, withdrawHash, refundHash).then((tx) => {
            tradeId = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })
        var expectedTradeId = Web3Utils.soliditySha3(sender, recipient, token.address, value, withdrawHash, refundHash);
        // Assert successful trade initiation
        assert.equal(tradeId.toString(), expectedTradeId);
    })

    it('Deposit', async () => {
        console.log("\n==== Test: Deposit ====");
        const token = await Token.new();
        const hydrogen = await Hydrogen.new();

        sender = accounts[0];
        recipient = accounts[1];

        value = 1000;
        // Mint token to be used
        await token.mint(value);
        balance = await token.balanceOf.call(sender)
        assert.equal(balance.toString(), value.toString());

        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeId;
        // 1. Initiate Trade Agreement
        await hydrogen.InitiateTradeAgreement(token.address, recipient, value, withdrawHash, refundHash).then((tx) => {
            tradeId = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })
        var expectedTradeId = Web3Utils.soliditySha3(sender, recipient, token.address, value, withdrawHash, refundHash);
        assert.equal(tradeId.toString(), expectedTradeId);

        // 2. Deposit funds to Hydrogen contract
        await token.metadataTransfer(hydrogen.address, value, tradeId, {from: sender, gas: 900000}).then((tx) => {
            console.log("Deposit: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful deduction from account
        balance = await token.balanceOf.call(sender);
        assert.equal(balance.toString(), "0");

        // Assert successful credit to contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());
    })

    it('Deposit And Refund', async () => {
        console.log("\n==== Test: Deposit And Refund ====");
        const token = await Token.deployed();
        const hydrogen = await Hydrogen.new();

        sender = accounts[0];
        recipient = accounts[1];

        value = 1000;
        // Mint token to be used
        await token.mint(value);
        balance = await token.balanceOf.call(sender)
        assert.equal(balance.toString(), value.toString());
        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeId;
        // 1. Initiate Trade Agreement
        await hydrogen.InitiateTradeAgreement(token.address, recipient, value, withdrawHash, refundHash).then((tx) => {
            tradeId = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })
        var expectedTradeId = Web3Utils.soliditySha3(sender, recipient, token.address, value, withdrawHash, refundHash);
        assert.equal(tradeId.toString(), expectedTradeId);

        // 2. Deposit funds to Hydrogen contract
        await token.metadataTransfer(hydrogen.address, value, tradeId, {from: sender, gas: 900000}).then((tx) => {
            console.log("Deposit: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful deduction from account
        balance = await token.balanceOf.call(sender);
        assert.equal(balance.toString(), "0");

        // Assert successful credit to contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());

        // 3. Refund funds from contract back to owner account
        await hydrogen.Refund(tradeId, refundRef, {from: sender, gas: 900000}).then((tx) => {
            console.log("Refund: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful refund to original owner
        balance = await token.balanceOf.call(sender);
        assert.equal(balance.toString(), value.toString());

        // Assert successful removal from contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), "0");
    })

    it('Fail Refund with incorrect reference', async () => {
        console.log("\n==== Test: Fail Refund with incorrect reference ====");
        const token = await Token.new();
        const hydrogen = await Hydrogen.new();

        sender = accounts[0];
        recipient = accounts[1];

        value = 1000;
        // Mint token to be used
        await token.mint(value);
        balance = await token.balanceOf.call(sender)
        assert.equal(balance.toString(), value.toString());

        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeId;
        // 1. Initiate Trade Agreement
        await hydrogen.InitiateTradeAgreement(token.address, recipient, value, withdrawHash, refundHash).then((tx) => {
            tradeId = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })

        var expectedTradeId = Web3Utils.soliditySha3(sender, recipient, token.address, value, withdrawHash, refundHash);
        assert.equal(tradeId.toString(), expectedTradeId);

        // 2. Deposit funds to Hydrogen contract
        await token.metadataTransfer(hydrogen.address, value, tradeId, {from: sender, gas: 900000}).then((tx) => {
            console.log("Deposit: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful deduction from account
        balance = await token.balanceOf.call(sender);
        assert.equal(balance.toString(), "0");

        // Assert successful credit to contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());

        // 3. Provide incorrect reference which should revert
        await hydrogen.Refund(tradeId, "Refund2", {from: sender, gas: 900000}).should.be.rejectedWith(REVERT_MSG);

        // Assert that funds have not been returned to owner
        balance = await token.balanceOf.call(sender);
        assert.equal(balance.toString(), "0");

        // Assert that funds remain on hydrogen contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());
    })

    it('Deposit And Withdraw', async () => {
        console.log("\n==== Test: Deposit And Withdraw ====");
        const token = await Token.new();
        const hydrogen = await Hydrogen.new();

        sender = accounts[0];
        recipient = accounts[1];

        value = 1000;
        // Mint token to be used
        await token.mint(value);
        balance = await token.balanceOf.call(sender)
        assert.equal(balance.toString(), value.toString());

        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeId;
        // 1. Initiate Trade Agreement
        await hydrogen.InitiateTradeAgreement(token.address, recipient, value, withdrawHash, refundHash).then((tx) => {
            tradeId = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })

        var expectedTradeId = Web3Utils.soliditySha3(sender, recipient, token.address, value, withdrawHash, refundHash);
        assert.equal(tradeId.toString(), expectedTradeId);

        // 2. Deposit funds to Hydrogen contract
        await token.metadataTransfer(hydrogen.address, value, tradeId, {from: sender, gas: 900000}).then((tx) => {
            console.log("Deposit: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful deduction from account
        balance = await token.balanceOf.call(sender);
        assert.equal(balance.toString(), "0");

        // Assert successful credit to contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());

        // 3. Withdraw funds from hydrogen successfully with withdraw reference
        await hydrogen.Withdraw(tradeId, withdrawRef, {from: recipient, gas: 900000}).then((tx) => {
            console.log("Withdraw: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful withdrawal to recipient
        balance = await token.balanceOf.call(recipient);
        assert.equal(balance.toString(), value.toString());

        // Assert successful withdrawal from contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), "0");
    })

    it('Fail Withdraw with incorrect reference', async () => {
        console.log("\n==== Test: Fail Withdraw with incorrect reference ====");
        const token = await Token.new();
        const hydrogen = await Hydrogen.new();

        sender = accounts[0];
        recipient = accounts[1];

        value = 1000;
        // Mint token to be used
        await token.mint(value);
        balance = await token.balanceOf.call(sender)
        assert.equal(balance.toString(), value.toString());

        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeId;
        // 1. Initiate Trade Agreement
        await hydrogen.InitiateTradeAgreement(token.address, recipient, value, withdrawHash, refundHash).then((tx) => {
            tradeId = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })

        var expectedTradeId = Web3Utils.soliditySha3(sender, recipient, token.address, value, withdrawHash, refundHash);
        assert.equal(tradeId.toString(), expectedTradeId);

        // 2. Deposit funds to Hydrogen contract
        await token.metadataTransfer(hydrogen.address, value, tradeId, {from: sender, gas: 900000}).then((tx) => {
            console.log("Deposit: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful deduction from account
        balance = await token.balanceOf.call(sender);
        assert.equal(balance.toString(), "0");

        // Assert successful credit to contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());

        // 3. Provide incorrect withdraw reference which should revert
        await hydrogen.Withdraw(tradeId, "Withdraw2", {from: recipient, gas: 900000}).should.be.rejectedWith(REVERT_MSG);

        // Assert failure to withdraw to recipient
        balance = await token.balanceOf.call(recipient);
        assert.equal(balance.toString(), "0");

        // Assert that funds remain on contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());

    })
});
