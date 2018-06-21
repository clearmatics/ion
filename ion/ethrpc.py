"""
This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <http://unlicense.org/>
"""

import json
import requests
import time
import warnings
from binascii import hexlify, unhexlify
from io import IOBase

from collections import namedtuple
from ethereum.abi import encode_abi, decode_abi
from requests.adapters import HTTPAdapter
from requests.exceptions import ConnectionError as RequestsConnectionError

from .crypto import keccak_256
from .utils import CustomJSONEncoder, require, big_endian_to_int, zpad, encode_int, normalise_address

GETH_DEFAULT_RPC_PORT = 8545
ETH_DEFAULT_RPC_PORT = 8545
PARITY_DEFAULT_RPC_PORT = 8545
PYETHAPP_DEFAULT_RPC_PORT = 4000
MAX_RETRIES = 3
JSON_MEDIA_TYPE = 'application/json'


BLOCK_TAG_EARLIEST = 'earliest'
BLOCK_TAG_LATEST   = 'latest'
BLOCK_TAG_PENDING  = 'pending'
BLOCK_TAGS = (
    BLOCK_TAG_EARLIEST,
    BLOCK_TAG_LATEST,
    BLOCK_TAG_PENDING,
)


class EthJsonRpcError(Exception):
    pass


class ConnectionError(EthJsonRpcError):
    pass


class BadStatusCodeError(EthJsonRpcError):
    pass


class BadJsonError(EthJsonRpcError):
    pass


class BadResponseError(EthJsonRpcError):
    pass


def hex_to_dec(x):
    '''
    Convert hex to decimal
    '''
    return int(x, 16)


def clean_hex(d):
    '''
    Convert decimal to hex and remove the "L" suffix that is appended to large
    numbers
    '''
    return hex(d).rstrip('L')

def validate_block(block):
    if isinstance(block, str):
        if block not in BLOCK_TAGS:
            raise ValueError('invalid block tag')
    if isinstance(block, int):
        block = hex(block)
    return block


def wei_to_ether(wei):
    '''
    Convert wei to ether
    '''
    return 1.0 * wei / 10**18


def ether_to_wei(ether):
    '''
    Convert ether to wei
    '''
    return ether * 10**18


class EthTransaction(namedtuple('_TxStruct', ('rpc', 'txid'))):
    def details(self):
        txid = self.txid
        if txid[:2] != '0x':
            txid = '0x' + txid
        return self.rpc.eth_getTransactionByHash(txid)

    def wait(self):
        return self.receipt(wait=True)

    def receipt(self, wait=False, tick_fn=None):
        # TODO: add `timeout` param
        first = True
        txid = self.txid
        if txid[:2] != '0x':
            txid = '0x' + txid
        while True:
            receipt = self.rpc.eth_getTransactionReceipt(txid)
            # TODO: turn into asynchronous notification / future
            if receipt:
                return receipt
            if not wait:
                break
            try:
                if first:
                    if hasattr(wait, '__call__'):
                        wait()
                    first = False
                elif tick_fn:
                    tick_fn(self)
                time.sleep(1)
            except KeyboardInterrupt:
                break

    def __str__(self):
        return self.txid


