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
    *   @description    persists the last submitted block of a chain being validated
    */
	struct BlockHeader {
        uint256 blockNumber;
		bytes32 blockHash;
		bytes32 prevBlockHash;
        bytes32 txRootHash;
        bytes32 receiptRootHash;
	}

    struct Metadata {
        address[] validators;
        mapping (address => bool) m_validators;
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
    mapping (bytes32 => mapping (bytes32 => bool)) public m_blockhashes;
	mapping (bytes32 => mapping (bytes32 => BlockHeader)) public m_blockheaders;
	mapping (bytes32 => mapping (bytes32 => Metadata)) public m_blockmetadata;
    mapping (bytes32 => bytes32[]) public heads;

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
            require(!m_blockhashes[_chainId][_genesisBlockHash], "Chain already exists with identical genesis");
        } else {
            chains[_chainId] = true;
            ion.addChain(_storeAddr, _chainId);
        }

        addGenesisBlock(_chainId, _validators, _genesisBlockHash);
    }

	/*
    * SubmitBlock
    * param: _chainId (bytes32) Unique id of chain submitting block from
    * param: _rlpBlockHeader (bytes) RLP-encoded byte array of the block header from other chain including all proposal and commit seals
    * param: _storeAddr (address) Address of block store contract to store block to
    *
    * Submission of block headers from another chain. Signatures held in the extraData field of _rlpSignedBlockHeader is recovered
    * and if valid the block is persisted as BlockHeader structs defined above.
    */
    function SubmitBlock(bytes32 _chainId, bytes _rlpUnsignedBlockHeader, bytes _rlpSignedBlockHeader, bytes _commitSeals, address _storageAddr) onlyRegisteredChains(_chainId) public {
        RLP.RLPItem[] memory header = _rlpSignedBlockHeader.toRLPItem().toList();

        // Check the parent hash is the same as the previous block submitted
		bytes32 parentBlockHash = SolUtils.BytesToBytes32(header[0].toBytes(), 1);
		require(m_blockhashes[_chainId][parentBlockHash], "Not child of previous block!");

        // Verify that validator and sealers are correct
        require(checkSignature(_chainId, header[12].toData(), keccak256(_rlpUnsignedBlockHeader), parentBlockHash), "Signer is not validator");
        require(checkSeals(_chainId, _commitSeals, _rlpSignedBlockHeader, parentBlockHash), "Sealer(s) not valid");

        // Append new block to the struct
        addValidators(_chainId, header[12].toData(), keccak256(_rlpSignedBlockHeader), parentBlockHash);
        storeBlock(_chainId, keccak256(_rlpSignedBlockHeader), parentBlockHash, SolUtils.BytesToBytes32(header[4].toBytes(), 1), SolUtils.BytesToBytes32(header[5].toBytes(), 1), header[8].toUint(), _rlpSignedBlockHeader, _storageAddr);

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

        Metadata storage metadata = m_blockmetadata[_chainId][_genesisBlockHash];
        metadata.validators = _validators;

        // Append validators and vote threshold
        for (uint256 i = 0; i < _validators.length; i++) {
            metadata.m_validators[_validators[i]] = true;
        }
        metadata.threshold = 2*(_validators.length/3) + 1;

        m_blockhashes[_chainId][_genesisBlockHash] = true;

        emit GenesisCreated(_chainId, _genesisBlockHash);
    }

    /*
    * checkSignature
    * param: _chainId (bytes32) Unique id of interoperating chain
    * param: _extraData (bytes) Byte array of the extra data containing signature
    * param: _rlpBlockHeader (bytes) Byte array of RLP encoded unsigned block header
    * param: _parentBlockHash (bytes32) Parent block hash of current block being checked
    *
    * Checks that the submitted block has actually been signed, recovers the signer and checks if they are validator in
    * parent block
    */
    function checkSignature(bytes32 _chainId, bytes _extraData, bytes32 _hash, bytes32 _parentBlockHash) internal view returns (bool) {
        // Retrieve Istanbul Extra
        bytes memory istanbulExtra = new bytes(_extraData.length - 32);
        SolUtils.BytesToBytes(istanbulExtra, _extraData, 32);

        RLP.RLPItem[] memory signature = istanbulExtra.toRLPItem().toList();

        bytes memory extraDataSig = new bytes(65);
        SolUtils.BytesToBytes(extraDataSig, signature[1].toBytes(), signature[1].toBytes().length-65);

        // Recover the signature
        address sigAddr = ECVerify.ecrecovery(keccak256(_hash), extraDataSig);
        Metadata storage parentMetadata = m_blockmetadata[_chainId][_parentBlockHash];

        // Check if signature is a validator that exists in previous block
		return parentMetadata.m_validators[sigAddr];
    }

    /*
    * checkSeals
    * param: _chainId (bytes32) Unique id of interoperating chain
    * param: _seals (bytes) RLP encoded list of 65 byte seals
    * param: _rlpBlock (bytes) Byte array of RLP encoded unsigned block header
    * param: _parentBlockHash (bytes32) Parent block hash of current block being checked
    *
    * Checks that the submitted block has enough seals to be considered valid as per the IBFT Soma rules
    */
    function checkSeals(bytes32 _chainId, bytes _seals, bytes _rlpBlock, bytes32 _parentBlockHash) internal view returns (bool) {
        bytes32 signedHash = keccak256(abi.encodePacked(keccak256(_rlpBlock), 0x02));
        Metadata storage parentMetadata = m_blockmetadata[_chainId][_parentBlockHash];
        uint256 validSeals = 0;

        // Check if signature is a validator that exists in previous block
        RLP.RLPItem[] memory seals = _seals.toRLPItem().toList();
        for (uint i = 0; i < seals.length; i++) {
            
            // Recover the signature
            address sigAddr = ECVerify.ecrecovery(signedHash, seals[i].toData());
            if (!parentMetadata.m_validators[sigAddr]) {
                return false;
            }
            validSeals++;
        }

        if (validSeals < parentMetadata.threshold) {
            return false;
        }

		return true;
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
    function addValidators(bytes32 _chainId, bytes _extraData, bytes32 _blockHash, bytes32 _parentBlockHash) internal {
        // Metadata storage parentMetadata = m_blockmetadata[_chainId][_parentBlockHash];
        Metadata storage metadata = m_blockmetadata[_chainId][_blockHash];

        address[] storage newValidators = metadata.validators;

        // Retrieve Istanbul Extra
        bytes memory rlpIstanbulExtra = new bytes(_extraData.length - 32);
        SolUtils.BytesToBytes(rlpIstanbulExtra, _extraData, 32);

        RLP.RLPItem[] memory istanbulExtra = rlpIstanbulExtra.toRLPItem().toList();
        RLP.RLPItem[] memory decodedExtra = istanbulExtra[0].toBytes().toRLPItem().toList();

        for (uint i = 0; i < decodedExtra.length; i++) {
            address validator = decodedExtra[i].toAddress();
            newValidators.push(validator);
            metadata.m_validators[validator] = true;
        }

        metadata.validators = newValidators;
        metadata.threshold = 2*(newValidators.length/3) + 1;

    }

    /*
    * storeBlock
    * param: _chainId (bytes32) Unique id of interoperating chain
    * param: _hash (address) Byte array of the extra data containing signature
    * param: _parentHash (bytes32) Current block hash being checked
    * param: _txRootHash (bytes32) Parent block hash of current block being checked
    * param: _receiptRootHash (bytes32) Parent block hash of current block being checked
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
        bytes32 _txRootHash,
        bytes32 _receiptRootHash,
        uint256 _height,
        bytes _rlpBlockHeader,
        address _storageAddr
    ) internal {
        m_blockhashes[_chainId][_hash] = true;

        BlockHeader storage header = m_blockheaders[_chainId][_hash];
        header.blockNumber = _height;
        header.blockHash = _hash;
        header.prevBlockHash = _parentHash;
        header.txRootHash = _txRootHash;
        header.receiptRootHash = _receiptRootHash;

        // Add block to Ion
        ion.storeBlock(_storageAddr, _chainId, _rlpBlockHeader);
    }

    function getValidators(bytes32 _chainId, bytes32 _blockHash) public view returns (address[]) {
        return m_blockmetadata[_chainId][_blockHash].validators;
    }

}
