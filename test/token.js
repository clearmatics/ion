'use strict';

const BigNumber = web3.BigNumber;

const should = require('chai')
    .use(require('chai-as-promised'))
    .use(require('chai-bignumber')(BigNumber))
    .should();

const Token = artifacts.require("Token");

contract('Token', (accounts) => {
	let token;
    let token_owner;

    beforeEach(async function() {
		token = await Token.new();
    });

    it("should return the correct totalSupply after construction", async function()
    {
        const totalSupply = await token.totalSupply();
        totalSupply.should.be.bignumber.equal(new BigNumber(0));
    });


    it('should throw an error when trying to transfer to 0x0', async function()
    {
    	await token.mint(500);

    	await token.transfer(0x0, 100, {from:accounts[1]}).should.be.rejected;
    });

    it('should throw an error when trying to transfer more than balance', async function() {
		(await token.balanceOf(accounts[0], {from:accounts[0]})).should.be.bignumber.equal(new BigNumber(0));
		await token.mint(500, {from:accounts[0]});

		(await token.balanceOf(accounts[0], {from:accounts[0]})).should.be.bignumber.equal(new BigNumber(500));
    	await token.transfer(accounts[1], 501, {from:accounts[0]}).should.be.rejected;

    	await token.transfer(accounts[1], 500, {from:accounts[0]});
    	(await token.balanceOf(accounts[0], {from:accounts[0]})).should.be.bignumber.equal(new BigNumber(0));
    });
});
