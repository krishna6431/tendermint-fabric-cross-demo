# Demo

This directory is used for the playground to run the demo and several temporary directories are generated after the demo is executed by scripts.

## Directory structure

This directory consists of the following files/directories

- chains
  - fabric ... fabric-related files including docker-compose.yaml for network
    - chaincode
      - fabibc ... fabric chaincode application
  - tendermint ... tendermint simapp working as CLI and tendermint application including docker-compose.yaml for network
- configs ... config files for alpha cli, beta cli, fabric app, relayer
- scripts ... shell scripts for alpha cli, beta cli, relayer, scenario

## Setup Network

```Makefile
# install fabric tools and build Fabric/Tendermint CLI, Relayer
make build -j5

# down related containers, remove volumes
make network-down

# remove any generated data
make clean

# prepare fabric network and tendermint network using docker containers
make network
```

## About sample-scenario

Execute the Cross-chain swap as below

- Alice transfers 10 tokens to Bob on Tendermint α
- Bob transfers 10 tokens to Alice on Fabric β

#### Initial Balance

|       | Tendermint α | Fabric β |
| ----- | ------------ | -------- |
| Alice | 10           | 0        |
| Bob   | 0            | 10       |

#### Expected Result of Balance after executed

|       | Tendermint α | Fabric β |
| ----- | ------------ | -------- |
| Alice | 0            | 10       |
| Bob   | 10           | 1        |

### Actor

| Actor                  | Tendermint α | Fabric β             |
| ---------------------- | ------------ | -------------------- |
| Alice                  | account      | AliceMSP member      |
| Bob                    | account      | BobMSP member        |
| TokenOwner for α Chain | account      | -                    |
| TokenOwner for β Chain | -            | TokenOwnerMSP member |

### Flow and Commands in [sample-scenario](https://github.com/datachainlab/fabric-tendermint-cross-demo/blob/main/demo/scripts/scenario/sample-scenario)

1. Run `erc20 mint` command for initial balance

```bash
infoln "Mint token to Alice by TokenOwner on the Alpha Chain"
printAlpha "${ALPHACLI_TOKEN_OWNER} erc20 mint \
--receiver-address ${ALICE_ALPHA_ID} \
--amount ${AMOUNT}"

infoln "Mint token to Bob by TokenOwner on the Beta Chain"
printBeta "${BETACLI_TOKEN_OWNER} erc20 mint \
--receiver-id ${BOB_BETA_ID} \
--amount ${AMOUNT}"
```

2. Confirm balance

```bash
printAlpha "${ALPHACLI_ALICE} erc20 balance-of --owner-address ${ALICE_ALPHA_ID}"
alice_alpha_amount=${LATEST_RESULT}
expected_alice_alpha=$((${alice_alpha_amount} - ${AMOUNT}))
printAlpha "${ALPHACLI_BOB} erc20 balance-of --owner-address ${BOB_ALPHA_ID}"
bob_alpha_amount=${LATEST_RESULT}
expected_bob_alpha=$((${bob_alpha_amount} + ${AMOUNT}))
printBeta "${BETACLI_ALICE} erc20 balance-of --owner-id ${ALICE_BETA_ID}"
alice_beta_amount=${LATEST_RESULT}
expected_alice_beta=$((${alice_beta_amount} + ${AMOUNT}))
printBeta "${BETACLI_BOB} erc20 balance-of --owner-id ${BOB_BETA_ID}"
bob_beta_amount=${LATEST_RESULT}
expected_bob_beta=$((${bob_beta_amount} - ${AMOUNT}))
```

3. Create 2 `ContractTx`s including 2 function calls. See `call_info` variable.

```bash
infoln "Create a ContractTx by TokenOwner on the Off Chain"
call_info="{\"method\":\"transfer\",\"args\":[\"${BOB_ALPHA_ID}\",\"${AMOUNT}\"]}"
printAlpha "${ALPHACLI_BINARY} --home ${ALPHA_DATA} cross create-contract-tx
--signer-address ${ALICE_ALPHA_ID} \
--call-info $call_info \
--output-document ${DATA_DIR}/alpha-tx-1.json"

call_info="{\"method\":\"transfer\",\"args\":[\"${ALICE_BETA_ID}\",\"${AMOUNT}\"]}"
chan="channel-0:cross"
printBeta "${BETACLI_TOKEN_OWNER} cross create-contract-tx
--signer-id ${BOB_BETA_ID} \
--call-info $call_info \
--initiator-chain-channel $chan \
--output-document ${DATA_DIR}/beta-tx-1.json"
```

4. Create `InitiateTx` with created `ContractTxs` and send it on Tendermint α Chain. Then `transaction ID` can be returned.

