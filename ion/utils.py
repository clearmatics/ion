## Copyright (c) 2016-2018 Clearmatics Technologies Ltd
## Copyright (c) 2018 Harry Roberts.
## SPDX-License-Identifier: LGPL-3.0+

import sys
import json
from base64 import b64encode, b64decode
from binascii import hexlify, unhexlify
from functools import reduce

from rlp.sedes import big_endian_int
from rlp.utils import decode_hex, str_to_bytes

if sys.version_info.major == 2:
    def bytearray_to_bytestr(value):
        return bytes(''.join(chr(c) for c in value))
else:
    def bytearray_to_bytestr(value):
        return bytes(value)


TT256 = 2 ** 256


safe_ord = ord if sys.version_info.major == 2 else lambda x: x if isinstance(x, int) else ord(x)

bytes_to_int = lambda x: reduce(lambda o, b: (o << 8) + safe_ord(b), [0] + list(x))


bit_clear = lambda n, b: n ^ (1<<(b-1)) if n & 1<<(b-1) else n

bit_set = lambda n, b: n | (1<<(b-1))

bit_test = lambda n, b: 0 != (n & (1<<(b-1)))


def packl(lnum):
    assert lnum >= 0
    if lnum == 0:
        return b'\0'
    s = hex(lnum)[2:].rstrip('L')
    if len(s) & 1:
        s = '0' + s
    return unhexlify(s)

int_to_big_endian = packl


def big_endian_to_int(x): return big_endian_int.deserialize(
    str_to_bytes(x).lstrip(b'\x00'))

zpad = lambda x, l: b'\x00' * max(0, l - len(x)) + x


u256be = lambda v: zpad(int_to_big_endian(v), 32)

flatten = lambda l: [item for sublist in l for item in sublist]

# Turns a `defaultdict(defaultdict)` into a flat dictionary
dict_dump = lambda diff: {c: dict(d.items()) for c, d in diff.items()}


def is_numeric(x):
    return isinstance(x, int)


def encode_int(v):
    """encodes an integer into serialization"""
    if not is_numeric(v) or v < 0 or v >= TT256:
        raise Exception("Integer invalid or out of range: %r" % v)
    return int_to_big_endian(v)


def scan_bin(v):
    if v[:2] in ('0x', b'0x'):
        return decode_hex(v[2:])
    else:
        return decode_hex(v)


def require(arg, msg=None):
    if not arg:
        raise RuntimeError(msg or "Requirement failed")


def normalise_address(addr):
    if len(addr) == 20:
        addr = hexlify(addr).decode('ascii')
    if addr[:2] == '0x':
        addr = addr[2:]
    require(len(addr) == 40, "Invalid address: " + str(addr))
    return addr


class Marshalled(object):
    def tojson(self):
        return tojson(self)

    def marshal(self):
        return marshal(list(self))

    @classmethod
    def unmarshal(cls, args):
        return cls(*map(unmarshal, args))


def tojson(x):
    return json.dumps(marshal(x), cls=CustomJSONEncoder)


def marshal(x):
    if isinstance(x, (int, type(None))):
        return x
    if isinstance(x, (str, bytes)):
        return b64encode(x)
    if isinstance(x, (tuple, list)):
        return map(marshal, x)
    if isinstance(x, Marshalled):
        return x.marshal()
    raise ValueError("Cannot marshal type: %r - %r" % (type(x), x))


def unmarshal(x):
    if x is None or isinstance(x, int):
        return x
    if isinstance(x, (str, bytes)):
        return b64decode(x)
    if isinstance(x, (tuple, list)):
        return map(unmarshal, x)
    if isinstance(x, Marshalled):
        return x.unmarshal(x)
    raise ValueError("Cannot unmarshal type: %r - %r" % (type(x), x))


class CustomJSONEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, bytes):
            return obj.decode('utf-8', 'backslashreplace')
        return json.JSONEncoder.default(self, obj)


# XXX: about about tojson?
def json_dumps(obj):
    return json.dumps(obj, cls=CustomJSONEncoder)