class EthJsonRpc(object):
    '''
    Ethereum JSON-RPC client class
    '''

    DEFAULT_GAS_PER_TX = 900000
    DEFAULT_GAS_PRICE = 50 * 10**9  # 50 gwei

    def __init__(self, host='localhost', port=GETH_DEFAULT_RPC_PORT, tls=False):
        self.host = host
        self.port = port
        self.tls = tls
        self.session = requests.Session()
        self.session.mount(self.host, HTTPAdapter(max_retries=MAX_RETRIES))

    def _call(self, method, params=None, _id=1):

        params = params or []
        data = {
            'jsonrpc': '2.0',
            'method':  method,
            'params':  params,
            'id':      _id,
        }
        scheme = 'http'
        if self.tls:
            scheme += 's'
        url = '{}://{}:{}'.format(scheme, self.host, self.port)
        headers = {'Content-Type': JSON_MEDIA_TYPE}
        try:
            encoded_data = json.dumps(data, cls=CustomJSONEncoder) 
            r = self.session.post(url, headers=headers, data=encoded_data)
        except RequestsConnectionError:
            raise ConnectionError(url)
        if r.status_code / 100 != 2:
            raise BadStatusCodeError(r.status_code)
        try:
            response = r.json()
        except ValueError:
            raise BadJsonError(r.text)
        try:
            return response['result']
        except KeyError:
            raise BadResponseError(response)

    def _encode_function(self, signature, param_values):

        prefix = big_endian_to_int(keccak_256(signature.encode('utf-8')).digest()[:4])

        if signature.find('(') == -1:
            raise RuntimeError('Invalid function signature. Missing "(" and/or ")"...')

        if signature.find(')') - signature.find('(') == 1:
            return encode_int(prefix)

        types = signature[signature.find('(') + 1: signature.find(')')].split(',')
        encoded_params = encode_abi(types, param_values)
        return zpad(encode_int(prefix), 4) + encoded_params

    def _solproxy_bind(self, method, address, account):
        ins = [_['type'] for _ in method['inputs']]
        outs = [_['type'] for _ in method['outputs']]
        sig = method['name'] + '(' + ','.join(ins) + ')'

        if method['constant']:
            # XXX: document len(outs) and different behaviour...
            if len(outs) > 1:
                return lambda *args, **kwa: self.call(address, sig, args, outs, **kwa)
            return lambda *args, **kwa: self.call(address, sig, args, outs, **kwa)[0]
        if account is None:
            return None
        return lambda *args, **kwa: self.call_with_transaction(account, address, sig, args, **kwa)

    def proxy(self, abi, address, account=None):
        """
        Provides a Python proxy object which exposes the contract ABI as 
        callable methods, allowing for seamless use of contracts from Python... 
        """
        # XXX: specific to Ethereum addresses, 20 octets
        address = normalise_address(address)

        if account is not None:
            account = normalise_address(account)

        if isinstance(abi, IOBase):
            abi = json.load(abi)
        elif isinstance(abi, str):
            with open(abi) as jsonfile:
                abi = json.load(jsonfile)
        require(isinstance(abi, list))

        proxy = dict()
        for method in abi:
            if method['type'] != 'function':
                continue

            handler = self._solproxy_bind(method, address, account)
            if handler is None:
                continue

            sig = "%s(%s)" % (method['name'], ','.join([i['type'] for i in method['inputs']]))
            sig_hash = keccak_256(sig.encode('utf-8')).hexdigest()[:8]

            # Provide an alternate, where the explicit function signature
            proxy[method['name']] = handler
            proxy[method['name'] + '_' + sig_hash] = handler

        return namedtuple('SolProxy', proxy.keys())(*proxy.values())

################################################################################
# high-level methods
################################################################################

    def receipt(self, txid, wait=False, raise_on_error=False):
        # TODO: add `timeout` param
        transaction = EthTransaction(self, txid)
        receipt = transaction.receipt(wait=wait)
        if raise_on_error:
            if int(receipt['status'], 16) == 0:
                raise EthJsonRpcError("Transaction was aborted")
        return receipt

    def receipt_wait(self, txid, raise_on_error=True):
        """
        Wait for the transaction to be mined, then return receipt
        """
        return self.receipt(txid, raise_on_error)

    def transfer(self, from_, to, amount):
        '''
        Send wei from one address to another
        '''
        return self.eth_sendTransaction(from_address=from_, to_address=to, value=amount)

    def create_contract(self, from_, code, gas, sig=None, args=None):
        '''
        Create a contract on the blockchain from compiled EVM code. Returns the
        transaction hash.
        '''
        from_ = from_ or self.eth_coinbase()
        if sig is not None and args is not None:
             types = sig[sig.find('(') + 1: sig.find(')')].split(',')
             encoded_params = encode_abi(types, args)
             code += hexlify(encoded_params)
        return self.eth_sendTransaction(from_address=from_, gas=gas, data=code)

    def get_contract_address(self, tx):
        '''
        Get the address for a contract from the transaction that created it
        '''
        receipt = self.eth_getTransactionReceipt(tx)
        return receipt['contractAddress']

    def call(self, address, sig, args, result_types):
        '''
        Call a contract function on the RPC server, without sending a
        transaction (useful for reading data)
        '''
        data = self._encode_function(sig, args)
        data_hex = hexlify(data)
        response = self.eth_call(to_address=address, data=data_hex)
        # XXX: horrible hack for when RPC returns '0x0'...
        if (len(result_types) == 0 or result_types[0] == 'uint256') and response == '0x0':
            response = '0x' + ('0' * 64)
        return decode_abi(result_types, unhexlify(response[2:]))

    def call_with_transaction(self, from_, address, sig, args, gas=None, gas_price=None, value=None):
        '''
        Call a contract function by sending a transaction (useful for storing
        data)
        '''
        gas = gas or self.DEFAULT_GAS_PER_TX
        gas_price = gas_price or self.DEFAULT_GAS_PRICE
        data = self._encode_function(sig, args)
        data_hex = hexlify(data)
        return self.eth_sendTransaction(from_address=from_, to_address=address, data=data_hex, gas=gas,
                                        gas_price=gas_price, value=value)

