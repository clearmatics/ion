pragma solidity ^0.4.23;

contract EventConsuming {
    uint256[] consumedNonces;
    address ionCompatible;

    constructor(address _ionCompAddress) public {
        ionCompatible = _ionCompAddress;
    }

    modifier IonCompatibleOnly(){
        require( msg.sender == ionCompatible,
        "Caller is not IonCompatible contract."
        );
        _;
    }

    function consumeEvent(uint256 nonce) internal {
        consumedNonces.push(nonce);
    }
}