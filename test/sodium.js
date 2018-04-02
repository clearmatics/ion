const BN = require('bignumber.js')

const merkle = require('./merkle')
const Sodium = artifacts.require("Sodium");

contract.only('Sodium', (accounts) => {

  it.only('test JS Merkle', () => {

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
    //testData.forEach((leaf,idx) => assert(merkle.proofMerkle(leaf,path[idx],tree[1]), 'proof failed'))
    //testData.forEach((leaf,idx) => assert(!merkle.proofMerkle('10',path[idx],tree[1]), 'proof failed'))

    // console.log(JSON.stringify(tree,2,2))
  })

  it.only('data test', async () => {
    const sodium = await Sodium.deployed();

    const testData = ["1","2","3","4","5","6","7"]
    const tree = merkle.createMerkle(testData)
    const path = testData.map(value => merkle.pathMerkle(value,tree[0]))

    console.log(path.length)
    console.log(merkle.proofMerkle(testData[0],path[0],tree[1],false,true))

    testData.forEach((leaf,idx) => assert(merkle.proofMerkle(leaf,path[idx],tree[1])))

    const leafHash = merkle.rawDataMerkleHash(testData[0])
    const rootArg = tree[1]

    const nextBlock = await sodium.NextBlock()
    //console.log('0x'+nextBlock.toString(16))
    const receiptUpdate = await sodium.Update(nextBlock,[rootArg])
    //console.log(receiptUpdate)
    //const valid = await sodium.Verify(nextBlock.toString(),leafArg,pathArg)
    const valid = await sodium.Verify(nextBlock,leafHash,path[0])

    /*
    const valid = await sodium.Verify2(nextBlock,leafHash,path[0])
    valid.logs.map(l=>{
      //console.log(l.event,Object.keys(l.args).map(k=>[k,l.args[k].toString(16)]))
      if(l.event === 'HashEvent') {
        const aHex = l.args.a.toString(16).padStart(64,'0')
        const bHex = l.args.b.toString(16).padStart(64,'0')
        const abHex = new BN('0x' + aHex + bHex)
        const expectedResult = l.args.result.toString(16)
        const bitSet = l.args.bitSet
        console.log('Contract: a:',aHex,'b:',bHex,'result:',expectedResult)
        console.log('Truffle:',merkle.merkleHash(abHex).toString(16))
      }
    });
    path[0].map((p,idx)=> console.log('p('+idx+'):',p.toString(16)))
    console.log('leaf:',leafHash.toString(16))
    console.log('root:',rootArg.toString(16))
    //console.log(valid.map(v => '>> 0x'+v.toString(16)))
    */
     console.log(valid)
    assert(valid,'Sodium.verify() failed!')
  })

  /*
   *
   * expected = [ '0x673620737675e2755ce8269a99904022d15da8d5843f5aec205cd243ff80240a',
  '0x228f6de7a6bb38a9720976366c195e4fe678de171770eaed2adbc6e1f45fecbd',
  '0x60f539fe715f17c20861225c41494f3196f57de3d586dc95865a4a5fa9bb5722',
  '0x1a792cf089bfa56eae57ffe87e9b22f9c9bfe52c1ac300ea1f43f4ab53b4b794' ]
   *
   *> web3Utils.soliditySha3(path[0],leaf)
'0x673620737675e2755ce8269a99904022d15da8d5843f5aec205cd243ff80240a'
>  web3Utils.soliditySha3(path[1],expected[0])
'0xa28f6de7a6bb38a9720976366c195e4fe678de171770eaed2adbc6e1f45fecbd'
> cutBit(new BN(web3Utils.soliditySha3(path[1],expected[0])),0xFF).toString(16)
'228f6de7a6bb38a9720976366c195e4fe678de171770eaed2adbc6e1f45fecbd'
> cutBit(new BN(web3Utils.soliditySha3(expected[1],cutBit(path[2],0xFF))),0xFF).toString(16)
'60f539fe715f17c20861225c41494f3196f57de3d586dc95865a4a5fa9bb5722'
> cutBit(new BN(web3Utils.soliditySha3(expected[2],cutBit(path[3],0xFF))),0xFF).toString(16)
'1a792cf089bfa56eae57ffe87e9b22f9c9bfe52c1ac300ea1f43f4ab53b4b794'
   *
   */

  it('check harry\'s data with js merkle implementation',async () => {

    const expected = {
      root: "0x1a792cf089bfa56eae57ffe87e9b22f9c9bfe52c1ac300ea1f43f4ab53b4b794",
      leaf: "0x2584db4a68aa8b172f70bc04e2e74541617c003374de6eb4b295e823e5beab01",
      path: [
        "0x1ab0c6948a275349ae45a06aad66a8bd65ac18074615d53676c09b67809099e0"
        ,"0x093fd25755220b8f497d65d2538c01ed279c131f63e42b2942867f2bd6622486"
        ,"0xb1d101d9a9d27c3a8ed9d1b6548626eacf3d19546306117eb8af547d1e97189e"
        ,"0xcb431dd627bc8dcfd858eae9304dc71a8d3f34a8de783c093188bb598eeafd04"
      ]
    }

    const obj = await Sodium.deployed();
    const nextBlock = await obj.NextBlock()
    const receiptUpdate = await obj.Update(nextBlock.toString(),[expected.root])
    const expectedHashResult = await obj.Verify2(nextBlock.toString(),expected.leaf,expected.path)

    //console.log('Event:',expectedHashResult.logs.map(l=>l.args.bitTest+', '+l.args.node.toString(16)))
    console.log('Event:',expectedHashResult.logs.map(l => l.event + ': ' + Object.keys(l.args).map(k => k+': '+l.args[k].toString(16))))

    /*
    expectedHashResult
      .map(v => new BN(v))
      .filter(v => !v.equals(0))
      .map((v,idx) => console.log(idx+':','0x'+v.toString(16)))
      */

    const valid = merkle.proofMerkle(new BN(expected.leaf),expected.path.map(p => new BN(p)),new BN(expected.root),true)
    assert(valid,'data is not valid in JS merkle implementation')

  })

  it('harry data test', async () => {
    const obj = await Sodium.deployed();

    const root = "0x1a792cf089bfa56eae57ffe87e9b22f9c9bfe52c1ac300ea1f43f4ab53b4b794"
    const leafHash = "0x2584db4a68aa8b172f70bc04e2e74541617c003374de6eb4b295e823e5beab01"
    const path = [
      "0x1ab0c6948a275349ae45a06aad66a8bd65ac18074615d53676c09b67809099e0"
      ,"0x093fd25755220b8f497d65d2538c01ed279c131f63e42b2942867f2bd6622486"
      ,"0xb1d101d9a9d27c3a8ed9d1b6548626eacf3d19546306117eb8af547d1e97189e"
      ,"0xcb431dd627bc8dcfd858eae9304dc71a8d3f34a8de783c093188bb598eeafd04"
    ]
    const nextBlock = await obj.NextBlock()
    //console.log('0x'+nextBlock.toString(16))

    const receiptUpdate = await obj.Update(nextBlock.toString(),[root])
    //console.log(receiptUpdate)

    const valid = await obj.Verify(nextBlock.toString(),leafHash,path)
    //console.log(valid)
    assert(valid)
  });
});
