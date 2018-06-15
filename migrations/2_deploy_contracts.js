const Ion = artifacts.require("Ion");

module.exports = async (deployer) => {
  try {
    deployer.deploy(Ion, "0xab830ae0774cb20180c8b463202659184033a9f30a21550b89a2b406c3ac8075")

  } catch(err) {
    console.log('ERROR on deploy:',err);
  }

};