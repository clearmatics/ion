pragma solidity ^0.4.23;

import "./RLP.sol";

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

    function CheckTxProof() public {
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

