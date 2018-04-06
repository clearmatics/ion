const web3Utils = require('web3-utils')
const web3Abi = require('web3-eth-abi');

const merkle = require('./merkle');

const Token = artifacts.require("Token");
const Sodium = artifacts.require("Sodium");
const Fluoride = artifacts.require("Fluoride");

contract('Fluoride', (accounts) => {
  const joinHex = arr => '0x' + arr.map(el => el.slice(2)).join('')
  const generateState = (expire,tokenAddr,amount) => web3Utils.sha3(joinHex([expire, tokenAddr, amount]))

  it('startA startB', async () => {
    const token = await Token.new();
    const sodium = await Sodium.new();
		const fluoride = await Fluoride.new(sodium.address);

    const fluorideAddrA = fluoride.address
    const fluorideAddrB = fluoride.address

    const unixTimestamp = Math.round((new Date()).getTime() / 1000)
    const expireHex = min => '0x' + web3Utils.toHex(unixTimestamp + min).slice(2).padStart(64,'0')

    const expireA = expireHex(60) // need to add some proper time to here
    const expireB = expireHex(60) // need to add some proper time to here

    const tokenA = token.address
    const tokenB = token.address

    const amountA = '0x' + web3Utils.toHex(100).slice(2).padStart(64,'0')
    const amountB = '0x' + web3Utils.toHex(100).slice(2).padStart(64,'0')

    const addrA = accounts[0]
    const addrB = accounts[1]

    const stateA = generateState(expireA, tokenA, amountA)
    const stateB = generateState(expireB, tokenB, amountB)

    const hashA = web3Utils.sha3(joinHex([fluorideAddrA,stateA]))
    const hashB = web3Utils.sha3(joinHex([fluorideAddrB,stateB]))

    const sigA = web3.eth.sign(addrA,hashA)

    const hashAB = web3Utils.sha3(joinHex([hashA,addrA,hashB]))
    const sigB = web3.eth.sign(addrB,hashAB)

    const hashABC = web3Utils.sha3(joinHex([hashAB,addrB]))
    const sigC = web3.eth.sign(addrA,hashABC)

    assert.equal(hashABC,web3Utils.soliditySha3(hashAB,addrB),'hashing badly when compared with solidity')


    // the two start functions should be done on different chains?
    const receiptStartOnAbyA = await fluoride.Start_OnAbyA(
      fluorideAddrA, expireA, tokenA, amountA, sigA,
      fluorideAddrB, stateB, sigB, sigC)
    const tradeIdA = receiptStartOnAbyA.logs.find(l => l.event === 'OnDeposit').args.trade_id

    const receiptStartOnBbyB = await fluoride.Start_OnBbyB(
      fluorideAddrA, stateA, sigA,
      fluorideAddrB, expireB, tokenB, amountB, sigB, sigC)
    const tradeIdB = receiptStartOnBbyB.logs.find(l => l.event === 'OnDeposit').args.trade_id
    assert.equal(tradeIdA,tradeIdB,'tradeIds are different')
  })

  it('withdraw', async () => {
    // ChainA A -PAY-> B
    // ChainB B -PAY-> A
    const tokenA = await Token.new();
    const sodiumA = await Sodium.new();
		const fluorideA = await Fluoride.new(sodiumA.address);
    const tokenB = await Token.new();
    const sodiumB = await Sodium.new();
		const fluorideB = await Fluoride.new(sodiumB.address);

    const fluorideAddrA = fluorideA.address
    const fluorideAddrB = fluorideB.address

    const unixTimestamp = Math.round((new Date()).getTime() / 1000)

    const expireA = '0x' + web3Utils.toHex(unixTimestamp + 60).slice(2).padStart(64,'0') // need to add some proper time to here
    const expireB = '0x' + web3Utils.toHex(unixTimestamp + 60).slice(2).padStart(64,'0') // need to add some proper time to here

    const tokenAddrA = tokenA.address
    const tokenAddrB = tokenB.address

    const intAmountA = 123
    const intAmountB = 93

    const amountA = '0x' + web3Utils.toHex(intAmountA).slice(2).padStart(64,'0')
    const amountB = '0x' + web3Utils.toHex(intAmountB).slice(2).padStart(64,'0')

    const addrA = accounts[1]
    const addrB = accounts[2]

    const stateA = generateState(expireA, tokenAddrA, amountA)
    const stateB = generateState(expireB, tokenAddrB, amountB)

    const hashA = web3Utils.sha3(joinHex([fluorideAddrA,stateA]))
    const hashB = web3Utils.sha3(joinHex([fluorideAddrB,stateB]))

    const sigA = web3.eth.sign(addrA,hashA)

    const hashAB = web3Utils.sha3(joinHex([hashA,addrA,hashB]))
    const sigB = web3.eth.sign(addrB,hashAB)

    const hashABC = web3Utils.sha3(joinHex([hashAB,addrB]))
    const sigC = web3.eth.sign(addrA,hashABC)

    // the two start functions should be done on different chains?
    const receiptStartOnAbyA = await fluorideA.Start_OnAbyA(
      fluorideAddrA, expireA, tokenAddrA, amountA, sigA,
      fluorideAddrB, stateB, sigB, sigC, { from: addrA })

    const receiptStartOnBbyB = await fluorideB.Start_OnBbyB(
      fluorideAddrA, stateA, sigA,
      fluorideAddrB, expireB, tokenAddrB, amountB, sigB, sigC, { from: addrB })
    const tradeId = receiptStartOnAbyA.logs.find(l => l.event === 'OnDeposit').args.trade_id

    // mint and transfer token
    const tokenOwner = accounts[0]
    const totalSupply = 1000

    const receiptMintA = await tokenA.mint(totalSupply,{ from: tokenOwner })
    const receiptTansferA = await tokenA.transfer(addrA,totalSupply/2,{ from: tokenOwner })
    const receiptMintB = await tokenB.mint(totalSupply,{ from: tokenOwner })
    const receiptTansferB = await tokenB.transfer(addrB,totalSupply/2,{ from: tokenOwner })

    const tokenABalA0 = await tokenA.balanceOf(addrA)
    const tokenABalB0 = await tokenA.balanceOf(addrB)
    const tokenBBalA0 = await tokenB.balanceOf(addrA)
    const tokenBBalB0 = await tokenB.balanceOf(addrB)

    //transfer value to contract on each chain
    const overloadedTransferAbi = {
      "constant": false,
      "inputs": [
        { "name": "_to", "type": "address" },
        { "name": "_value", "type": "uint256" },
        { "name": "_data", "type": "bytes" }
      ],
      "name": "transfer",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    }
    const transferMethodTransactionDataA = web3Abi
      .encodeFunctionCall( overloadedTransferAbi, [ fluorideAddrA,intAmountA, tradeId ]);
    const transferMethodTransactionDataB = web3Abi
      .encodeFunctionCall( overloadedTransferAbi, [ fluorideAddrB,intAmountB, tradeId ]);
    const receiptPayA = await web3.eth
      .sendTransaction({from: addrA, to: tokenA.address, data: transferMethodTransactionDataA, value: 0});
    const receiptPayB = await web3.eth
      .sendTransaction({from: addrB, to: tokenB.address, data: transferMethodTransactionDataB, value: 0});

    const tokenABalA1 = await tokenA.balanceOf(addrA)
    const tokenABalB1 = await tokenA.balanceOf(addrB)
    const tokenBBalA1 = await tokenB.balanceOf(addrA)
    const tokenBBalB1 = await tokenB.balanceOf(addrB)

    // UPDATE SODIUM (this should be done by event relay/lithium)
    const expectEventA = joinHex([fluorideAddrB/*, topic*/,tradeId])
    const expectEventB = joinHex([fluorideAddrA/*, topic*/,tradeId])

    const testDataA = [expectEventA,"2","3","4","5","6","7"]
    const treeA = merkle.createMerkle(testDataA)
    const pathA = testDataA.map(value => merkle.pathMerkle(value,treeA[0]))

    const leafHashA = merkle.merkleHash(testDataA[0])
    const rootArgA = treeA[1]

    const nextBlockA = await sodiumA.NextBlock()
    const receiptUpdateA = await sodiumA.Update(nextBlockA,[rootArgA])
    const validA = await sodiumA.Verify(nextBlockA,leafHashA,pathA[0])
    assert(validA,'SodiumA.verify() failed!')

    const testDataB = [expectEventB,"2","3","4","5","6","7"]
    const treeB = merkle.createMerkle(testDataB)
    const pathB = testDataB.map(value => merkle.pathMerkle(value,treeB[0]))

    const leafHashB = merkle.merkleHash(testDataB[0])
    const rootArgB = treeB[1]

    const nextBlockB = await sodiumB.NextBlock()
    const receiptUpdateB = await sodiumB.Update(nextBlockB,[rootArgB])
    const validB = await sodiumB.Verify(nextBlockB,leafHashB,pathB[0])
    assert(validB,'SodiumB.verify() failed!')

    // WITHDRAW VALUE
    const receiptWithDrawA = await fluorideA.Withdraw(tradeId, nextBlockA, pathA[0], { from: addrB })
    const receiptWithDrawB = await fluorideB.Withdraw(tradeId, nextBlockB, pathB[0], { from: addrA })

    const tokenABalA2 = await tokenA.balanceOf(addrA)
    const tokenABalB2 = await tokenA.balanceOf(addrB)
    const tokenBBalA2 = await tokenB.balanceOf(addrA)
    const tokenBBalB2 = await tokenB.balanceOf(addrB)

    // asserts about the state of the tokens

    /*
    console.log('TokenA Balance A',
      'Initial:',tokenABalA0.toString(10),
      'Deposit:',tokenABalA1.toString(10),
      'Withdraw:',tokenABalA2.toString(10))
    console.log('TokenA Balance B',
      'Initial:',tokenABalB0.toString(10),
      'Deposit:',tokenABalB1.toString(10),
      'Withdraw:',tokenABalB2.toString(10))
    console.log('TokenB Balance A',
      'Initial:',tokenBBalA0.toString(10),
      'Deposit:',tokenBBalA1.toString(10),
      'Withdraw:',tokenBBalA2.toString(10))
    console.log('TokenB Balance B',
      'Initial:',tokenBBalB0.toString(10),
      'Deposit:',tokenBBalB1.toString(10),
      'Withdraw:',tokenBBalB2.toString(10))
      */

    assert(tokenABalA0 == totalSupply/2)
    assert(tokenBBalA0 == 0)
    assert(tokenABalB0 == 0)
    assert(tokenBBalB0 == totalSupply/2)

    assert(tokenABalA1 == ((totalSupply/2) - intAmountA))
    assert(tokenBBalA1 == 0)
    assert(tokenABalB1 == 0)
    assert(tokenBBalB1 == ((totalSupply/2) - intAmountB))

    assert(tokenABalA2 == ((totalSupply/2) - intAmountA))
    assert(tokenBBalA2 == intAmountB)
    assert(tokenABalB2 == intAmountA)
    assert(tokenBBalB2 == ((totalSupply/2) - intAmountB))


    /*
    console.log('tokenOwner:',tokenOwner)
    console.log('addrA:',addrA)
    console.log('addrB:',addrB)
    console.log('addrA:',addrA)
    console.log('tokenA:',tokenA)
    console.log('amountA:',amountA)
    console.log('expireA:',expireA)
    console.log('stateA:',stateA)
    console.log('hashA:',hashA)
    console.log('hashB:',hashB)
    console.log('hashAB:',hashAB)
    console.log('hashABC:',hashABC)
    */
    //fluoride.Start_OnBbyB(fluorideAddrB, expireB, tokenB, amountB, sigB, fluorideAddrA, stateB, sigA, sigC)

    // fluoride.Start_OnAbyA(
    //     address a_contract, uint a_expire, address a_token, uint256 a_amount, bytes a_sig,
    //     address b_contract, bytes32 b_state, bytes b_sig,
    //     bytes c_sig )
    //     a_hash = keccak256(a_contract, keccak256(a_expire, a_token, a_amount))
    //     b_hash = keccak256(b_contract, b_state)
    //
    // fluoride.Start_OnBbyB(
    //     address a_contract, bytes32 a_state, bytes a_sig,
    //     address b_contract, uint256 b_expire, address b_token, uint256 b_amount, bytes b_sig,
    //     bytes c_sig )
		//     a_hash = keccak256(a_contract, a_state),
		//     b_hash = keccak256(b_contract, keccak256(b_expire, b_token, b_amount)),

    //const tradeAgreeResult = await fluoride.VerifyTradeAgreement(hashA,sigA,hashB,sigB,sigC)

    // assert.fail('reached the END!!! this is here to check EVENTS')
  })

  // TODO: TEST CANCEL
  // TODO: TEST REFUND
});

