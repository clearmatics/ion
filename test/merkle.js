const Web3Utils = require('web3-utils');
const BN = require('bignumber.js')

// TODO: too complicated for what it is, simplify it!!
const bitClear = (n,b) => {
  let bitStr = n.toString(2).padStart(b,'0')
  const firstBitIdx = bitStr.length-b
  const firstBit = Number(bitStr[firstBitIdx])
  if(firstBit) {
    bitStr = bitStr.substr(0,firstBitIdx) + (firstBit ^ 1) + bitStr.substr(firstBitIdx+1)
  }
  return new BN(bitStr,2).toString(10)
}
const hash = (data) => new BN(Web3Utils.sha3(data)).toString(10)
const merkleHash = (data) => bitClear(new BN(hash(data)), 0xFF)

// TODO: doesn't solve odd sized item arrays (it ignores the last element)
const treeLevel = (items) => items.reduce((prev,el, idx, arr) => {
  if(idx % 2) {
    const h1 = new BN(arr[idx-1])
    const h2 = new BN(el)
    const h1hex = h1.toString(16).padStart(64,'0')
    const h2hex = h2.toString(16).padStart(64,'0')
    return prev.concat(merkleHash('0x' + h1hex + h2hex))
  }
  return prev
} , [])

const createMerkle = (items) => {
  const leafHash = items.map((leaf) => merkleHash(leaf))
    .sort((a,b) => (new BN(a)) - (new BN(b)))
  if(leafHash.length % 2 !== 0) leafHash.push(merkleHash('merkle-tree-extra'))
  const tree = [ leafHash ]
  while (tree[0].length !== 1) {
    tree.unshift(treeLevel(tree[0]))
  }
  const root = tree[0];
  return [tree.sort(l => 1)].concat(root)
}

const bitSet = (n,b) => {
  const bnN = new BN(n)
  let resStr = bnN.toString(2).padStart(b,'0')
  const idx = resStr.length - b
  resStr = resStr.substr(0,idx) + '1' + resStr.substr(idx+1)
  const bnRes = new BN(resStr,2)
  return bnRes.toString(10)
}

const bitTest = (n,b) => {
  const bnN = new BN(n)
  let resStr = bnN.toString(2).padStart(b,'0')
  const idx = resStr.length - b
  return (Number(resStr[idx]) === 1)
}

// expects a leaf node
const pathMerkle = (item,tree) => {
  const itemHash = merkleHash(item)
  let idx = tree[0].findIndex(h => h === itemHash)

  const result = tree.slice(0,-1).reduce((path,level) => {
    const hash = (idx % 2) ? level[idx-1] : bitSet(level[idx+1],0xFF)
    idx = idx >> 1 // divide by 2
    return path.concat(hash)
  },[])

  return result
}

const proofMerkle = (leaf, path, root) => {
  const leafHash = merkleHash(leaf)
  const result = path.reduce((prev,item) => {
    const isOdd = bitTest(item, 0xFF)
    const hashStr1 = isOdd ?  prev : item
    const hashStr2 = isOdd ?  bitClear(new BN(item), 0xFF) : prev
    const h1 = new BN(hashStr1)
    const h2 = new BN(hashStr2)
    const h1hex = h1.toString(16).padStart(64,'0')
    const h2hex = h2.toString(16).padStart(64,'0')
    return merkleHash('0x' + h1hex + h2hex)
  }, leafHash)
  return (result === root)
}

module.exports = { createMerkle, treeLevel, hash, merkleHash, bitClear, bitSet, pathMerkle, proofMerkle }
