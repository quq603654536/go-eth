package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

// 声明全局变量，用于存储以太坊客户端连接
var client *ethclient.Client

// printBlockInfo 函数用于获取和打印区块信息
func printBlockInfo() {
	// 获取最新的区块头信息
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// 打印最新区块号
	fmt.Println("最新区块号:", header.Number.String())

	// 将区块号转换为 big.Int 类型
	blockNumber := big.NewInt(header.Number.Int64())
	// 获取指定区块号的完整区块信息
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	// 打印区块详细信息
	fmt.Println("区块号:", block.Number().Uint64())
	fmt.Println("区块时间戳:", block.Time())
	fmt.Println("区块难度值:", block.Difficulty().Uint64())
	fmt.Println("区块哈希值:", block.Hash().Hex())
	fmt.Println("区块中的交易数量:", len(block.Transactions()))

	// 获取区块中的交易总数（使用另一种方式）
	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("区块交易总数:", count)
}

// sendEth 函数用于发送以太币交易
func sendEth() {
	// 从环境变量获取私钥并转换为ECDSA格式
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// 从私钥获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	// 从公钥获取发送方地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取账户当前的nonce值（交易序号）
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// 设置转账金额（0.0001 ETH，转换为wei单位）
	value := big.NewInt(0.0001 * 1e18)
	fmt.Println("转账金额(wei):", value)

	// 设置gas限制（标准ETH转账需要21000 gas）
	gasLimit := uint64(21000)

	// 获取当前建议的gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 设置接收方地址
	toAddress := common.HexToAddress("0x5A0615f4091bc72a02b96568d7d7377885bD6eA0")

	// 创建新的交易
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// 获取链ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 使用私钥对交易进行签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 在发送交易前检查账户余额
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 计算交易所需的总金额（转账金额 + gas费用）
	gasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	totalCost := new(big.Int).Add(value, gasCost)

	// 检查余额是否足够
	if balance.Cmp(totalCost) < 0 {
		log.Fatalf("地址%s 余额不足: 需要 %s wei, 但只有 %s wei", fromAddress, totalCost.String(), balance.String())
	}

	// 发送交易到网络
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("交易已发送，交易哈希: %s", signedTx.Hash().Hex())
}

func main() {
	var err error
	// 加载环境变量配置文件
	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("测试网地址:", os.Getenv("SEPOLIA_RPC_URL"))

	// 连接到以太坊测试网节点
	client, err = ethclient.Dial(os.Getenv("SEPOLIA_RPC_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// 执行区块信息查询
	printBlockInfo()
	// 执行ETH转账
	sendEth()
}
