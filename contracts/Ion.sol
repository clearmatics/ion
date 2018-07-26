// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.23;

import "./libraries/ECVerify.sol";
import "./libraries/RLP.sol";
import "./libraries/PatriciaTrie.sol";

contract Ion {
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using RLP for bytes;

    struct BlockHeader {
        uint256 blockHeight;
        bytes32 prevBlockHash;
        bytes32 txRootHash;
        bytes32 receiptRootHash;
    }


    address[] public validators;
    bytes32 public blockHash;
    bytes32 public chainId;
    // bytes32[] public chains;
    uint256 public blockHeight;

    mapping (bytes32 => bool) public chains;
    mapping (bytes32 => mapping (bytes32 => bool)) public m_blockhashes;
    mapping (bytes32 => mapping (bytes32 => BlockHeader)) public m_blockheaders;

    // input chainId to validators address returning bool
    mapping (bytes32 => mapping (address => bool)) public m_validators;


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


    enum ProofType { TX, RECEIPT, ROOTS }

    event VerifiedProof(bytes32 chainId, bytes32 blockHash, uint proofType);
    event broadcastSignature(address signer);
	event broadcastHash(bytes32 blockHash);
/*
========================================================================================================================

    Modifiers

========================================================================================================================
*/

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

    modifier onlyExistingBlocks(bytes32 _id, bytes32 _hash) {
        require(m_blockhashes[_id][_hash], "Block does not exist for chain");
        _;
    }

/*
========================================================================================================================

    Functions

========================================================================================================================
*/

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

        for (uint256 i = 0; i < _validators.length; i++) {
            m_validators[_id][_validators[i]] = true;
    	}

        m_blockhashes[_id][_genesisHash] = true;
		m_blockheaders[_id][_genesisHash].blockHeight = 0;
    }

    /*
    * SubmitBlock
    * param: _id (bytes32) Unique id of chain submitting block from
    * param: _rlpBlockHeader (bytes) RLP-encoded byte array of the block header from other chain without the signature in extraData
    * param: _rlpSignedBlockHeader (bytes) RLP-encoded byte array of the block header from other chain with the signature in extraData
    *
    * Submission of block headers from another chain. Signatures held in the extraData field of _rlpSignedBlockHeader is recovered
    * and if valid the block is persisted as BlockHeader structs defined above.
    */
    function SubmitBlock(bytes32 _id, bytes _rlpBlockHeader, bytes _rlpSignedBlockHeader) public onlyRegisteredChains(_id) {
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
		bytes32 _parentBlockHash = bytesToBytes32(header[0].toBytes(), 1);
		require(m_blockhashes[_id][_parentBlockHash]==true, "Not child of previous block!");

        // Check the blockhash
        bytes32 _blockHash = keccak256(_rlpSignedBlockHeader);
        emit broadcastHash(_blockHash);

        recoverSignature(_id, signedHeader[12].toBytes(), _rlpBlockHeader);

        // Append the new block to the struct
		blockHeight++;
		m_blockheaders[_id][_blockHash].blockHeight = blockHeight;
		m_blockheaders[_id][_blockHash].prevBlockHash = _parentBlockHash;
        m_blockheaders[_id][_blockHash].txRootHash = bytesToBytes32(header[4].toBytes(), 1);
        m_blockheaders[_id][_blockHash].receiptRootHash = bytesToBytes32(header[5].toBytes(), 1);

        addBlockHashToChain(_id, _blockHash);

    }

    function recoverSignature(bytes32 _id, bytes signedHeader, bytes _rlpBlockHeader) internal {
        bytes memory extraDataSig = new bytes(65);
        uint256 length = signedHeader.length;
        bytesToBytes(extraDataSig, signedHeader, length-65);

        // Recover the signature of 
        address sigAddr = ECVerify.ecrecovery(keccak256(_rlpBlockHeader), extraDataSig);
		require(m_validators[_id][sigAddr]==true, "Signer not a validator!");

        emit broadcastSignature(sigAddr);
    }

    /*
    * CheckTxProof
    * param: _id (bytes32) Unique id of chain submitting block from
    * param: _blockHash (bytes32) Block hash of block being submitted
    * param: _value (bytes) RLP-encoded transaction object array with fields defined as: https://github.com/ethereumjs/ethereumjs-tx/blob/0358fad36f6ebc2b8bea441f0187f0ff0d4ef2db/index.js#L50
    * param: _parentNodes (bytes) RLP-encoded array of all relevant nodes from root node to node to prove
    * param: _path (bytes) Byte array of the path to the node to be proved
    *
    * emits: VerifiedTxProof(chainId, blockHash, proofType)
    *        chainId: (bytes32) hash of the chain verifying proof against
    *        blockHash: (bytes32) hash of the block verifying proof against
    *        proofType: (uint) enum of proof type
    *
    * All data associated with the proof must be constructed and provided to this function. Modifiers restrict execution
    * of this function to only allow if the chain the proof is for is registered to this contract and if the block that
    * the proof is for has been submitted.
    */
    function CheckTxProof(
        bytes32 _id,
        bytes32 _blockHash,
        bytes _value,
        bytes _parentNodes,
        bytes _path
    )
        onlyRegisteredChains(_id)
        onlyExistingBlocks(_id, _blockHash)
        public
        returns (bool)
    {
        BlockHeader storage blockHeader = m_blockheaders[_id][_blockHash];
        assert( PatriciaTrie.verifyProof(_value, _parentNodes, _path, blockHeader.txRootHash) );

        emit VerifiedProof(_id, _blockHash, uint(ProofType.TX));
        return true;
    }

    /*
    * CheckReceiptProof
    * param: _id (bytes32) Unique id of chain submitting block from
    * param: _blockHash (bytes32) Block hash of block being submitted
    * param: _value (bytes) RLP-encoded receipt object array with fields defined as: https://github.com/ethereumjs/ethereumjs-tx/blob/0358fad36f6ebc2b8bea441f0187f0ff0d4ef2db/index.js#L50
    * param: _parentNodes (bytes) RLP-encoded array of all relevant nodes from root node to node to prove
    * param: _path (bytes) Byte array of the path to the node to be proved
    *
    * emits: VerifiedTxProof(chainId, blockHash, proofType)
    *        chainId: (bytes32) hash of the chain verifying proof against
    *        blockHash: (bytes32) hash of the block verifying proof against
    *        proofType: (uint) enum of proof type
    *
    * All data associated with the proof must be constructed and provided to this function. Modifiers restrict execution
    * of this function to only allow if the chain the proof is for is registered to this contract and if the block that
    * the proof is for has been submitted.
    */
    function CheckReceiptProof(
        bytes32 _id,
        bytes32 _blockHash,
        bytes _value,
        bytes _parentNodes,
        bytes _path
    )
    onlyRegisteredChains(_id)
    onlyExistingBlocks(_id, _blockHash)
    public
    returns (bool)
    {
        BlockHeader storage blockHeader = m_blockheaders[_id][_blockHash];
        assert( PatriciaTrie.verifyProof(_value, _parentNodes, _path, blockHeader.receiptRootHash) );

        emit VerifiedProof(_id, _blockHash, uint(ProofType.RECEIPT));
        return true;
    }

    /*
    * CheckRootsProof
    * param: _id (bytes32) Unique id of chain submitting block from
    * param: _blockHash (bytes32) Block hash of block being submitted
    * param: _txNodes (bytes) RLP-encoded relevant nodes of the Tx trie
    * param: _receiptNodes (bytes) RLP-encoded relevant nodes of the Receipt trie
    *
    * emits: VerifiedTxProof(chainId, blockHash, proofType)
    *        chainId: (bytes32) hash of the chain verifying proof against
    *        blockHash: (bytes32) hash of the block verifying proof against
    *        proofType: (uint) enum of proof type
    *
    * All data associated with the proof must be constructed and provided to this function. Modifiers restrict execution
    * of this function to only allow if the chain the proof is for is registered to this contract and if the block that
    * the proof is for has been submitted.
    */
    function CheckRootsProof(
        bytes32 _id,
        bytes32 _blockHash,
        bytes _txNodes,
        bytes _receiptNodes
    )
        onlyRegisteredChains(_id)
        onlyExistingBlocks(_id, _blockHash)
        public
        returns (bool)
    {
        BlockHeader storage blockHeader = m_blockheaders[_id][_blockHash];

        assert( blockHeader.txRootHash == getRootNodeHash(_txNodes) );
        assert( blockHeader.receiptRootHash == getRootNodeHash(_receiptNodes) );

        emit VerifiedProof(_id, _blockHash, uint(ProofType.ROOTS));
        return true;
    }

    /*
    * @description      when a block is submitted the root hash must be added to a mapping of chains to hashes
    * @param _chainId   unique identifier of the chain from which the block hails     
    * @param _hash      root hash of the block being added
    */
    function addBlockHashToChain(bytes32 _chainId, bytes32 _hash) internal {
        m_blockhashes[_chainId][_hash] = true;
    }

    // XXX: Requires addition documentation
    function getBlockHeader(bytes32 _id, bytes32 _blockHash) public view returns (bytes32[3]) {
        BlockHeader storage header = m_blockheaders[_id][_blockHash];

        return [header.prevBlockHash, header.txRootHash, header.receiptRootHash];
    }

