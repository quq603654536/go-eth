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

var client *ethclient.Client

func printBlockInfo() {
	// 获取最新的区块头信息
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// 打印最新区块号
	fmt.Println(header.Number.String())

	// 将区块号转换为 big.Int 类型
	blockNumber := big.NewInt(header.Number.Int64())
	// 获取指定区块号的完整区块信息
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	// 打印区块信息：
	// 1. 区块号
	fmt.Println(block.Number().Uint64())
	// 2. 区块时间戳
	fmt.Println(block.Time())
	// 3. 区块难度值
	fmt.Println(block.Difficulty().Uint64())
	// 4. 区块哈希值
	fmt.Println(block.Hash().Hex())
	// 5. 区块中的交易数量
	fmt.Println(len(block.Transactions()))

	// 获取区块中的交易总数（另一种方式）
	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}

	// 打印交易总数
	fmt.Println(count)
}

func sendEth() {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(0.0001 * 1e18)

	fmt.Println("转账wei:\n", value)

	gasLimit := uint64(21000)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x5A0615f4091bc72a02b96568d7d7377885bD6eA0")

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// 在发送交易前添加余额检查
	balance, err := client.BalanceAt(context.Background(), fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

	// 计算所需总金额 = 转账金额 + gas费用
	gasCost := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	totalCost := new(big.Int).Add(value, gasCost)

	if balance.Cmp(totalCost) < 0 {
		log.Fatalf("地址%s 余额不足: 需要 %s wei, 但只有 %s wei", fromAddress, totalCost.String(), balance.String())
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}

func main() {
	var err error
	// 加载 .env 文件
	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("测试网地址", os.Getenv("SEPOLIA_RPC_URL"))

	// 连接到以太坊节点（这里使用的是 Scroll Sepolia 测试网）
	client, err = ethclient.Dial(os.Getenv("SEPOLIA_RPC_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// printBlockInfo()
	sendEth()
}
