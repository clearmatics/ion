const Ion = artifacts.require("Ion");
const Recover = artifacts.require("Recover");
const Validation = artifacts.require("Validation");
const PatriciaTrie = artifacts.require("PatriciaTrie");
const EventFunction = artifacts.require("Function");
const EventVerifier = artifacts.require("EventVerifier");

module.exports = async (deployer) => {
  try {
    deployer.deploy(Recover)
      .then(() => Recover.deployed)
      .then(() => deployer.deploy(Validation, "0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177"))
      .then(() => Validation.deployed)
      .then(() => deployer.deploy(PatriciaTrie))
      .then(() => PatriciaTrie.deployed)
      .then(() => deployer.link(PatriciaTrie, Ion))
      .then(() => deployer.deploy(Ion, "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075"))
      .then(() => Ion.deployed)
      .then(() => deployer.deploy(EventVerifier))
      .then(() => EventVerifier.deployed)
      .then(() => deployer.deploy(EventFunction, Ion.address, EventVerifier.address))
      .then(() => EventFunction.deployed)
  } catch(err) {
    console.log('ERROR on deploy:',err);
  }

};