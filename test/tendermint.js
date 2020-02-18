// Copyright (c) 2020 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+

/*
    Tendermint Validation contract test

    Tests here are standalone unit tests for tendermint module functionality.
    Other contracts have been mocked to simulate basic behaviour.

    Tests the tenderming scheme for block submission, validator signature verification and more.
*/


const MockIon = artifacts.require("MockIon");
const MockStorage = artifacts.require("MockStorage");
const Tendermint = artifacts.require("TendermintAutonity");

contract("Tendermint", (accounts) => {
    beforeEach("Setup contract for each test", async () => {
        ion = await MockIon.new(DEPLOYEDCHAINID);
        tendermint = await Tendermint.new(ion.address);
        storage = await MockStorage.new(ion.address);
    })

    describe("Register Chain", () => {
        it("Succesful Register Chain", async() => {
            
        })
    })
})