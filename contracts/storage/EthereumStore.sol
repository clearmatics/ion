// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.24;

import "../libraries/ECVerify.sol";
import "../libraries/RLP.sol";
import "../libraries/PatriciaTrie.sol";
import "../libraries/SolidityUtils.sol";
import "./BlockStore.sol";

contract EthereumStore is BlockStore {
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using RLP for bytes;

    /*
    *   @description    BlockHeader struct containing trie root hashes for tx verifications
    */
    struct BlockHeader {
        bytes32 txRootHash;
        bytes32 receiptRootHash;
    }

    mapping (bytes32 => bool) public m_blockhashes;
    mapping (bytes32 => BlockHeader) public m_blockheaders;

    enum ProofType { TX, RECEIPT, ROOTS }

    event BlockAdded(bytes32 chainId, bytes32 blockHash);
    event VerifiedProof(bytes32 chainId, bytes32 blockHash, uint proofType);

    constructor(address _ionAddr) BlockStore(_ionAddr) public {}

    /*
    * onlyExistingBlocks
    * param: _id (bytes32) Unique id of chain supplied to function
    * param: _hash (bytes32) Block hash which needs validation
    *
    * Modifier that checks if the provided block hash has been verified by the validation contract
    */
    modifier onlyExistingBlocks(bytes32 _hash) {
        require(m_blockhashes[_hash], "Block does not exist for chain");
        _;
    }


    /*
    * @description          when a block is submitted the header must be added to a mapping of blockhashes and m_chains to blockheaders
    * @param _chainId       ID of the chain the block is from
    * @param _blockHash     Block hash of the block being added
    * @param _blockBlob     Bytes blob of the RLP-encoded block header being added
    */
    function addBlock(bytes32 _chainId, bytes _blockBlob)
        public
        onlyIon
        onlyRegisteredChains(_chainId)
    {
        bytes32 blockHash = keccak256(_blockBlob);
        require(!m_blockhashes[blockHash], "Block already exists" );

        RLP.RLPItem[] memory header = _blockBlob.toRLPItem().toList();
        require(header.length == 15, "Block Header parameter mismatch");

        m_blockhashes[blockHash] = true;
        m_blockheaders[blockHash].txRootHash = header[4].toBytes32();
        m_blockheaders[blockHash].receiptRootHash = header[5].toBytes32();

        emit BlockAdded(_chainId, blockHash);
    }

    function CheckProofs(bytes32 _chainId, bytes32 _blockHash, bytes _proof) public returns (bytes memory) {
        RLP.RLPItem[] memory proof = _proof.toRLPItem().toList();

        require(proof.length == 5, "Malformed proof");

        assert(CheckRootsProof(_chainId, _blockHash, proof[2].toBytes(), proof[4].toBytes()));
        assert(CheckTxProof(_chainId, _blockHash, proof[1].toBytes(), proof[2].toBytes(), proof[0].toBytes()));
        assert(CheckReceiptProof(_chainId, _blockHash, proof[3].toBytes(), proof[4].toBytes(), proof[0].toBytes()));

        return proof[3].toBytes();
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
        bytes32 _chainId,
        bytes32 _blockHash,
        bytes _value,
        bytes _parentNodes,
        bytes _path
    )
        onlyRegisteredChains(_chainId)
        onlyExistingBlocks(_blockHash)
        internal
        returns (bool)
    {
        verifyProof(_value, _parentNodes, _path, m_blockheaders[_blockHash].txRootHash);

        emit VerifiedProof(_chainId, _blockHash, uint(ProofType.TX));
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
    * All data associated with the proof must be constructed and paddChainrovided to this function. Modifiers restrict execution
    * of this function to only allow if the chain the proof is for is registered to this contract and if the block that
    * the proof is for has been submitted.
    */
    function CheckReceiptProof(
        bytes32 _chainId,
        bytes32 _blockHash,
        bytes _value,
        bytes _parentNodes,
        bytes _path
    )
        onlyRegisteredChains(_chainId)
        onlyExistingBlocks(_blockHash)
        internal
        returns (bool)
    {
        verifyProof(_value, _parentNodes, _path, m_blockheaders[_blockHash].receiptRootHash);

        emit VerifiedProof(_chainId, _blockHash, uint(ProofType.RECEIPT));
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
        bytes32 _chainId,
        bytes32 _blockHash,
        bytes _txNodes,
        bytes _receiptNodes
    )
        onlyRegisteredChains(_chainId)
        onlyExistingBlocks(_blockHash)
        internal
        returns (bool)
    {
        assert( m_blockheaders[_blockHash].txRootHash == getRootNodeHash(_txNodes) );
        assert( m_blockheaders[_blockHash].receiptRootHash == getRootNodeHash(_receiptNodes) );

        emit VerifiedProof(_chainId, _blockHash, uint(ProofType.ROOTS));
        return true;
    }

    function verifyProof(bytes _value, bytes _parentNodes, bytes _path, bytes32 _hash) internal {
        assert( PatriciaTrie.verifyProof(_value, _parentNodes, _path, _hash) );
    }

/*
========================================================================================================================

    Helper Functions

========================================================================================================================
*/

    /*
    * @description      returns the root node of an RLP encoded Patricia Trie
	* @param _rlpNodes  RLP encoded trie
	* @returns          root hash
	*/
    function getRootNodeHash(bytes _rlpNodes) private view returns (bytes32) {
        RLP.RLPItem[] memory nodeList = _rlpNodes.toRLPItem().toList();

        bytes memory b_nodeRoot = RLP.toBytes(nodeList[0]);

        return keccak256(b_nodeRoot);
    }


}

