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
    address payable _to = 0xFFcf8FDEE72ac11b5c542428B35EEF5769C409f0;

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

    /// @notice  2-8 分账
    /// @dev
    function transfer()  public payable {
        // transfer 转账给_ownerpublish
        _ownerpublish.transfer(msg.value * 20 / 100);
        _to.transfer(msg.value * 80 / 100 );
    }

    /// @notice
    /// @dev
    /// @return 发布者余额
    function acc() public view returns(uint) {
        return _ownerpublish.balance; // 合约发布者余额
    }


}
