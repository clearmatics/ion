// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.23;

import "./RLP.sol";

library PatriciaTrie {

    function verifyProof(bytes _value, bytes _parentNodes, bytes _path, bytes32 _root) internal returns (bool) {
        RLP.RLPItem memory nodes = RLP.toRLPItem(_parentNodes);
        RLP.RLPItem[] memory parentNodes = RLP.toList(nodes);

        bytes32 currentNodeKey = _root;

        uint traversedNibbles = 0;
        bytes memory path = toNibbleArray(_path, false);

        for (uint i = 0; i < parentNodes.length; i++) {
            if (currentNodeKey != keccak256(RLP.toBytes(parentNodes[i]))) {
                return false;
            }

            RLP.RLPItem[] memory currentNode = RLP.toList(parentNodes[i]);

            if (currentNode.length == 17) {
                // Branch Node
                (currentNodeKey, traversedNibbles) = processBranchNode(currentNode, traversedNibbles, path, _value);
            } else if (currentNode.length == 2) {
                // Extension/Leaf Node
                (currentNodeKey, traversedNibbles) = processExtensionLeafNode(currentNode, traversedNibbles, path, _value);
            } else {
                return false;
            }

            // Read comment block below for explanation of this
            if (currentNodeKey == 0x0) {
                return traversedNibbles == 1;
            }
        }
    }

    /**
    Node Processing

    processBranchNodes returns (bytes32 currentNodeKey, uint traversedNibbles)
    processExtensionLeafNode returns (bytes32 currentNodeKey, uint traversedNibbles)

    Due to the dual nature of how a branch node may be processed where the next node in the path could be either
    referenced by hash or nested in the branch node is the total RLP-encoded node is less than 32 bytes (nested node),
    we required separation of logic due to "stack-too-deep" issues and opted for a messy returning of reused variables.
    These returned variables now hold two purposes:

    * currentNodeKey (bytes32): Holds value of the hash of the next node to be processed. If processing is finished this
                                value is 0x0.
    * traversedNibbles (uint):  Tracks how many nibbles have been traversed. If processing is finished this value will
                                be 0 if verification failed, and 1 if verification succeeded.

    The dual-functionality of these variables is the crux of how I avoided stack issues which makes the code somewhat
    unreadable. If there is an improvement to this algorithm that can make it more readable please share.

    */

    function processBranchNode(RLP.RLPItem[] memory _currentNode, uint _traversedNibbles, bytes memory _path, bytes _value) private returns (bytes32, uint) {
        if (_traversedNibbles == _path.length) {
            return (0x0, checkNodeValue(_value, RLP.toBytes(_currentNode[16])) ? 1 : 0);
        }

        uint16 nextPathNibble = uint16(_path[_traversedNibbles]);
        RLP.RLPItem memory nextNode = _currentNode[nextPathNibble];
        _traversedNibbles += 1;

        bytes32 currentNodeKey;
        if (RLP.toBytes(nextNode).length < 32) {
            //Nested 'Node'
            (currentNodeKey, _traversedNibbles) = processNestedNode(nextNode, _traversedNibbles, _path, _value);
        } else {
            currentNodeKey = RLP.toBytes32(_currentNode[nextPathNibble]);
        }
        return (currentNodeKey, _traversedNibbles);
    }

    function processExtensionLeafNode(RLP.RLPItem[] memory _currentNode, uint _traversedNibbles, bytes memory _path, bytes _value) private returns (bytes32, uint) {
        bytes memory nextPathNibbles = RLP.toData(_currentNode[0]);
        _traversedNibbles += toNibbleArray(nextPathNibbles, true).length;

        if (_traversedNibbles == _path.length) {
            return (0x0, checkNodeValue(_value, RLP.toData(_currentNode[1])) ? 1 : 0);
        }

        // Reached a leaf before end of the path. Proof false.
        if (toNibbleArray(nextPathNibbles, true).length == 0) {
            return (0x0, 0);
        }

        bytes memory nextNodeKey = RLP.toData(_currentNode[1]);
        bytes32 currentNodeKey = bytesToBytes32(nextNodeKey, 0);

        return (currentNodeKey, _traversedNibbles);
    }

    function processNestedNode(RLP.RLPItem memory _nextNode, uint _traversedNibbles, bytes memory _path, bytes _value) private returns (bytes32, uint) {
        RLP.RLPItem[] memory currentNode = RLP.toList(_nextNode);
        if (currentNode.length == 17) {
            // Branch Node
            return processBranchNode(currentNode, _traversedNibbles, _path, _value);
        } else if (currentNode.length == 2) {
            // Leaf Node
            return processExtensionLeafNode(currentNode, _traversedNibbles, _path, _value);
        } else {
            return (0x0, 0);
        }
    }

    function checkNodeValue(bytes _expected, bytes _nodeValue) private returns (bool) {
        return keccak256(_expected) == keccak256(_nodeValue);
    }

    function toNibbleArray(bytes b, bool hexPrefixed) private pure returns (bytes) {
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
        return byte(uint8(i) * uint8(2) ** uint8(bits));
    }

    function rightShift(byte i, uint8 bits) private pure returns (byte) {
        return byte(uint8(i) / uint8(2) ** uint8(bits));
    }

    function bytesToBytes32(bytes b, uint offset) private pure returns (bytes32) {
        bytes32 out;

        for (uint i = 0; i < 32; i++) {
            out |= bytes32(b[offset + i] & 0xFF) >> (i * 8);
        }
        return out;
    }
}
