pragma solidity ^0.4.24;

import "../Ion.sol";
import "../storage/BlockStore.sol";

contract MockIon is Ion {

    constructor(bytes32 _id) Ion(_id) {}

    function registerValidationModule() public {
        require( isContract(msg.sender), "Caller address is not a valid contract. Please inherit the BlockStore contract for proper usage." );
        require( !m_registered_validation[msg.sender], "Validation module has already been registered." );

        m_registered_validation[msg.sender] = true;
        validation_modules.push(msg.sender);
    }

    function addChain(address _storageAddress, bytes32 _chainId) {
        BlockStore store = BlockStore(_storageAddress);
        store.addChain(_chainId);
    }

    function storeBlock(address _storageAddress, bytes32 _chainId, bytes32 _blockHash, bytes _blockBlob) public {
        BlockStore store = BlockStore(_storageAddress);
        store.addBlock(_chainId, _blockHash, _blockBlob);
    }
}