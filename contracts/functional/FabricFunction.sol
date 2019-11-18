pragma solidity ^0.5.12;

import "../storage/FabricStore.sol";

contract FabricFunction {
    FabricStore blockStore;

    constructor(address _storeAddr) public {
        blockStore = FabricStore(_storeAddr);
    }

    event State(uint blockNo, uint txNo, string mvalue);

    function execute(uint _blockNo, uint _txNo, string memory _value) internal {
        emit State(_blockNo, _txNo, _value);
    }

    function retrieveAndExecute(bytes32 _chainId, string memory _channelId, string memory _key) public {
        uint blockVersion;
        uint txVersion;
        string memory value;

        (blockVersion, txVersion, value) = blockStore.getState(_chainId, _channelId, _key);

        execute(blockVersion, txVersion, value);
    }
}
