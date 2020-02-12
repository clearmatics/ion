pragma solidity ^0.5.12;

import "../libraries/MerkleTree.sol";

contract MerkleTreeTest {
    event Result(bytes32 result);
    event Verify(bool result);

    function testRoot(bytes32[] memory elements) public returns (bytes32) {
        bytes32 result = MerkleTree.generateRoot(elements);
        emit Result(result);
        return result;
    }

    function testVerify(bytes32[] memory proof, bytes32 root, bytes32 leaf) public returns (bool) {
        bool result = MerkleTree.verify(proof, root, leaf);
        emit Verify(result);
        return result;
    }
}