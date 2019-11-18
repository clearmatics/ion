const assert = require('assert');
const path = require('path');
const utils = require('./helpers/utils.js');
const async = require('async');

const bytecode = "0x608060405234801561001057600080fd5b5061113f806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80638c9d0f9e14610030575b600080fd5b6102216004803603608081101561004657600080fd5b810190808035906020019064010000000081111561006357600080fd5b82018360208201111561007557600080fd5b8035906020019184600183028401116401000000008311171561009757600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290803590602001906401000000008111156100fa57600080fd5b82018360208201111561010c57600080fd5b8035906020019184600183028401116401000000008311171561012e57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192908035906020019064010000000081111561019157600080fd5b8201836020820111156101a357600080fd5b803590602001918460018302840111640100000000831117156101c557600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192908035906020019092919050505061023b565b604051808215151515815260200191505060405180910390f35b600061024985858585610253565b9050949350505050565b600061025d611094565b61026685610386565b90506060610273826103e3565b905060008490506000809050606061028c88600061049a565b905060008090505b8451811015610377576102b98582815181106102ac57fe5b60200260200101516107b4565b8051906020012084146102d5576000965050505050505061037e565b60606102f38683815181106102e657fe5b60200260200101516103e3565b90506011815114156103185761030b8185858f610813565b809550819650505061034c565b60028151141561033b5761032e8185858f610918565b809550819650505061034b565b600097505050505050505061037e565b5b6000801b851415610369576001841497505050505050505061037e565b508080600101915050610294565b5050505050505b949350505050565b61038e611094565b60008251905060008114156103bc5760405180604001604052806000815260200160008152509150506103de565b6000602084019050604051806040016040528082815260200183815250925050505b919050565b60606103ee82610a0e565b6103f757600080fd5b600061040283610a40565b90508060405190808252806020026020018201604052801561043e57816020015b61042b6110ae565b8152602001906001900390816104235790505b5091506104496110c8565b61045284610ab8565b905060005b61046082610afe565b156104925761046e82610b28565b84828151811061047a57fe5b60200260200101819052508080600101915050610457565b505050919050565b60608060ff6040519080825280601f01601f1916602001820160405280156104d15781602001600182028038833980820191505090505b509050600080905060008090505b85518110156106fb576104f06110e8565b61050f8783815181106104ff57fe5b602001015160f81c60f81b610b85565b905085801561051e5750600082145b1561064357600160f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168160006002811061055757fe5b60200201517effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614806105e25750600360f81b7effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916816000600281106105b957fe5b60200201517effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916145b1561063e57806001600281106105f457fe5b6020020151848460ff168151811061060857fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506001830192505b6106ed565b8060006002811061065057fe5b6020020151848460ff168151811061066457fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350806001600281106106a057fe5b6020020151846001850160ff16815181106106b757fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053506002830192505b5080806001019150506104df565b5060608160ff166040519080825280601f01601f1916602001820160405280156107345781602001600182028038833980820191505090505b50905060008090505b8260ff168110156107a75783818151811061075457fe5b602001015160f81c60f81b82828151811061076b57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350808060010191505061073d565b5080935050505092915050565b6060600082602001519050806040519080825280601f01601f1916602001820160405280156107f25781602001600182028038833980820191505090505b5091506000811461080d5761080c83600001518383610c4d565b5b50919050565b600080835185141561086a576000610847846108428960108151811061083557fe5b60200260200101516107b4565b610c8f565b610852576000610855565b60015b8160001b91508060ff1690509150915061090f565b600061088b85878151811061087b57fe5b602001015160f81c60f81b610caa565b9050610895611094565b878261ffff16815181106108a557fe5b60200260200101519050600187019650600060206108c2836107b4565b5110156108e2576108d582898989610cbd565b8099508192505050610905565b610902898461ffff16815181106108f557fe5b6020026020010151610d28565b90505b8088945094505050505b94509492505050565b600080606061093a8760008151811061092d57fe5b6020026020010151610d3d565b905061094781600161049a565b518601955084518614156109a157600061097d856109788a60018151811061096b57fe5b6020026020010151610d3d565b610c8f565b61098857600061098b565b60015b8160001b91508060ff1690509250925050610a05565b60006109ae82600161049a565b5114156109cb576000808160001b91508090509250925050610a05565b60606109ea886001815181106109dd57fe5b6020026020010151610d3d565b905060006109f9826000610dad565b90508088945094505050505b94509492505050565b60008082602001511415610a255760009050610a3b565b60008260000151905060c0815160001a10159150505b919050565b6000610a4b82610a0e565b610a585760009050610ab3565b60008083600001519050805160001a91506000610a7485610e22565b82019050600060018660200151840103905060005b818311610aaa57610a9983610eb4565b830192508080600101915050610a89565b80955050505050505b919050565b610ac06110c8565b610ac982610a0e565b610ad257600080fd5b6000610add83610e22565b83600001510190508282600001819052508082602001818152505050919050565b6000610b08611094565b826000015190508060200151816000015101836020015110915050919050565b610b30611094565b610b3982610afe565b15610b7b576000826020015190506000610b5282610eb4565b905081836000018181525050808360200181815250508082018460200181815250505050610b80565b600080fd5b919050565b610b8d6110e8565b6000610b9a836004610f4e565b90506000600f60f81b841690506040518060400160405280837effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff191681525092505050919050565b6020601f820104836020840160005b83811015610c7c5760208102808401518184015250600181019050610c5c565b5060008551602001860152505050505050565b60008180519060200120838051906020012014905092915050565b6000806002830151905080915050919050565b6000806060610ccb876103e3565b9050601181511415610ced57610ce381878787610813565b9250925050610d1f565b600281511415610d0d57610d0381878787610918565b9250925050610d1f565b6000808160001b915080905092509250505b94509492505050565b6000610d3382610f77565b60001b9050919050565b6060610d4882610fd6565b610d5157600080fd5b600080610d5d84611007565b8092508193505050806040519080825280601f01601f191660200182016040528015610d985781602001600182028038833980820191505090505b509250610da6828483610c4d565b5050919050565b60008060008090505b6020811015610e17576008810260ff60f81b8683870181518110610dd657fe5b602001015160f81c60f81b167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916901c821791508080600101915050610db6565b508091505092915050565b60008082602001511415610e395760009050610eaf565b60008083600001519050805160001a91506080821015610e5e57600092505050610eaf565b60b8821080610e7a575060c08210158015610e79575060f882105b5b15610e8a57600192505050610eaf565b60c0821015610ea357600160b783030192505050610eaf565b600160f7830301925050505b919050565b600080825160001a90506080811015610ed05760019150610f48565b60b8811015610ee757600160808203019150610f47565b60c0811015610f115760b78103806020036101000a60018501510480820160010193505050610f46565b60f8811015610f2857600160c08203019150610f45565b60f78103806020036101000a600185015104808201600101935050505b5b5b5b50919050565b60008160ff16600260ff160a60ff168360f81c60ff1681610f6b57fe5b0460f81b905092915050565b6000610f8282610fd6565b610f8b57600080fd5b600080610f9784611007565b80925081935050506020811115610fad57600080fd5b6000811415610fc157600092505050610fd1565b806020036101000a825104925050505b919050565b60008082602001511415610fed5760009050611002565b60008260000151905060c0815160001a109150505b919050565b60008061101383610fd6565b61101c57600080fd5b60008084600001519050805160001a9150608082101561104a5780935060019250838393509350505061108f565b60b88210156110685760018560200151039250600181019350611086565b600060b7830390508060018760200151030393506001818301019450505b83839350935050505b915091565b604051806040016040528060008152602001600081525090565b604051806040016040528060008152602001600081525090565b60405180604001604052806110db6110ae565b8152602001600081525090565b604051806040016040528060029060208202803883398082019150509050509056fea265627a7a7231582005a69cc9513dbd6b7b562bda66b69949785d93e43d5be8022b7666ca34e5a06e64736f6c634300050d0032";
const ABI = [{
    "constant":true,
    "inputs":[{
        "internalType":"bytes",
        "name":"_value",
        "type":"bytes"
    }, {
        "internalType":"bytes",
        "name":"_parentNodes",
        "type":"bytes"
    }, {
        "internalType":"bytes",
        "name":"_path",
        "type":"bytes"
    }, {
        "internalType":"bytes32",
        "name":"_root",
        "type":"bytes32"
    }],
    "name":"testVerify",
    "outputs":[{"internalType":"bool","name":"","type":"bool"}],
    "payable":false,
    "stateMutability":"view",
    "type":"function"
}];

