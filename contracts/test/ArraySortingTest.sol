pragma solidity ^0.5.12;

import "../libraries/SortArray.sol";

contract ArraySortingTest {
    event Root(bytes32 result);

    function testSort(address[] memory elements) public returns (address[] memory) {
        return SortArray.sortAddresses(elements);
    }

    function testSortAndHash(address[] memory elements) public returns (bytes32) {
        return keccak256(abi.encode(SortArray.sortAddresses(elements)));
    }
    
}