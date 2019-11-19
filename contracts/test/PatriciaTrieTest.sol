pragma solidity ^0.5.12;

import "../libraries/PatriciaTrie.sol";

contract PatriciaTrieTest {
    event Result(bool result);
    function testVerify(bytes memory _value, bytes memory _parentNodes, bytes memory _path, bytes32 _root) public returns (bool) {
        bool result = PatriciaTrie.verifyProof(_value, _parentNodes, _path, _root);
        emit Result(result);
        return result;
    }
}