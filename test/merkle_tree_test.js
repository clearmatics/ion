const assert = require('assert');
const path = require('path');
const utils = require('./helpers/utils.js');
const { MerkleTree } = require('./helpers/merkleTree.js');
const { keccak256, bufferToHex } = require('ethereumjs-util');
require('chai')
    .use(require('chai-as-promised'))
    .should();

const MerkleTreeContract = artifacts.require("MerkleTreeTest");

contract("Merkle Tree", (accounts) => {

    const completeElements = [accounts[0], accounts[1], accounts[2], accounts[3]]
    const sparseBinaryTree = [accounts[0], accounts[1], accounts[2]]
    const largeSparse = [accounts[0], accounts[1], accounts[2], accounts[3], accounts[4], accounts[5], accounts[6], accounts[7]];

    describe("Generate root", () => {
        it("Should correctly calculate the merkle root", async () => {
            const merkleTree = new MerkleTree(completeElements);
            const expectedRoot = merkleTree.getHexRoot();

            const merkleContract = await MerkleTreeContract.new()

            root = await merkleContract.testRoot.call(merkleTree.elements)

            assert.equal(root, expectedRoot);
        })

        it("Calculates root from sparse leaves", async () => {
            const merkleTree = new MerkleTree(sparseBinaryTree);
            const expectedRoot = merkleTree.getHexRoot();

            const merkleContract = await MerkleTreeContract.new()

            root = await merkleContract.testRoot.call(merkleTree.elements)

            assert.equal(root, expectedRoot);
        })

        it("Calculates root from many sparse leaves", async () => {
            const merkleTree = new MerkleTree(largeSparse);
            const expectedRoot = merkleTree.getHexRoot();

            const merkleContract = await MerkleTreeContract.new()

            root = await merkleContract.testRoot.call(merkleTree.elements)
            assert.equal(root, expectedRoot);
        })

        it("Returns single leaf without hashing", async () => {
            const merkleTree = new MerkleTree(sparseBinaryTree);
            const expectedRoot = merkleTree.elements[0].toString('hex');

            const merkleContract = await MerkleTreeContract.new()

            root = await merkleContract.testRoot.call(merkleTree.elements.slice(0,1))

            assert.equal(root, "0x" + expectedRoot);
        })

        it("Reverts on empty list", async () => {
            const merkleTree = new MerkleTree(sparseBinaryTree);
            const expectedRoot = merkleTree.elements[0].toString('hex');

            const merkleContract = await MerkleTreeContract.new()

            root = await merkleContract.testRoot.call().should.be.rejected;
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

                const anotherMerkleTree = new MerkleTree(sparseBinaryTree);
                const badProof = anotherMerkleTree.getHexProof(sparseBinaryTree[0]);

                const merkleContract = await MerkleTreeContract.new()
                verification = await merkleContract.testVerify.call(badProof, correctRoot, correctLeaf)
                assert.equal(verification, false)
            });

        })
})