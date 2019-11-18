pragma solidity ^0.5.12;

import "../contracts/libraries/PatriciaTrie.sol";

contract PatriciaTrieTest {
    function testVerify(bytes memory _value, bytes memory _parentNodes, bytes memory _path, bytes32 _root) public view returns (bool) {
        return PatriciaTrie.verifyProof(_value, _parentNodes, _path, _root);
    }
}