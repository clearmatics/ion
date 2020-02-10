// Copyright (c) 2016-2020 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.5.12;

import "./RLP.sol";

library MerkleTree {

    // elements are the hashed pre-images
    function generateRoot(bytes32[] memory elements) internal pure returns (bytes32) {
        require(elements.length > 0, "You passed an empty array");

        uint padding = calculatePadding(elements);
        uint i;
        uint finalLength = elements.length + padding

        // add padding
        for (i = elements.length - 1; i < finalLength; i++) {
            elements[i] = bytes32(0);
        }
        
        // hash pairs of consecutive values and append the resulting hash
        for (i = 0; i < finalLength - 1; i++) {
            elements[elements.length + i + 1] = hashPair(elements[2*i], elements[2*i+1]);
        }

        // last element should be root 
        return elements[elements.length - 1];
    }

    function hashPair(bytes32 elementA, bytes32 elementB) internal pure returns (bytes32) {

        // return hash of pair or one of the two if the other is 0x0..
        if (elementA == bytes32(0)) {
            return elementB;
        } else if (elementB == bytes32(0)) {
            return elementA;
        }

        // sort the two (for verification purpose), rlp encode and hash
        if (elementA > elementB) {
            return keccak256(abi.encodePacked(elementA, elementB));
        } else {
            return keccak256(abi.encodePacked(elementB, elementA));
        }
    }

    // find the number of zero elements to add to create a complete binary tree
    function calculatePadding(bytes32[] memory elements) internal pure returns (uint) {

        for (uint8 i = 1; i <= elements.length; i++) {

            uint values;

            // TODO hardcode values instead of do exp
            assembly { values := exp(2,i) }

            if (values == elements.length) {
                // already complete
                return 0;
            } else if (values > elements.length) {
                return values - elements.length;
            }
        }
    }
}