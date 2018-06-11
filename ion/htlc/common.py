import os
import time

ONE_MINUTE = 60
ONE_HOUR = ONE_MINUTE * 60
ONE_DAY = ONE_HOUR * 24
ONE_YEAR = ONE_DAY * 365

DEFAULT_EXPIRY_DURATION = 10 * ONE_MINUTE
MINIMUM_EXPIRY_DURATION = 2 * ONE_MINUTE
DURATION_OR_EPOCH_SPLIT = ONE_YEAR


def make_htlc_proxy(rpc, contract, account):
    # TODO: embed 'abi/HTLC.abi' file in package resources?
    return rpc.proxy('abi/HTLC.abi', contract, account)


def get_default_expiry():
    return int(time.time()) + DEFAULT_EXPIRY_DURATION


def get_random_secret_32():
    return '0x' + os.urandom(32).encode('hex')
