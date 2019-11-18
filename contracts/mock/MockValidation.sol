pragma solidity ^0.5.12;

import "../IonCompatible.sol";

contract MockValidation is IonCompatible {
    constructor (address _ionAddr) IonCompatible(_ionAddr) public {}

    function register() public returns (bool) {
        ion.registerValidationModule();
        return true;
    }

    function SubmitBlock(address _storageAddress, bytes32 _chainId, bytes memory _blockBlob) public {
        ion.storeBlock(_storageAddress, _chainId, _blockBlob);
    }
}
