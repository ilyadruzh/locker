// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {LockerService} from "../src/LockerService.sol";

contract LockerSetUp is Test {
    LockerService public lockerService;
    address owner = address(0xFF);
    uint256 _fee = 100;
    address customer = address(0x01);
    address[] _users;
    address user1 = address(0x11);
    address user2 = address(0x12);
    address user3 = address(0x13);
    address user4 = address(0x14);
    bytes32 public createdLockerId;

    function setUp() public {
        vm.prank(owner);
        lockerService = new LockerService(_fee);

        _users.push(user1);
        _users.push(user2);
        _users.push(user3);

        vm.prank(customer);
        vm.deal(customer, 1 ether);
        createdLockerId = lockerService.createLocker{value: 1 ether}(_users);

        bytes32[] memory storedLockerIds = lockerService.getOrdersByCustomer(customer);
        assertEqUint(uint256(createdLockerId), uint256(storedLockerIds[0]));
        // LockerService.Locker memory locker = lockerService.getLockerById(storedLockerIds[0]);
    }
}
