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

import "./ERC223_receiving_contract.sol";
import "./openzeppelin-token/ERC20/StandardToken.sol";
/// @title an interface for Token.
/// @notice it contains various utility functions that will be used by a token.
interface TokensNetworkInterface {
    function chain_id() external view returns (uint256); //必须大于0
    function contract_version() external view returns (string); //必须是0.6_
    function signature_prefix() external view returns (string); //必须是'\x19Spectrum Signed Message:\n'

}
contract SMTToken is  StandardToken {

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
    )
    public
    {
        name = "SMTToken";            //名字必须是这个,临时通过名字来检测功能支持                        // Set the name for display purposes
        decimals = 18;                            // Amount of decimals for display purposes
        symbol = _tokenSymbol;                               // Set the symbol for display purposes
        totalSupply_ = 1; //这个参数没啥用,专用于token network,维持为0
        tokenNetwork=_tokenNetwork;
        //要检测指定的这个tokennetwork地址是有效的
        TokensNetworkInterface tn=TokensNetworkInterface(_tokenNetwork);
        require(tn.chain_id()>0);
        require(keccak256(tn.contract_version())==keccak256("0.6._"));
        require(keccak256(tn.signature_prefix())==keccak256('\x19Spectrum Signed Message:\n'));
    }
    function () external { revert(); }


    //调用这个合约存钱,创建channel,调用方式和调用  只能给指定合约充值
    function buyAndTransfer( bytes _data) public payable {
        require(msg.value>0); //不能充值0
        //只允许充值到tokensNetwork合约中
        ERC223ReceivingContract receiver = ERC223ReceivingContract(tokenNetwork);
        receiver.tokenFallback(msg.sender, msg.value, _data);
    }
    //钱退回账户的过程,特殊处理
    /**
    * @dev Transfer token for a specified address
    * @param _to The address to transfer to.
    * @param _value The amount to be transferred.
    */
    function transfer(address _to, uint256 _value) public returns (bool) {
        require(msg.sender==tokenNetwork); //只能是token network转回来,其他没用.
//        require(_value <= balances[msg.sender]);
        require(_to != address(0));
//        balances[msg.sender] = balances[msg.sender].sub(_value);
        _to.transfer(_value); // 退回给这个账户
        return true;
    }

    //以下函数都应该禁用,只是为了符合erc223的接口罢了
    function transferFrom(
        address _from,
        address _to,
        uint256 _value
    )
    public
    returns (bool){
        revert();
    }
    function approve(address _spender, uint256 _value) public returns (bool) {
        revert();
    }
    function increaseApproval(address _spender, uint256 _addedValue) public returns (bool) {
        revert();
    }
    function decreaseApproval(address _spender, uint256 _subtractedValue) public returns (bool)
    {
        revert();
    }
}
