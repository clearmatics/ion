import json
from collections import namedtuple
from ethjsonrpc import EthJsonRpc


def solproxy_bind(rpc, method, address, account):
    ins = [_['type'] for _ in method['inputs']]
    outs = [_['type'] for _ in method['outputs']]
    sig = method['name'] + '(' + ','.join(ins) + ')'
    if method['constant']:
        return lambda *args, **kwa: (rpc.call(address, sig, args, outs, **kwa)
                                     if len(outs) > 1 else
                                     rpc.call(address, sig, args, outs, **kwa)[0])
    return lambda *args, **kwa: (rpc.call_with_transaction(account, address, sig, args, **kwa)
                                 if len(outs) > 1 else
                                 rpc.call_with_transaction(account, address, sig, args, **kwa)[0])


def solproxy(rpc, abi, address, account):
    if isinstance(rpc, tuple):
        EthJsonRpc(*rpc)
    if isinstance(abi, str):
        with open(abi) as handle:
            abi = json.load(handle)
    proxy = {method['name']: solproxy_bind(rpc, method, address, account)
             for method in abi
             if method['type'] == 'function'}
    return namedtuple('SolProxy', proxy.keys())(*proxy.values())
