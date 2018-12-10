// Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.4.0;
import "./Ion.sol";

contract IonCompatible {
    /*  The Ion contract that proofs would be made to. Ensure that prior to verification attempts that the relevant
        blocks have been submitted to the Ion contract. */
    Ion internal ion;

    constructor(address _ionAddr) public {
        ion = Ion(_ionAddr);
    }
}
