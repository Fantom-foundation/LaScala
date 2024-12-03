const {ethers, JsonRpcProvider, toQuantity, ZeroAddress, AbiCoder, getNumber} = require("ethers");
const fs = require("fs");
const NodeDriverAuthAddress = "0xd100ae0000000000000000000000000000000000";
const NodeDriverAuthAbi = [{
    "constant": false,
    "inputs": [{"internalType": "bytes", "name": "diff", "type": "bytes"}],
    "name": "updateNetworkRules",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "function"
}];
const PRIVATE_KEY = "163f5f0f9a621d72fedd85ffca3d08d131ab4e812181e0d30ffd1c885d20aac7"
const PUBLIC_ADDR = "0x239fA7623354eC26520dE878B52f13Fe84b06971"

function prepareUpdateNetworkRulesCall(rulesPath) {
    // load contents of rulesPath file and convert it to bytes
    const rules = fs.readFileSync(rulesPath);
    var rulesBytes = new Uint8Array(rules)
    var iface = new ethers.Interface(NodeDriverAuthAbi)
    var calldata = iface.encodeFunctionData("updateNetworkRules", [rulesBytes]);
    return calldata;
}

/**
 * @param {string} url
 * @param {string} rulesPath
 * @returns {Promise<void>}
 */
async function main(url, rulesPath) {
    return new Promise(async (resolve) => {
        const provider = new ethers.JsonRpcProvider(url);
        const signer = new ethers.Wallet(PRIVATE_KEY, provider)

        var calldata = prepareUpdateNetworkRulesCall(rulesPath);
        const tx = await signer.sendTransaction({
            to: NodeDriverAuthAddress,
            data: calldata,
        });
        console.log(tx);

        // confirm status of tx receipt
        const receipt = await provider.waitForTransaction(tx.hash);
        if (receipt.status !== 1) {
            console.log(receipt);
        } else {
            console.log("Transaction receipt confirmed");
        }

        console.log("Rules updated - wait for next epoch for rules to take an effect...");

        resolve(true);
    })
}

const argv = require("minimist")(process.argv.slice(2), {string: ["url", "rulesPath"]});
if (argv.url === undefined) {
    console.error("Missing required RPC endpoint URL address --url");
    process.exit(1);
}

if (argv.rulesPath === undefined) {
    console.error("Missing required new rulesPath --rulesPath");
    process.exit(1);
}

main(argv.url, argv.rulesPath).catch(error => {
    console.error(error);
}).then(() => console.log('done'));
