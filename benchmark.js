const fs = require("fs-extra")
const Debug = require('web3-eth-debug').Debug
const Web3 = require('web3');

benchmarkTx = (txHash, name, options, web3) => {
    web3.currentProvider.send({
       "jsonrcp":"2.0",
       "method":"debug_traceTransaction",
       "params":[txHash, options]
    }, (err, res) => {
        if(err===null && !res.error){
            aggregate(res.result, name, options)
        } else {
            console.log(err, name + "-" + res.error.message)
        }
    })

}

aggregate = (txTrace, name, options) => {
    aggregateObj = {}
    aggregateObj.gas = txTrace.gas
    aggregateObj.maxStackDepth = 0

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

    // word count the opcodes
    for (log of txTrace.structLogs) {
        aggregateObj[log.op] = aggregateObj[log.op] ? aggregateObj[log.op] += 1 : 1
    }

    fs.ensureFile(FILEPATH)
    .then(() => {
      data = fs.readFileSync(FILEPATH)
      data = data.length === 0 ? {} : JSON.parse(data)
      data[name] = aggregateObj
      
      fs.writeJsonSync(FILEPATH, data)
    })
    .catch(err => {
      console.error(err)
    })
}

// ENTRYPOINT 
const web3 = new Web3();
web3.setProvider(new web3.providers.HttpProvider('http://localhost:8545', {timeout:100000}));

const FILEPATH = "./gas-public.json"

options = {disableStorage:true, disableStack:true, disableMemory:true}

// data = fs.readJsonSync("./txMaps.json")

// for (var key of Object.keys(data)) {
//     benchmarkTx(data[key], key, options, web3)
// }

benchmarkTx("0x22eb8a81dd3d949b2e36e38a9f5221c399c3faf0b4d07e56ca3ee8673a97fa7b", "register", options, web3)

module.exports = benchmarkTx