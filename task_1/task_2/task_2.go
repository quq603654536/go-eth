package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	Counter "github.com/wkm/go-eth/task_1/task_2/contract"
)

// 声明全局变量，用于存储以太坊客户端连接
var client *ethclient.Client

func main() {
	var err error

	// 加载环境变量配置文件
	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("测试网地址", os.Getenv("SEPOLIA_RPC_URL"))

	// 连接到以太坊测试网节点
	client, err = ethclient.Dial(os.Getenv("SEPOLIA_RPC_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// 从环境变量中获取私钥并转换为ECDSA格式
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

	// 从公钥获取钱包地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取账户当前的nonce值（交易序号）
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	// 获取当前建议的gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 创建交易认证器，用于签名交易
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce)) // 设置nonce
	auth.Value = big.NewInt(0)            // 设置转账金额（这里为0）
	auth.GasLimit = uint64(300000)        // 设置gas限制
	auth.GasPrice = gasPrice              // 设置gas价格

	// 部署智能合约
	address, tx, instance, err := Counter.DeployContract(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	// 打印部署的合约地址和交易哈希
	fmt.Println("合约地址:", address.Hex())
	fmt.Println("交易哈希:", tx.Hash().Hex())

	_ = instance // 保留合约实例以供后续使用
}
