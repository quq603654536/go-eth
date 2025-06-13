package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	Counter "github.com/wkm/go-eth/task_1/task_2/contract"
)

// 声明全局变量，用于存储以太坊客户端连接
var client *ethclient.Client

// 测试网地址 https://sepolia.infura.io/v3/**
// 合约地址: 0x67d12A1a25B0Bf4e3D9D319c970A818D39dAB128
// 交易哈希: 0xe1d555e71bf9171b2da250aaaceb5123895ae2e159d1ab80addd3ab1c150d7a2
// 等待交易确认...
// 合约部署成功！
// 等待 Incr 交易确认...
// Incr 交易成功！
// 当前计数: 1

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
	address, tx, _, err := Counter.DeployContract(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	// 打印部署的合约地址和交易哈希
	fmt.Println("合约地址:", address.Hex())
	fmt.Println("交易哈希:", tx.Hash().Hex())

	// 等待交易被确认
	fmt.Println("等待交易确认...")
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status == 0 {
		log.Fatal("合约部署失败")
	}
	fmt.Println("合约部署成功！")

	doSoming(address.Hex(), auth)
}

func doSoming(address string, auth *bind.TransactOpts) {
	newAddress := eth_common.HexToAddress(address)
	instance, err := Counter.NewContract(newAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// 更新 nonce 值
	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	// 更新 gas 价格，使用更高的价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 将 gas 价格提高 50%
	auth.GasPrice = new(big.Int).Mul(gasPrice, big.NewInt(15))
	auth.GasPrice = new(big.Int).Div(auth.GasPrice, big.NewInt(10))

	// 调用合约的 Incr 方法
	tx, err := instance.Incr(auth)
	if err != nil {
		log.Fatal(err)
	}

	// 等待交易被确认
	fmt.Println("等待 Incr 交易确认...")
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status == 0 {
		log.Fatal("Incr 交易失败")
	}
	fmt.Println("Incr 交易成功！")

	// 创建调用选项
	callOpts := &bind.CallOpts{}

	// 调用合约的 GetCounter 方法
	var counter *big.Int
	counter, err = instance.GetCounter(callOpts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("当前计数:", counter.Int64())
}
