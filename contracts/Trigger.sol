pragma solidity ^0.4.23;

contract Trigger {
    event Triggered(address caller);

    function fire() public {
        emit Triggered(msg.sender);
    }
}
