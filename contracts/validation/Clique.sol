// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.23;

import "../libraries/ECVerify.sol";
import "../libraries/RLP.sol";
import "../libraries/SolidityUtils.sol";
import "../IonCompatible.sol";
import "../storage/BlockStore.sol";

/*
    Smart contract for validation of blocks that use the Clique PoA consensus algorithm
    Blocks must be submitted sequentially due to the voting mechanism of Clique.
*/

contract Clique is IonCompatible {
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
        mapping (address => uint256) m_proposals;
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
        require( _chainId != ion.chainId(), "Cannot add this chain id to chain register" );

        if (chains[_chainId]) {
            require( !m_blockhashes[_chainId][_genesisBlockHash], "Chain already exists with identical genesis" );
        } else {
            chains[_chainId] = true;
            ion.addChain(_storeAddr, _chainId);
        }

        addGenesisBlock(_chainId, _validators, _genesisBlockHash, _storeAddr);
    }

	/*
    * SubmitBlock
    * param: _chainId (bytes32) Unique id of chain submitting block from
    * param: _rlpBlockHeader (bytes) RLP-encoded byte array of the block header from other chain without the signature in extraData
    * param: _rlpSignedBlockHeader (bytes) RLP-encoded byte array of the block header from other chain with the signature in extraData
    * param: _storeAddr (address) Address of block store contract to store block to
    *
    * Submission of block headers from another chain. Signatures held in the extraData field of _rlpSignedBlockHeader is recovered
    * and if valid the block is persisted as BlockHeader structs defined above.
    */
    function SubmitBlock(bytes32 _chainId, bytes _rlpBlockHeader, bytes _rlpSignedBlockHeader, address _storageAddr) onlyRegisteredChains(_chainId) public {
        RLP.RLPItem[] memory header = _rlpBlockHeader.toRLPItem().toList();
        RLP.RLPItem[] memory signedHeader = _rlpSignedBlockHeader.toRLPItem().toList();
        require( header.length == signedHeader.length, "Header properties length mismatch" );

        // Check header and signedHeader contain the same data
        for (uint256 i=0; i<signedHeader.length; i++) {
            // Skip extra data field
            if (i==12) {
                bytes memory extraData = new bytes(32);
                bytes memory extraDataSigned = new bytes(32);
                SolUtils.BytesToBytes(extraData, signedHeader[i].toBytes(), 2);
                SolUtils.BytesToBytes(extraDataSigned, header[i].toBytes(), 1);
                require(keccak256(extraDataSigned) == keccak256(extraData), "Header data doesn't match!");
            } else {
                require(keccak256(header[i].toBytes()) == keccak256(signedHeader[i].toBytes()), "Header data doesn't match!");
            }
        }

        // Check the parent hash is the same as the previous block submitted
		bytes32 parentBlockHash = SolUtils.BytesToBytes32(header[0].toBytes(), 1);
		require( m_blockhashes[_chainId][parentBlockHash], "Not child of previous block!" );
        require( checkSignature(_chainId, signedHeader[12].toBytes(), _rlpBlockHeader, parentBlockHash), "Signer is not validator" );

        // Append the new block to the struct
        addProposal(_chainId, SolUtils.BytesToAddress(header[2].toBytes(), 1), keccak256(_rlpSignedBlockHeader), parentBlockHash);
        storeBlock(_chainId, keccak256(_rlpSignedBlockHeader), parentBlockHash, SolUtils.BytesToBytes32(header[4].toBytes(), 1), SolUtils.BytesToBytes32(header[5].toBytes(), 1), header[8].toUint(), _rlpSignedBlockHeader, _storageAddr);
        shiftHead(_chainId, keccak256(_rlpSignedBlockHeader), parentBlockHash);

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
    * param: _storeAddr (address) Address of block store contract to register chain to
    *
    * Adds a genesis block with the validators and other metadata for this genesis block
    */
    function addGenesisBlock(bytes32 _chainId, address[] _validators, bytes32 _genesisBlockHash, address _storeAddr) internal {
        BlockHeader storage header = m_blockheaders[_chainId][_genesisBlockHash];
        header.blockNumber = 0;
        header.blockHash = _genesisBlockHash;

        Metadata storage metadata = m_blockmetadata[_chainId][_genesisBlockHash];
        metadata.validators = _validators;

        // Append validators and vote threshold
        for (uint256 i = 0; i < _validators.length; i++) {
            metadata.m_validators[_validators[i]] = true;
        }
        metadata.threshold = (_validators.length/2) + 1;

        m_blockhashes[_chainId][_genesisBlockHash] = true;
        shiftHead(_chainId, _genesisBlockHash, 0x0);

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
    function checkSignature(bytes32 _chainId, bytes _extraData, bytes _rlpBlockHeader, bytes32 _parentBlockHash) internal returns (bool) {
        bytes memory extraDataSig = new bytes(65);
        uint256 length = _extraData.length;
        SolUtils.BytesToBytes(extraDataSig, _extraData, length-65);

        // Recover the signature of 
        address sigAddr = ECVerify.ecrecovery(keccak256(_rlpBlockHeader), extraDataSig);
        Metadata storage parentMetadata = m_blockmetadata[_chainId][_parentBlockHash];

        // Check if signature is a validator that exists in previous block
		return parentMetadata.m_validators[sigAddr];
    }

    /*
    * addProposal
    * param: _chainId (bytes32) Unique id of interoperating chain
    * param: _candidate (address) Byte array of the extra data containing signature
    * param: _blockHash (bytes32) Current block hash being checked
    * param: _parentBlockHash (bytes32) Parent block hash of current block being checked
    *
    * Modifies the proposal/validator set via votes collated from the block. Checks parent block for latest state.
    */
    function addProposal(bytes32 _chainId, address _candidate, bytes32 _blockHash, bytes32 _parentBlockHash) internal {
        Metadata storage parentMetadata = m_blockmetadata[_chainId][_parentBlockHash];
        Metadata storage metadata = m_blockmetadata[_chainId][_blockHash];

        if (_candidate != 0x0) {
            uint newVoteCount;
            uint newThreshold = metadata.threshold;
            address[] storage newValidators = metadata.validators;

            // If votes pass threshold, add validator if exists, remove validator if not exists. Else metadata equal to parent
            if ( (parentMetadata.m_proposals[_candidate] + 1) >= parentMetadata.threshold && !parentMetadata.m_validators[_candidate]) {
                newVoteCount = 0;

                for (uint i = 0; i < parentMetadata.validators.length; i++) {
                    newValidators.push(parentMetadata.validators[i]);
                }
                newValidators.push(_candidate);
            } else if ( (parentMetadata.m_proposals[_candidate] + 1) >= parentMetadata.threshold && parentMetadata.m_validators[_candidate]) {
                newVoteCount = 0;

                for (uint j = 0; j < parentMetadata.validators.length; j++) {
                    if (parentMetadata.validators[j] != _candidate) {
                        newValidators.push(parentMetadata.validators[j]);
                    }
                }
            } else {
                newVoteCount = parentMetadata.m_proposals[_candidate] + 1;

                for (uint k = 0; k < parentMetadata.validators.length; k++) {
                    newValidators.push(parentMetadata.validators[k]);
                }
            }

            metadata.m_proposals[_candidate] = newVoteCount;
            newThreshold = (newValidators.length/2) + 1;

            for (uint vi = 0; vi < newValidators.length; vi++) {
                metadata.m_validators[newValidators[vi]] = true;
                if (newValidators[vi] != _candidate) {
                    metadata.m_proposals[newValidators[vi]] = parentMetadata.m_proposals[newValidators[vi]];
                }
            }
        } else {
            // If no vote, set current block metadata equal to parent block
            metadata.validators = parentMetadata.validators;
            metadata.threshold = parentMetadata.threshold;

            for (uint pi = 0; pi < parentMetadata.validators.length; pi++) {
                metadata.m_validators[parentMetadata.validators[pi]] = true;
                metadata.m_proposals[parentMetadata.validators[pi]] = parentMetadata.m_proposals[parentMetadata.validators[pi]];
            }
        }
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
        ion.storeBlock(_storageAddr, _chainId, _hash, _rlpBlockHeader);
    }

    /*
    * shiftHead
    * param: _chainId (bytes32) Unique id of chain
    * param: _childHash (bytes32) New block hash
    * param: _parentHash (bytes32) Previous block hash
    *
    * Updates set of current open chain heads per chain. Open chain heads are blocks that do not have a child that can
    * be built upon.
    */
    function shiftHead(bytes32 _chainId, bytes32 _childHash, bytes32 _parentHash) public {
        int index = -1;
        bytes32[] storage chainHeads = heads[_chainId];

        // Check if parent hash is an open head and replace with child
        for (uint i = 0; i < chainHeads.length; i++) {
            if (chainHeads[i] == _parentHash) {
                index = int(i);

                delete chainHeads[uint(index)];
                chainHeads[uint(index)] = _childHash;

                return;
            }
        }

        // If parent is not an open head, child is, so append to heads
        chainHeads.push(_childHash);
    }

    function getValidators(bytes32 _chainId, bytes32 _blockHash) constant returns (address[]) {
        return m_blockmetadata[_chainId][_blockHash].validators;
    }

    function getProposal(bytes32 _chainId, bytes32 _blockHash, address _candidate) constant returns (uint256) {
        return m_blockmetadata[_chainId][_blockHash].m_proposals[_candidate];
    }
}
