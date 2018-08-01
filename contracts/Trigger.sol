pragma solidity ^0.4.23;

/*
    Trigger

    Example contract that emits an event to be consumed.

    Currently an instance deployed to:
    Rinkeby @: 0x61621bcf02914668f8404c1f860e92fc1893f74c
    Deployment Tx Hash: 0xc9500e84af2394e1d91b43e40c9c89f105636748f95ae05c11c73f2fd755795e
    Deployed Block Number: 2657325
    `fire()` call Tx Hash 0xafc3ab60059ed38e71c7f6bea036822abe16b2c02fcf770a4f4b5fffcbfe6e7e

    The current tests are running against generated proofs from Rinkeby for the above data and consumes the event
    emitted in the transaction executed.
*/

contract Trigger {
    event Triggered(address caller);

    function fire() public {
        emit Triggered(msg.sender);
    }
}