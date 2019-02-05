pragma solidity ^0.4.24;

import "./BlockStore.sol";
import "../libraries/RLP.sol";
import "../libraries/SolidityUtils.sol";

contract FabricStore is BlockStore {
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using RLP for bytes;

    struct Chain {
        bytes32 id;
        mapping (string => Channel) m_channels;
    }

    struct Channel {
        string id;
        mapping (string => bool) blocks;
        mapping (string => Block) m_blocks;
        mapping (string => Transaction) m_transactions;
        mapping (string => bool) m_transactions_exist;
        mapping (string => State) m_state;
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
        string blockHash;
        string[] namespaces;
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

    struct State {
        string key;
        RSVersion version;
        string value;
    }

    mapping (bytes32 => Chain) public m_networks;

    constructor(address _ionAddr) BlockStore(_ionAddr) public {}

    event BlockAdded(bytes32 chainId, string channelId, string blockHash);

    function addChain(bytes32 _chainId) onlyIon public returns (bool) {
        require(super.addChain(_chainId), "Storage addChain parent call failed");

        Chain storage chain = m_networks[_chainId];
        chain.id = _chainId;

        return true;
    }

    // Function name is inaccurate for Fabric due to blocks being a sub-structure to a channel
    // Will need refactoring
    function addBlock(bytes32 _chainId, bytes _blockBlob)
        public
        onlyIon
        onlyRegisteredChains(_chainId)
    {
        RLP.RLPItem[] memory data = _blockBlob.toRLPItem().toList();

        // Iterate all channel objects in the data structure
        for (uint i = 0; i < data.length; i++) {
            decodeChannelObject(_chainId, data[i].toBytes());
        }
    }

    function decodeChannelObject(bytes32 _chainId, bytes _channelRLP) internal {
        RLP.RLPItem[] memory channelRLP = _channelRLP.toRLPItem().toList();

        string memory channelId = channelRLP[0].toAscii();
        Channel storage channel = m_networks[_chainId].m_channels[channelId];

        // Currently adds the channel if it does not exist. This may need changing.
        if (keccak256(abi.encodePacked(channel.id)) == keccak256(abi.encodePacked(""))) {
            channel.id = channelId;
        }

//        RLP.RLPItem[] memory blocksRLP = channelRLP[1].toList();
//
//        // Iterate all blocks in the channel structure. Currently not used as we only focus on parsing single blocks
//        for (uint i = 0; i < blocksRLP.length; i++) {
//            Block memory block = decodeBlockObject(_chainId, channelId, channelRLP[1].toBytes());
//            require(!channel.blocks[block.hash], "Block with identical hash already exists");
//            channel.blocks[block.hash] = true;
//            channel.m_blocks[block.hash] = block;
//
//            emit BlockAdded(_chainId, channelId, block.hash);
//        }

        Block memory blk = decodeBlockObject(_chainId, channelId, channelRLP[1].toBytes());
        require(!channel.blocks[blk.hash], "Block with identical hash already exists");

        mutateState(_chainId, channelId, blk);

        channel.blocks[blk.hash] = true;
        channel.m_blocks[blk.hash] = blk;

        emit BlockAdded(_chainId, channelId, blk.hash);
    }

    function decodeBlockObject(bytes32 _chainId, string _channelId, bytes _blockRLP) internal returns (Block memory) {
        RLP.RLPItem[] memory blockRLP = _blockRLP.toRLPItem().toList();

        string memory blockHash = blockRLP[0].toAscii();

        Block memory blk;

        blk.number = blockRLP[1].toUint();
        blk.hash = blockHash;
        blk.prevHash = blockRLP[2].toAscii();
        blk.dataHash = blockRLP[3].toAscii();
        blk.timestamp_s = blockRLP[4].toUint();
        blk.timestamp_nanos = blockRLP[5].toUint();

        RLP.RLPItem[] memory txnsRLP = blockRLP[6].toList();

        blk.transactions = new string[](txnsRLP.length);

        // Iterate all transactions in the block
        for (uint i = 0; i < txnsRLP.length; i++) {
            string memory txId = decodeTxObject(txnsRLP[i].toBytes(), _chainId, _channelId);
            require(!isTransactionExists(_chainId, _channelId, txId), "Transaction already exists");
            blk.transactions[i] = txId;
            injectBlockHashToTx(_chainId, _channelId, txId, blockHash);
            flagTx(_chainId, _channelId, txId);
        }

        return blk;
    }

    function decodeTxObject(bytes _txRLP, bytes32 _chainId, string _channelId) internal returns (string) {
        RLP.RLPItem[] memory txRLP = _txRLP.toRLPItem().toList();

        Transaction storage txn = m_networks[_chainId].m_channels[_channelId].m_transactions[txRLP[0].toAscii()];
        txn.id = txRLP[0].toAscii();

        RLP.RLPItem[] memory namespacesRLP = txRLP[1].toList();

        // Iterate all namespace rwsets in the transaction
        for (uint i = 0; i < namespacesRLP.length; i++) {
            RLP.RLPItem[] memory nsrwRLP = namespacesRLP[i].toList();

            Namespace storage namespace = txn.m_nsrw[nsrwRLP[0].toAscii()];
            namespace.namespace = nsrwRLP[0].toAscii();
            txn.namespaces.push(nsrwRLP[0].toAscii());

            // Iterate all read sets in the namespace
            RLP.RLPItem[] memory readsetsRLP = nsrwRLP[1].toList();
            for (uint j = 0; j < readsetsRLP.length; j++) {
                namespace.reads.push(decodeReadset(readsetsRLP[j].toBytes()));
            }

            // Iterate all write sets in the namespace
            RLP.RLPItem[] memory writesetsRLP = nsrwRLP[2].toList();
            for (uint k = 0; k < writesetsRLP.length; k++) {
                namespace.writes.push(decodeWriteset(writesetsRLP[k].toBytes()));
            }
        }

        return txRLP[0].toAscii();
    }

    function mutateState(bytes32 _chainId, string _channelId, Block memory _blk) internal {
        string[] memory txIds = _blk.transactions;

        // Iterate across all transactions
        for (uint i = 0; i < txIds.length; i++) {
            Transaction storage txn = m_networks[_chainId].m_channels[_channelId].m_transactions[txIds[i]];

            // Iterate across all namespaces
            for (uint j = 0; j < txn.namespaces.length; j++) {
                string storage namespace = txn.namespaces[j];

                // Iterate across all writesets and check readset version of each write key against stored version
                for (uint k = 0; k < txn.m_nsrw[namespace].writes.length; k++) {
                    State storage state = m_networks[_chainId].m_channels[_channelId].m_state[txn.m_nsrw[namespace].writes[k].key];

                    if (keccak256(abi.encodePacked(state.key)) == keccak256(abi.encodePacked(txn.m_nsrw[namespace].writes[k].key))) {
                        if (!isExpectedReadVersion(txn.m_nsrw[namespace], state.version, state.key))
                            continue;
                    }

                    state.key = txn.m_nsrw[namespace].writes[k].key;
                    state.version = RSVersion(_blk.number, i);
                    state.value = txn.m_nsrw[namespace].writes[k].value;
                }
            }
        }
    }

    function injectBlockHashToTx(bytes32 _chainId, string _channelId, string _txId, string _blockHash) internal {
        Transaction storage txn = m_networks[_chainId].m_channels[_channelId].m_transactions[_txId];
        txn.blockHash = _blockHash;
    }

    function flagTx(bytes32 _chainId, string _channelId, string _txId) internal {
        m_networks[_chainId].m_channels[_channelId].m_transactions_exist[_txId] = true;
    }

    function decodeReadset(bytes _readsetRLP) internal view returns (ReadSet memory) {
        RLP.RLPItem[] memory readsetRLP = _readsetRLP.toRLPItem().toList();

        string memory key = readsetRLP[0].toAscii();

        RLP.RLPItem[] memory rsv = readsetRLP[1].toList();

        uint blockNo = rsv[0].toUint();
        uint txNo = 0;

        if (rsv.length > 1) {
            txNo = rsv[1].toUint();
        }
        RSVersion memory version = RSVersion(blockNo, txNo);

        return ReadSet(key, version);
    }

    function decodeWriteset(bytes _writesetRLP) internal view returns (WriteSet memory){
        RLP.RLPItem[] memory writesetRLP = _writesetRLP.toRLPItem().toList();

        string memory key = writesetRLP[0].toAscii();
        string memory value = writesetRLP[2].toAscii();

        bool isDelete = false;
        string memory isDeleteStr = writesetRLP[1].toAscii();
        if (keccak256(abi.encodePacked(isDeleteStr)) == keccak256(abi.encodePacked("true"))) {
            isDelete = true;
        }

        return WriteSet(key, isDelete, value);
    }

    function isExpectedReadVersion(Namespace memory _namespace, RSVersion memory _version, string _key) internal pure returns (bool) {
        ReadSet[] memory reads = _namespace.reads;

        for (uint i = 0; i < reads.length; i++) {
            ReadSet memory readset = reads[i];

            if (keccak256(abi.encodePacked(readset.key)) == keccak256(abi.encodePacked(_key)))
                return isSameVersion(readset.version, _version);
        }

        return false;
    }

    function isSameVersion(RSVersion memory _v1, RSVersion memory _v2) internal pure returns (bool) {
        if (_v1.blockNo != _v2.blockNo)
            return false;

        if (_v1.txNo != _v2.txNo)
            return false;

        return true;
    }

    function getBlock(bytes32 _chainId, string _channelId, string _blockHash) public view returns (uint, string, string, string, uint, uint, string) {
        Block storage blk = m_networks[_chainId].m_channels[_channelId].m_blocks[_blockHash];

        require(keccak256(abi.encodePacked(blk.hash)) != keccak256(abi.encodePacked("")), "Block does not exist.");

        string memory txs = blk.transactions[0];

        for (uint i = 1; i < blk.transactions.length; i++) {
            txs = string(abi.encodePacked(txs, ",", blk.transactions[i]));
        }

        return (blk.number, blk.hash, blk.prevHash, blk.dataHash, blk.timestamp_s, blk.timestamp_nanos, txs);
    }

    function getTransaction(bytes32 _chainId, string _channelId, string _txId) public view returns (string, string) {
        Transaction storage txn = m_networks[_chainId].m_channels[_channelId].m_transactions[_txId];

        require(isTransactionExists(_chainId, _channelId, _txId), "Transaction does not exist.");

        string memory ns = txn.namespaces[0];

        for (uint i = 1; i < txn.namespaces.length; i++) {
            ns = string(abi.encodePacked(ns, ",", txn.namespaces[i]));
        }

        return (txn.blockHash, ns);
    }

    function isTransactionExists(bytes32 _chainId, string _channelId, string _txId) public view returns (bool) {
        return m_networks[_chainId].m_channels[_channelId].m_transactions_exist[_txId];
    }

    function getNSRW(bytes32 _chainId, string _channelId, string _txId, string _namespace) public view returns (string, string) {
        Namespace storage ns = m_networks[_chainId].m_channels[_channelId].m_transactions[_txId].m_nsrw[_namespace];

        require(keccak256(abi.encodePacked(ns.namespace)) != keccak256(abi.encodePacked("")), "Namespace does not exist.");

        string memory reads;
        for (uint i = 0; i < ns.reads.length; i++) {
            RSVersion storage version = ns.reads[i].version;
            reads = string(abi.encodePacked(reads, "{ key: ", ns.reads[i].key, ", version: { blockNo: ", SolUtils.UintToString(version.blockNo), ", txNo: ", SolUtils.UintToString(version.txNo), " } } "));
        }

        string memory writes;
        for (uint j = 0; j < ns.writes.length; j++) {
            writes = string(abi.encodePacked(writes, "{ key: ", ns.writes[j].key, ", isDelete: ", SolUtils.BoolToString(ns.writes[j].isDelete), ", value: ", ns.writes[j].value, " } "));
        }

        return (reads, writes);
    }

    function getState(bytes32 _chainId, string _channelId, string _key) public view returns (uint, uint, string) {
        State storage state = m_networks[_chainId].m_channels[_channelId].m_state[_key];

        require(keccak256(abi.encodePacked(state.key)) != keccak256(abi.encodePacked("")), "Key unrecognised.");

        return (state.version.blockNo, state.version.txNo, state.value);
    }
}