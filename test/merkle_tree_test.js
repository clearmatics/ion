const assert = require('assert');
const path = require('path');
const utils = require('./helpers/utils.js');          
const { MerkleTree } = require('./helpers/merkleTree.js');
const { keccak256, bufferToHex } = require('ethereumjs-util');

const MerkleTreeContract = artifacts.require("MerkleTreeTest");

contract("Merkle Tree", (accounts) => {

    const completeElements = [accounts[0], accounts[1], accounts[2], accounts[3]]
    const allButOneElements = [accounts[0], accounts[1], accounts[2]]
    const sparseElements = [accounts[0], accounts[1], accounts[2], accounts[3], accounts[4]]

    describe("Generate root", () => {
        it("Should correctly calculate the merkle root", async () => {
            const merkleTree = new MerkleTree(completeElements);
            const expectedRoot = merkleTree.getHexRoot();

            const merkleContract = await MerkleTreeContract.new()

            root = await merkleContract.testRoot.call(merkleTree.elements)
            // console.log(expectedRoot, root)
            assert.equal(root, expectedRoot)
        })

        it("Calculates the root with all but one leaves", async () => {
            const merkleTree = new MerkleTree(allButOneElements);
            const expectedRoot = merkleTree.getHexRoot();

            const merkleContract = await MerkleTreeContract.new()

            root = await merkleContract.testRoot.call(merkleTree.elements)
            // console.log(expectedRoot, root)
            assert.equal(root, expectedRoot)
        })

        it("Calculates the root with sparse leaves", async () => {
            const merkleTree = new MerkleTree(sparseElements);
            const expectedRoot = merkleTree.getHexRoot();

            const merkleContract = await MerkleTreeContract.new()

            root = await merkleContract.testRoot.call(merkleTree.elements)
            // console.log(expectedRoot, root)
            assert.equal(root, expectedRoot)
        })


        it("Calculates the root with one element", async () => {
            const merkleTree = new MerkleTree([accounts[0]]);
            const expectedRoot = merkleTree.getHexRoot();

            const merkleContract = await MerkleTreeContract.new()

            root = await merkleContract.testRoot.call(merkleTree.elements)
            // console.log(expectedRoot, root)
            assert.equal(root, expectedRoot)
        })


        it("Returns zero with no elements", async () => {
            const merkleTree = new MerkleTree([]);
            const expectedRoot = "0x0000000000000000000000000000000000000000000000000000000000000000";

            const merkleContract = await MerkleTreeContract.new()

            root = await merkleContract.testRoot.call(merkleTree.elements)
            // console.log(expectedRoot, root)
            assert.equal(root, expectedRoot)
        })
    })

    describe("Verify root", () => {

        it("Should return true for a valid Merkle proof", async () => {
            const merkleTree = new MerkleTree(completeElements);
            const expectedRoot = merkleTree.getHexRoot();

            const merkleContract = await MerkleTreeContract.new()
            const contractRoot = await merkleContract.testRoot.call(merkleTree.elements)

            assert.equal(contractRoot, expectedRoot)

            const proof = merkleTree.getHexProof(completeElements[0])
            const leaf = bufferToHex(keccak256(completeElements[0]))

            verification = await merkleContract.testVerify.call(proof, contractRoot, leaf)
            assert.equal(verification, true)
        })

        it('should return false for an invalid Merkle proof', async function () {
            const correctMerkleTree = new MerkleTree(completeElements);
            const correctRoot = correctMerkleTree.getHexRoot();
            const correctLeaf = bufferToHex(keccak256(completeElements[0]));

            const anotherMerkleTree = new MerkleTree(allButOneElements);
            const badProof = anotherMerkleTree.getHexProof(allButOneElements[0]);

            const merkleContract = await MerkleTreeContract.new()
            verification = await merkleContract.testVerify.call(badProof, correctRoot, correctLeaf)
            assert.equal(verification, false)
        });

    })
})