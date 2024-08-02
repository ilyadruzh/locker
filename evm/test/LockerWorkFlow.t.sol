// SPDX-License-Identifier: UNLICENSED
pragma solidity >=0.7.0 <0.9.0;
pragma abicoder v2;

import "forge-std/console.sol";
import {LockerSetUp} from "./LockerSetUp.t.sol";

// forge test --match-contract WorkFlow
contract WorkFlow is LockerSetUp {
    address customer = address(0x01);

    address user1 = address(0x11);
    address user2 = address(0x12);
    address user3 = address(0x13);

    // createLockerWithAmount(address[] memory users, uint256 amount)
    function test_createLockerWithAmount() public {
        address[] memory _users;
        _users[0] = user1;
        _users[1] = user2;
        _users[2] = user3;

        vm.prank(customer);
        bytes32 createdLockerId = locker.createLockerWithAmount(_users, 1000);
        bytes32[] memory storedLockerIds = locker.getOrdersByCustomer(customer);

        assertEqUint(uint256(createdLockerId), uint256(storedLockerIds[0]));
    }
}
