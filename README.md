# go-eth
go-eth 开发

# abigen工具安装
sudo apt update
sudo apt install golang-go

在 ~/.bashrc添加
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
source ~/.bashrc

go install github.com/ethereum/go-ethereum/cmd/abigen@latest

abigen --version

## 命令solc
solc --bin Counter.sol -o .
solc --abi Counter.sol -o .

abigen --abi=Counter_sol_Counter.abi --pkg=counter --out=Counter.go
abigen --bin=Counter.bin --abi=Counter_sol_Counter.abi --pkg=counter --out=Counter_2.go