/*
========================================================================================================================

    Helper Functions

========================================================================================================================
*/

    /*
    * @description  copies 32 bytes from input into the output
	* @param output	memory allocation for the data you need to extract
	* @param input  array from which the data should be extracted
	* @param buf	index which the data starts within the byte array needs to have 32 bytes appended
	*/
	function bytesToBytes32(bytes input, uint256 buf) internal pure returns (bytes32 output) {
		buf = buf + 32;
        assembly {
			output := mload(add(input, buf))
		}
	}

    /*
    * @description  copies output.length bytes from the input into the output
	* @param output	memory allocation for the data you need to extract
	* @param input  array from which the data should be extracted
	* @param buf	index which the data starts within the byte array
	*/
	function bytesToBytes(bytes output, bytes input, uint256 buf) constant internal {
		uint256 outputLength = output.length;
		buf = buf + 32; // Append 32 as we need to point past the variable type definition
		assembly {
           let ret := staticcall(3000, 4, add(input, buf), outputLength, add(output, 32), outputLength)
	    }
	}

    function getRootNodeHash(bytes _rlpNodes) private view returns (bytes32) {
        RLP.RLPItem memory nodes = RLP.toRLPItem(_rlpNodes);
        RLP.RLPItem[] memory nodeList = RLP.toList(nodes);

        bytes memory b_nodeRoot = RLP.toBytes(nodeList[0]);

        return keccak256(b_nodeRoot);
    }


}

