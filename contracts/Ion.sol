// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.23;

import "./libraries/ECVerify.sol";
import "./libraries/RLP.sol";
import "./libraries/PatriciaTrie.sol";
import "./libraries/SolidityUtils.sol";
import "./storage/BlockStore.sol";

contract Ion {

    bytes32 public chainId;
    mapping (address => bool) public m_registered_validation;
    address[] public validation_modules;

    /*
    * Constructor
    * param: id (bytes32) Unique id to identify this chain that the contract is being deployed to.
    *
    * Supplied with a unique id to identify this chain to others that may interoperate with it.
    * The deployer must assert that the id is indeed public and that it is not already being used
    * by another chain
    */
    constructor(bytes32 _id) public {
        chainId = _id;
    }

    /*
    * onlyRegisteredValidation
    * param: _addr (address) Address of the Validation module being registered
    *
    * Modifier that checks if the provided chain id has been registered to this contract
    */
    modifier onlyRegisteredValidation() {
        require( isContract(msg.sender), "Caller address is not a valid contract. Please inherit the BlockStore contract for proper usage." );
        require( m_registered_validation[msg.sender], "Validation module is not registered");
        _;
    }

    // Pseudo-modifier returns boolean, used with different 'require's to input custom revert messages
    function isContract(address _addr) internal returns (bool) {
        uint size;
        assembly { size := extcodesize(_addr) }
        return (size > 0);
    }


    function registerValidationModule() public {
        require( isContract(msg.sender), "Caller address is not a valid contract. Please inherit the BlockStore contract for proper usage." );
        require( !m_registered_validation[msg.sender], "Validation module has already been registered." );

        m_registered_validation[msg.sender] = true;
        validation_modules.push(msg.sender);
    }

    function addChain(address _storageAddress, bytes32 _chainId) onlyRegisteredValidation public {
        BlockStore store = BlockStore(_storageAddress);
        store.addChain(_chainId);
    }

    /*
    * storeBlock
    * param:
    *
    */
    function storeBlock(address _storageAddress, bytes32 _chainId, bytes32 _blockHash, bytes _blockBlob) onlyRegisteredValidation public {
        require( isContract(_storageAddress), "Storage address provided is not contract.");
        BlockStore store = BlockStore(_storageAddress);

        store.addBlock(_chainId, _blockHash, _blockBlob);
    }
}