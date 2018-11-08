pragma solidity ^0.4.24;

import "../storage/BlockStore.sol";

/*
    Mock Block Store contract

    This mocking contract is used to simulate interactions and asserting certain return data from interaction via other
    contracts being tested. Use as a tool for testing ONLY.

    This is not an accurate representation of a block store contract and should not be used in any way as a
    representation of a block store contract. Please refer to BlockStore.sol and inherit functionality from that base
    contract and see EthereumStore.sol for more implementation details.

*/

contract MockStorage is BlockStore {

    constructor(address _ionAddr) BlockStore(_ionAddr) public {}

    event AddedBlock(bytes32 blockHash);
    function addBlock(bytes32 _chainId, bytes32 _blockHash, bytes _blockBlob) {
        emit AddedBlock(_blockHash);
    }
}
