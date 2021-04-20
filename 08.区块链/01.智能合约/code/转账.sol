pragma solidity ^0.5.16;


/// @title Account
/// @notice
/// @dev
contract Account {

    /*
        address(this)   合约地址
        msg.sender   from地址
    */

    address public _owner;
    /// @notice
    address payable public _ownerpublish;
    // address payable _acc = 0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1;

    constructor() public {
        _owner = address(this); // 合约地址
        _ownerpublish = msg.sender; // 合约发布者
    }


    /// @notice
    /// @dev
    /// @return 当前帐号钱包地址
    function getSend() public view returns(address) {
        return msg.sender; // 当前账号
    }

    /// @notice
    /// @dev
    function transfer()  public payable {
        // transfer 转账给_ownerpublish
        _ownerpublish.transfer(msg.value);
    }

    /// @notice
    /// @dev
    /// @return 发布者余额
    function acc() public view returns(uint) {
        return _ownerpublish.balance; // 合约发布者余额
    }


}
