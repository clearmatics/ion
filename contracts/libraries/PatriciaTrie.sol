pragma solidity ^0.4.23;

import "./RLP.sol";

library PatriciaTrie {

    function verifyProof(bytes _value, bytes _parentNodes, bytes _path, bytes32 _root) returns (bool) {
        RLP.RLPItem memory nodes = RLP.toRLPItem(_parentNodes);
        RLP.RLPItem[] memory parentNodes = RLP.toList(nodes);

        bytes32 currentNodeKey = _root;
        bytes memory b_currentNode;
        RLP.RLPItem[] memory currentNode;

        uint traversedNibbles = 0;
        bytes memory path = toNibbleArray(_path, false);

        for (uint i = 0; i < parentNodes.length; i++) {
            b_currentNode = RLP.toBytes(parentNodes[i]);

            if (currentNodeKey != keccak256(b_currentNode)) {
                return false;
            }

            currentNode = RLP.toList(parentNodes[i]);

            if (currentNode.length == 17) {
                // Branch Node
                if (traversedNibbles == path.length) {

                    if(keccak256(RLP.toBytes(currentNode[16])) == keccak256(_value)) {
                        return true;
                    } else {
                        return false;
                    }
                }

                uint16 nextPathNibble = uint16(path[traversedNibbles]);
                currentNodeKey = RLP.toBytes32(currentNode[nextPathNibble]);
                traversedNibbles += 1;

            } else if (currentNode.length == 2) {
                // Extension/Leaf Node
                bytes memory nextPathNibbles = RLP.toData(currentNode[0]);
                traversedNibbles += toNibbleArray(nextPathNibbles, true).length;

                if (traversedNibbles == path.length) {
                    if(keccak256(RLP.toData(currentNode[1])) == keccak256(_value)) {
                        return true;
                    } else {
                        return false;
                    }
                }

                // Reached a leaf before end of the path. Proof false.
                if (toNibbleArray(nextPathNibbles, true).length == 0) {
                    return false;
                }

                bytes memory nextNodeKey = RLP.toData(currentNode[1]);
                currentNodeKey = bytesToBytes32(nextNodeKey, 0);
            } else {
                return false;
            }
        }
    }

    function toNibbleArray(bytes b, bool hexPrefixed) private returns (bytes) {
        bytes memory nibbleArray = new bytes(255);

        uint8 nibblesFound = 0;
        for (uint i = 0; i < b.length; i++) {
            byte[2] memory nibbles = byteToNibbles(b[i]);

            if (hexPrefixed && i == 0) {
                if (nibbles[0] == 1 || nibbles[0] == 3) {
                    nibbleArray[nibblesFound] = nibbles[1];
                    nibblesFound += 1;
                }
            } else {
                nibbleArray[nibblesFound] = nibbles[0];
                nibbleArray[nibblesFound + 1] = nibbles[1];
                nibblesFound += 2;
            }
        }

        bytes memory finiteNibbleArray = new bytes(nibblesFound);
        for (uint j = 0; j < nibblesFound; j++) {
            finiteNibbleArray[j] = nibbleArray[j];
        }
        return finiteNibbleArray;
    }

    function byteToNibbles(byte b) private pure returns (byte[2]) {
        byte firstNibble = rightShift(b, 4);
        byte secondNibble = b & 0xf;

        return [firstNibble, secondNibble];
    }

    function leftShift(byte i, uint8 bits) private pure returns (byte) {
        return byte(uint8(i) * 2 ** bits);
    }

    function rightShift(byte i, uint8 bits) private pure returns (byte) {
        return byte(uint8(i) / 2 ** bits);
    }

    function bytesToBytes32(bytes b, uint offset) private pure returns (bytes32) {
        bytes32 out;

        for (uint i = 0; i < 32; i++) {
            out |= bytes32(b[offset + i] & 0xFF) >> (i * 8);
        }
        return out;
    }
}
