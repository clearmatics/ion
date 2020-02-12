// Copyright (c) 2016-2020 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.5.12;

library MerkleTree {

    // elements are the hashed pre-images
    function generateRoot(bytes32[] memory elements) internal returns (bytes32) {
        require(elements.length > 0, "You passed an empty array");

        uint elementsToHash = elements.length;

        uint hashedElements = 0;
        while (elementsToHash > 1) {
            uint largest2Power = calculateLargest2Power(elementsToHash);

            bytes32 hash = hashArrayElements(elements, hashedElements, largest2Power);

            hashedElements += largest2Power - 1;
            elements[hashedElements] = hash;
            elementsToHash -= hashedElements;
        }

        // last element should be root
        return elements[elements.length - 1];
    }

    // Hashes power-of-2 number of array elements down to 1
    // Takes an array of elements and hashes a subarray from index + number of elements
    function hashArrayElements(bytes32[] memory elements, uint fromIndex, uint numOfElements) internal pure returns (bytes32) {
        require(numOfElements % 2 == 0, "Number of elements to hash must be a power of 2");
        require(elements.length >= fromIndex + numOfElements, "Cannot hash so many elements from array. Too short.");

        // Find the tree depth for the number of elements needed to be hashed in the subarray
        int depth = calculateLargest2PowerIndex(numOfElements);
        for (int i = depth; i > 0; i--) {
            // Calculate how many elements to hash for each depth of the tree
            uint depthElements;
            assembly { depthElements := exp(2, i)}

            // Hash elements in the current tree depth and assign the result to the beginning of the subarray
            for (uint j = 0; j < depthElements; j += 2) {
                bytes32 hash = hashPair(elements[fromIndex + j], elements[fromIndex + j + 1]);
                elements[fromIndex + j/2] = hash;
            }
        }

        // Root should be assigned to the beginning of the subarray
        return elements[fromIndex];
    }

    function hashPair(bytes32 elementA, bytes32 elementB) internal pure returns (bytes32) {
        // sort the two (for verification purpose), rlp encode and hash
        if (elementA < elementB) {
            return keccak256(abi.encodePacked(elementA, elementB));
        } else {
            return keccak256(abi.encodePacked(elementB, elementA));
        }
    }

    // Returns the largest power of 2 contained in a number
    function calculateLargest2Power(uint number) internal pure returns (uint) {
        if (number >= 1024) {
            //Unsupported
            return 0;
        } else if (number >= 512) {
            return 512;
        } else if (number >= 256) {
            return 256;
        } else if (number >= 128) {
            return 128;
        } else if (number >= 64) {
            return 64;
        } else if (number >= 32) {
            return 32;
        } else if (number >= 16) {
            return 16;
        } else if (number >= 8) {
            return 8;
        } else if (number >= 4) {
            return 4;
        } else if (number >= 2){
            return 2;
        } else if (number >= 1){
            return 1;
        } else {
            return 0;
        }
    }

    // Returns the largest power of 2 contained in a number
    function calculateLargest2PowerIndex(uint number) internal pure returns (int) {
        if (number < 1) {
            return -1;
        }

        if (number > 1024) {
            //Unsupported
            return -1;
        } else if (number >= 512) {
            return 9;
        } else if (number >= 256) {
            return 8;
        } else if (number >= 128) {
            return 7;
        } else if (number >= 64) {
            return 6;
        } else if (number >= 32) {
            return 5;
        } else if (number >= 16) {
            return 4;
        } else if (number >= 8) {
            return 3;
        } else if (number >= 4) {
            return 2;
        } else if (number >= 2) {
            return 1;
        } else if (number >= 1) {
            return 0;
        } else {
            return 0;
        }
    }
}