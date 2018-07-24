pragma solidity ^0.4.23;

import "./RLP.sol";

contract EventVerifier {
    function verify(bytes20 _contractEmittedAddress, bytes _rlpReceipt, bytes20 _expectedAddress) public returns (bool) {

        /* Here we decode the receipt into it's consituents and grab the logs as we know it's position in the receipt
        object and proceed to decode the logs also.
        */
        RLP.RLPItem[] memory receipt = RLP.toList(RLP.toRLPItem(_rlpReceipt));
        RLP.RLPItem[] memory logs = RLP.toList(receipt[3]);


        /*
        In this example there should be only one set of logs but in other cases there may be multiple sets of logs,
        one set per event emitted per transaction. It will be the writer of the event verifier to locate the correct
        log in the receipts through iteration.
        */
        assert( logs.length == 1 );


        // The log decoded into a list: log = [ contractAddress (bytes32), topics (array), data (bytes) ]
        RLP.RLPItem[] memory log = RLP.toList(logs[0]);


        /*
        In the case where there are multiple logs, the below code should iterate through them to find the correct one:

        bytes32 expectedEventSignature = 0xsomespecifichash
        RLP.RLPItem[] memory log;
        for (uint i = 0; i < logs.length; i++) {
            log = RLP.toList(logs[i]);
            RLP.RLPItem[] memory topics = RLP.toList(log[1]);
            bytes32 eventSignature = RLP.toBytes32(topics[0]);
            if (eventSignature == expectedEventSignature) {
                break;
            }
        }

        And then separated into it's constituents for processing below:
        */
        bytes memory contractEmittedEvent = RLP.toData(log[0]);
        RLP.RLPItem[] memory topics = RLP.toList(log[1]);
        bytes memory data = RLP.toData(log[2]);

        bytes20 emissionSource = bytesToBytes20(contractEmittedEvent, 0);
        assert( emissionSource == _contractEmittedAddress);

        /*
        This section below is specific to this event and checks the relevant data.
        In this event we only expect a single un-indexed address parameter which will be present in the data field.
        The data field pads it's contents if they are less than 32 bytes. Therefore we know that our address parameter
        exists in the 20 least significant bytes of the data field.

        We copy the last 20 bytes of our data field to a bytes20 variable to compare against the supplied expected
        parameter in the event from our function call. This acts as our conditional check that the event called is what
        the user expects.
        */
        bytes20 b20_address = bytesToBytes20(data, data.length - 20);
        assert( b20_address == _expectedAddress );

        /*
        Once verified, the logs of this specific event are proven as true and returns as such. Else, the execution
        reverts.
        */
        return true;
    }

    function bytesToBytes20(bytes _b, uint _offset) private pure returns (bytes20) {
        bytes20 out;

        for (uint i = 0; i < 20; i++) {
            out |= bytes20(_b[_offset + i] & 0xFF) >> (i * 8);
        }
        return out;
    }
}
