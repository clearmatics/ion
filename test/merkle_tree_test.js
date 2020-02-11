const assert = require('assert');
const path = require('path');
const utils = require('./helpers/utils.js');
const { MerkleTree } = require('./helpers/merkleTree.js');

const MerkleTreeContract = artifacts.require("MerkleTreeTest");

contract("Merkle Tree", (accounts) => {

    const completeElements = [accounts[0], accounts[1], accounts[2], accounts[3]]
    const sparseBinaryTree = [accounts[0], accounts[1], accounts[2]]

    describe("Generate root", () => {
        it("Should correctly calculate the merkle root", async () => {
            const merkleTree = new MerkleTree(completeElements);
            const root = merkleTree.getHexRoot();
            console.log(root)

            const merkleContract = await MerkleTreeContract.new()

            receipt = await merkleContract.testRoot(merkleTree.elements)
            console.log(receipt)

            assert(false);
        })
    })
})