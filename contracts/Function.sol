pragma solidity ^0.4.24;

import "./EventConsuming.sol";

contract FunctionConsumesEvent is EventConsuming {
    constructor(address _ionCompAddress) EventConsuming(_ionCompAddress) {}

    event Executed();

    function execute() IonCompatibleOnly {
        emit Executed();
    }

    function CheckReceipt() IonCompatibleOnly {

        ion.CheckReceiptProof();
    }
}
