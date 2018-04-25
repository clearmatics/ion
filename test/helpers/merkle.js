// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

const Web3Utils = require('web3-utils');
const BN = require('bignumber.js')

/*
const bnjs = require('bn.js')
const bitTest = (n,b) => (new bnjs(n.toString(16),16)).testn(b)
const bitClear = (n,b) => bitTest(n,b) ? new BN((new bnjs(n.toString(16),16)).xor((new bnjs(1)).bincn(b))) : n
const bitSet = (n,b) => new BN((new bnjs(n.toString(16),16)).setn(b))
*/

const toggleBit = (n,b,bitValue) => {
  //b += 1
  let resStr = n.toString(2).padStart(b,'0')
  const idx = resStr.length - b
  resStr = resStr.substr(0,idx) + bitValue + resStr.substr(idx+1)
  const bnRes = new BN(resStr,2)
  return bnRes
}


const bitTest = (n,b) => {
  //b += 1
  const resStr = n.toString(2).padStart(b,'0')
  const idx = resStr.length - b
  return Number(resStr[idx]) === 1
}


const bitClear = (n,b) => toggleBit(n,b,'0')


const bitSet = (n,b) => toggleBit(n,b,'1')


const toHex = n => {
  let nHex = n.toString(16)
  if(nHex.length <= 64)
    nHex = '0x' + nHex.padStart(64,'0')
  else
    nHex = '0x' + nHex.padStart(64*2,'0')
  return nHex
}


const joinHex2BN = (a,b) => new BN('0x' + toHex(a) + toHex(b).substring(2))


const hash = n => new BN(Web3Utils.sha3(Web3Utils.isBigNumber(n) ? toHex(n) : n))


const merkleHash = n => bitClear(hash(n), 0xFF)


const treeLevel = items => items.reduce((prev,el, idx, arr) => {
  if(idx % 2)
    return prev.concat(merkleHash(joinHex2BN(arr[idx-1],el)))
  return prev
} , [])


const createMerkle = (items) => {
  const extraHash = merkleHash('merkle-tree-extra')
  const leafHash = items
    .map((leaf) => merkleHash(leaf))
    .sort((a,b) => (new BN(a)) - (new BN(b)))
  if(leafHash.length % 2 !== 0) leafHash.push(extraHash)
  const tree = [ leafHash ]
  while (tree[0].length !== 1) {
    const level = treeLevel(tree[0])
    if(level.length !== 1 && level.length % 2 !== 0) level.push(extraHash) //levels need to be pair
    tree.unshift(level)
  }
  const root = tree[0];
  return [tree.sort(l => 1)].concat(root)
}


const pathMerkle = (leaf,tree) => {
  const leafHash = merkleHash(leaf)
  let idx = tree[0].findIndex(h => h.equals(leafHash))
  //console.log(tree.map(l => l.map(b=>(b||'0x0').toString(16))))

  const result = tree
    .slice(0,-1)
    .reduce((path,level) => {
      const hash = (idx % 2) ? level[idx-1] : bitSet(level[idx+1],0xFF)
      idx = idx >> 1 // divide by 2
      return path.concat(hash)
    },[])

  return result
}


const proofMerkle = (leaf, path, root, hashLeaf,debug) => {
  const leafHash = hashLeaf ? leaf : merkleHash(leaf)
  const result = path
    .reduce((prev,item) => {
      const bitSide = bitTest(item, 0xFF)
      const h1 = bitSide ?  prev : item
      const h2 = bitSide ?  bitClear(item, 0xFF) : prev
      const hashValue = merkleHash(joinHex2BN(h1,h2))
      if(debug) {
        console.log(bitSide)
        console.log(h1.toString(16).length,h2.toString(16).length, h1.toString(16),h2.toString(16))
        console.log(joinHex2BN(h1,h2).toString(16).length)
        console.log(hashValue.toString(16).length,hashValue.toString(16))
        console.log('=======================================================')
      }
      return hashValue
    }, leafHash)
  return (result.equals(root))
}


const merkle = {
  createMerkle,
  treeLevel,
  hash,
  merkleHash,
  bitClear,
  bitSet,
  pathMerkle,
  proofMerkle,
  bitTest,
}
module.exports = merkle
