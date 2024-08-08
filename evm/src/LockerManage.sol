// SPDX-License-Identifier: MIT

pragma solidity ^0.8.24;

contract LockerManage {
    address public whitelist;
    address public owner;
    uint256 public feePrice;
    uint256 public feeAmount;

    modifier onlyOwner() {
        require(msg.sender == owner, "You are not owner!");
        _;
    }

    function changeFeePrice(uint256 newFeePrice) public onlyOwner {
        feePrice = newFeePrice;
    }

    function changeOwner(address newOwner) public onlyOwner {
        owner = newOwner;
    }
}
