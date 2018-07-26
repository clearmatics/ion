const Ion = artifacts.require("Ion");
const PatriciaTrie = artifacts.require("PatriciaTrie");
const PatriciaTrieTest = artifacts.require("PatriciaTrieTest");
const Function = artifacts.require("Function");
const EventVerifier = artifacts.require("EventVerifier");

module.exports = async (deployer) => {
  try {
    await deployer.deploy(PatriciaTrie);
    await deployer.link(PatriciaTrie, Ion);
    await deployer.deploy(Ion, "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075");
    eventVerifier = await deployer.deploy(EventVerifier);
    ion = await Ion.deployed();
    eventVerifier = await EventVerifier.deployed();
    await deployer.deploy(Function, ion.address, eventVerifier.address);

  } catch(err) {
    console.log('ERROR on deploy:',err);
  }

};