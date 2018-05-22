pragma solidity ^0.4.18;

contract Ion {
    event OnTradeInitiated(bytes32 _tradeId, address _owner);
    event OnDeposit(bytes32 _tradeId, address _owner);
    event OnRefund(bytes32 _tradeId, address _owner, bytes _refundRef);
    event OnWithdraw(bytes32 _tradeId, address _recipient, bytes _withdrawRef);

    struct TradeAgreement {
        address owner;
        address recipient;
        address token;
        uint256 value;
        bytes32 withdrawHash;
        bytes32 refundHash;
    }
}
