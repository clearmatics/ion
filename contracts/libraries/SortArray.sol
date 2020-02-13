// Copyright (c) 2016-2020 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.5.12;

library SortArray {

    function sortAddresses(address[] memory addresses) internal returns (address[] memory) {

        quickSort(addresses, 0, addresses.length - 1);

        return addresses;
    }

    function quickSort(address[] memory elements, uint left, uint right) internal {
        uint i = left;
        uint j = right;

        if (i == j) {
            return;
        }

        address pivot = elements[right];

        while (i <= j) {

            while (elements[i] < pivot) {
                i++;
            }

            while (pivot < elements[j]) {
                j--;
            }

            if (i <= j) {
                (elements[i], elements[j]) = (elements[j], elements[i]);
                i++;
                j--;
            }
        }

        if (left < j)
            quickSort(elements, left, j);
        if (i < right)
            quickSort(elements, i, right);
    }   
}