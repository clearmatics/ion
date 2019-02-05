pragma solidity ^0.4.24;

import "../storage/FabricStore.sol";

contract FabricFunction {
    FabricStore blockStore;

    constructor(address _storeAddr) public {
        blockStore = FabricStore(_storeAddr);
    }

    event State(uint blockNo, uint txNo, string value);

    function execute(uint _blockNo, uint _txNo, string _value) internal {
        emit State(_blockNo, _txNo, _value);
    }

    function retrieveAndExecute(bytes32 _chainId, string _channelId, string _key) public {
        uint blockVersion;
        uint txVersion;
        string memory value;

        (blockVersion, txVersion, value) = blockStore.getState(_chainId, _channelId, _key);

        execute(blockVersion, txVersion, value);
    }
}
