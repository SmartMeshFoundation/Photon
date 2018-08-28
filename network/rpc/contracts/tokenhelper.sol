pragma solidity ^0.4.24;
import "./Token.sol";

/// @title
contract ApproveAndCallFallBack {
    function receiveApproval(address from, uint256 _amount, address _token, bytes _data) public returns (bool success);
}

/// @title
contract TokenHelper {

    function approveAndCall(address token,address _spender, uint256 _amount, bytes _extraData
    ) public returns (bool success) {
        require(Token(token).approve(_spender, _amount));
        require(ApproveAndCallFallBack(_spender).receiveApproval(
                msg.sender,
                _amount,
                this,
                _extraData
            ));
        return true;
    }
}
