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
  return tree
}

module.exports = { createMerkle, treeLevel, hash, merkleHash, bitClear }
/*
 * hashsResult = new bn(web3.sha3('7'))
 * merkl_hash??
 *
 *

 the extra leaf is extra = merkle_hash("merkle-tree-extra")

 merkle.merkle_tree(["1","2","3","4","5","6","7"])



 [
  [
    [
      8568612641526826488487436752726739043287191320122540356069953783894380777505,
      8763638472773768691201326883407021568462294246273894496415427229083082408032,
      19224855404247632006917173431419498680506051063941070371722880450128577361118,
      61795459977501490647348212754130855970016313872340374962921336716751708851142,
      64645341593328157176709656265449880868558868673380425455960412802858937540801,
      74330811247603495249613868516695563873247293176611122272199330092769797099053,
      78469846343542442363028680824980501212021332975324075417961003849793346933925,
      75317570447191171753008806478868650352148013528306361601609880810432714200529
    ],
    [
      6560824545851281876686151142367952893930617484325436481370811303698242675212,
      14094329272021934754728783365468382816047630355461653340632553426278198853241,
      25919299780512511508061958642305261009583198324725036212440752482930702519878,
      11791415309425995046749154607832041856871129882141188736462372751874115368248
    ],
    [
      22114525030336665972036957912787127870644756898138077124815002206627656645846,
      74561778027252859083209130121920474961655350982938755244738788717578708084930
    ],
    [
      5587813875922595628752214729735723034111050560116231646359963981668986135460
    ]
  ],
  5587813875922595628752214729735723034111050560116231646359963981668986135460
]

[
  [
    [sort([ merkle_hash("1"),merkle_hash("2"),merkle_hash("3"),...,merkle_hash("7"),merkle_hash("merkle-tree-extra") ])] // H11,H12,H13,...,H18
    [merkle_hash(H11,H12),..,merkle_hash(H17,H18)] // H21, H22, H23,H24
    [merkle_hash(H21,H22),merkle_hash(H23,H24)] // H31, H32
    merkle_hash(H31,H32) // H41
  ],
  merkle_hash(H31,H32) // ROOT
]



merkle.merkle_path("7",tree[0])
[
  37711660782102817547094073135578998531779790412684035506279823231061364818016,
  43042351581350983610621529617640359779365126521871794350496949428256481263225,
  103509800336581907939101876374092451924972847149348896254603184719556990494914
]

 */
