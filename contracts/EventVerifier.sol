pragma solidity ^0.4.23;

import "./RLP.sol";

contract EventVerifier {

    function verify(bytes _rlpReceipt, bytes20 _expectedAddress) public returns (bool) {
        RLP.RLPItem[] memory receipt = RLP.toList(RLP.toRLPItem(_rlpReceipt));

        RLP.RLPItem[] memory logs = RLP.toList(receipt[3]);

        require( logs.length == 1 );

        bytes memory b_address = RLP.toBytes(logs[0]);
        bytes20 b20_address;
        assembly {
            b20_address := mload(add(b_address, 0))
        }

        require( b20_address == _expectedAddress );

        return true;
    }
}