contract('Patricia Trie', (accounts) => {
    describe('VerifyProof', async function () {
        it('should successfully verify all proofs', async function () {

            let contract = new web3.eth.Contract(ABI);
            await contract.deploy({data: bytecode})
            .send({from: accounts[0], gas: "0xFFFFFFFFFFFF"})
            .then((contractInstance) => {
                patriciatrietest = contractInstance;

                testData['success'].forEach( (data) => {
                    patriciatrietest.methods.testVerify(data.value, data.nodes, data.path, data.rootHash).send({from: accounts[0], gas: "0xFFFFFFFFFFFF"}, function (result) {
                        assert.ifError(result);
                        console.log(result);
                        assert.equal(result, true);
                    });
                })
            })
        });

        it('should fail verifying all proofs with incompatible data', async function () {
            let contract = new web3.eth.Contract(ABI);
            await contract.deploy({data: bytecode})
            .send({from: accounts[0], gas: "0xFFFFFFFFFFFF"})
            .then((contractInstance) => {
                patriciatrietest = contractInstance;

                testData['fail'].forEach( async (data) => {
                    await patriciatrietest.methods.testVerify(data.value, data.nodes, data.path, data.rootHash).send({from: accounts[0], gas: "0xFFFFFFFFFFFF"}, function (result) {
                        assert.ifError(result);
                        console.log(result);
                        assert.equal(result, false);
                    });
                })
            })
        });

    });
});

