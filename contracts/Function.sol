pragma solidity ^0.4.23;

import "./EventConsuming.sol";

contract Function is EventConsuming {
    constructor(address _ionCompAddress) EventConsuming(_ionCompAddress) {}

    event Executed();

    function execute() IonCompatibleOnly {
        emit Executed();
    }

    function CheckReceipt() IonCompatibleOnly {

    }
}
