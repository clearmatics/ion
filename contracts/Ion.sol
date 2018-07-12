pragma solidity ^0.4.23;

import "./libraries/RLP.sol";
import "./libraries/PatriciaTrie.sol";

contract Ion {
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using RLP for bytes;

    struct BlockHeader {
        bytes32 prevBlockHash;
        bytes32 txRootHash;
        bytes32 receiptRootHash;
    }

    bytes32 public chainId;
    bytes32[] public chains;
    mapping (bytes32 => bytes32[]) public m_blockhashes;
    mapping (bytes32 => BlockHeader) public m_blockheaders;

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

    event VerifiedTxProof(bytes32 chainId, bytes32 blockHash);
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
        bool chainRegistered = false;
        for (uint i = 0; i < chains.length; i++) {
            if (_id == chains[i]) {
                chainRegistered = true;
                break;
            }
        }
        require(chainRegistered, "Chain is not registered");
        _;
    }

    modifier onlyExistingBlocks(bytes32 _id, bytes32 _hash) {
        bool blockExists = false;
        bytes32[] hashes = m_blockhashes[_id];
        for (uint i = 0; i < hashes.length; i++) {
            if (_hash == hashes[i]) {
                blockExists = true;
                break;
            }
        }
        require(blockExists, "Block does not exist for chain");
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
    function RegisterChain(bytes32 _id) public {
        require( _id != chainId, "Cannot add this chain id to chain register" );
        for (uint i = 0; i < chains.length; i++) {
            require( chains[i] != _id, "Chain already exists" );
        }
        chains.push(_id);
    }

    /*
    * SubmitBlock
    * param: _id (bytes32) Unique id of chain submitting block from
    * param: _blockHash (bytes32) Block hash of block being submitted
    * param: _rlpBlockHeader (bytes) RLP-encoded byte array of the block header from other chain
    *
    * Submission of block headers from another chain, deconstructed and persisted as BlockHeader structs defined above
    * and adds it to the list of known block hashes and headers of specified chain.
    */
    function SubmitBlock(bytes32 _id, bytes32 _blockHash, bytes _rlpBlockHeader) onlyRegisteredChains(_id) public {
        RLP.RLPItem[] memory header = _rlpBlockHeader.toRLPItem().toList();

        bytes32 hashedHeader = keccak256(_rlpBlockHeader);
        require( hashedHeader == _blockHash );

        BlockHeader storage blockHeader = m_blockheaders[_blockHash];
        blockHeader.prevBlockHash = bytesToBytes32(header[0].toBytes(), 1);
        blockHeader.txRootHash = bytesToBytes32(header[4].toBytes(), 1);
        blockHeader.receiptRootHash = bytesToBytes32(header[5].toBytes(), 1);

        addBlockHashToChain(_id, _blockHash);
    }

    /*
    * CheckTxProof
    * param: _id (bytes32) Unique id of chain submitting block from
    * param: _blockHash (bytes32) Block hash of block being submitted
    * param: _value (bytes) RLP-encoded transaction object array with fields defined as: https://github.com/ethereumjs/ethereumjs-tx/blob/0358fad36f6ebc2b8bea441f0187f0ff0d4ef2db/index.js#L50
    * param: _parentNodes (bytes) RLP-encoded array of all relevant nodes from root node to node to prove
    * param: _path (bytes) Byte array of the path to the node to be proved
    *
    * emits: VerifiedTxProof(chainId, blockHash)
    *        chainId: (bytes32) hash of the chain verifying proof against
    *        blockHash: (bytes32) hash of the block verifying proof against
    *
    * All data associated with the proof must be constructed and provided to this function. Modifiers restrict execution
    * of this function to only allow if the chain the proof is for is registered to this contract and if the block that
    * the proof is for has been submitted.
    */
    function CheckTxProof(bytes32 _id, bytes32 _blockHash, bytes _value, bytes _parentNodes, bytes _path) onlyRegisteredChains(_id) onlyExistingBlocks(_id, _blockHash) public returns (bool) {
        BlockHeader storage blockHeader = m_blockheaders[_blockHash];
        assert( PatriciaTrie.verifyProof(_value, _parentNodes, _path, blockHeader.txRootHash) );
        emit VerifiedTxProof(_id, _blockHash);
        return true;
    }

    function CheckReceiptProof() public {
    }

    function CheckRootsProof() public {
    }


    function addBlockHashToChain(bytes32 _chainId, bytes32 _hash) internal {
        bytes32[] storage blockHashes = m_blockhashes[_chainId];

        for (uint i = 0; i < blockHashes.length; i++) {
            require( blockHashes[i] != _hash );
        }
        blockHashes.push(_hash);
    }

    function getBlockHeader(bytes32 _blockHash) public view returns (bytes32[3]) {
        BlockHeader storage header = m_blockheaders[_blockHash];

        return [header.prevBlockHash, header.txRootHash, header.receiptRootHash];
    }
/*
========================================================================================================================

    Helper Functions

========================================================================================================================
*/
    function bytesToBytes32(bytes b, uint offset) private pure returns (bytes32) {
        bytes32 out;

        for (uint i = 0; i < 32; i++) {
            out |= bytes32(b[offset + i] & 0xFF) >> (i * 8);
        }
        return out;
    }
}