const testData = {
    "success": [{
        "rootHash": "0xda2e968e25198a0a41e4dcdc6fcb03b9d49274b3d44cb35d921e4ebe3fb5c54c",
        "path": "0x61",
        "value": "0x857465737431",
        "nodes": "0xf83bf839808080808080c8318685746573743180a0207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a978080808080808080"
    }, {
        "rootHash": "0xda2e968e25198a0a41e4dcdc6fcb03b9d49274b3d44cb35d921e4ebe3fb5c54c",
        "path": "0x826162",
        "value": "0x74",
        "nodes": "0xf87ff839808080808080c8318685746573743180a0207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a978080808080808080f8428080c58320616274cc842061626386857465737433a05d495bd9e35ab0dab60dec18b21acc860829508e7df1064fce1f0b8fa4c0e8b2808080808080808080808080"
    }, {
        "rootHash": "0xda2e968e25198a0a41e4dcdc6fcb03b9d49274b3d44cb35d921e4ebe3fb5c54c",
        "path": "0x83616263",
        "value": "0x857465737433",
        "nodes": "0xf87ff839808080808080c8318685746573743180a0207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a978080808080808080f8428080c58320616274cc842061626386857465737433a05d495bd9e35ab0dab60dec18b21acc860829508e7df1064fce1f0b8fa4c0e8b2808080808080808080808080"
    }, {
        "rootHash": "0xda2e968e25198a0a41e4dcdc6fcb03b9d49274b3d44cb35d921e4ebe3fb5c54c",
        "path": "0x8461626564",
        "value": "0x857465737435",
        "nodes": "0xf8cbf839808080808080c8318685746573743180a0207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a978080808080808080f8428080c58320616274cc842061626386857465737433a05d495bd9e35ab0dab60dec18b21acc860829508e7df1064fce1f0b8fa4c0e8b2808080808080808080808080e583161626a06b1a1127b4c489762c8259381ff9ecf51b7ef0c2879b89e72c993edc944f1ccce5808080ca8220648685746573743480ca822064868574657374358080808080808080808080"
    }, {
        "rootHash": "0xda2e968e25198a0a41e4dcdc6fcb03b9d49274b3d44cb35d921e4ebe3fb5c54c",
        "path": "0x8461626364",
        "value": "0x857465737434",
        "nodes": "0xf8cbf839808080808080c8318685746573743180a0207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a978080808080808080f8428080c58320616274cc842061626386857465737433a05d495bd9e35ab0dab60dec18b21acc860829508e7df1064fce1f0b8fa4c0e8b2808080808080808080808080e583161626a06b1a1127b4c489762c8259381ff9ecf51b7ef0c2879b89e72c993edc944f1ccce5808080ca8220648685746573743480ca822064868574657374358080808080808080808080"
    }],
    "fail": [{
        "rootHash": "0xda2e968e25198a0a41e4dcdc6fcb03b9d49274b3d44cb35d921e4ebe3fb5c54c",
        "path": "0x61",
        "value": "0x857465737432",
        "nodes": "0xf83bf839808080808080c8318685746573743180a0207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a978080808080808080"
    }, {
        "rootHash": "0xda2e968e25198a0a41e4dcdc6fcb03b9d49274b3d44cb35d921e4ebe3fb5c54c",
        "path": "0x826163",
        "value": "0x75",
        "nodes": "0xf87ff839808080808080c8318685746573743180a0207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a978080808080808080f8428080c58320616274cc842061626386857465737433a05d495bd9e35ab0dab60dec18b21acc860829508e7df1064fce1f0b8fa4c0e8b2808080808080808080808080"
    }, {
        "rootHash": "0xda2e968e25198a0a41e4dcdc6fcb03b9d49274b3d44cb35d921e4ebe3fb5c54c",
        "path": "0x83616263",
        "value": "0x857465737434",
        "nodes": "0xf87ff839808080808080c8318685746573743180a0207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a978080808080808080f8428080c58320616274cc842061626386857465737433a05d495bd9e35ab0dab60dec18b21acc860829508e7df1064fce1f0b8fa4c0e8b2808080808080808080808080"
    }, {
        "rootHash": "0xda2e968e25198a0a41e4dcdc6fcb03b9d49274b3d44cb35d921e4ebe3fb5c54c",
        "path": "0x8461626564",
        "value": "0x857465737435",
        "nodes": "0xf8cbf839808080808080c8318685746573743180a0207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a978080808080808080f8428080c58320616274cc842061626386857465737433a05d495bd9e35ab0dab60dec18b21acc860829508e7df1064fce1f0b8fa4c0e8b2808080808080808080808080e583161626a06b1a1127b4c489762c8259381ff9ecf51b7ef0c2879b89e72c993edc944f1ccce5808080ca8220648685746573743480ca822064868574657374358080808080808080808085"
    }, {
        "rootHash": "0xda2e968e25198a0a41e4dcdc6fcb03b9d49274b3d44cb35d921e4ebe3fb5c54c",
        "path": "0x8461626364",
        "value": "0x857465737435",
        "nodes": "0xf8cbf839808080808080c8318685746573743180a0207947cf85c03bd3d9f9ff5119267616318dcef0e12de2f8ca02ff2cdc720a978080808080808080f8428080c58320616274cc842061626386857465737433a05d495bd9e35ab0dab60dec18b21acc860829508e7df1064fce1f0b8fa4c0e8b2808080808080808080808080e583161626a06b1a1127b4c489762c8259381ff9ecf51b7ef0c2879b89e72c993edc944f1ccce5808080ca8220648685746573743480ca822064868574657374358080808080808080808080"
    }]
}