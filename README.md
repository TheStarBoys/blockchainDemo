# blockchain Demo
> The blockchain demo by TheStarBoys

## Website
http://www.blockchaindemo.top/

## Block
Block shows the general appearance of a block.
1. `BLOCK`'s `HASH`, `Timestamp` and `Nonce` are changed by change of `BLOCK`'s `DATA`.

## Blockchain
Blockchain shows the connection process between blocks.
1. Genesis block's `DATA` read-only.
2. Clicking `ADD BLOCK` will directly dig out a block.
3. To try clicking `PREVIOUS HASH`, you'll turn to previous block.
4. If you change DATA of the block excluded `GENESIS BLOCK`, it's `HASH` will be change also and `BLOCK` which is after it will turn red to warn you.

## Transactions
Transactions presents the concept of `UTXO` in `Bitcoin`.
1. `BLOCK` shows the details of the transactions.
2. The table of accounts' balance put on the left of the page.
3. Refusing the illegal transaction.
4. Clicking one of `TXID`s of the `Transaction`'s Input, the page will turn to where is the `Transaction` that corresponds to the `TXID`.

## Merkle tree
Merkle tree demonstrates an important data structure in `Bitcoin` that can be used for simple payment verification scenarios.
1. If you change `TX`'s `DATA`, `HASH` that corresponds to the `TX` will be change, which will turn red to warn you. In the end, `Merkel Root` will be change also.

## Address
Address shows how an address in `Bitcoin` is obtained.
1. Clicking `GENERATE` will generate a address of Suiting `Bitcoin`'s norm.