// SPDX-License-Identifier: MIT

pragma solidity >=0.8.2 <0.9.0;

import "./ILocker.sol";

contract Locker is ILocker {
    address public owner;
    uint256 public feePrice;
    uint256 public feeAmount;

    // cutomers => lockerId
    mapping(address => bytes32[]) public customers;

    // lockerId => (user => amount)
    mapping(bytes32 => mapping(address => uint256)) public lockers;

    event LockerCreated(bytes32 indexed lockerId, address indexed customer, address[] users, uint256 amount);

    error NotEnoughValue(uint256 excepted, uint256 actual);

    modifier onlyOwner() {
        require(msg.sender == owner, "you are not owner");
        _;
    }

    constructor(uint256 _fee) {
        owner = msg.sender;
        feePrice = _fee;
    }

    function createLockerWithAmount(address[] memory users, uint256 amount)
        external
        payable
        override
        returns (bytes32 lockerId)
    {
        uint256 allAmount = users.length * amount;

        if (msg.value >= allAmount + feePrice) {
            feeAmount += msg.value - allAmount;

            // TODO: если автор решит два раза вызвать однотипную раздачу?
            // TODO: если автор отправит два адреса с разными количесставми?
            lockerId = keccak256(abi.encodePacked(msg.sender, users, amount));

            for (uint8 count = 0; count < users.length; count++) {
                lockers[lockerId][users[count]] = amount;
            }

            customers[msg.sender].push(lockerId);

            emit LockerCreated(lockerId, msg.sender, users, amount);
        } else {
            revert NotEnoughValue(allAmount, msg.value);
        }
    }

    function createLockerWithoutAmount(address[] memory users) public payable returns (bytes32 lockerId) {
        uint256 amount = msg.value / users.length;

        lockerId = keccak256(abi.encodePacked(msg.sender, users, amount));

        for (uint8 count = 0; count < users.length; count++) {
            lockers[lockerId][users[count]] = amount;
        }

        customers[msg.sender].push(lockerId);

        emit LockerCreated(lockerId, msg.sender, users, amount);
    }

    function claim(bytes32 lockerId) external override returns (bool result) {
        require(lockers[lockerId][msg.sender] != 0, "already claimed");

        lockers[lockerId][msg.sender] = 0;

        payable(msg.sender).transfer(lockers[lockerId][msg.sender]);

        return true;
    }

    function changeFeePrice(uint256 newFeePrice) public onlyOwner {
        feePrice = newFeePrice;
    }

    function changeOwner(address newOwner) public onlyOwner {
        owner = newOwner;
    }

    function withdraw() public onlyOwner {
        payable(msg.sender).transfer(feeAmount);
    }

    function getOrdersByCustomer(address customer) external view returns (bytes32[] memory orders) {
        return customers[customer];
    }
}
