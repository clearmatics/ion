// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.23;

import "./libraries/ECVerify.sol";
import "./libraries/RLP.sol";
import "./libraries/SolidityUtils.sol";

contract Validation {
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using RLP for bytes;

    bytes32 public chainId;

    /*
    *   @description    persists the last submitted block of a chain being validated
    */
	struct BlockHeader {
        uint256 blockHeight;
		bytes32 latestHash; 
		bytes32 prevBlockHash;
        bytes32 txRootHash;
        bytes32 receiptRootHash;
	}

    mapping (bytes32 => bool) public chains;
    mapping (bytes32 => mapping (bytes32 => bool)) public m_blockhashes;
	mapping (bytes32 => mapping (bytes32 => BlockHeader)) public m_blockheaders;
	mapping (bytes32 => mapping (address => bool)) public m_validators;

	event broadcastSignature(address signer);
	event broadcastHash(bytes32 blockHash);

	/*
	*	@param _id		genesis block of the blockchain where the contract is deployed
	*/
	constructor (bytes32 _id) public {
		chainId = _id;
	}


    /*
    * RegisterChain
    * param: chainId (bytes32) Unique id of another chain to interoperate with
    *
    * Supplied with an id of another chain, checks if this id already exists in the known set of ids
    * and adds it to the list of known chains.
    */
    function RegisterChain(bytes32 _id, address[] _validators, bytes32 _genesisHash) public {
        require( _id != chainId, "Cannot add this chain id to chain register" );
        require(!chains[_id], "Chain already exists" );
        chains[_id] = true;

        // Append validators
        for (uint256 i = 0; i < _validators.length; i++) {
            m_validators[_id][_validators[i]] = true;
    	}

		m_blockheaders[_id][_genesisHash].blockHeight = 0;
		m_blockhashes[_id][_genesisHash] = true;
    }

	/*
    * ValidateBlock
    * param: _id (bytes32) Unique id of chain submitting block from
    * param: _rlpBlockHeader (bytes) RLP-encoded byte array of the block header from other chain without the signature in extraData
    * param: _rlpSignedBlockHeader (bytes) RLP-encoded byte array of the block header from other chain with the signature in extraData
    *
    * Submission of block headers from another chain. Signatures held in the extraData field of _rlpSignedBlockHeader is recovered
    * and if valid the block is persisted as BlockHeader structs defined above.
    */
    function SubmitBlock(bytes32 _id, bytes _rlpBlockHeader, bytes _rlpSignedBlockHeader) onlyRegisteredChains(_id) public {
        RLP.RLPItem[] memory header = _rlpBlockHeader.toRLPItem().toList();
        RLP.RLPItem[] memory signedHeader = _rlpSignedBlockHeader.toRLPItem().toList();

        // Check header and signedHeader contain the same data
        for (uint256 i=0; i<signedHeader.length; i++) {
            // Skip extra data field
            if (i==12) {
                continue;
            } else{
                require(keccak256(header[i].toBytes())==keccak256(signedHeader[i].toBytes()), "Header data doesn't match!");
            }
        }

        // Check the parent hash is the same as the previous block submitted
		bytes32 _parentBlockHash = SolUtils.BytesToBytes32(header[0].toBytes(), 1);
		require(m_blockhashes[_id][_parentBlockHash], "Not child of previous block!");

        // Check the blockhash
        bytes32 _blockHash = keccak256(_rlpSignedBlockHeader);
        emit broadcastHash(_blockHash);

        recoverSignature(_id, signedHeader[12].toBytes(), _rlpBlockHeader);

        // Append the new block to the struct
		m_blockheaders[_id][_blockHash].blockHeight = header[8].toUint();
		m_blockheaders[_id][_blockHash].latestHash = _blockHash;
		m_blockheaders[_id][_blockHash].prevBlockHash = _parentBlockHash;
        m_blockheaders[_id][_blockHash].txRootHash = SolUtils.BytesToBytes32(header[4].toBytes(), 1);
        m_blockheaders[_id][_blockHash].receiptRootHash = SolUtils.BytesToBytes32(header[5].toBytes(), 1);

        addBlockHashToChain(_id, _blockHash);

    }

    function recoverSignature(bytes32 _id, bytes signedHeader, bytes _rlpBlockHeader) internal {
        bytes memory extraDataSig = new bytes(65);
        uint256 length = signedHeader.length;
        SolUtils.BytesToBytes(extraDataSig, signedHeader, length-65);

        // Recover the signature of 
        address sigAddr = ECVerify.ecrecovery(keccak256(_rlpBlockHeader), extraDataSig);
		require(m_validators[_id][sigAddr], "Signer not a validator!");

        emit broadcastSignature(sigAddr);
    }

    /*
    * @description      when a block is submitted the root hash must be added to a mapping of chains to hashes
    * @param _id        unique identifier of the chain from which the block hails     
    * @param _hash      root hash of the block being added
    */
    function addBlockHashToChain(bytes32 _id, bytes32 _hash) internal {
        m_blockhashes[_id][_hash] = true;
    }

    /*
    * @description  returns the transaction root hash of a specific block
    * @param _id    unique identifier of the chain from which the block hails     
    * @param _hash  root hash of the block being queried
    */
    function getTxRootHash(bytes32 _id, bytes32 _hash) public returns(bytes32) {
        return(m_blockheaders[_id][_hash].txRootHash);
    }

    /*
    * @description  returns the receipt root hash of a specific block
    * @param _id    unique identifier of the chain from which the block hails     
    * @param _hash  root hash of the block being queried
    */
    function getReceiptRootHash(bytes32 _id, bytes32 _hash) public returns(bytes32) {
        return(m_blockheaders[_id][_hash].receiptRootHash);
    }

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


}
