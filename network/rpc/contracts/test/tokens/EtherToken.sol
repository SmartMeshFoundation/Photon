pragma solidity ^0.4.24;

import "./openzeppelin-token/ERC20/StandardToken.sol";


contract EtherToken is StandardToken {

    function buy() public payable {
        balances[msg.sender] += msg.value;
        totalSupply_ += msg.value;

        emit Transfer(address(0), msg.sender, msg.value);
    }

    function sell(uint256 amount) public {
        require(balances[msg.sender] >= amount, "Insufficient balance.");

        balances[msg.sender] -= amount;
        totalSupply_ -= amount;
        msg.sender.transfer(amount);

        emit Transfer(msg.sender, address(0), amount);
    }
}