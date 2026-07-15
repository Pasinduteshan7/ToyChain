# ⛓ ToyChain

> A minimal blockchain simulator built from scratch in pure Go — no third-party SDKs, no networking, just the standard library.

ToyChain is an educational implementation of a blockchain that demonstrates core concepts like **SHA-256 hashing**, **Proof-of-Work mining**, **tamper detection**, and **ledger validation** — all in a single-process CLI application.

---

## ✨ Features

- 🧱 **Block Structure** — Height, timestamp, transactions, prev-hash link, nonce, and SHA-256 hash
- 🏁 **Genesis Block** — Deterministic starting block with all-zero prev-hash
- 🔐 **SHA-256 Hashing** — Deterministic fingerprinting with documented field order
- 💰 **Transaction Ledger** — Balance tracking with overdraft protection and faucet minting
- ⛏️ **Proof-of-Work Mining** — Brute-force nonce search for leading-zero hash targets
- 🛡️ **Tamper Detection** — 5-rule chain validation catches any post-mining edits
- 💾 **JSON Persistence** — Full chain state saved/loaded automatically via `chain.json`
- ⚙️ **Configurable Parameters** — Difficulty, block size, and data file via CLI flags
- 🖥️ **Interactive Menu** — Colored, menu-driven interface for demos and exploration
- 📊 **Experiment Mode** — Built-in difficulty-vs-effort sweep for research data

---

## 🚀 Getting Started

### Prerequisites
- **Go 1.22+** installed on your machine

### Build & Run

```bash
git clone https://github.com/Pasinduteshan8/ToyChain-.git
cd ToyChain-
go build -o toychain.exe .
```

---

## 🖥️ Interactive Mode

Run with no arguments to launch the interactive menu:

```bash
./toychain
```

```
╔══════════════════════════════════════════╗
║         ⛓  TOYCHAIN BLOCKCHAIN  ⛓        ║
║     A Minimal Blockchain Simulator       ║
╚══════════════════════════════════════════╝

  Blocks: 1    Difficulty: 3    Pending: 0

  1. Add Transaction
  2. Mine Block
  3. View Blockchain
  4. View Balances
  5. Validate Blockchain
  6. Run Experiment
  7. Exit

Choice: _
```

---

## 🛠️ Command-Line Mode

You can also use one-shot commands with flags:

```bash
toychain [-data path] [-difficulty N] [-maxblock N] <command> [args]
```

### Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-data` | `chain.json` | Path to the persisted chain file |
| `-difficulty` | `3` | PoW difficulty (leading hex zeros) — only for new chains |
| `-maxblock` | `0` | Max transactions per mined block (0 = unlimited) |

### Commands

| Command | Arguments | Description |
|---------|-----------|-------------|
| `add-tx` | `<sender> <recipient> <amount>` | Queue a transaction (use `-` as sender for faucet) |
| `mine` | — | Mine a block from the pending pool |
| `print` | — | Print every block in the chain |
| `validate` | — | Validate the chain and report tampering |
| `balances` | — | Show current account balances |
| `experiment` | `[maxDifficulty]` | Run difficulty-vs-effort sweep (default: 5) |

### Example Session

```bash
# Step 1: Mint 100 coins to Alice (faucet transaction)
./toychain add-tx - alice 100

# Step 2: Alice sends 30 to Bob
./toychain add-tx alice bob 30

# Step 3: Mine the pending transactions into a block
./toychain mine

# Step 4: View the chain, balances, and validate
./toychain print
./toychain balances      # alice: 70, bob: 30
./toychain validate      # chain is VALID

# Step 5: Run the difficulty experiment
./toychain experiment 6
```

---

## 📊 Difficulty vs. Effort Experiment

The built-in experiment command demonstrates exponential mining growth:

```bash
./toychain experiment 6
```

```
| Difficulty | Attempts (nonces tried) | Time (ms) | Hash prefix    |
|------------|-------------------------|-----------|----------------|
| 1          | 7                       | 0.000     | `07c1cc07ec…`  |
| 2          | 164                     | 0.000     | `0066805642…`  |
| 3          | 7,236                   | 8.170     | `000cfb346c…`  |
| 4          | 87,273                  | 77.011    | `000065d0c0…`  |
| 5          | 407,727                 | 352.408   | `00000b860d…`  |
| 6          | 1,956,130               | 1,720.075 | `0000000648…`  |
```

