// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {Locker} from "../src/Locker.sol";

contract LockerSetUp is Test {
    Locker public locker;
    address owner = address(0xFF);

    uint256 _fee = 100;

    function setUp() public {
        vm.prank(owner);
        locker = new Locker(_fee);
    }
}