################################################################################
# JSON-RPC methods
################################################################################

    def web3_clientVersion(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_clientversion

        TESTED
        '''
        return self._call('web3_clientVersion')

    def web3_sha3(self, data):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#web3_sha3

        TESTED
        '''
        data = hexlify(str(data))
        return self._call('web3_sha3', [data])

    def net_version(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#net_version

        TESTED
        '''
        return self._call('net_version')

    def net_listening(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#net_listening

        TESTED
        '''
        return self._call('net_listening')

    def net_peerCount(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#net_peercount

        TESTED
        '''
        return hex_to_dec(self._call('net_peerCount'))

    def eth_protocolVersion(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_protocolversion

        TESTED
        '''
        return self._call('eth_protocolVersion')

    def eth_syncing(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_syncing

        TESTED
        '''
        return self._call('eth_syncing')

    def eth_coinbase(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_coinbase

        TESTED
        '''
        return self._call('eth_coinbase')

    def eth_mining(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_mining

        TESTED
        '''
        return self._call('eth_mining')

    def eth_hashrate(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_hashrate

        TESTED
        '''
        return hex_to_dec(self._call('eth_hashrate'))

    def eth_gasPrice(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gasprice

        TESTED
        '''
        return hex_to_dec(self._call('eth_gasPrice'))

    def eth_accounts(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_accounts

        TESTED
        '''
        return self._call('eth_accounts')

    def eth_blockNumber(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_blocknumber

        TESTED
        '''
        return hex_to_dec(self._call('eth_blockNumber'))

    def eth_getBalance(self, address=None, block=BLOCK_TAG_LATEST):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getbalance

        TESTED
        '''
        address = address or self.eth_coinbase()
        block = validate_block(block)
        return hex_to_dec(self._call('eth_getBalance', [address, block]))

    def eth_getStorageAt(self, address=None, position=0, block=BLOCK_TAG_LATEST):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getstorageat

        TESTED
        '''
        block = validate_block(block)
        return self._call('eth_getStorageAt', [address, hex(position), block])

    def eth_getTransactionCount(self, address, block=BLOCK_TAG_LATEST):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactioncount

        TESTED
        '''
        block = validate_block(block)
        return hex_to_dec(self._call('eth_getTransactionCount', [address, block]))

    def eth_getBlockTransactionCountByHash(self, block_hash):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblocktransactioncountbyhash

        TESTED
        '''
        return hex_to_dec(self._call('eth_getBlockTransactionCountByHash', [block_hash]))

    def eth_getBlockTransactionCountByNumber(self, block=BLOCK_TAG_LATEST):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblocktransactioncountbynumber

        TESTED
        '''
        block = validate_block(block)
        return hex_to_dec(self._call('eth_getBlockTransactionCountByNumber', [block]))

    def eth_getUncleCountByBlockHash(self, block_hash):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclecountbyblockhash

        TESTED
        '''
        return hex_to_dec(self._call('eth_getUncleCountByBlockHash', [block_hash]))

    def eth_getUncleCountByBlockNumber(self, block=BLOCK_TAG_LATEST):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclecountbyblocknumber

        TESTED
        '''
        block = validate_block(block)
        return hex_to_dec(self._call('eth_getUncleCountByBlockNumber', [block]))

    def eth_getCode(self, address, default_block=BLOCK_TAG_LATEST):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getcode

        NEEDS TESTING
        '''
        if isinstance(default_block, str):
            if default_block not in BLOCK_TAGS:
                raise ValueError
        return self._call('eth_getCode', [address, default_block])

    def eth_sign(self, address, data):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign

        NEEDS TESTING
        '''
        return self._call('eth_sign', [address, data])

    def eth_sendTransaction(self, to_address=None, from_address=None, gas=None, gas_price=None, value=None, data=None,
                            nonce=None):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendtransaction

        NEEDS TESTING
        '''
        if len(to_address) == 20:
            to_address = hexlify(to_address)
        if len(from_address) == 20:
            from_address = hexlify(from_address)
        params = {}
        params['from'] = normalise_address(from_address) or self.eth_coinbase()
        if to_address is not None:
            params['to'] = normalise_address(to_address)
        if gas is not None:
            params['gas'] = hex(gas)
        if gas_price is not None:
            params['gasPrice'] = clean_hex(gas_price)
        if value is not None:
            params['value'] = clean_hex(value)
        if data is not None:
            params['data'] = data.decode('utf-8')
        if nonce is not None:
            params['nonce'] = hex(nonce)
        txid = self._call('eth_sendTransaction', [params])
        return EthTransaction(self, txid)

    def eth_sendRawTransaction(self, data):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sendrawtransaction

        NEEDS TESTING
        '''
        return self._call('eth_sendRawTransaction', [data])

    def eth_call(self, to_address, from_address=None, gas=None, gas_price=None, value=None, data=None,
                 default_block=BLOCK_TAG_LATEST):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_call

        NEEDS TESTING
        '''
        if isinstance(default_block, str):
            if default_block not in BLOCK_TAGS:
                raise ValueError
        if from_address is not None and len(from_address) == 20:
            from_address = hexlify(from_address)
        if len(to_address) == 20:
            to_address = hexlify(to_address)
        obj = {}
        obj['to'] = normalise_address(to_address)
        if from_address is not None:
            obj['from'] = normalise_address(from_address)
        if gas is not None:
            obj['gas'] = hex(gas)
        if gas_price is not None:
            obj['gasPrice'] = clean_hex(gas_price)
        if value is not None:
            obj['value'] = value
        if data is not None:
            obj['data'] = data.decode('utf-8')
        return self._call('eth_call', [obj, default_block])

    def eth_estimateGas(self, to_address=None, from_address=None, gas=None, gas_price=None, value=None, data=None,
                        default_block=BLOCK_TAG_LATEST):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_estimategas

        NEEDS TESTING
        '''
        if isinstance(default_block, str):
            if default_block not in BLOCK_TAGS:
                raise ValueError
        obj = {}
        if to_address is not None:
            obj['to'] = normalise_address(to_address)
        if from_address is not None:
            obj['from'] = normalise_address(from_address)
        if gas is not None:
            obj['gas'] = hex(gas)
        if gas_price is not None:
            obj['gasPrice'] = clean_hex(gas_price)
        if value is not None:
            obj['value'] = value
        if data is not None:
            obj['data'] = data
        return hex_to_dec(self._call('eth_estimateGas', [obj, default_block]))

    def eth_getBlockByHash(self, block_hash, tx_objects=True):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbyhash

        TESTED
        '''
        return self._call('eth_getBlockByHash', [block_hash, tx_objects])

    def eth_getBlockByNumber(self, block=BLOCK_TAG_LATEST, tx_objects=True):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getblockbynumber

        TESTED
        '''
        block = validate_block(block)
        return self._call('eth_getBlockByNumber', [block, tx_objects])

    def eth_getTransactionByHash(self, tx_hash):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyhash

        TESTED
        '''
        return self._call('eth_getTransactionByHash', [tx_hash])

    def eth_getTransactionByBlockHashAndIndex(self, block_hash, index=0):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyblockhashandindex

        TESTED
        '''
        return self._call('eth_getTransactionByBlockHashAndIndex', [block_hash, hex(index)])

    def eth_getTransactionByBlockNumberAndIndex(self, block=BLOCK_TAG_LATEST, index=0):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionbyblocknumberandindex

        TESTED
        '''
        block = validate_block(block)
        return self._call('eth_getTransactionByBlockNumberAndIndex', [block, hex(index)])

    def eth_getTransactionReceipt(self, tx_hash):
        # type: (string) -> dict
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_gettransactionreceipt

        TESTED
        '''
        return self._call('eth_getTransactionReceipt', [str(tx_hash)])

    def eth_getUncleByBlockHashAndIndex(self, block_hash, index=0):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclebyblockhashandindex

        TESTED
        '''
        return self._call('eth_getUncleByBlockHashAndIndex', [block_hash, hex(index)])

    def eth_getUncleByBlockNumberAndIndex(self, block=BLOCK_TAG_LATEST, index=0):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getunclebyblocknumberandindex

        TESTED
        '''
        block = validate_block(block)
        return self._call('eth_getUncleByBlockNumberAndIndex', [block, hex(index)])

    def eth_getCompilers(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getcompilers

        TESTED
        '''
        return self._call('eth_getCompilers')

    def eth_compileSolidity(self, code):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_compilesolidity

        TESTED
        '''
        return self._call('eth_compileSolidity', [code])

    def eth_compileLLL(self, code):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_compilelll

        N/A
        '''
        return self._call('eth_compileLLL', [code])

    def eth_compileSerpent(self, code):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_compileserpent

        N/A
        '''
        return self._call('eth_compileSerpent', [code])

    def eth_newFilter(self, from_block=BLOCK_TAG_LATEST, to_block=BLOCK_TAG_LATEST, address=None, topics=None):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newfilter

        NEEDS TESTING
        '''
        _filter = {
            'fromBlock': from_block,
            'toBlock':   to_block,
            'address':   address,
            'topics':    topics,
        }
        return self._call('eth_newFilter', [_filter])

    def eth_newBlockFilter(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newblockfilter

        TESTED
        '''
        return self._call('eth_newBlockFilter')

    def eth_newPendingTransactionFilter(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_newpendingtransactionfilter

        TESTED
        '''
        return hex_to_dec(self._call('eth_newPendingTransactionFilter'))

    def eth_uninstallFilter(self, filter_id):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_uninstallfilter

        NEEDS TESTING
        '''
        return self._call('eth_uninstallFilter', [filter_id])

    def eth_getFilterChanges(self, filter_id):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getfilterchanges

        NEEDS TESTING
        '''
        return self._call('eth_getFilterChanges', [filter_id])

    def eth_getFilterLogs(self, filter_id):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getfilterlogs

        NEEDS TESTING
        '''
        return self._call('eth_getFilterLogs', [filter_id])

    def eth_getLogs(self, filter_object):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getlogs

        NEEDS TESTING
        '''
        return self._call('eth_getLogs', [filter_object])

    def eth_getWork(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_getwork

        TESTED
        '''
        return self._call('eth_getWork')

    def eth_submitWork(self, nonce, header, mix_digest):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_submitwork

        NEEDS TESTING
        '''
        return self._call('eth_submitWork', [nonce, header, mix_digest])

    def eth_submitHashrate(self, hash_rate, client_id):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_submithashrate

        TESTED
        '''
        return self._call('eth_submitHashrate', [hex(hash_rate), client_id])

    def db_putString(self, db_name, key, value):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#db_putstring

        TESTED
        '''
        warnings.warn('deprecated', DeprecationWarning)
        return self._call('db_putString', [db_name, key, value])

    def db_getString(self, db_name, key):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#db_getstring

        TESTED
        '''
        warnings.warn('deprecated', DeprecationWarning)
        return self._call('db_getString', [db_name, key])

    def db_putHex(self, db_name, key, value):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#db_puthex

        TESTED
        '''
        if not value.startswith('0x'):
            value = '0x{}'.format(value)
        warnings.warn('deprecated', DeprecationWarning)
        return self._call('db_putHex', [db_name, key, value])

    def db_getHex(self, db_name, key):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#db_gethex

        TESTED
        '''
        warnings.warn('deprecated', DeprecationWarning)
        return self._call('db_getHex', [db_name, key])

    def shh_version(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#shh_version

        N/A
        '''
        return self._call('shh_version')

    def shh_post(self, topics, payload, priority, ttl, from_=None, to=None):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#shh_post

        NEEDS TESTING
        '''
        whisper_object = {
            'from':     from_,
            'to':       to,
            'topics':   topics,
            'payload':  payload,
            'priority': hex(priority),
            'ttl':      hex(ttl),
        }
        return self._call('shh_post', [whisper_object])

    def shh_newIdentity(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#shh_newidentity

        N/A
        '''
        return self._call('shh_newIdentity')

    def shh_hasIdentity(self, address):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#shh_hasidentity

        NEEDS TESTING
        '''
        return self._call('shh_hasIdentity', [address])

    def shh_newGroup(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#shh_newgroup

        N/A
        '''
        return self._call('shh_newGroup')

    def shh_addToGroup(self):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#shh_addtogroup

        NEEDS TESTING
        '''
        return self._call('shh_addToGroup')

    def shh_newFilter(self, to, topics):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#shh_newfilter

        NEEDS TESTING
        '''
        _filter = {
            'to':     to,
            'topics': topics,
        }
        return self._call('shh_newFilter', [_filter])

    def shh_uninstallFilter(self, filter_id):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#shh_uninstallfilter

        NEEDS TESTING
        '''
        return self._call('shh_uninstallFilter', [filter_id])

    def shh_getFilterChanges(self, filter_id):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#shh_getfilterchanges

        NEEDS TESTING
        '''
        return self._call('shh_getFilterChanges', [filter_id])

    def shh_getMessages(self, filter_id):
        '''
        https://github.com/ethereum/wiki/wiki/JSON-RPC#shh_getmessages

        NEEDS TESTING
        '''
        return self._call('shh_getMessages', [filter_id])


class ParityEthJsonRpc(EthJsonRpc):
    '''
    EthJsonRpc subclass for Parity-specific methods
    '''

    def __init__(self, host='localhost', port=PARITY_DEFAULT_RPC_PORT, tls=False):
        EthJsonRpc.__init__(self, host=host, port=port, tls=tls)

    def trace_filter(self, from_block=None, to_block=None, from_addresses=None, to_addresses=None):
        '''
        https://github.com/ethcore/parity/wiki/JSONRPC-trace-module#trace_filter

        TESTED
        '''
        params = {}
        if from_block is not None:
            from_block = validate_block(from_block)
            params['fromBlock'] = from_block
        if to_block is not None:
            to_block = validate_block(to_block)
            params['toBlock'] = to_block
        if from_addresses is not None:
            if not isinstance(from_addresses, list):
                from_addresses = [from_addresses]
            params['fromAddress'] = from_addresses
        if to_addresses is not None:
            if not isinstance(to_addresses, list):
                to_addresses = [to_addresses]
            params['toAddress'] = to_addresses
        return self._call('trace_filter', [params])

    def trace_get(self, tx_hash, positions):
        '''
        https://github.com/ethcore/parity/wiki/JSONRPC-trace-module#trace_get

        NEEDS TESTING
        '''
        if not isinstance(positions, list):
            positions = [positions]
        return self._call('trace_get', [tx_hash, positions])

    def trace_transaction(self, tx_hash):
        '''
        https://github.com/ethcore/parity/wiki/JSONRPC-trace-module#trace_transaction

        TESTED
        '''
        return self._call('trace_transaction', [tx_hash])

    def trace_block(self, block=BLOCK_TAG_LATEST):
        '''
        https://github.com/ethcore/parity/wiki/JSONRPC-trace-module#trace_block

        TESTED
        '''
        block = validate_block(block)
        return self._call('trace_block', [block])
