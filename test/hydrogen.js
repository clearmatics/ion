const crypto = require('crypto');
const web3 = require('web3');
const Web3Utils = require('web3-utils');

const Hydrogen = artifacts.require("Hydrogen");
const Token = artifacts.require("Token");

require('chai')
 .use(require('chai-as-promised'))
 .should();

const REVERT_MSG = 'VM Exception while processing transaction: revert';

contract('Hydrogen', (accounts) => {
    it('Initiate Trade', async () => {
        console.log("\n==== Test: Initiate Trade ====");
        const token = await Token.new();
        const tokenB = await Token.new();
        const hydrogen = await Hydrogen.new();

        initiator = accounts[0];
        counterparty = accounts[1];

        value = 1000;

        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeId;
        console.log("Initiating trade...");
        // 1. Initiate Trade Agreement
        await hydrogen.InitiateTradeAgreement(token.address, tokenB.address, initiator, counterparty, value, withdrawHash, refundHash).then((tx) => {
            tradeId = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })
        var expectedTradeId = Web3Utils.soliditySha3(initiator, counterparty, token.address, tokenB.address, value, withdrawHash, refundHash);
        // Assert successful trade initiation
        assert.equal(tradeId.toString(), expectedTradeId);
    })

    it('Deposit', async () => {
        console.log("\n==== Test: Deposit ====");
        const token = await Token.new();
        const hydrogen = await Hydrogen.new();
        const tokenB = await Token.new();

        initiator = accounts[0];
        counterparty = accounts[1];

        value = 1000;
        // Mint token to be used
        await token.mint(value);
        balance = await token.balanceOf.call(initiator)
        assert.equal(balance.toString(), value.toString());

        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeId;
        // 1. Initiate Trade Agreement
        await hydrogen.InitiateTradeAgreement(token.address, tokenB.address, initiator, counterparty, value, withdrawHash, refundHash).then((tx) => {
            tradeId = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })
        var expectedTradeId = Web3Utils.soliditySha3(initiator, counterparty, token.address, tokenB.address, value, withdrawHash, refundHash);
        assert.equal(tradeId.toString(), expectedTradeId);

        // 2. Deposit funds to Hydrogen contract
        await token.metadataTransfer(hydrogen.address, value, tradeId, {from: initiator, gas: 900000}).then((tx) => {
            console.log("Deposit: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful deduction from account
        balance = await token.balanceOf.call(initiator);
        assert.equal(balance.toString(), "0");

        // Assert successful credit to contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());

        // 3. Check Deposit
        deposited = await hydrogen.CheckDeposit.call(tradeId, {from: initiator, gas: 900000})
        assert(deposited);
    })

    it('Deposit And Refund', async () => {
        console.log("\n==== Test: Deposit And Refund ====");
        const token = await Token.new();
        const hydrogen = await Hydrogen.new();
        const tokenB = await Token.new();

        initiator = accounts[0];
        counterparty = accounts[1];

        value = 1000;
        // Mint token to be used
        await token.mint(value);
        balance = await token.balanceOf.call(initiator)
        assert.equal(balance.toString(), value.toString());
        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeId;
        // 1. Initiate Trade Agreement
        await hydrogen.InitiateTradeAgreement(token.address, tokenB.address, initiator, counterparty, value, withdrawHash, refundHash).then((tx) => {
            tradeId = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })
        var expectedTradeId = Web3Utils.soliditySha3(initiator, counterparty, token.address, tokenB.address, value, withdrawHash, refundHash);
        assert.equal(tradeId.toString(), expectedTradeId);

        // 2. Deposit funds to Hydrogen contract
        await token.metadataTransfer(hydrogen.address, value, tradeId, {from: initiator, gas: 900000}).then((tx) => {
            console.log("Deposit: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful deduction from account
        balance = await token.balanceOf.call(initiator);
        assert.equal(balance.toString(), "0");

        // Assert successful credit to contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());

        // 3. Refund funds from contract back to owner account
        await hydrogen.Refund(tradeId, refundRef, {from: initiator, gas: 900000}).then((tx) => {
            console.log("Refund: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful refund to original owner
        balance = await token.balanceOf.call(initiator);
        assert.equal(balance.toString(), value.toString());

        // Assert successful removal from contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), "0");
    })

    it('Fail Refund with incorrect reference', async () => {
        console.log("\n==== Test: Fail Refund with incorrect reference ====");
        const token = await Token.new();
        const hydrogen = await Hydrogen.new();
        const tokenB = await Token.new();

        initiator = accounts[0];
        counterparty = accounts[1];

        value = 1000;
        // Mint token to be used
        await token.mint(value);
        balance = await token.balanceOf.call(initiator)
        assert.equal(balance.toString(), value.toString());

        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeId;
        // 1. Initiate Trade Agreement
        await hydrogen.InitiateTradeAgreement(token.address, tokenB.address, initiator, counterparty, value, withdrawHash, refundHash).then((tx) => {
            tradeId = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })

        var expectedTradeId = Web3Utils.soliditySha3(initiator, counterparty, token.address, tokenB.address, value, withdrawHash, refundHash);
        assert.equal(tradeId.toString(), expectedTradeId);

        // 2. Deposit funds to Hydrogen contract
        await token.metadataTransfer(hydrogen.address, value, tradeId, {from: initiator, gas: 900000}).then((tx) => {
            console.log("Deposit: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful deduction from account
        balance = await token.balanceOf.call(initiator);
        assert.equal(balance.toString(), "0");

        // Assert successful credit to contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());

        // 3. Provide incorrect reference which should revert
        await hydrogen.Refund(tradeId, "Refund2", {from: initiator, gas: 900000}).should.be.rejectedWith(REVERT_MSG);

        // Assert that funds have not been returned to owner
        balance = await token.balanceOf.call(initiator);
        assert.equal(balance.toString(), "0");

        // Assert that funds remain on hydrogen contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());
    })

    it('Fail to Unlock Funds by incorrect address', async () => {
        console.log("\n==== Test: Fail to Unlock Funds by incorrect address ====");
        const token = await Token.new();
        const hydrogen = await Hydrogen.new();
        const tokenB = await Token.new();

        initiator = accounts[0];
        counterparty = accounts[1];

        value = 1000;
        // Mint token to be used
        await token.mint(value);
        balance = await token.balanceOf.call(initiator)
        assert.equal(balance.toString(), value.toString());

        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeId;
        // 1. Initiate Trade Agreement
        await hydrogen.InitiateTradeAgreement(token.address, tokenB.address, initiator, counterparty, value, withdrawHash, refundHash).then((tx) => {
            tradeId = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })

        var expectedTradeId = Web3Utils.soliditySha3(initiator, counterparty, token.address, tokenB.address, value, withdrawHash, refundHash);
        assert.equal(tradeId.toString(), expectedTradeId);

        // 2. Deposit funds to Hydrogen contract
        await token.metadataTransfer(hydrogen.address, value, tradeId, {from: initiator, gas: 900000}).then((tx) => {
            console.log("Deposit: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful deduction from account
        balance = await token.balanceOf.call(initiator);
        assert.equal(balance.toString(), "0");

        // Assert successful credit to contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());

        // 3. Fail Unlocking of funds
        await hydrogen.UnlockFunds(tradeId, {from: counterparty, gas: 900000}).should.be.rejectedWith(REVERT_MSG);
    })

    it('Deposit And Withdraw', async () => {
        console.log("\n==== Test: Deposit And Withdraw ====");
        const token = await Token.new();
        const hydrogen = await Hydrogen.new();
        const tokenB = await Token.new();

        initiator = accounts[0];
        counterparty = accounts[1];

        value = 1000;
        // Mint token to be used
        await token.mint(value);
        balance = await token.balanceOf.call(initiator)
        assert.equal(balance.toString(), value.toString());

        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeId;
        // 1. Initiate Trade Agreement
        await hydrogen.InitiateTradeAgreement(token.address, tokenB.address, initiator, counterparty, value, withdrawHash, refundHash).then((tx) => {
            tradeId = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })

        var expectedTradeId = Web3Utils.soliditySha3(initiator, counterparty, token.address, tokenB.address, value, withdrawHash, refundHash);
        assert.equal(tradeId.toString(), expectedTradeId);

        // 2. Deposit funds to Hydrogen contract
        await token.metadataTransfer(hydrogen.address, value, tradeId, {from: initiator, gas: 900000}).then((tx) => {
            console.log("Deposit: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful deduction from account
        balance = await token.balanceOf.call(initiator);
        assert.equal(balance.toString(), "0");

        // Assert successful credit to contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());

        // 3. Unlock funds for withdrawal
        await hydrogen.UnlockFunds(tradeId, {from: initiator, gas: 900000}).then((tx) => {
            console.log("Unlock Funds: Gas used = " + tx.receipt.gasUsed);
        });
        fundsUnlocked = await hydrogen.getFundsLockedForTrade.call(tradeId);
        assert(fundsUnlocked);

        // 4. Withdraw funds from hydrogen successfully with withdraw reference
        await hydrogen.Withdraw(tradeId, withdrawRef, {from: counterparty, gas: 900000}).then((tx) => {
            console.log("Withdraw: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful withdrawal to counterparty
        balance = await token.balanceOf.call(counterparty);
        assert.equal(balance.toString(), value.toString());

        // Assert successful withdrawal from contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), "0");
    })

    it('Fail Withdraw', async () => {
        console.log("\n==== Test: Fail Withdraw ====");
        const token = await Token.new();
        const hydrogen = await Hydrogen.new();
        const tokenB = await Token.new();

        initiator = accounts[0];
        counterparty = accounts[1];

        value = 1000;
        // Mint token to be used
        await token.mint(value);
        balance = await token.balanceOf.call(initiator)
        assert.equal(balance.toString(), value.toString());

        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeId;
        // 1. Initiate Trade Agreement
        await hydrogen.InitiateTradeAgreement(token.address, tokenB.address, initiator, counterparty, value, withdrawHash, refundHash).then((tx) => {
            tradeId = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })

        var expectedTradeId = Web3Utils.soliditySha3(initiator, counterparty, token.address, tokenB.address, value, withdrawHash, refundHash);
        assert.equal(tradeId.toString(), expectedTradeId);

        // 2. Deposit funds to Hydrogen contract
        await token.metadataTransfer(hydrogen.address, value, tradeId, {from: initiator, gas: 900000}).then((tx) => {
            console.log("Deposit: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful deduction from account
        balance = await token.balanceOf.call(initiator);
        assert.equal(balance.toString(), "0");

        // Assert successful credit to contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());


        // 3. Fail withdrawals

        console.log("\n==== Test: Fail Withdraw with incorrect trade id ====");
        await hydrogen.Withdraw(tradeId & 0x00, withdrawRef, {from: counterparty, gas: 900000}).should.be.rejectedWith(REVERT_MSG);
        console.log("Success\n")

        console.log("\n==== Test: Fail Withdraw by depositor ====");
        await hydrogen.Withdraw(tradeId, withdrawRef, {from: initiator, gas: 900000}).should.be.rejectedWith(REVERT_MSG);
        console.log("Success\n");

        console.log("\n==== Test: Fail Withdraw by locked funds ====");
        await hydrogen.Withdraw(tradeId, withdrawRef, {from: counterparty, gas: 900000}).should.be.rejectedWith(REVERT_MSG);
        console.log("Success\n");

        // Unlock funds for withdrawal
        await hydrogen.UnlockFunds(tradeId, {from: initiator, gas: 900000}).then((tx) => {
            console.log("Unlock Funds: Gas used = " + tx.receipt.gasUsed);
        });

        console.log("\n==== Test: Fail Withdraw with incorrect reference ====");
        await hydrogen.Withdraw(tradeId, withdrawRef + "2", {from: counterparty, gas: 900000}).should.be.rejectedWith(REVERT_MSG);
        console.log("Success\n");

        // Assert failure to withdraw to counterparty
        balance = await token.balanceOf.call(counterparty);
        assert.equal(balance.toString(), "0");

        // Assert that funds remain on contract
        balance = await token.balanceOf.call(hydrogen.address);
        assert.equal(balance.toString(), value.toString());

    })

    it('Successful Full Flow', async () => {
        console.log("\n==== Test: Successful Full Flow ====");
        const tokenA = await Token.new();
        const hydrogenA = await Hydrogen.new();

        const tokenB = await Token.new();
        const hydrogenB = await Hydrogen.new();

        initiator = accounts[0];
        counterparty = accounts[1];

        value = 1000;
        // Mint token to be used
        await tokenA.mint(value, {from: initiator, gas: 900000});
        balance = await tokenA.balanceOf.call(initiator)
        assert.equal(balance.toString(), value.toString());

        // Mint token to be used
        await tokenB.mint(value, {from: counterparty, gas: 900000});
        balance = await tokenB.balanceOf.call(counterparty)
        assert.equal(balance.toString(), value.toString());

        refundRef = "Refund";
        withdrawRef = "Withdraw";
        refundHash = Web3Utils.sha3(refundRef);
        withdrawHash = Web3Utils.sha3(withdrawRef);

        var tradeIdA;
        var tradeIdB;
        // 1. Initiate Trade Agreements
        await hydrogenA.InitiateTradeAgreement(tokenA.address, tokenB.address, initiator, counterparty, value, withdrawHash, refundHash, {from: initiator, gas: 900000}).then((tx) => {
            tradeIdA = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })

        await hydrogenB.InitiateTradeAgreement(tokenA.address, tokenB.address, initiator, counterparty, value, withdrawHash, refundHash, {from: counterparty, gas: 900000}).then((tx) => {
            tradeIdB = tx.logs[0].args._tradeId;
            console.log("InitiateTradeAgreement: Gas used = " + tx.receipt.gasUsed);
        })

        // Assert that both trade agreements created on both chains are equal by hash
        assert.equal(tradeIdA, tradeIdB);

        var expectedTradeId = Web3Utils.soliditySha3(initiator, counterparty, tokenA.address, tokenB.address, value, withdrawHash, refundHash);
        assert.equal(tradeIdA.toString(), expectedTradeId);

        // 2. Deposit funds to Hydrogen contract
        await tokenA.metadataTransfer(hydrogenA.address, value, tradeIdA, {from: initiator, gas: 900000}).then((tx) => {
            console.log("Deposit: Gas used = " + tx.receipt.gasUsed);
        });

        await tokenB.metadataTransfer(hydrogenB.address, value, tradeIdB, {from: counterparty, gas: 900000}).then((tx) => {
            console.log("Deposit: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful deduction from account
        balance = await tokenA.balanceOf.call(initiator);
        assert.equal(balance.toString(), "0");

        // Assert successful credit to contract
        balance = await tokenA.balanceOf.call(hydrogenA.address);
        assert.equal(balance.toString(), value.toString());


        // Assert successful deduction from account
        balance = await tokenB.balanceOf.call(counterparty);
        assert.equal(balance.toString(), "0");

        // Assert successful credit to contract
        balance = await tokenB.balanceOf.call(hydrogenB.address);
        assert.equal(balance.toString(), value.toString());

        // 3. Unlock funds for withdrawal
        await hydrogenA.UnlockFunds(tradeIdA, {from: initiator, gas: 900000}).then((tx) => {
            console.log("Unlock Funds: Gas used = " + tx.receipt.gasUsed);
        });

        // 4. Attempt withdrawal
        await hydrogenA.Withdraw(tradeIdA, withdrawRef, {from: counterparty, gas: 900000}).then((tx) => {
            console.log("Withdraw from chain A: Gas used = " + tx.receipt.gasUsed);
        });
        await hydrogenB.Withdraw(tradeIdB, withdrawRef, {from: initiator, gas: 900000}).then((tx) => {
            console.log("Withdraw from chain B: Gas used = " + tx.receipt.gasUsed);
        });

        // Assert successful withdrawal to counterparty
        balance = await tokenA.balanceOf.call(counterparty);
        assert.equal(balance.toString(), value.toString());
        balance = await tokenB.balanceOf.call(initiator);
        assert.equal(balance.toString(), value.toString());

        // Assert successful withdrawal from contract
        balance = await tokenA.balanceOf.call(hydrogenA.address);
        assert.equal(balance.toString(), "0");
        balance = await tokenB.balanceOf.call(hydrogenB.address);
        assert.equal(balance.toString(), "0");

    })
});
