const Ion = artifacts.require("Ion");
const Recover = artifacts.require("Recover");
const Validation = artifacts.require("Validation");
const PatriciaTrie = artifacts.require("PatriciaTrie");

module.exports = async (deployer) => {
  try {
    deployer.deploy(Recover)
      .then(() => Recover.deployed)
      .then(() => deployer.deploy(Validation, ["0x2be5ab0e43b6dc2908d5321cf318f35b80d0c10d", "0x8671e5e08d74f338ee1c462340842346d797afd3"], "0xc3bac257bbd04893316a76d41b6ff70de5f65c9f24db128864a6322d8e0e2f28"))
      .then(() => Validation.deployed)
      .then(() => deployer.deploy(PatriciaTrie))
      .then(() => PatriciaTrie.deployed)
      .then(() => deployer.link(PatriciaTrie, Ion))
      .then(() => deployer.deploy(Ion, "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075"))
      .then(() => Ion.deployed)
  } catch(err) {
    console.log('ERROR on deploy:',err);
  }

};