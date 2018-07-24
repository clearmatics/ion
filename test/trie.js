var Trie = require('merkle-patricia-tree')
levelup = require('levelup')
const rlp = require('rlp');
const HDWalletProvider = require('truffle-hdwallet-provider');
const mnemonic = "select they inform invite result believe equal daughter front arrest wagon miss same menu twenty";

const Web3 = require('web3');
const EP = require('eth-proof');
web3 = new Web3(new HDWalletProvider(mnemonic, 'https://rinkeby.infura.io/v3/973b00227ca84ced8266b2ab6d7592cb'));

function processLogs(logs) {
    var rlpLogs = [];
    for (var i = 0; i < logs.length; i++) {
        var address = Buffer.from(logs[i].address.slice(2), 'hex')
        var topics = logs[i].topics.map(topic => Buffer.from(topic.slice(2), 'hex'))
        var data = Buffer.from(logs[i].data.slice(2), 'hex')
        rlpLogs.push([address, topics, data]);
    }
    return rlpLogs;
}

async function getReceipt(txHash) {
    receipt = await web3.eth.getTransactionReceipt(txHash);

    var cummulativeGas = Buffer.from(receipt.cumulativeGasUsed.toString('16'), 'hex')
    var bloomFilter = Buffer.from(receipt.logsBloom.slice(2), 'hex')
    var setOfLogs = processLogs(receipt.logs);

    if (receipt.status !== undefined && receipt.status != null){
        var status = receipt.status ? Buffer.from('01', 'hex') : Buffer.from('00', 'hex');
        var rawReceipt = rlp.encode([status,cummulativeGas,bloomFilter,setOfLogs]);
        return rawReceipt;
      } else {
        var postTransactionState = strToBuf(receipt.root)
        var rawReceipt = rlp.encode([postTransactionState, cummulativeGas,bloomFilter,setOfLogs])
        return rawReceipt;
    }
}

async function addReceiptToTrie(trie, txHash, prevTrieRoot) {
    console.log("Inserting receipt for hash: " + txHash);
    receipt = await getReceipt(txHash);
    console.log("RECEIPT RLP:" + receipt.toString('hex'))

    tx = await web3.eth.getTransaction(txHash);
    path = rlp.encode(tx.transactionIndex);
    console.log("Key: " + path.toString('hex'));

    console.log(trie.root)
    console.log(prevTrieRoot)
    do {
        await trie.put(path, receipt, (error) => {console.log("ERROR:" + error)});
    } while (trie.root == prevTrieRoot)
}

async function generateReceiptProof(block) {
    block = await web3.eth.getBlock(block);

    txs = block.transactions;

    var trie = new Trie();
    var lastTrieRoot = ""
    for (let i = 0; i < 2; i++) {
        await addReceiptToTrie(trie, txs[i], lastTrieRoot);
        lastTrieRoot = trie.root;
    }
}

//generateReceiptProof("0x694752333dd1bd0f806cc6ef1063162f4f330c88f9dcd9e61174fcf5e4927eb7");

let root = Buffer.from("f871a012d378fe6800bc18f22e715a31971ef7e73ac5d1d85384f4b66ac32036ae43dea004d6e2678656a957ac776dbef512a04d266c1af3e2c5587fd233261a3d423213808080808080a05fac317a4d6d78181319fbc7e2cae4a9260f1a6afb5c6fea066e2308eed416818080808080808080", 'hex')
let second = Buffer.from("f9016b20b90167f901640183252867b9010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000010000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000f85af8589461621bcf02914668f8404c1f860e92fc1893f74ce1a027a9902e06885f7c187501d61990eae923b37634a8d6dda55a04dc7078395340a0000000000000000000000000279884e133f9346f2fad9cc158222068221b613e", 'hex')
let leaf = Buffer.from("f90151a03da235c6dd0fbdaf208c60cbdca0d609dee2ba107495aa7adaa658362616c8aaa09ebf378a9064aa4da0512c55c790a5e007ac79d2713e4533771cd2c95be47a4da0c06fed36ffe1f2ec164ba88f73b353960448d2decbb65355c5298a33555de742a0e057afe423ee17e5499c570a56880b0f5b5c1884b90ff9b9b5baa827f72fc816a093e06093cd2fdb67e0f87cfcc35ded2f445cc1309a0ff178e59f932aeadb6d73a0193e4e939fbc5d34a570bea3fff7c6d54adcb1c3ab7ef07510e7bd5fcef2d4b3a0a17a0c71c0118092367220f65b67f2ba2eb9068ff5270baeabe8184a01a37f14a03479a38e63123d497588ad5c31d781276ec8c11352dd3895c8add34f9a2b786ba042254728bb9ab94b58adeb75d2238da6f30382969c00c65e55d4cc4aa474c0a6a03c088484aa1c73b8fb291354f80e9557ab75a01c65d046c2471d19bd7f2543d880808080808080", 'hex')

nodes = rlp.encode([root, second, leaf]);

console.log(nodes.toString('hex'))