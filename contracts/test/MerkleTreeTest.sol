pragma solidity ^0.5.12;

import "../libraries/MerkleTree.sol";

contract MerkleTreeTest {
    event Result(bytes32 result);

    function testRoot(bytes32[] memory elements) public returns (bytes32) {
        bytes32 result = MerkleTree.generateRoot(elements);
        emit Result(result);
        return result;
    }
}