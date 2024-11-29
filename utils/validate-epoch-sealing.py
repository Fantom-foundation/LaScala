#!/usr/bin/python
import asyncio
import math
import sys

from web3 import Web3

zero_addr = '0x0000000000000000000000000000000000000000'
cache = {}
latest_block = {}


def load_epoch(web3, blk):
    try:
        block = hex(blk)
        if block in cache:
            return cache[block]
        j = web3.eth.get_block(block)
        if j is None:
            raise Exception('ERROR block: ', int(block, 16),
                            " doesn't exists.")
        r = int(j['epoch'], 16)
        cache[block] = r
        return r
    except Exception as e:
        raise Exception(f"Error: {e}")


def find_first(web3, epoch, min, max):
    last_low = -1
    last_high = -1

    while True:
        if max - min <= 1:
            if epoch == load_epoch(web3, min):
                return min, last_low, last_high
            elif epoch == load_epoch(web3, max):
                return max, last_low, last_high

        # choosing the bigger block in middle of the interval
        half = math.ceil((max - min) / 2) + min
        e = load_epoch(web3, half)
        if e < epoch:
            min = half
        else:
            # labeling the boundaries of last block
            # first time when the epoch is correct then the last block is in between half and max
            if e == epoch and last_low == -1:
                last_low = half
                last_high = max
            if max == half:
                raise Exception('ERROR epoch: ', epoch, " not found")
            max = half


def find_last(web3, epoch, min, max):
    while True:
        # choosing the smaller block in the middle of the interval
        half = math.floor((max - min) / 2) + min
        e = load_epoch(web3, half)
        if max - min == 0:
            return max
        if max - min == 1:
            if epoch != e:
                raise Exception(e, ' != ', epoch)
            if e + 1 != load_epoch(web3, half + 1):
                return half + 1
            return half
        if e <= epoch:
            min = half
        else:
            max = half


def load_block_range(web3, epoch):
    max = latest_block['result']['number']
    current_epoch = int(latest_block['result']['epoch'], 16)
    if current_epoch < epoch:
        raise Exception('ERROR epoch: ', epoch,
                        " isn't finished. Lastest finished epoch is: ",
                        current_epoch - 1)
    r = find_first(web3, epoch, 0, max)
    last_min_bound = r[1]
    last_max_bound = r[2]
    if r[1] == -1 or r[2] == -1:
        last_min_bound = r[0]
        last_max_bound = max
    l = find_last(web3, epoch, last_min_bound, last_max_bound)
    return [r[0], l]


def get_zero_addr_nonce_at_block(web3, i):
    return web3.eth.get_transaction_count(zero_addr, block_identifier=i)


def validate_epoch_seal_calls_locations(web3, first_block, last_block):
    # get reference nonce from first block
    expected_nonce = get_zero_addr_nonce_at_block(web3, first_block)
    i = load_epoch(web3, first_block)
    last_epoch = load_epoch(web3, last_block)

    if i == last_epoch:
        pass

    # check if the first block isn't the epoch seal block
    rf = load_block_range(web3, i)
    # check if the last block isn't the epoch seal block
    rl = load_block_range(web3, last_epoch)
    if latest_block['result']['number'] != rl[1]:
        # last epoch can be checked only if we know for sure that the last block is really epoch seal
        # this needs to be done because running node doesn't know
        # if the latest block is or isn't epoch seal from load_block_range call
        last_epoch += 1
    if first_block == rf[1]:
        # the first checked block is epoch seal block so expected nonce is decremented by 2
        # (state before SealEpochStats and SealEpochValidators calls)
        expected_nonce -= 2
        print(
            f"Block {first_block} is epoch seal block, so block {first_block - 1} needs to be checked for nonce as well")

    #  check nonces of all epochs within the range
    while i < last_epoch:
        r = load_block_range(web3, i)
        ln = get_zero_addr_nonce_at_block(web3, r[1])

        # check nonce is correct before sealing epoch
        before_seal_nonce = get_zero_addr_nonce_at_block(web3, r[1] - 1)
        if before_seal_nonce != expected_nonce and r[1] - 1 != 0:
            raise Exception(f"Error: block {r[1] - 1} nonce expected: {expected_nonce} got: {before_seal_nonce}")

        # check nonce is correct in block that is sealing the epoch (SealEpochStats and SealEpochValidators were called)
        expected_nonce += 2
        if ln != expected_nonce:
            raise Exception(f"Error: block {r[1]} nonce expected: {expected_nonce} got: {ln}")
        if i % 1 == 0:
            print(
                f"Progress: {round((r[1] - first_block) / (last_block - first_block) * 100, 2)}% - epoch {i} - block {r[1]}")
        i += 1


async def main():
    fb = int(sys.argv[1])
    lb = sys.argv[2]
    url = str(sys.argv[3])

    web3 = Web3(Web3.HTTPProvider(url))
    latest_block['result'] = web3.eth.get_block('latest')

    if lb == "last":
        lb = latest_block['result']['number']

    print("Starting from block: ", fb, " to block: ", lb)
    validate_epoch_seal_calls_locations(web3, fb, int(lb))

    print("Finished successfully")


if __name__ == '__main__':
    asyncio.run(main())
