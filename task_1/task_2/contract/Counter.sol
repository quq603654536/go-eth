// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Counter {
    uint256 private counter;

    function incr() external {
        ++counter;
    }

    function getCounter() external view returns (uint256) {
        return counter;
    }
}
