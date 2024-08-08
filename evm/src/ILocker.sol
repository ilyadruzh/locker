// SPDX-License-Identifier: MIT

pragma solidity ^0.8.24;

interface ILocker {
    function createLocker(address[] calldata users) external payable returns (bytes32 orderId);

    function createTimeLocker(address[] calldata users, uint256 timelock)
        external
        payable
        returns (bytes32 orderId);

    function claim(bytes32 orderId) external returns (bool result);

    function createERC20Locker(address[] memory users, address token, uint256 amount)
        external
        payable
        returns (bytes32 lockerId);

    function createERC20TimeLocker(address[] memory users, address token, uint256 amount, uint256 timelock)
        external
        payable
        returns (bytes32 lockerId);
}
