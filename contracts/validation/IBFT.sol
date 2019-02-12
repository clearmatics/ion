// Copyright (c) 2016-2019 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.23;

import "../libraries/ECVerify.sol";
import "../libraries/RLP.sol";
import "../libraries/SolidityUtils.sol";
import "../IonCompatible.sol";
import "../storage/BlockStore.sol";

/*
    Smart contract for validation of blocks that use the IBFT-Soma consensus algorithm
    Blocks must be submitted sequentially due to the voting mechanism of IBFT-Soma.
*/

contract IBFT is IonCompatible {
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using RLP for bytes;

    /*
    * @description    persists the last submitted block of a chain being validated
    */
	struct BlockHeader {
        uint256 blockNumber;
		bytes32 blockHash;
		bytes32 prevBlockHash;
        address[] validators;
        uint256 threshold;
	}

    event GenesisCreated(bytes32 chainId, bytes32 blockHash);
    event BlockSubmitted(bytes32 chainId, bytes32 blockHash);

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

    mapping (bytes32 => bool) public chains;
    mapping (bytes32 => bytes32) public m_chainHeads;
	mapping (bytes32 => mapping (bytes32 => BlockHeader)) public m_blockheaders;

	constructor (address _ionAddr) IonCompatible(_ionAddr) public {}

/* =====================================================================================================================

        Public Functions

   =====================================================================================================================
*/
    function register() public returns (bool) {
        ion.registerValidationModule();
        return true;
    }

    /*
    * RegisterChain
    * param: _chainId (bytes32) Unique id of another chain to interoperate with
    * param: _validators (address[]) Array containing the validators at the genesis block
    * param: _genesisHash (bytes32) Hash of the genesis block for the chain being registered with Ion
    * param: _storeAddr (address) Address of block store contract to register chain to
    *
    * Registers knowledge of the id of another interoperable chain requiring the genesis block metadata. Allows
    * the initialising of genesis blocks and their validator sets for chains. Multiple may be submitted and built upon
    * and is not opinionated on how they are used.
    */
    function RegisterChain(bytes32 _chainId, address[] _validators, bytes32 _genesisBlockHash, address _storeAddr) public {
        require(_chainId != ion.chainId(), "Cannot add this chain id to chain register");

        if (chains[_chainId]) {
            require(m_chainHeads[_chainId] == bytes32(0x0), "Chain already exists");
        } else {
            chains[_chainId] = true;
            ion.addChain(_storeAddr, _chainId);
        }

        addGenesisBlock(_chainId, _validators, _genesisBlockHash);
    }

	/*
    * SubmitBlock
    * param: _chainId (bytes32) Unique id of chain submitting block from
    * param: _rlpUnsignedBlockHeader (bytes) RLP-encoded byte array of the block header from IBFT-Soma chain containing only validator set in IstanbulExtra field
    * param: _rlpSignedBlockHeader (bytes) RLP-encoded byte array of the block header from other chain including all proposal seal in the IstanbulExtra field
    * param: _commitSeals (bytes) RLP-encoded commitment seals that are typically contained in the last element of the IstanbulExtra field
    * param: _storeAddr (address) Address of block store contract to store block to
    *
    * Submission of block headers from another chain. 
    */
    function SubmitBlock(bytes32 _chainId, bytes _rlpUnsignedBlockHeader, bytes _rlpSignedBlockHeader, bytes _commitSeals, address _storageAddr) onlyRegisteredChains(_chainId) public {
        RLP.RLPItem[] memory header = _rlpSignedBlockHeader.toRLPItem().toList();

        // Check the parent hash is the same as the previous block submitted
		bytes32 parentBlockHash = SolUtils.BytesToBytes32(header[0].toBytes(), 1);
		require(m_chainHeads[_chainId] == parentBlockHash, "Not child of previous block!");

        // Verify that validator and sealers are correct
        require(checkSignature(_chainId, header[12].toData(), keccak256(_rlpUnsignedBlockHeader), parentBlockHash), "Signer is not validator");
        require(checkSeals(_chainId, _commitSeals, _rlpSignedBlockHeader, parentBlockHash), "Sealer(s) not valid");

        // Append new block to the struct
        addValidators(_chainId, header[12].toData(), keccak256(_rlpSignedBlockHeader));
        storeBlock(_chainId, keccak256(_rlpSignedBlockHeader), parentBlockHash, header[8].toUint(), _rlpSignedBlockHeader, _storageAddr);

        emit BlockSubmitted(_chainId, keccak256(_rlpSignedBlockHeader));
    }


/* =====================================================================================================================

        Internal Functions

   =====================================================================================================================
*/

    /*
    * addGenesisBlock
    * param: _chainId (bytes32) Unique id of another chain to interoperate with
    * param: _validators (address[]) Array containing the validators at the genesis block
    * param: _genesisHash (bytes32) Hash of the genesis block for the chain being registered with Ion
    *
    * Adds a genesis block with the validators and other metadata for this genesis block
    */
    function addGenesisBlock(bytes32 _chainId, address[] _validators, bytes32 _genesisBlockHash) internal {
        BlockHeader storage header = m_blockheaders[_chainId][_genesisBlockHash];
        header.blockNumber = 0;
        header.blockHash = _genesisBlockHash;
        header.validators = _validators;
        header.threshold = 2*(_validators.length/3) + 1;

        m_chainHeads[_chainId] = _genesisBlockHash;
        emit GenesisCreated(_chainId, _genesisBlockHash);
    }

    /*
    * checkSignature
    * param: _chainId (bytes32) Unique id of interoperating chain
    * param: _extraData (bytes) Byte array of the extra data containing signature
    * param: _hash (bytes32) Hash of the unsigned block header
    * param: _parentBlockHash (bytes32) Parent block hash of current block being checked
    *
    * Checks that the submitted block has actually been signed, recovers the signer and checks if they are validator in
    * parent block
    */
    function checkSignature(bytes32 _chainId, bytes _extraData, bytes32 _hash, bytes32 _parentBlockHash) internal view returns (bool) {
        // Retrieve Istanbul Extra Data
        bytes memory istanbulExtra = new bytes(_extraData.length - 32);
        SolUtils.BytesToBytes(istanbulExtra, _extraData, 32);

        RLP.RLPItem[] memory signature = istanbulExtra.toRLPItem().toList();

        bytes memory extraDataSig = new bytes(65);
        SolUtils.BytesToBytes(extraDataSig, signature[1].toBytes(), signature[1].toBytes().length-65);

        // Recover the signature
        address sigAddr = ECVerify.ecrecovery(keccak256(_hash), extraDataSig);
        BlockHeader storage parentBlock = m_blockheaders[_chainId][_parentBlockHash];

        // Check if signature is a validator that exists in previous block
		return isValidator(parentBlock.validators, sigAddr);
    }

    /*
    * checkSeals
    * param: _chainId (bytes32) Unique id of interoperating chain
    * param: _seals (bytes) RLP-encoded list of 65 byte seals
    * param: _rlpBlock (bytes) Byte array of RLP encoded unsigned block header
    * param: _parentBlockHash (bytes32) Parent block hash of current block being checked
    *
    * Checks that the submitted block has enough seals to be considered valid as per the IBFT Soma rules
    */
    function checkSeals(bytes32 _chainId, bytes _seals, bytes _rlpBlock, bytes32 _parentBlockHash) internal view returns (bool) {
        bytes32 signedHash = keccak256(abi.encodePacked(keccak256(_rlpBlock), 0x02));
        BlockHeader storage parentBlock = m_blockheaders[_chainId][_parentBlockHash];
        uint256 validSeals = 0;

        // Check if signature is a validator that exists in previous block
        RLP.RLPItem[] memory seals = _seals.toRLPItem().toList();
        for (uint i = 0; i < seals.length; i++) {
            // Recover the signature
            address sigAddr = ECVerify.ecrecovery(signedHash, seals[i].toData());
            if (!isValidator(parentBlock.validators, sigAddr))
                return false;
            validSeals++;
        }

        if (validSeals < parentBlock.threshold)
            return false;

		return true;
    }

    function isValidator(address[] _validators, address _validator) internal pure returns (bool) {
        for (uint i = 0; i < _validators.length; i++) {
            if (_validator == _validators[i])
                return true;
        }
        return false;
    }

    /*
    * addValidators
    * param: _chainId (bytes32) Unique id of interoperating chain
    * param: _extraData (bytes) Byte array of the extra data containing signature
    * param: _blockHash (bytes32) Current block hash being checked
    * param: _parentBlockHash (bytes32) Parent block hash of current block being checked
    *
    * Updates the validators from the RLP encoded extradata
    */
    function addValidators(bytes32 _chainId, bytes _extraData, bytes32 _blockHash) internal {
        BlockHeader storage newBlock = m_blockheaders[_chainId][_blockHash];

        // Retrieve Istanbul Extra Data
        bytes memory rlpIstanbulExtra = new bytes(_extraData.length - 32);
        SolUtils.BytesToBytes(rlpIstanbulExtra, _extraData, 32);

        RLP.RLPItem[] memory istanbulExtra = rlpIstanbulExtra.toRLPItem().toList();
        RLP.RLPItem[] memory decodedExtra = istanbulExtra[0].toBytes().toRLPItem().toList();

        for (uint i = 0; i < decodedExtra.length; i++) {
            address validator = decodedExtra[i].toAddress();
            newBlock.validators.push(validator);
        }

        newBlock.threshold = 2*(newBlock.validators.length/3) + 1;
    }

    /*
    * storeBlock
    * param: _chainId (bytes32) Unique id of interoperating chain
    * param: _hash (address) Byte array of the extra data containing signature
    * param: _parentHash (bytes32) Current block hash being checked
    * param: _height (bytes32) Parent block hash of current block being checked
    * param: _rlpBlockHeader (bytes32) Parent block hash of current block being checked
    * param: _storageAddr (bytes32) Parent block hash of current block being checked
    *
    * Takes the submitted block to propagate to the storage contract.
    */
    function storeBlock(
        bytes32 _chainId,
        bytes32 _hash,
        bytes32 _parentHash,
        uint256 _height,
        bytes _rlpBlockHeader,
        address _storageAddr
    ) internal {
        m_chainHeads[_chainId] = _hash;

        BlockHeader storage header = m_blockheaders[_chainId][_hash];
        header.blockNumber = _height;
        header.blockHash = _hash;
        header.prevBlockHash = _parentHash;

        delete m_blockheaders[_chainId][_parentHash];

        // Add block to Ion
        ion.storeBlock(_storageAddr, _chainId, _rlpBlockHeader);
    }

    function getValidators(bytes32 _chainId) public view returns (address[]) {
        return m_blockheaders[_chainId][m_chainHeads[_chainId]].validators;
    }

}
