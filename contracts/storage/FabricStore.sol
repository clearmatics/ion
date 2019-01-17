pragma solidity ^0.4.24;

import "./BlockStore.sol";
import "../libraries/RLP.sol";

contract FabricStore is BlockStore {
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using RLP for bytes;

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
        string[] transactions;
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
    mapping (string => Transaction) internal m_transactions;
    mapping (string => bool) internal m_transactions_exist;

    constructor(address _ionAddr) BlockStore(_ionAddr) public {}

    // Function name is inaccurate for Fabric due to blocks being a sub-structure to a channel
    // Will need refactoring
    function addBlock(bytes32 _chainId, bytes32 _blockHash, bytes _blockBlob)
        onlyIon
        onlyRegisteredChains(_chainId)
    {
        RLP.RLPItem[] memory data = _blockBlob.toRLPItem().toList();

        // Iterate all channel objects in the data structure
        for (uint i = 0; i < data.length; i++) {
            decodeChannelObject(data[i].toBytes());
        }

//        emit BlockAdded(_chainId, _blockHash);
    }

    function decodeChannelObject(bytes _channelRLP) internal {
        RLP.RLPItem[] memory channelRLP = _channelRLP.toRLPItem().toList();

        string memory channelId = channelRLP[0].toAscii();
        Channel storage channel = m_channels[channelId];

        // Currently adds the channel if it does not exist. This may need changing.
        if (keccak256(channel.id) == keccak256("")) {
            channel.id = channelId;
        }

        RLP.RLPItem[] memory blocksRLP = channelRLP[1].toBytes().toRLPItem().toList();

        // Iterate all blocks in the channel structure.
        for (uint i = 0; i < blocksRLP.length; i++) {
            Block memory block = decodeBlockObject(channelId, blocksRLP[i].toBytes());
            require(!channel.blocks[block.hash], "Block with identical hash already exists");
            channel.blocks[block.hash] = true;
            channel.m_blocks[block.hash] = block;
        }
    }

    function decodeBlockObject(string _channelId, bytes _blockRLP) internal returns (Block memory){
        RLP.RLPItem[] memory blockRLP = _blockRLP.toRLPItem().toList();

        string memory blockHash = blockRLP[0].toAscii();

        Block memory block;

        block.number = blockRLP[1].toUint();
        block.hash = blockHash;
        block.prevHash = blockRLP[2].toAscii();
        block.dataHash = blockRLP[3].toAscii();
        block.timestamp_s = blockRLP[4].toUint();
        block.timestamp_nanos = blockRLP[5].toUint();

        RLP.RLPItem[] memory txnsRLP = blockRLP[6].toBytes().toRLPItem().toList();

        block.transactions = new string[](txnsRLP.length);

        // Iterate all transactions in the block
        for (uint i = 0; i < txnsRLP.length; i++) {
            block.transactions[i] = decodeTxObject(txnsRLP[i].toBytes());
        }

        return block;
    }

    function decodeTxObject(bytes _txRLP) internal returns (string) {
        RLP.RLPItem[] memory txRLP = _txRLP.toRLPItem().toList();

        string memory txId = txRLP[0].toAscii();

        require(!m_transactions_exist[txId], "Transaction ID already exists");
        m_transactions_exist[txId] = true;

        Transaction storage tx = m_transactions[txId];
        tx.id = txId;

        RLP.RLPItem[] memory namespacesRLP = txRLP[0].toBytes().toRLPItem().toList();

        // Iterate all namespace rwsets in the transaction
        for (uint i = 0; i < namespacesRLP.length; i++) {
//            Namespace memory namespace = decodeNamespaceRW(namespacesRLP[i].toBytes());
//            tx.m_nsrw[namespace.namespace] = namespace;

            Namespace storage namespace = tx.m_nsrw[namespace.namespace];
            namespace.namespace = nsrwRLP[0].toAscii();

            RLP.RLPItem[] memory nsrwRLP = namespacesRLP[i].toBytes().toRLPItem().toList();

            // Iterate all read sets in the namespace
            RLP.RLPItem[] memory readsetsRLP = nsrwRLP[1].toBytes().toRLPItem().toList();
            for (uint j = 0; j < readsetsRLP.length; j++) {
                namespace.reads.push(decodeReadset(readsetsRLP[j].toBytes()));
            }

            // Iterate all write sets in the namespace
            RLP.RLPItem[] memory writesetsRLP = nsrwRLP[2].toBytes().toRLPItem().toList();
            WriteSet[] memory writesets = new WriteSet[](writesetsRLP.length);
            for (uint k = 0; k < writesetsRLP.length; k++) {
                namespace.writes.push(decodeWriteset(writesetsRLP[k].toBytes()));
            }
        }

        return txId;
    }

//    function decodeNamespaceRW(bytes _nsrwRLP) internal returns (Namespace memory) {
//        RLP.RLPItem[] memory nsrwRLP = _nsrwRLP.toRLPItem().toList();
//
//        string memory namespace = nsrwRLP[0].toAscii();
//
//        // Iterate all read sets in the namespace
//        RLP.RLPItem[] memory readsetsRLP = nsrwRLP[1].toBytes().toRLPItem().toList();
//        ReadSet[] memory readsets = new ReadSet[](readsetsRLP.length);
//        for (uint i = 0; i < readsetsRLP.length; i++) {
//            readsets[i] = decodeReadset(readsetsRLP[i].toBytes());
//        }
//
//        // Iterate all write sets in the namespace
//        RLP.RLPItem[] memory writesetsRLP = nsrwRLP[2].toBytes().toRLPItem().toList();
//        WriteSet[] memory writesets = new WriteSet[](writesetsRLP.length);
//        for (uint j = 0; j < writesetsRLP.length; j++) {
//            writesets[i] = decodeWriteset(writesetsRLP[i].toBytes());
//        }
//
//        return Namespace(namespace, readsets, writesets);
//    }

    function decodeReadset(bytes _readsetRLP) internal returns (ReadSet memory) {
        RLP.RLPItem[] memory readsetRLP = _readsetRLP.toRLPItem().toList();

        string memory key = readsetRLP[0].toAscii();

        RLP.RLPItem[] memory rsv = readsetRLP[1].toBytes().toRLPItem().toList();

        uint blockNo = rsv[0].toUint();
        uint txNo = 0;

        if (rsv.length > 1) {
            txNo = rsv[1].toUint();
        }
        RSVersion memory version = RSVersion(blockNo, txNo);

        return ReadSet(key, version);
    }

    function decodeWriteset(bytes _writesetRLP) internal returns (WriteSet memory){
        RLP.RLPItem[] memory writesetRLP = _writesetRLP.toRLPItem().toList();

        string memory key = writesetRLP[0].toAscii();
        bool isDelete = writesetRLP[1].toBool();
        string memory value = writesetRLP[2].toAscii();

        return WriteSet(key, isDelete, value);
    }
}