// Copyright (c) 2020 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.5.12;

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
        bytes32 validatorsHash; // hash of the sorted validators list
        uint256 votingThreshold; // threshold of the voting power needed to consider a block valid 
    }

    event GenesisCreated(bytes32 chainId, bytes32 blockHash);
    event BlockSubmitted(bytes32 chainId, bytes32 blockHash);
    event Header(bytes header, bytes header2); 

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
    modifier onlyRegisteredChains(bytes32 id) {
        require(supportedChains[id], "Chain is not registered");
        _;
    }

    // initialize ION hub contract
    constructor (address ionAddr) IonCompatible(ionAddr) public {}

/* =====================================================================================================================

        View Functions

   =====================================================================================================================
*/

    /*
    * getValidatorsRoot
    * 
    * @param: chainId (bytes32) id of the chain of the block containing the validators hash
    * @param: blockHash (bytes32) hash of the block containing the validators hash
    *
    * @description: Returns the validators hash of a specific chain head
    */
    function getValidatorsRoot(bytes32 chainId, bytes32 blockHash) external view returns (bytes32) {
        return id_chainHeaders[chainId][blockHash].validatorsHash;
    }

/* =====================================================================================================================

        External Functions

   =====================================================================================================================
*/
    
    /*
    * RegisterValidationModule 
    * 
    * @description: Register this validation module to ION hub in order to submit blocks to it
    */
    function RegisterValidationModule() public returns (bool) {
        ion.registerValidationModule();
        return true;
    }

    /*
    * RegisterChain
    * 
    * @param: chainId (bytes32) Unique id of another chain to interoperate with
    * @param: validators (address[]) Array containing the validators at the genesis block
    * @param: initialThreshold (uint256) Voting power threshold at the genesis block
    * @param: genesisHash (bytes32) Hash of the genesis block for the chain being registered with Ion
    *
    * @description: Adds a genesis block with the validators and other metadata for this genesis block
    */
    function RegisterChain(
        bytes32 chainId, 
        address[] calldata validators, 
        uint256 initialTreshold, 
        bytes32 genesisBlockHash, 
        address storeAddr
    ) external {

        require(chainId != ion.chainId(), "Cannot add this chain id to chain register");
        require(id_chainHeaders[chainId][genesisBlockHash].blockHash == bytes32(0), "This chain already exists");

        if (!supportedChains[chainId]) {
            // initialize the chain if it's not yet 
            supportedChains[chainId] = true;
        }

        // register this chain to ion hub
        ion.addChain(storeAddr, chainId);

        // store genesis block data needed to validate further blocks
        BlockHeader storage header = id_chainHeaders[chainId][genesisBlockHash];
        header.blockHash = genesisBlockHash;
        header.validatorsHash = keccak256(abi.encode(SortArray.sortAddresses(validators)));
        header.votingThreshold = initialTreshold;

        emit GenesisCreated(chainId, genesisBlockHash);
    }

    function SubmitBlock(
        bytes32 chainId, 
        bytes memory rlpUnsignedBlockHeader, 
        bytes memory rlpSignedBlockHeader, 
        bytes memory commitSeals, 
        address[] memory validatorsPreviousBlock,
        address storageAddr
    ) onlyRegisteredChains (chainId) public {

        // unmarshal rlp block header
        RLP.RLPItem[] memory header = rlpSignedBlockHeader.toRLPItem().toList();

        bytes32 expectedParentHash = SolUtils.BytesToBytes32(header[0].toBytes(), 1);

        // storage pointer to parent block
        BlockHeader storage parentBlock = id_chainHeaders[chainId][expectedParentHash];

        // check parent hash is correct
        require(parentBlock.blockHash == expectedParentHash, "Not child of previous block!");

        // verify the passed set of validators is the one that signed the previous block
        require(parentBlock.validatorsHash == keccak256(abi.encode(SortArray.sortAddresses(validatorsPreviousBlock))), "This is not the set of validators of the parent block");

        // use that set of validators to verify signatures of this block
        require(checkSignature(chainId, header[12].toData(), keccak256(rlpUnsignedBlockHeader), expectedParentHash, validatorsPreviousBlock), "Signer is not a validator");

        // and to verify sealers
        require(checkSeals(chainId, commitSeals, rlpSignedBlockHeader, expectedParentHash, validatorsPreviousBlock), "Sealer(s) not valid");

        // valid block - store it with the new set of validators
        storeBlock(chainId, header[12].toData(), expectedParentHash, rlpSignedBlockHeader, storageAddr);

        emit BlockSubmitted(chainId, keccak256(rlpSignedBlockHeader));
    }

/* =====================================================================================================================

        Internal Functions

   =====================================================================================================================
*/  

    function checkSignature(
        bytes32 chainId, 
        bytes memory extraData, 
        bytes32 hashUnsignedHeader, 
        bytes32 parentHash, 
        address[] memory validators
    ) internal returns (bool) {
        return true;
    }

    function checkSeals(
        bytes32 chainId, 
        bytes memory seals, 
        bytes memory rlpBlock, 
        bytes32 parentHash, 
        address[] memory validators
    ) internal view returns (bool) {
        return true;
    }

    function storeBlock(
        bytes32 chainId, 
        bytes memory extraData,
        bytes32 parentHash, 
        bytes memory rlpSignedBlockHeader,
        address storageAddr
    ) internal {

        // point to parent block and overwrite it with new one
        BlockHeader storage header = id_chainHeaders[chainId][parentHash];
        header.blockHash = keccak256(rlpSignedBlockHeader);
        header.parentHash = parentHash;
        (header.validatorsHash, header.votingThreshold) = calculateRootAndThreshold(extraData);

        // add block to storage module through ION hub 
        ion.storeBlock(storageAddr, chainId, rlpSignedBlockHeader);
    }

    function calculateRootAndThreshold(bytes memory extraData) internal returns (bytes32, uint256) {
        
        // retrieve committee and voting power from extraData header field 
        bytes memory rlpExtraData = new bytes(extraData.length - 32);
        SolUtils.BytesToBytes(rlpExtraData, extraData, 32);

        
        RLP.RLPItem[] memory extraDecoded = rlpExtraData.toRLPItem().toList()[0].toBytes().toRLPItem().toList();
    
        address[] memory newValidators = new address[](extraDecoded.length);

        for (uint i = 0; i < extraDecoded.length; i++) {
            newValidators[i] = extraDecoded[i].toAddress();
        }

        return (keccak256(abi.encode(SortArray.sortAddresses(newValidators))), 1);
    }

}