```bash
infoln "Create and send an InitiateTx to start Cross-chain tx by Alice on the Alpha Chain"
printAlpha "${ALPHACLI_ALICE} cross create-initiate-tx --contract-txs=${DATA_DIR}/alpha-tx-1.json,${DATA_DIR}/beta-tx-1.json"
tx_id=$(echo ${LATEST_RESULT} | jq -r '.logs[0].events[1].attributes[0].value')
```

5. Create `IBCSignTx` with `transaction ID` and send it on Fabric β Chain.

```bash
infoln "Send IBCSignTx for the InitiateTx by Bob on the Beta Chain"
printBeta "${BETACLI_BOB} cross ibc-signtx
--tx-id ${tx_id} \
--initiator-chain-channel ${chan}"
```

6. Relay packets on the given path in both direcions for Authentication.

```bash
printRelay "${RLY} tx relay ${PATH_NAME}"
sleep 5
```

7. Relay acknowledgments on the given path in both directions for Authentication.

```bash
printRelay "${RLY} tx relay-acknowledgements ${PATH_NAME}"
sleep 5
```

8. Authentication status `remaining_signers` as array is returned by the below command.

```bash
printAlpha "${ALPHACLI_BINARY} --home ${ALPHA_DATA} cross tx-auth-state $tx_id"
```

9. Status at Coordinator chain is returned by the below command.

```bash
printAlpha "${ALPHACLI_BINARY} --home ${ALPHA_DATA} cross coordinator-state $tx_id"
```

10. Make sure that balance is not changed yet in this time

```bash
printAlpha "${ALPHACLI_ALICE} erc20 balance-of --owner-address ${ALICE_ALPHA_ID}"
actual_alice_alpha=${LATEST_RESULT}
printAlpha "${ALPHACLI_BOB} erc20 balance-of --owner-address ${BOB_ALPHA_ID}"
actual_bob_alpha=${LATEST_RESULT}
assertEqual $actual_alice_alpha $alice_alpha_amount "Alice has still $actual_alice_alpha token during Prepare phase."
assertEqual $actual_bob_alpha $bob_alpha_amount "Bob has still $actual_bob_alpha token during Prepare phase."
```

11. Relay packets and acknowledgements on the given path in both directions until ends.

```bash
infoln "Relay packets"
set +e
for i in {1..5}; do
  printRelay "${RLY} tx relay ${PATH_NAME}"
  printRelay "${RLY} tx relay-acknowledgements ${PATH_NAME}"
  printRelay "${RLY} query unrelayed-packets ${PATH_NAME}"
  packet=${LATEST_RESULT}
  printRelay "${RLY} query unrelayed-acknowledgements ${PATH_NAME}"
  ack=${LATEST_RESULT}
  printAlpha "${ALPHACLI_BINARY} --home ${ALPHA_DATA} cross tx-auth-state $tx_id"
  printAlpha "${ALPHACLI_BINARY} --home ${ALPHA_DATA} cross coordinator-state $tx_id"
  if [ "$packet" = '{"src":[],"dst":[]}' ] && [ "$ack" = '{"src":[],"dst":[]}' ]; then
    break
  fi
  sleep 5
done
```

12. Get current account's balance

```bash
infoln "Assert result"
printAlpha "${ALPHACLI_ALICE} erc20 balance-of --owner-address ${ALICE_ALPHA_ID}"
actual_alice_alpha=${LATEST_RESULT}
printAlpha "${ALPHACLI_BOB} erc20 balance-of --owner-address ${BOB_ALPHA_ID}"
actual_bob_alpha=${LATEST_RESULT}
printBeta "${BETACLI_ALICE} erc20 balance-of --owner-id ${ALICE_BETA_ID}"
actual_alice_beta=${LATEST_RESULT}
printBeta "${BETACLI_BOB} erc20 balance-of --owner-id ${BOB_BETA_ID}"
actual_bob_beta=${LATEST_RESULT}
```

13. Make sure that current status of coordinator chain, commit status, expected balance.

```bash
printAlpha "${ALPHACLI_BINARY} --home ${ALPHA_DATA} cross coordinator-state $tx_id"
tx_res=${LATEST_RESULT}
isCommit "$tx_res" "$tx_id"
assertEqual $actual_alice_alpha $expected_alice_alpha \
"Alice balance on the Alpha from $alice_alpha_amount to $actual_alice_alpha"
assertEqual $actual_bob_alpha $expected_bob_alpha \
"Bob balance on the Alpha from $bob_alpha_amount to $actual_bob_alpha"
assertEqual $actual_alice_beta $expected_alice_beta \
"Alice balance on the Beta from $alice_beta_amount to $actual_alice_beta"
assertEqual $actual_bob_beta $expected_bob_beta \
"Bob balance on the Beta from $bob_beta_amount to $actual_bob_beta"
```

### Run sample scenario

```Makefile
# initialize relayer, fabric CLI, tendermint CLI, and run handshake for IBC between fabric and tendermint by creating transactions.
make run-init

# run ./scripts/scenario/sample-scenario. See `About sample-scenario` section for more detail.
make run
```
