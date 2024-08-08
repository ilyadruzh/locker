// SPDX-License-Identifier: UNLICENSED
pragma solidity >=0.7.0 <0.9.0;
pragma abicoder v2;

import "forge-std/console.sol";
import {LockerSetUp} from "./LockerSetUp.t.sol";
import {LockerService} from "../src/LockerService.sol";

// forge test --match-contract LockerWorkFlow
contract LockerWorkFlow is LockerSetUp {
    // createLockerWithAmount(address[] memory users, uint256 amount)

    function testFail_createLockerWithoutUsers() public {
        address[] memory epmtyUsers;

        vm.prank(customer);
        vm.deal(customer, 1 ether);
        createdLockerId = lockerService.createLocker{value: 1 ether}(epmtyUsers);

        bytes32[] memory storedLockerIds = lockerService.getOrdersByCustomer(customer);
        assertEqUint(uint256(createdLockerId), uint256(storedLockerIds[0]));
    }

    function testFail_createLockerWithoutFund() public {
        vm.prank(customer);
        createdLockerId = lockerService.createLocker{value: 1 ether}(_users);

        bytes32[] memory storedLockerIds = lockerService.getOrdersByCustomer(customer);
        assertEqUint(uint256(createdLockerId), uint256(storedLockerIds[0]));
    }

    function test_createLockerWithSmallFund_101() public {
        vm.prank(customer);
        vm.deal(customer, 101);
        bytes32 lockerWithSmallFund = lockerService.createLocker{value: 101}(_users);

        LockerService.Locker memory locker = lockerService.getLockerById(lockerWithSmallFund);
        

        assertEqUint(locker.amount, 1);
    }

    function test_authorizedClaim() public {
        vm.prank(user3);
        bool res = lockerService.claim(createdLockerId);
        assertTrue(res);

        bool claimed = lockerService.isClaimed(createdLockerId, user3);
        assertTrue(claimed, "Not claimed");
    }

    function test_unauthorizedClaim() public {
        vm.prank(user4);
        vm.expectRevert("Not in set of members");
        lockerService.claim(createdLockerId);
    }

    // function test_createTimeLocker() public {
    //     address[] memory _users;
    //     _users[0] = user1;
    //     _users[1] = user2;
    //     _users[2] = user3;

    //     uint256 timelock;

    //     vm.prank(customer);
    //     bytes32 createdLockerId = lockerService.createTimeLocker(_users, timelock);
    //     bytes32[] memory storedLockerIds = lockerService.getOrdersByCustomer(customer);

    //     assertEqUint(uint256(createdLockerId), uint256(storedLockerIds[0]));
    // }
}
