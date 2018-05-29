pragma solidity ^0.4.18;

contract Ion {
    event OnTradeInitiated(bytes32 _tradeId, address _owner);
    event OnDeposit(bytes32 _tradeId, address _owner);
    event OnRefund(bytes32 _tradeId, address _owner, bytes _refundRef);
    event OnFundsUnlocked(bytes32 _tradeId, bool unlocked);
    event OnWithdraw(bytes32 _tradeId, address _recipient, bytes _withdrawRef);

    struct TradeAgreement {
        address initiator;
        address counterparty;
        address initiatorToken;
        address counterpartyToken;
        uint256 value;
        bytes32 withdrawHash;
        bytes32 refundHash;
        bool fundsUnlocked;
        address owner;
    }
}
