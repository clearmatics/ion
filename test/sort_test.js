const SortingContract = artifacts.require("ArraySortingTest");
const { keccak256, bufferToHex } = require('ethereumjs-util');
const Web3EthAbi = require('web3-eth-abi');

contract("Sorting library", async (accounts) => {
    let expectedSortedAccounts, sortContract;

    describe("Sort", async () => {
        before(async () => {
            sortContract = await SortingContract.new()
        })

        it("Sorts an even number of elements", async () => {
            // without the upper case some weird things happen
            expectedSortedAccounts = accounts.map(x => x.toUpperCase()).sort()
            sorted = await sortContract.testSort.call(accounts)
            
            sorted = sorted.map(x => x.toUpperCase())
            expect(sorted).to.eql(expectedSortedAccounts) // deep equality - assert.equal fails
        })

        it("Sorts an odd number of elements", async () => {
            // without the upper case some weird things happen
            expectedSortedAccounts = accounts.slice(0,5).map(x => x.toUpperCase()).sort()
            sorted = await sortContract.testSort.call(accounts.slice(0,5))
            
            sorted = sorted.map(x => x.toUpperCase())
            expect(sorted).to.eql(expectedSortedAccounts.slice(0,5)) // deep equality - assert.equal fails
        })

        it("Sorts and hash an even number of elements", async () => {
            expectedHash = bufferToHex(keccak256(Web3EthAbi.encodeParameter("address[]", accounts.map(x => x.toLowerCase()).sort())))

            hash = await sortContract.testSortAndHash.call(accounts)
            assert.equal(expectedHash, hash)
        })


        it("Sorts and hash an odd number of elements", async () => {
            expectedHash = bufferToHex(keccak256(Web3EthAbi.encodeParameter("address[]", accounts.slice(0,5).map(x => x.toLowerCase()).sort())))

            hash = await sortContract.testSortAndHash.call(accounts.slice(0,5))
            assert.equal(expectedHash, hash)
        })
    })

})