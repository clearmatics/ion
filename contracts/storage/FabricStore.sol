pragma solidity ^0.5.12;

import "./BlockStore.sol";
import "../libraries/RLPReader.sol";
import "../libraries/SolidityUtils.sol";

contract FabricStore is BlockStore {
    using RLPReader for RLPReader.RLPItem;
    using RLPReader for RLPReader.Iterator;
    using RLPReader for bytes;

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
    function addBlock(bytes32 _chainId, bytes memory _blockBlob)
        public
        onlyIon
        onlyRegisteredChains(_chainId)
    {
        RLPReader.RLPItem[] memory data = _blockBlob.toRLPItem().toList();

        // Iterate all channel objects in the data structure
        for (uint i = 0; i < data.length; i++) {
            decodeChannelObject(_chainId, data[i].toBytes());
        }
    }

    function decodeChannelObject(bytes32 _chainId, bytes memory _channelRLP) internal {
        RLPReader.RLPItem[] memory channelRLP = _channelRLP.toRLPItem().toList();

        string memory channelId = string(channelRLP[0].toBytes());
        Channel storage channel = m_networks[_chainId].m_channels[channelId];

        // Currently adds the channel if it does not exist. This may need changing.
        if (keccak256(abi.encodePacked(channel.id)) == keccak256(abi.encodePacked(""))) {
            channel.id = channelId;
        }

//        RLPReader.RLPItem[] memory blocksRLP = channelRLP[1].toList();
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

    function decodeBlockObject(bytes32 _chainId, string memory _channelId, bytes memory _blockRLP) internal returns (Block memory) {
        RLPReader.RLPItem[] memory blockRLP = _blockRLP.toRLPItem().toList();

        string memory blockHash = string(blockRLP[0].toBytes());

        Block memory blk;

        blk.number = blockRLP[1].toUint();
        blk.hash = blockHash;
        blk.prevHash = string(blockRLP[2].toBytes());
        blk.dataHash = string(blockRLP[3].toBytes());
        blk.timestamp_s = blockRLP[4].toUint();
        blk.timestamp_nanos = blockRLP[5].toUint();

        RLPReader.RLPItem[] memory txnsRLP = blockRLP[6].toList();

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

    function decodeTxObject(bytes memory _txRLP, bytes32 _chainId, string memory _channelId) internal returns (string memory) {
        RLPReader.RLPItem[] memory txRLP = _txRLP.toRLPItem().toList();

        Transaction storage txn = m_networks[_chainId].m_channels[_channelId].m_transactions[string(txRLP[0].toBytes())];
        txn.id = string(txRLP[0].toBytes());

        RLPReader.RLPItem[] memory namespacesRLP = txRLP[1].toList();

        // Iterate all namespace rwsets in the transaction
        for (uint i = 0; i < namespacesRLP.length; i++) {
            RLPReader.RLPItem[] memory nsrwRLP = namespacesRLP[i].toList();

            Namespace storage namespace = txn.m_nsrw[string(nsrwRLP[0].toBytes())];
            namespace.namespace = string(nsrwRLP[0].toBytes());
            txn.namespaces.push(string(nsrwRLP[0].toBytes()));

            // Iterate all read sets in the namespace
            RLPReader.RLPItem[] memory readsetsRLP = nsrwRLP[1].toList();
            for (uint j = 0; j < readsetsRLP.length; j++) {
                namespace.reads.push(decodeReadset(readsetsRLP[j].toBytes()));
            }

            // Iterate all write sets in the namespace
            RLPReader.RLPItem[] memory writesetsRLP = nsrwRLP[2].toList();
            for (uint k = 0; k < writesetsRLP.length; k++) {
                namespace.writes.push(decodeWriteset(writesetsRLP[k].toBytes()));
            }
        }

        return string(txRLP[0].toBytes());
    }

    function mutateState(bytes32 _chainId, string memory _channelId, Block memory _blk) internal {
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

    function injectBlockHashToTx(bytes32 _chainId, string memory _channelId, string memory _txId, string memory _blockHash) internal {
        Transaction storage txn = m_networks[_chainId].m_channels[_channelId].m_transactions[_txId];
        txn.blockHash = _blockHash;
    }

    function flagTx(bytes32 _chainId, string memory _channelId, string memory _txId) internal {
        m_networks[_chainId].m_channels[_channelId].m_transactions_exist[_txId] = true;
    }

    function decodeReadset(bytes memory _readsetRLP) internal pure returns (ReadSet memory) {
        RLPReader.RLPItem[] memory readsetRLP = _readsetRLP.toRLPItem().toList();

        string memory key = string(readsetRLP[0].toBytes());

        RLPReader.RLPItem[] memory rsv = readsetRLP[1].toList();

        uint blockNo = rsv[0].toUint();
        uint txNo = 0;

        if (rsv.length > 1) {
            txNo = rsv[1].toUint();
        }
        RSVersion memory version = RSVersion(blockNo, txNo);

        return ReadSet(key, version);
    }

    function decodeWriteset(bytes memory _writesetRLP) internal pure returns (WriteSet memory){
        RLPReader.RLPItem[] memory writesetRLP = _writesetRLP.toRLPItem().toList();

        string memory key = string(writesetRLP[0].toBytes());
        string memory value = string(writesetRLP[2].toBytes());

        bool isDelete = false;
        string memory isDeleteStr = string(writesetRLP[1].toBytes());
        if (keccak256(abi.encodePacked(isDeleteStr)) == keccak256(abi.encodePacked("true"))) {
            isDelete = true;
        }

        return WriteSet(key, isDelete, value);
    }

    function isExpectedReadVersion(Namespace memory _namespace, RSVersion memory _version, string memory _key) internal pure returns (bool) {
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

    function getBlock(bytes32 _chainId, string memory _channelId, string memory _blockHash) public view returns (uint, string memory, string memory, string memory, uint, uint, string memory) {
        Block storage blk = m_networks[_chainId].m_channels[_channelId].m_blocks[_blockHash];

        require(keccak256(abi.encodePacked(blk.hash)) != keccak256(abi.encodePacked("")), "Block does not exist.");

        string memory txs = blk.transactions[0];

        for (uint i = 1; i < blk.transactions.length; i++) {
            txs = string(abi.encodePacked(txs, ",", blk.transactions[i]));
        }

        return (blk.number, blk.hash, blk.prevHash, blk.dataHash, blk.timestamp_s, blk.timestamp_nanos, txs);
    }

    function getTransaction(bytes32 _chainId, string memory _channelId, string memory _txId) public view returns (string memory, string memory) {
        Transaction storage txn = m_networks[_chainId].m_channels[_channelId].m_transactions[_txId];

        require(isTransactionExists(_chainId, _channelId, _txId), "Transaction does not exist.");

        string memory ns = txn.namespaces[0];

        for (uint i = 1; i < txn.namespaces.length; i++) {
            ns = string(abi.encodePacked(ns, ",", txn.namespaces[i]));
        }

        return (txn.blockHash, ns);
    }

    function isTransactionExists(bytes32 _chainId, string memory _channelId, string memory _txId) public view returns (bool) {
        return m_networks[_chainId].m_channels[_channelId].m_transactions_exist[_txId];
    }

    function getNSRW(bytes32 _chainId, string memory _channelId, string memory _txId, string memory _namespace) public view returns (string memory, string memory) {
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

    function getState(bytes32 _chainId, string memory _channelId, string memory _key) public view returns (uint, uint, string memory) {
        State storage state = m_networks[_chainId].m_channels[_channelId].m_state[_key];

        require(keccak256(abi.encodePacked(state.key)) != keccak256(abi.encodePacked("")), "Key unrecognised.");

        return (state.version.blockNo, state.version.txNo, state.value);
    }
}