const fs = require("fs-extra")
const Debug = require('web3-eth-debug').Debug
const Web3 = require('web3');

benchmarkTx = (txHash, options) => {
    const web3 = new Web3();
    web3.setProvider(new web3.providers.HttpProvider('http://localhost:8545', {timeout:100000}));

    web3.currentProvider.send({
       "jsonrcp":"2.0",
       "method":"debug_traceTransaction",
       "params":[txHash, options]
    }, (err, res) => console.log(err, res))

}

aggregate = (txTrace, options) => {
    console.log(txTrace)
    aggregateObj = {}
    aggregateObj.name = "checkProof"
    aggregateObj.gas = txTrace.gas
    aggregateObj.maxStackDepth = 0
    aggregateObj.maxMemoryDepth = 0

    if(!options.disableStack){
        // stack max depth
        for (log of txTrace.structLogs) {
            if (log.stack.length > aggregateObj.maxStackDepth)
                aggregateObj.maxStackDepth = log.stack.length
        }
    } else if(!options.disableMemory){
        // memory max depth
        for (log of txTrace.structLogs) {
            if (log.memory.length > aggregateObj.maxMemoryDepth)
                aggregateObj.maxMemoryDepth = log.memory.length
        }
    }

    // write final object containing aggregate
    fs.writeJson("./gas.json", aggregateObj, {flag: "w+"})
}


options = {disableStorage:true, disableStack:true, disableMemory:false, timeout:"1m"}

benchmarkTx("0x6035502a3dcc23be3b7ff9275674db8bd6ff9711255caac878d4faf7fd44b568", options)

module.exports = benchmarkTx