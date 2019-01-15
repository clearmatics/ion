pragma solidity ^0.4.24;

import "./BlockStore.sol";
import "../libraries/RLP.sol";

contract FabricStore is BlockStore {

    struct Channel {
        string id;
        mapping (string => bool) blocks;
        mapping (string => Block) m_blocks;
    }

    struct Block {
        uint number;
        string hash;
        string prevHash;
        string dataHash;
        uint timestamp_s;
        uint timestamp_nanos;
        mapping (string => bool) transactions;
        mapping (string => Transaction) m_transactions;
    }

    struct Transaction {
        string id;
        mapping (string => Namespace) m_nsrw;
    }

    struct Namespace {
        string namespace;
        ReadSet[] reads;
        WriteSet[] writes;
    }

    struct ReadSet {
        string key;
        RSVersion version;
    }

    struct RSVersion {
        uint blockNo;
        uint txNo;
    }

    struct WriteSet {
        string key;
        bool isDelete;
        string value;
    }

    mapping (string => string) internal m_ledgerState;
    mapping (string => Channel) internal m_channels;

    constructor(address _ionAddr) BlockStore(_ionAddr) public {}

    function addBlock(bytes32 _chainId, bytes32 _blockHash, bytes _blockBlob)
        onlyIon
        onlyRegisteredChains(_chainId)
    {
        require(!m_blockhashes[_blockHash], "Block already exists" );
        RLP.RLPItem[] memory data = _blockBlob.toRLPItem().toList();

        // Iterate all channel objects in the data structure
        for (uint i = 0; i < data.length; i++) {
            decodeChannelObject(data[i]);
        }

        emit BlockAdded(_chainId, _blockHash);
    }

    function decodeAndAddChannelObject(bytes _channelRLP) internal {
        RLP.RLPItem[] memory channelRLP = _channelRLP.toRLPItem().toList();

        string channelId = channelRLP[0].toAscii();
        Channel storage channel = m_channels[channelId];

        // Currently adds the channel if it does not exist. This may need changing.
        if (channel.id == "") {
            channel.id = channelId;
        }

        RLP.RLPItem[] memory blocksRLP = channelRLP[1].toRLPItem().toList();

        // Iterate all blocks in the channel structure.
        for (uint i = 0; i < blocksRLP.length; i++) {
            decodeAndAddBlockObject(channel, blocksRLP[i]);
        }
    }

    function decodeAndAddBlockObject(Channel storage _channel, bytes _blockRLP) internal {
        RLP.RLPItem[] memory blockRLP = _blockRLP.toRLPItem().toList();

        string blockHash = blockRLP[0].toAscii();

        require(!_channel.blocks[blockHash], "Block hash already exists");
        _channel.blocks[blockHash] = true;
        Block storage block = _channel.m_blocks[blockHash];

        block.number = blockRLP[1].toUint();
        block.hash = blockHash;
        block.prevHash = blockRLP[2].toAscii();
        block.dataHash = blockRLP[3].toAscii();
        block.timestamp_s = blockRLP[4].toUint();
        block.timestamp_nanos = blockRLP[5].toUint();

        RLP.RLPItem[] memory txnsRLP = blockRLP[6].toRLPItem().toList();

        // Iterate all transactions in the block
        for (uint i = 0; i < txnsRLP.length; i++) {
            decodeAndAddTxObject(block, txnsRLP[i]);
        }
    }

    function decodeAndAddTxObject(Block storage _block, bytes _txRLP) internal {
        RLP.RLPItem[] memory txRLP = _txRLP.toRLPItem().toList();

        string txId = txRLP[0].toAscii();

        require(!_block.transactions[txId], "Transaction ID already exists");
        _block.transactions[txId] = true;

        Transaction storage tx = _block.transactions[txId];
        tx.id = txId;

        RLP.RLPItem[] memory namespacesRLP = txRLP[0].toRLPItem().toList();

        // Iterate all namespace rwsets in the transaction
        for (uint i = 0; i < namespacesRLP.length; i++) {
            decodeAndAddNamespaceRW(tx, namespacesRLP[i]);
        }
    }

    function decodeAndAddNamespaceRW(Transaction storage _tx, bytes _nsrwRLP) internal {
        RLP.RLPItem[] memory nsrwRLP = _nsrwRLP.toRLPItem().toList();

        string namespace = nsrwRLP[0].toAscii();

        Namespace storage ns = _tx.m_nsrw[namespace];
        ns.namespace = namespace;

        // Iterate all read sets in the namespace
        RLP.RLPItem[] memory readsetsRLP = nsrwRLP[1].toRLPItem().toList();
        for (uint i = 0; i < readsetsRLP.length; i++) {
            ns.reads.push(decodeReadset(readsetsRLP[i]));
        }

        // Iterate all write sets in the namespace
        RLP.RLPItem[] memory writesetsRLP = nsrwRLP[2].toRLPItem().toList();
        for (uint j = 0; j < writesetsRLP.length; j++) {
            ns.writes.push(decodeAndAddWriteset(ns, writesetsRLP[i]));
        }
    }

    function decodeAndAddReadset(bytes _readsetRLP) internal returns (ReadSet storage) {
        RLP.RLPItem[] memory readsetRLP = _readsetRLP.toRLPItem().toList();

        string key = readsetRLP[0].toAscii();

        RLP.RLPItem[] storage rsv = readsetRLP[1].toRLPItem().toList();

        uint blockNo = rsv[0].toUint();
        uint txNo = rsv[1].toUint();
        RSVersion storage version = RSVersion(blockNo, txNo);

        return Readset(key, version);
    }

    function decodeAndAddWriteset(Namespace storage _ns, bytes _writesetRLP) internal returns (WriteSet storage){
        RLP.RLPItem[] memory writesetRLP = _writesetRLP.toRLPItem().toList();

        string key = writesetRLP[0].toAscii();
        bool isDelete = writesetRLP[1].toBool();
        string value = writesetRLP[2].toAscii();

        return WriteSet(key, isDelete, value);
    }
}