Each extra leading zero multiplies the search space by **16×** (one hex digit = 4 bits).

---

## 🏗️ Architecture

```
toychain/
├── main.go                 # Entry point — delegates to cli.Run()
├── block/
│   ├── block.go            # Block struct, SHA-256 hashing, genesis block
│   └── block_test.go       # Hash determinism, genesis correctness
├── chain/
│   ├── chain.go            # Blockchain type, mining loop, mempool
│   ├── validate.go         # 5-rule chain validation (tamper detection)
│   ├── persist.go          # JSON save/load (chain.json)
│   ├── chain_test.go       # Mining, chain linking tests
│   ├── validate_test.go    # Tamper detection, honest chain tests
│   └── persist_test.go     # Save/load round-trip tests
├── ledger/
│   ├── transaction.go      # Transaction struct (Sender, Recipient, Amount)
│   ├── balances.go         # Balance tracking, overdraft protection
│   └── balances_test.go    # Overspend rejection, faucet tests
├── cli/
│   ├── cli.go              # Command-line argument parsing & dispatch
│   └── interactive.go      # Interactive menu-based CLI
├── report.md               # Research report template
└── demo.json               # Sample chain data
```

---

## 🧠 Design Decisions

- **Hashing:** SHA-256 over `Height | Timestamp | Transactions(JSON) | PrevHash | Nonce` with `|` separators. The `Hash` field itself is excluded (circular dependency). Manual concatenation makes the hash input unambiguous.

- **Derived Balances:** Balances are replayed from genesis — no separate balance table. This guarantees balances can never drift out of sync with the chain.

- **Double Validation:** Transactions are checked against projected balances when added to the mempool. Chain-level `Validate()` independently checks hash/link/PoW integrity for tamper detection.

- **Tamper Detection:** Two independent checks: (1) recomputed hash must match stored hash, (2) each block's `PrevHash` must match the previous block's actual hash — catching even re-mined altered blocks.

---

## 🛡️ How Validation Works (The Immune System)

Validation is the core of blockchain security. The `Validate()` function runs **5 strict checks** on every single block from Genesis to the tip. If a block fails even one check, the chain is rejected.

1. **Hash Re-computation (Content Tampering)**
   * **The Method:** Recalculates the SHA-256 hash of the block's data from scratch.
   * **What it checks:** Compares the new hash to the saved hash. If they don't match, someone edited the block's contents (e.g., changing a transaction) after it was mined.

2. **Previous Hash Link (Chain Continuity)**
   * **The Method:** Checks if Block N's `PrevHash` exactly matches Block N-1's actual `Hash`.
   * **What it checks:** Ensures the blocks are chained together properly. If a hacker tampers with an old block, its hash changes, breaking the link to the next block.

3. **Proof-of-Work Compliance (Mining Legitimacy)**
   * **The Method:** Checks if the block's hash starts with the required number of zeros.
   * **What it checks:** Proves the miner actually did the heavy mathematical lifting to find a valid Nonce, preventing hackers from instantly fabricating fake blocks.

4. **Structural Integrity (Order & Time)**
   * **The Method:** Checks the `Height` and `Timestamp`.
   * **What it checks:** Ensures heights increase by exactly 1 (no gaps or duplicates) and timestamps only move forward (no time travel).

5. **Ledger Replay (Economic Rules)**
   * **The Method:** Replays every transaction on a blank ledger from the beginning.
   * **What it checks:** Ensures the "laws of physics" of money are obeyed. Catches illegal transactions (e.g., negative amounts or overdrafts) even if the block's hash is mathematically perfect.

---

## ⚠️ Known Limitations

| Limitation | Description |
|------------|-------------|
| No Networking | Single-process simulator — no peer consensus |
| No Digital Signatures | Anyone can construct a transaction from any account |
| No Merkle Tree | Transactions hashed as flat JSON, not a tree with proof paths |
| Basic Persistence | Single JSON file, not designed for concurrent access |

---

## 🧪 Testing

```bash
go test ./... -v
```

Tests cover: hash determinism, genesis correctness, PoW target satisfaction, chain linking, honest-chain validation, tamper detection, overspend rejection, mempool limits, and save/load round-tripping.

---

## 📜 License

This project is for educational purposes as part of a Backend Engineering Internship assessment.
