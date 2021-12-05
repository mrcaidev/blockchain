package utils

import "math"

// 挖矿的奖励。
const Subsidy = 10

// 钱包集数据库。
const WalletsFile = "wallets.dat"

// 当前版本号。
const Version = byte(0x00)

// 难度系数。
const Difficulty = 24
const MaxNonce = math.MaxInt64

// 数据库路径。
const DBFile = "blockchain.db"
const BlocksBucket = "blocks"

const GenesisCoinbase = "Genesis Coinbase"
