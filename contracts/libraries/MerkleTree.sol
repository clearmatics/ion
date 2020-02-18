// Copyright (c) 2016-2020 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.5.12;

library MerkleTree {

    // elements are the hashed pre-images
    function generateRoot(bytes32[] memory elements) internal pure returns (bytes32) {
        
        // corner cases
        if (elements.length == 0) {
            return bytes32(0);
        } else if (elements.length == 1) {
            return elements[0];
        }

        uint pos = 0;
        uint i = 0;
        uint upperBoundIndex = elements.length - 1; // index at which stop hashing

        while(upperBoundIndex > 1) {

            if (i == upperBoundIndex){
                // last element of this layer - carry it over
                elements[pos] = elements[i];

                // reset index and pos 
                i = 0;
                pos = 0;

                // calculate new upper bound - solidity already rounds toward zero
                upperBoundIndex = upperBoundIndex / 2;

            } else {
                // i have two values to hash
                elements[pos] = hashPair(elements[i], elements[i + 1]);

                if (i == upperBoundIndex - 1) {
                    // those were last two elements of this layer
                
                    // reset index and pos 
                    i = 0;
                    pos = 0;

                    // calculate new upper bound - solidity already rounds toward zero
                    upperBoundIndex = upperBoundIndex / 2;
                } else {

                    // continue with this layer
                    i += 2;   
                    pos ++;
                }
            }
        }

        // i have last two elements to form root
        return hashPair(elements[0], elements[1]);
    }

    function verify(bytes32[] memory proof, bytes32 root, bytes32 leaf) internal pure returns (bool) {
        bytes32 computedHash = leaf;

        for (uint256 i = 0; i < proof.length; i++) {
            computedHash = hashPair(computedHash, proof[i]);
        }

        return computedHash == root;
    }

    function hashPair(bytes32 elementA, bytes32 elementB) internal pure returns (bytes32) {

        // return hash of pair or one of the two if the other is 0x0..
        if (elementA == bytes32(0)) {
            return elementB;
        } else if (elementB == bytes32(0)) {
            return elementA;
        }

        // sort the two (for verification purpose), rlp encode and hash
        if (elementA >= elementB) {
            return keccak256(abi.encodePacked(elementB, elementA));
        } else {
            return keccak256(abi.encodePacked(elementA, elementB));
        }
    }
}