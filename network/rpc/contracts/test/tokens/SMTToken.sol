/*
This Token Contract implements the standard token functionality (https://github.com/ethereum/EIPs/issues/20) as well as the following OPTIONAL extras intended for use by humans.

In other words. This is intended for deployment in something like a Token Factory or Mist wallet, and then used by humans.
Imagine coins, currencies, shares, voting weight, etc.
Machine-based, rapid creation of many tokens would not necessarily need these extra features or will be minted in other manners.

1) Initial Finite Supply (upon creation one specifies how much is minted).
2) In the absence of a token registry: Optional Decimal, Symbol & Name.
3) Optional approveAndCall() functionality to notify a contract if an approval() has occurred.

.*/

pragma solidity ^0.4.24;
import "./openzeppelin-token/ERC20/StandardToken.sol";

contract SMTToken is StandardToken {



    /*
    NOTE:
    The following variables are OPTIONAL vanities. One does not have to include them.
    They allow one to customise the token contract & in no way influences the core functionality.
    Some wallets/interfaces might not even bother to look at this information.
    */
    string public name;                   //fancy name: eg Simon Bucks
    uint8 public decimals;                //How many decimals to show. ie. There could 1000 base units with 3 decimals. Meaning 0.980 SBX = 980 base units. It's like comparing 1 wei to 1 ether.
    string public symbol;                 //An identifier: eg SBX
    string public version = 'H0.1';       //human 0.1 standard. Just an arbitrary versioning scheme.
    address public  tokenNetwork;  //特殊地址 tokens network


    constructor(
        string _tokenSymbol,
        address _tokenNetwork
    ) payable
    public
    {
        name = "SMTToken";            //名字必须是这个,临时通过名字来检测功能支持                        // Set the name for display purposes
        decimals = 18;                            // Amount of decimals for display purposes
        symbol = _tokenSymbol;                               // Set the symbol for display purposes
        require(msg.value>0);
        balances[msg.sender] = msg.value;
        totalSupply_ = msg.value;
        tokenNetwork=_tokenNetwork;
    }
    function () external { revert(); }
    /**
   * @dev Transfer the specified amount of tokens to the specified address.
   *      Invokes the `tokenFallback` function if the recipient is a contract.
   *      The token transfer fails if the recipient is a contract
   *      but does not implement the `tokenFallback` function
   *      or the fallback function to receive funds.
   *
   * @param _to    Receiver address.
   * @param _value Amount of tokens that will be transferred.
   * @param _data  Transaction metadata.
   */
    function transfer(address _to, uint256 _value, bytes _data) public {
        // Standard function transfer similar to ERC20 transfer with no _data .
        // Added due to backwards compatibility reasons .
        uint codeLength;

        assembly {
        // Retrieve the size of the code on target address, this needs assembly .
            codeLength := extcodesize(_to)
        }

        balances[msg.sender] = balances[msg.sender].sub(_value);
        balances[_to] = balances[_to].add(_value);
        if(codeLength>0) {
            ERC223ReceivingContract receiver = ERC223ReceivingContract(_to);
            receiver.tokenFallback(msg.sender, _value, _data);
        }
        emit Transfer(msg.sender, _to, _value, _data);
    }
    //调用这个合约存钱,创建channel,调用方式和调用  function transfer(address _to, uint256 _value, bytes _data) external  是一样的
    function buyAndTransfer(address _to, bytes _data) public payable {
        require(msg.value>0); //不能充值0
        buy(); //先充值,充值完成就转走
        transfer(_to,msg.value,_data);
    }
    //钱退回账户的过程,特殊处理一下.
    /**
* @dev Transfer token for a specified address
* @param _to The address to transfer to.
* @param _value The amount to be transferred.
*/
    function transfer(address _to, uint256 _value) public returns (bool) {
        require(_value <= balances[msg.sender]);
        require(_to != address(0));

        balances[msg.sender] = balances[msg.sender].sub(_value);
        balances[_to] = balances[_to].add(_value);
        emit Transfer(msg.sender, _to, _value);
        if (msg.sender==tokenNetwork) { //来自特殊地址的转账,直接提现到账户中
            sell(_value);
        }
        return true;
    }

    //将指定账户中的token,退回到相应账户
    function sellFrom(address _from, uint256 amount) internal {
        require(balances[_from] >= amount, "Insufficient balance.");

        balances[_from] -= amount; //先减去这个账户中token
        totalSupply_ -= amount;
        _from.transfer(amount); // 退回给这个账户
        //销毁token
        emit Transfer(_from, address(0), amount);
    }
    //允许充值到token中,但是不充值到channel中
    function buy() public payable {
        balances[msg.sender] += msg.value;
        totalSupply_ += msg.value;
        //增发token
        emit Transfer(address(0), msg.sender, msg.value);
    }
    //从token 退回到个人账户中
    function sell(uint256 amount) public {
        sell(msg.sender,amount);
    }


}
