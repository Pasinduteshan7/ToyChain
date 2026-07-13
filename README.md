
```markdown
# 🔗 ToyChain

> A minimal blockchain and ledger, built from scratch in pure Go (standard library only — no third-party chain SDKs, no networking).

ToyChain is an educational implementation of a blockchain, designed to demonstrate core concepts like Proof-of-Work (PoW), tamper detection, and ledger validation in a single-process simulator.

---

## 🚀 Getting Started

### Prerequisites
- **Go 1.22+** installed on your machine.

### Build & Run
Clone the repository and build the binary:

```bash
go build -o toychain .

```

---

## 🛠️ Usage & Commands

Global flags must go **before** the command:

```bash
toychain [-data path] [-difficulty N] [-maxblock N] <command> [args]

```

### Global Flags

| Flag | Default | Description |
| --- | --- | --- |
| `-data` | `chain.json` | Path to the persisted chain file. |
| `-difficulty` | `3` | PoW difficulty (leading hex zeros) — *only used when creating a new chain*. |
| `-maxblock` | `0` | Max transactions per mined block (0 = unlimited). |

### Commands

| Command | Arguments | Description |
| --- | --- | --- |
| `add-tx` | `<sender> <recipient> <amount>` | Queue a transaction. Use `-` as sender for a coinbase/faucet mint. |
| `mine` | — | Mine a block from the pending transaction pool. |
| `print` | — | Print every block in the chain to the console. |
| `validate` | — | Validate the whole chain, reporting the first bad block if any. |
| `balances` | — | Show current balances for every account seen. |

### Example Session

```bash
# Mint 100 coins to Alice
./toychain -data demo.json -difficulty 3 add-tx - alice 100 

# Alice sends 30 coins to Bob
./toychain -data demo.json add-tx alice bob 30

# Mine the pending transactions into a block
./toychain -data demo.json mine

# View the chain and verify balances
./toychain -data demo.json print
./toychain -data demo.json balances
./toychain -data demo.json validate

```

---

## 🧠 Design Decisions

* **Hashing:** A block's hash is SHA-256 over `Height | Timestamp | Transactions(JSON) | PrevHash | Nonce`, concatenated with `|` separators in that exact order (excluding the `Hash` field itself). Fields are concatenated manually rather than `json.Marshal`-ing the whole struct to make "which fields, what order" unambiguous.
* **Difficulty:** Defined as the required number of leading hex zeros in the hash. Mining brute-forces the nonce until the target is met.
* **Derived Balances:** `Blockchain.Balances()` replays every transaction in every mined block from genesis forward. There is no separate persisted balance table—this guarantees balances can never drift out of sync with the chain itself.
* **Double Validation:** A transaction is checked against projected balances when added to the pending pool (`AddTransaction`). Chain-level `Validate()` then re-checks hash/link/PoW integrity after the fact. `Validate()`'s job is strictly "has this data been tampered with," not "was this a legal transaction."
* **Tamper Detection:** Relies on two independent checks:
1. Recomputing a block's hash from its current fields must match the stored hash.
2. Each block's `PrevHash` must match the previous block's actual hash. (This catches attackers who re-mine an altered block to make it internally self-consistent).



---

## ⚠️ Known Limitations

* **No Networking:** This is explicitly a single-process simulator. There is no peer consensus.
* **No Digital Signatures:** Anyone can construct a transaction "from" any account. (Considered a future stretch goal).
* **No Merkle Tree:** A block's transaction list is hashed as a flat JSON array rather than summarized via a Merkle root.
* **Basic Persistence:** Uses a single flat JSON file, which is not designed for concurrent access from multiple processes.

---

## 🧪 Testing

Run the test suite with verbose output:

```bash
go test ./... -v

```

*Tests cover: hash determinism, genesis correctness, PoW target satisfaction, chain linking, honest-chain validation, tamper detection, overspend rejection, mempool limits, and save/load persistence round-tripping.*