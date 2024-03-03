pragma solidity 0.6.8;
pragma experimental ABIEncoderV2;

import "BN256G1.sol";
import "BN256G2.sol";

contract bn256test{
    struct Message {
        uint256 pid; // 作者ID
        uint256 mid; // 消息ID
    }

    Message[] messages; // 存储消息的数组

    // 事件用于通知消息添加
    event MessageAdded(uint256 indexed pid, uint256 indexed mid);

    // 添加消息到列表
    function addMessage(uint256 _pid, uint256 _mid) public {
        Message memory newMessage = Message(_pid, _mid);
        messages.push(newMessage);

        // 触发事件
        emit MessageAdded(_pid, _mid);
    }

    function searchMessagesByMid(uint256 _mid) public view returns (uint256[] memory) {
        uint256[] memory matchingPids = new uint256[](messages.length);
        uint256 count = 0;

        for (uint256 i = 0; i < messages.length; i++) {
            if (messages[i].mid == _mid) {
                matchingPids[count] = messages[i].pid;
                count++;
            }
        }

        // 去掉多余的零元素
        assembly {
            mstore(matchingPids, count)
        }

        return matchingPids;
    }
    
    // The prime q in the base field F_q for G1 and G2
    uint256 q = 21888242871839275222246405745257275088696311157297823662689037894645226208583;
    
    struct G1Point {
        uint256 x;
        uint256 y;
    }
    
    // Encoding of field elements is: X[0] * z + X[1]
    struct G2Point {
        uint256[2] x;
        uint256[2] y;
    }
    
    //@return the generator of G1
    function get_P1() internal pure returns (G1Point memory) {
        return G1Point(1,2);
    }
    
    //@return the generator of G2
    function get_P2() internal pure returns (G2Point memory) { // RE+IM
        return G2Point(
            [11559732032986387107991004021392285783925812861821192530917403151452391805634,
            10857046999023057135944570762232829481370756359578518086990519993285655852781],

            [4082367875863433681332203403145435568316851327593401208105741076214120093531,
            8495653923123431417604973247489272438418190587263600148770280649306958101930]
        );
    }

    // 我們論文用的參數
    uint256[8] data_ours =  [7295212943885535137761043108359461309601647987237926731738265796461598060035,4623803287399109611580097681536749355734855403003917189588428225284383471135,3745899144130829966074888229078829376511882025583224851938487878887116228575,19221218966914293210740534720977385257377733493369671293034520685030420978378,10331247486788799231356728707978625961226033367477117884161463113009346653253,15853000633514356373919958251715802396765840738295073096718922483158181535638,9562193070024042145110084860572806478440836752229285876189280312711117657447,14844603942051752802939650544486160631727010044955941074529948516911519590923];
    
    // 我們論文用來進行 verify 的 function
    function testOURS()public returns(bool)
    {
        G1Point memory PID_1 = G1Point(data_ours[0], data_ours[1]);
        G1Point memory addhash = G1Point(data_ours[2], data_ours[3]);
        G2Point memory s = G2Point([data_ours[4],data_ours[5]],[data_ours[6],data_ours[7]]);
        return pairing_check(get_P1(), s, g1add(PID_1,addhash), get_P2());
    }
    
    // 以下的 functions 可以在 smart contract 上進行定義在橢圓曲線上的運算，參考用
    function g2add(G2Point memory a, G2Point memory b)internal view returns(G2Point memory r)
    {
        (uint256 x_im, uint256 x_re, uint256 y_im, uint256 y_re) = BN256G2.ecTwistAdd(a.x[1], a.x[0], a.y[1], a.y[0], b.x[1], b.x[0], b.y[1], b.y[0]);
        return G2Point([x_re, x_im], [y_re, y_im]);
    }
    
    function g2mul(uint256 s, G2Point memory a)internal view returns(G2Point memory r)
    {
        (uint256 x_im, uint256 x_re, uint256 y_im, uint256 y_re) = BN256G2.ecTwistMul(s, a.x[1], a.x[0], a.y[1], a.y[0]);
        return G2Point([x_re, x_im], [y_re, y_im]);
    }
    
    function g1add(G1Point memory a, G1Point memory b)internal returns(G1Point memory r)
    {
        (uint256 x, uint256 y) = BN256G1.add([a.x, a.y, b.x, b.y]);
        return G1Point(x, y);
    }
    
    function g1mul(uint256 s, G1Point memory a)internal returns(G1Point memory r)
    {
        (uint256 x, uint256 y) = BN256G1.multiply([a.x, a.y, s]);
        return G1Point(x, y);
    }
    
    function g1negate(G1Point memory a) internal view returns(G1Point memory)
    {
        
        return G1Point(a.x, q-a.y);
    }
    
    function g2negate(G2Point memory a) internal view returns(G2Point memory)
    {
        
        return G2Point(
            [a.x[0], a.x[1]],
            [q - a.y[0], q - a.y[1]]);
    }
    
    function pairing_check(G1Point memory P, G2Point memory Q, G1Point memory R, G2Point memory S)internal returns(bool){
        P = g1negate(P);
        return BN256G1.bn256CheckPairing([P.x, P.y, Q.x[0], Q.x[1], Q.y[0], Q.y[1], R.x, R.y, S.x[0], S.x[1], S.y[0], S.y[1]]);
    }
    
}

