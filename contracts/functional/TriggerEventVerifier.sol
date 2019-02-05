// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.23;

import "../libraries/RLP.sol";
import "../libraries/SolidityUtils.sol";
import "../EventVerifier.sol";

/*
    TriggerEventVerifier

    Inherits from `EventVerifier` and verifies `Triggered` events.

    From the provided logs, we separate the data and define checks to assert certain information in the event and
    returns `true` if successful.

    Contracts similar to this that verify specific events should be designed only to verify the data inside the
    supplied events with similarly supplied expected outcomes. It is only meant to serve as a utility to perform defined
    checks against specific events.
*/
contract TriggerEventVerifier is EventVerifier {
    bytes32 eventSignature = keccak256("Triggered(address)");

    function verify(bytes20 _contractEmittedAddress, bytes _rlpReceipt, bytes20 _expectedAddress) public view returns (bool) {
        // Retrieve specific log for given event signature
        RLP.RLPItem[] memory log = retrieveLog(eventSignature, _contractEmittedAddress, _rlpReceipt);

        // Split logs into constituents. Not all constituents are used here
        bytes memory contractEmittedEvent = RLP.toData(log[0]);
        RLP.RLPItem[] memory topics = RLP.toList(log[1]);
        bytes memory data = RLP.toData(log[2]);

        /*
        This section below is specific to this event verifier and checks the relevant data.
        In this event we only expect a single un-indexed address parameter which will be present in the data field.
        The data field pads it's contents if they are less than 32 bytes. Therefore we know that our address parameter
        exists in the 20 least significant bytes of the data field.

        We copy the last 20 bytes of our data field to a bytes20 variable to compare against the supplied expected
        parameter in the event from our function call. This acts as our conditional check that the event called is what
        the user expects.
        */
        bytes20 b20_address = SolUtils.BytesToBytes20(data, data.length - 20);
        return b20_address == _expectedAddress;
    }

}
