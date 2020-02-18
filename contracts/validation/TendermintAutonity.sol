// Copyright (c) 2020 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.5.2;

import "../libraries/ECVerify.sol";
import "../libraries/RLP.sol";
import "../libraries/SolidityUtils.sol";
import "../libraries/SortArray.sol";
import "../storage/BlockStore.sol";
import "../IonCompatible.sol";

/* 
    Smart contract for validation of blocks that use the Autonity implementation of the 
    Tendermint consensus algorithm.
*/

contract TendermintAutonity is IonCompatible {
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using RLP for bytes;

    /*
    * @description    persists the last submitted block of a chain being validated
    */
    struct BlockHeader {
        bytes32 blockHash;
        bytes32 parentHash;
        bytes32 committeeHash; // hash of the sorted committee list 
        uint256 votingThreshold; // threshold of the voting power needed to consider a block valid 
    }

    event GenesisCreated(bytes32 chainId, bytes32 blockHash);
    event BlockSubmitted(bytes32 chainId, bytes32 blockHash);

    // ids of the chains registered to this module
    mapping (bytes32 => bool) public supportedChains; 
    
    // chainID to blockHash to blockHeader of the last submitted block
    // by not allowing to build multiple chains for the same chainID one could 
    // easily deny the service
    mapping (bytes32 => mapping(bytes32 => BlockHeader)) public id_chainHeaders;

    /*
    * onlyRegisteredChains
    * param: _id (bytes32) Unique id of chain supplied to function
    *
    * Modifier that checks if the provided chain id has been registered to this contract
    */
    modifier onlyRegisteredChains(bytes32 _id) {
        require(chains[_id], "Chain is not registered");
        _;
    }

    // initialize ION hub contract
    constructor (address _ionAddress) IonCompatible(_ionAddr) public {}


/* =====================================================================================================================

        Public Functions

   =====================================================================================================================
*/
    
    // register this validation module to ION hub so that this can send blocks to be stored
    function Register() public returns (bool) {
        ion.registervalidationModule();
        return true
    }

    function RegisterChain(address[] memory committee, bytes32 genesisBlockHash) {
        
    }


}