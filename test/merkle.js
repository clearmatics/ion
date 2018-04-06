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

// =======================================================================================================
// =======================================================================================================
// ==================================   TEST   ===========================================================
// =======================================================================================================
// =======================================================================================================

contract('Merkle.js', () => {
  it('test JS Merkle', () => {

    const testData = ["1","2","3","4","5","6","7"]

    const tree = merkle.createMerkle(testData)

    const expectedTree = [
      [
        [
          '8568612641526826488487436752726739043287191320122540356069953783894380777505',
          '8763638472773768691201326883407021568462294246273894496415427229083082408032',
          '19224855404247632006917173431419498680506051063941070371722880450128577361118',
          '61795459977501490647348212754130855970016313872340374962921336716751708851142',
          '64645341593328157176709656265449880868558868673380425455960412802858937540801',
          '74330811247603495249613868516695563873247293176611122272199330092769797099053',
          '78469846343542442363028680824980501212021332975324075417961003849793346933925',
          '75317570447191171753008806478868650352148013528306361601609880810432714200529'
        ],
        [
          '6560824545851281876686151142367952893930617484325436481370811303698242675212',
          '14094329272021934754728783365468382816047630355461653340632553426278198853241',
          '25919299780512511508061958642305261009583198324725036212440752482930702519878',
          '11791415309425995046749154607832041856871129882141188736462372751874115368248'
        ],
        [
          '22114525030336665972036957912787127870644756898138077124815002206627656645846',
          '74561778027252859083209130121920474961655350982938755244738788717578708084930'
        ],
        [
          '5587813875922595628752214729735723034111050560116231646359963981668986135460'
        ]
      ],
      '5587813875922595628752214729735723034111050560116231646359963981668986135460'
    ]

    const treeStr = [tree[0].map(i => i.map(j => j.toString(10))),tree[1].toString(10)]
    assert.deepEqual(treeStr,expectedTree)

    const expectedPaths = [
      [
        '19224855404247632006917173431419498680506051063941070371722880450128577361118',
        '6560824545851281876686151142367952893930617484325436481370811303698242675212',
        '103509800336581907939101876374092451924972847149348896254603184719556990494914'
      ],
      [
        '104265592756520220608901552731040627315465509694716502611474276812410996610513',
        '25919299780512511508061958642305261009583198324725036212440752482930702519878',
        '22114525030336665972036957912787127870644756898138077124815002206627656645846'
      ],
      [
        '90743482286830539503240959006302832933333810038750515972785732718729991261126',
        '6560824545851281876686151142367952893930617484325436481370811303698242675212',
        '103509800336581907939101876374092451924972847149348896254603184719556990494914'
      ],
      [
        '8568612641526826488487436752726739043287191320122540356069953783894380777505',
        '43042351581350983610621529617640359779365126521871794350496949428256481263225',
        '103509800336581907939101876374092451924972847149348896254603184719556990494914'
      ],
      [
        '103278833556932544105506614768867540836564789343021263282063726094748079509037',
        '40739437618755043902641900860004018820188626048551329746326768753852397778232',
        '22114525030336665972036957912787127870644756898138077124815002206627656645846'
      ],
      [
        '64645341593328157176709656265449880868558868673380425455960412802858937540801',
        '40739437618755043902641900860004018820188626048551329746326768753852397778232',
        '22114525030336665972036957912787127870644756898138077124815002206627656645846'
      ],
      [
        '37711660782102817547094073135578998531779790412684035506279823231061364818016',
        '43042351581350983610621529617640359779365126521871794350496949428256481263225',
        '103509800336581907939101876374092451924972847149348896254603184719556990494914'
      ]
    ]

    const path = testData.map(value => merkle.pathMerkle(value,tree[0]))
    assert.deepEqual(path.map(i => i.map(j => j.toString(10))),expectedPaths, 'paths badly created')

    const proof = testData.reduce((prev,leaf,idx) => (merkle.proofMerkle(leaf,path[idx],tree[1]) && prev), true)
    const negProof = testData.reduce((prev,leaf,idx) => !(merkle.proofMerkle('10',path[idx],tree[1]) && prev),true)
    assert(proof && negProof,'proof failed')
  })
})
