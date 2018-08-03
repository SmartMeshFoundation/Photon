pragma solidity ^0.4.24;
contract ERC223 {
    event Transfer(address indexed from, address indexed to, uint256 value, bytes indexed data);
    function transfer(address to, uint256 value, bytes data) external  ;
}