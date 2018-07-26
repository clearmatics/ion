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

	address public owner;
    bytes32 public chainId;

    /*
    *   @description    persists the last submitted block of a chain being validated
    */
	struct BlockHeader {
        uint256 blockHeight;
		bytes32 latestHash; 
		bytes32 prevBlockHash;
	}

    mapping (bytes32 => uint256) public blockHeight;
    mapping (bytes32 => bool) public chains;
    mapping (bytes32 => mapping (bytes32 => bool)) public m_blockhashes;
	mapping (bytes32 => BlockHeader) public m_blockheaders;
	mapping (bytes32 => mapping (address => bool)) public m_validators;

	event broadcastSignature(address signer);
	event broadcastHash(bytes32 blockHash);

	/*
	*	@param _id		genesis block of the blockchain where the contract is deployed
	*/
	constructor (bytes32 _id) public {
		owner = msg.sender;
		chainId = _id;
	}


    /*
    * InitChain
    * param: chainId (bytes32) Unique id of another chain to interoperate with
    *
    * Supplied with an id of another chain, checks if this id already exists in the known set of ids
    * and adds it to the list of known chains.
    */
    function InitChain(bytes32 _id, address[] _validators, bytes32 _genesisHash) public {
        require( _id != chainId, "Cannot add this chain id to chain register" );
        require(!chains[_id], "Chain already exists" );
        chains[_id] = true;

        // Append validators
        for (uint256 i = 0; i < _validators.length; i++) {
            m_validators[_id][_validators[i]] = true;
    	}

        // m_blockhashes[_id][_genesisHash] = true;
		m_blockheaders[_id].blockHeight = 0;
		m_blockheaders[_id].latestHash = _genesisHash;
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
    function ValidateBlock(bytes32 _id, bytes _rlpBlockHeader, bytes _rlpSignedBlockHeader) public onlyRegisteredChains(_id) {
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
		require(m_blockheaders[_id].latestHash==_parentBlockHash, "Not child of previous block!");

        // Check the blockhash
        bytes32 _blockHash = keccak256(_rlpSignedBlockHeader);
        emit broadcastHash(_blockHash);

        recoverSignature(_id, signedHeader[12].toBytes(), _rlpBlockHeader);

        // Append the new block to the struct
		m_blockheaders[_id].blockHeight++;
		m_blockheaders[_id].latestHash = _blockHash;
		m_blockheaders[_id].prevBlockHash = _parentBlockHash;

        addBlockHashToChain(_id, _blockHash);

    }

    function recoverSignature(bytes32 _id, bytes signedHeader, bytes _rlpBlockHeader) internal {
        bytes memory extraDataSig = new bytes(65);
        uint256 length = signedHeader.length;
        SolUtils.BytesToBytes(extraDataSig, signedHeader, length-65);

        // Recover the signature of 
        address sigAddr = ECVerify.ecrecovery(keccak256(_rlpBlockHeader), extraDataSig);
		require(m_validators[_id][sigAddr]==true, "Signer not a validator!");

        emit broadcastSignature(sigAddr);
    }

    /*
    * @description      when a block is submitted the root hash must be added to a mapping of chains to hashes
    * @param _chainId   unique identifier of the chain from which the block hails     
    * @param _hash      root hash of the block being added
    */
    function addBlockHashToChain(bytes32 _chainId, bytes32 _hash) internal {
        m_blockhashes[_chainId][_hash] = true;
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
