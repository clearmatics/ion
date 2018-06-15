pragma solidity ^0.4.0;

contract EventConsumable {
    event IonCompatibleEvent(
        bytes32 originChain,
        bytes32[] destinationChain,
        address emittedFromContract,
        uint256 nonce);
}
