pragma solidity ^0.4.24;

import "./BlockStore.sol";

contract FabricStore is BlockStore {

    struct Channel {
        string[] blocks;
        mapping (string => Block) m_blocks;
    }

    struct Block {
        uint number;
        string prevHash;
        string dataHash;
        string[] transactions;
        mapping (string => Transaction) m_transactions;
        uint timestamp;
        uint timestamp_nanos;
    }

    struct Transaction {
        mapping (string => Mutation) m_mutations;
        uint validationCode;
    }

    struct Mutation {
        ReadSet[] reads;
        WriteSet[] writes;
    }

    struct ReadSet {
        string key;
        string version;
    }

    struct WriteSet {
        string key;
        string value;
    }

    internal mapping (string => string) m_ledgerState;
    internal mapping (string => Channel) m_channels;

    constructor(address _ionAddr) BlockStore(_ionAddr) public {}

    function addBlock(bytes32 _chainId, bytes32 _blockHash, bytes _blockBlob)
        onlyIon
        onlyRegisteredChains(_chainId)
    {
        require(!m_blockhashes[_blockHash], "Block already exists" );
        RLP.RLPItem[] memory header = _blockBlob.toRLPItem().toList();

        bytes32 hashedHeader = keccak256(_blockBlob);
        require(hashedHeader == _blockHash, "Hashed header does not match submitted block hash!");

        m_blockhashes[_blockHash] = true;
        m_blockheaders[_blockHash].txRootHash = header[4].toBytes32();
        m_blockheaders[_blockHash].receiptRootHash = header[5].toBytes32();

        emit BlockAdded(_chainId, _blockHash);
    }

}