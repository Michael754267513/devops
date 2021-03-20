pragma solidity ^0.4.18;
/// @title hello
/// @notice
/// @dev
contract hello {
    /// @notice
    string greeting;

    /// @notice
    /// @dev
    /// @param _greeting
    /// @return
    function hello(string _greeting) public {
        greeting = _greeting;
    }

    /// @notice
    /// @dev
    /// @return
    function say() constant public returns (string) {
        return greeting;
    }
}
