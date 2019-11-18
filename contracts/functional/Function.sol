    // Copyright (c) 2016-2018 Clearmatics Technologies Ltd
// SPDX-License-Identifier: LGPL-3.0+
pragma solidity ^0.5.12;

import "../storage/EthereumStore.sol";

contract TriggerEventVerifier {
    function verify(bytes20 _contractEmittedAddress, bytes memory _rlpReceipt, bytes20 _expectedAddress) public returns (bool);
}

/*
    This function contract is the consumer of an event and performs some execution thereafter. In practice, this would
    be written by a contract designer that intends to consume specific events from another chain. As such all behaviour
    and dependence on such event must be defined here.

    Common custom behaviour may include:
    * Keeping an array of transaction hashes denoting the specific events from that transaction that
    have already been consumed to restrict multiple consumption or 'double spend' of events.
    * Extending the amount of expected event parameters above the stack limit. This might then require some other method
    of passing expected parameters to the contract possibly via RLP-encoding to compress all data to a single argument
    and decoding them within the `verifyAndExecute` function.
    * Including multiple event verifiers if a function requires proof of multiple state transitions from other chains.
    This would also bloat the local scope which is prone to 'stack too deep' issues which would require custom
    workarounds.
*/
contract Function {
    EthereumStore blockStore;

    /*  The event verifier for the specific event being consumed. Each event would require a different event verifier to
        be deployed and each consumer would reference the relevant verifier to prove logs. */
    TriggerEventVerifier verifier;

    /* Custom event that fires when execution is performed successfully. */
    event Executed();

    /*  Constructor. Requires Ion contract address and all used event verifier contract addresses. In this case we only
        use one verifier. */
    constructor(address _storeAddr, address _verifierAddr) public {
        blockStore = EthereumStore(_storeAddr);
        verifier = TriggerEventVerifier(_verifierAddr);
    }

    /* This is the function that is intended to be executed upon successful verification of proofs */
    function execute() internal {
        emit Executed();
    }

    /*  
        verifyAndExecute

        Core parameters for verification
        param: _chainId (bytes32)   Chain ID of the chain that the event being consumed was emitted on. This may require
                                    altering to (bytes) if proofs from multiple chains are needed.
        param: _blockHash (bytes32) Block hash of block with event to be consumed. This may require altering to (bytes)
                                    if proofs from multiple chains are needed.
        param: _contractEmittedAddress (bytes20)    Contract address of the source of event emission. This may require
                                                    altering to (bytes) if proofs from multiple chains are needed.
        param: _path (bytes)    RLP-encoded transaction index of the relevant transaction that emitted the event being
                                consumed. If multiple proofs are required, multiple paths supplied must be RLP-encoded
                                and an extra (bool) parameter provided to denote multiple paths included.
        param: _tx (bytes) RLP-encoded transaction object provided by proof generation.
        param: _txNodes (bytes) RLP_encoded transaction nodes provided by proof generation.
        param: _receipt (bytes) RLP-encoded receipt object provided by proof generation.
        param: _receiptNodes (bytes) RLP-encoded receipt nodes provided by proof generation.

        Custom parameters for verification
        param: _expectedAddress (bytes20) The expected address value in the event parameter being consumed.

        This is the only public function apart from the constructor and is the only interface to this contract. This 
        function wraps the verification and execution which only fires after a successful slew of verifications. As
        noted, stack restrictions will make it harder to implement multiple event consumption. Suggestions made here may
        not be the best way to achieve this but are possible methods. It may end up requiring separate functions for
        each event and persisting the consumption state of each event per tx hash and using that to allow or prevent
        verified execution. In our case, it is simple as we only consume a single event.
        */
    function verifyAndExecute(
        bytes32 _chainId,
        bytes32 _blockHash,
        bytes20 _contractEmittedAddress,
        bytes memory _proof,
        bytes20 _expectedAddress
    ) public {
        bytes memory receipt =  blockStore.CheckProofs(_chainId, _blockHash, _proof);

        require( verifier.verify(_contractEmittedAddress, receipt, _expectedAddress), "Event verification failed." );
        execute();
    }
}

