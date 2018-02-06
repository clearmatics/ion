import json
from collections import namedtuple
from .ethrpc import EthJsonRpc
from .crypto import keccak_256

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
    proxy = dict()
    for method in abi:
        if method['type'] != 'function':
            continue
        handler = solproxy_bind(rpc, method, address, account)
        sig = "%s(%s)" % (method['name'], ','.join([i['type'] for i in method['inputs']]))
        sig_hash = keccak_256(bytes(sig)).hexdigest()[:8]
        proxy[method['name']] = handler
        proxy[method['name'] + '_' + sig_hash] = handler
    return namedtuple('SolProxy', proxy.keys())(*proxy.values())
