// SPDX-License-Identifier: MIT

pragma solidity ^0.8.24;

import "./ILocker.sol";
import "./IERC20.sol";
import "./LockerManage.sol";
import "forge-std/console.sol";

contract LockerService is ILocker, LockerManage {
    struct Locker {
        address tokenERC20; // er
        uint256 amount;
        uint256 timelock;
        address[] users;
        uint256 part;
    }

    mapping(address => bytes32[]) public customers; // cutomers => lockerId
    mapping(bytes32 => Locker) public lockers; // lockerId => (user => amount)
    mapping(bytes32 => mapping(address => bool)) public claimed; // lockerId => (user => is_claimed)
    // mapping(bytes32 => address) public users; // lockerId => user
    mapping(address => uint256) public erc20Fees;

    event LockerCreated(bytes32 indexed lockerId, address indexed customer, Locker);
    event Claimed(bytes32 indexed lockerId, address indexed user, uint256 amount);

    error NotEnoughValue(uint256 excepted, uint256 actual);

    event ItIsNotTimeYet();

    constructor(uint256 _fee) {
        owner = msg.sender;
        feePrice = _fee;
    }

    // usual ETH locker
    function createLocker(address[] memory users) external payable returns (bytes32 lockerId) {
        require(msg.value > feePrice, "not enough ether");
        require(users.length > 0, "empty users list");

        uint256 amount = msg.value - feePrice;
        uint256 part = amount / users.length; // TODO: check

        Locker memory locker = Locker(address(0x0), amount, 0, users, part);

        lockerId = keccak256(abi.encodePacked(block.timestamp, msg.sender, abi.encode(locker)));

        lockerExists(lockerId);

        lockers[lockerId] = locker;

        customers[msg.sender].push(lockerId);
        emit LockerCreated(lockerId, msg.sender, locker);
    }

    function createTimeLocker(address[] memory users, uint256 timelock) external payable returns (bytes32 lockerId) {
        uint256 amount = msg.value - feePrice;
        uint256 part = users.length / amount; // TODO: проверить

        Locker memory locker = Locker(address(0x0), amount, timelock, users, part);

        lockerId = keccak256(abi.encodePacked(block.timestamp, msg.sender, abi.encode(locker)));
        lockerExists(lockerId);

        lockers[lockerId] = locker;

        customers[msg.sender].push(lockerId);
        emit LockerCreated(lockerId, msg.sender, locker);
    }

    function createERC20Locker(address[] memory users, address token, uint256 amount)
        external
        payable
        override
        returns (bytes32 lockerId)
    {
        uint256 lockerAmount = amount - feePrice;
        uint256 part = users.length / amount;

        _safeTransferFrom(IERC20(token), msg.sender, address(this), feePrice);

        Locker memory locker = Locker(token, lockerAmount, 0, users, part);

        lockerId = keccak256(abi.encodePacked(block.timestamp, msg.sender, abi.encode(locker)));
        lockerExists(lockerId);

        lockers[lockerId] = locker;
        customers[msg.sender].push(lockerId);
        erc20Fees[token] = erc20Fees[token] + feePrice;

        emit LockerCreated(lockerId, msg.sender, locker);
    }

    function createERC20TimeLocker(address[] memory users, address token, uint256 amount, uint256 timelock)
        external
        payable
        returns (bytes32 lockerId)
    {
        uint256 lockerAmount = amount - feePrice;
        uint256 part = users.length / amount;

        _safeTransferFrom(IERC20(token), msg.sender, address(this), feePrice);

        Locker memory locker = Locker(token, lockerAmount, timelock, users, part);

        lockerId = keccak256(abi.encodePacked(block.timestamp, msg.sender, abi.encode(locker)));
        lockerExists(lockerId);

        lockers[lockerId] = locker;
        customers[msg.sender].push(lockerId);
        erc20Fees[token] = erc20Fees[token] + feePrice;

        emit LockerCreated(lockerId, msg.sender, locker);
    }

    function claim(bytes32 lockerId) external override returns (bool result) {
        require(claimed[lockerId][msg.sender] != true, "already claimed");
        require(userInLocker(lockerId, msg.sender) == true, "Not in set of members");
        require(block.timestamp > lockers[lockerId].timelock, "It's not time yet");

        claimed[lockerId][msg.sender] = true;
        if (lockers[lockerId].tokenERC20 == address(0x0)) {
            payable(msg.sender).transfer(lockers[lockerId].part);
        } else {
            _safeTransferFrom(IERC20(lockers[lockerId].tokenERC20), address(this), msg.sender, lockers[lockerId].part);
        }
        emit Claimed(lockerId, msg.sender, lockers[lockerId].part);

        return true;
    }

    /// GETTERS
    function getOrdersByCustomer(address customer) external view returns (bytes32[] memory orders) {
        return customers[customer];
    }

    function getLockerById(bytes32 lockerId) external view returns (Locker memory locker) {
        return lockers[lockerId];
    }

    function isClaimed(bytes32 lockerId, address user) external view returns (bool) {
        return claimed[lockerId][user];
    }

    function lockerExists(bytes32 lockerId) internal view {
        if (lockers[lockerId].amount > 0) {
            revert("Locker already exists");
        }
    }

    function userInLocker(bytes32 lockerId, address user) internal view returns (bool) {
        for (uint256 i = 0; i < lockers[lockerId].users.length; i++) {
            if (lockers[lockerId].users[i] == user) {
                return true;
            }
        }
        return false;
    }

    function getLockerPart(bytes32 lockerId) external view returns (uint256 part) {
        part = lockers[lockerId].part;
        return part;
    }

    ///  WITHDRAW
    function withdraw() public onlyOwner {
        payable(owner).transfer(address(this).balance);
    }

    function withdrawERC20(IERC20 token) public onlyOwner {
        _safeTransferFrom(token, address(this), owner, erc20Fees[address(token)]);
        erc20Fees[address(token)] = 0;
    }

    function _safeTransferFrom(IERC20 token, address sender, address recipient, uint256 amount) private {
        bool sent = token.transferFrom(sender, recipient, amount);
        require(sent, "Token transfer failed");
    }
}
