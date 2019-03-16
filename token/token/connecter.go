package token

import (
	"context"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"fmt"
)

var (
	//  DAIToken 合约地址
	CoinAddr = common.HexToAddress("0x6cf944093b7828f4e345509b0400f56040aeae36")
	//  DAIToken 合约地址Hash
	CoinHash = CoinAddr.Hash()
)

// Connecter SuperCoin连接者
type Connecter struct {
	ctx  context.Context
	conn *ethclient.Client
	token *DAIToken
}

// NewConnecter 生成一个SuperCoin连接者
func NewConnecter() *Connecter {
	fmt.Println("conner")
	// Dial这里支持传入 ws、http、ipc的多种链接
	// 如果是经常需要调用最好还是使用 ws 方式保持通讯状态
	conn, err := ethclient.Dial("https://ropsten.infura.io/v3/81025a29acba47869578d1bd1934aa9e")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	token, err := NewDAIToken(CoinAddr, conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}
	return &Connecter{
		ctx:  context.Background(),
		conn: conn,
		token: token,
	}
}

// NewConnecterWithDeploy 部署合约，并创建一个connecter
func NewConnecterWithDeploy(ownerAuth *bind.TransactOpts) *Connecter {
	conn, err := ethclient.Dial("https://ropsten.infura.io/v3/81025a29acba47869578d1bd1934aa9e")
	if err != nil {
		panic(err)
	}
	_, tx, token, err := DeployDAIToken(ownerAuth, conn)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	CoinAddr, err = bind.WaitDeployed(ctx, conn, tx)
	if err != nil {
		panic(err)
	}
	CoinHash = CoinAddr.Hash()
	return &Connecter{
		ctx:  ctx,
		conn: conn,
		token: token,
	}
}

// BlockNumber 当前块高度
func (c *Connecter) BlockNumber() *big.Int {
	blockNumber, err := c.conn.BlockByNumber(c.ctx, nil)
	if err != nil {
		log.Fatal("Get block number error", err)
		return big.NewInt(0)
	}
	return blockNumber.Number()
}

// ContractName DAIToken名称
func (c *Connecter) ContractName() string {
	name, err := c.token.Name(nil)
	if err != nil {
		log.Fatal("Get contract name error", err)
		return ""
	}
	return name
}

// BalanceOfEth 以太币余额
func (c *Connecter) BalanceOfEth(addr common.Address) *big.Int {
	balance, err := c.conn.BalanceAt(c.ctx, addr, nil)
	if err != nil {
		log.Fatal("Get address eth balance error", err)
		return big.NewInt(0)
	}
	return balance
}

// BalanceOfCoin DAIToken余额
func (c *Connecter) BalanceOfToken(addr common.Address) *big.Int {
	balance, err := c.token.BalanceOf(nil, addr)
	if err != nil {
		log.Fatal("Get address SuperCoin balance error", err)
		return big.NewInt(0)
	}
	return balance
}

// TotalSupply DAIToken已铸币量
func (c *Connecter) TotalSupply() *big.Int {
	totalSupply, err := c.token.TotalSupply(nil)
	if err != nil {
		log.Fatal("Get SuperCoin totalSuply error", err)
		return big.NewInt(0)
	}
	return totalSupply
}

// AuthAccount 解锁账户
// 正式使用时候，此处应该限制GasPrice和GasLimit
func AuthAccount(key, passphrase string) *bind.TransactOpts {
	auth, err := bind.NewTransactor(strings.NewReader(key), passphrase)
	if err != nil {
		log.Fatalf("Auth account error: %v", err)
		return nil
	}
	log.Printf("Gas price is: %v, Gas limit is: %v", auth.GasPrice, auth.GasLimit)
	return auth
}

// TransferCoin DAIToken转账
func (c *Connecter) TransferToken(fromAuth *bind.TransactOpts, to common.Address, amount *big.Int) bool {
	tx, err := c.token.Transfer(fromAuth, to, amount)
	if err != nil {
		log.Fatalf("Transfer SuperCoin from %s to %s amount %s error: %v", fromAuth.From.String(), to.String(), amount.String(), err)
		return false
	}
	// 等待执行
	receipt, err := bind.WaitMined(c.ctx, c.conn, tx)
	if err != nil {
		log.Fatalf("Wait Transfer Mined error: %v", err)
		return false
	}
	return receipt.Status == 1
}

// TransferLogs SuperCoin转账记录
func (c *Connecter) TransferLogs(froms []common.Address, tos []common.Address) {
	iter, err := c.token.FilterTransfer(&bind.FilterOpts{}, froms, tos)
	defer iter.Close()
	if err != nil {
		panic(err)
	}
	for {
		if !iter.Next() {
			break
		}
		log.Printf("Transfer log: %s", iter.Event.String())
	}
}

func (s *DAITokenTransfer) String() string {
	return fmt.Sprintf("From: %s, To: %s, Amount: %s", s.From.Hex(), s.To.Hex(), s.Tokens)
}
