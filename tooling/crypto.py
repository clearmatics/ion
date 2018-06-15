## Copyright (c) 2016-2018 Clearmatics Technologies Ltd
## SPDX-License-Identifier: LGPL-3.0+

#!/usr/bin/env python
"""
Crypto: Has a load of useful crypto stuff
"""
from collections import namedtuple

from ethereum.utils import big_endian_to_int, encode_int32
from rlp.utils import ascii_chr

from sha3 import keccak_256

from .utils import Marshalled, u256be, safe_ord

try:
    import coincurve
except ImportError:
    from py_ecc.secp256k1 import ecdsa_raw_recover, ecdsa_raw_sign
    import warnings
    warnings.warn('could not import coincurve', ImportWarning)
    coincurve = None


# --------------------------------------------------------------------
# Datatypes


_EcdsaSignatureStruct = namedtuple('EcdsaSignature', ('v', 'r', 's'))


class EcdsaSignature(_EcdsaSignatureStruct, Marshalled):
    def __init__(self, *args, **kwa):
        _EcdsaSignatureStruct.__init__(*args, **kwa)

    def dump(self):
        # TODO: make same format as coincurve 65 byte str
        return ''.join([chr(self.v), self.r, self.s])

    def recover(self, rawhash):
        assert isinstance(self, EcdsaSignature)
        v, r, s = self
        if coincurve and hasattr(coincurve, "PublicKey"):
            try:
                pk = coincurve.PublicKey.from_signature_and_message(
                    ''.join([r, s, ascii_chr(v - 27)]),
                    rawhash,
                    hasher=None,
                )
                pub = pk.format(compressed=False)[1:]
            except BaseException:
                pub = b"\x00" * 64
        else:
            r = big_endian_to_int(r)
            s = big_endian_to_int(s)
            result = ecdsa_raw_recover(rawhash, (v, r, s))
            if result:
                x, y = result
                pub = encode_int32(x) + encode_int32(y)
            else:
                raise ValueError('Invalid VRS')
        assert len(pub) == 64

        # Convert to Ethereum address
        return keccak_256(pub).digest()[12:]


# --------------------------------------------------------------------
# ECDSA signature and address recovery


def ecdsa_sign(rawhash, key):
    # type: (bytes, bytes) -> EcdsaSignature
    if coincurve and hasattr(coincurve, 'PrivateKey'):
        pk = coincurve.PrivateKey(key)
        signature = pk.sign_recoverable(rawhash, hasher=None)
        v = safe_ord(signature[64]) + 27
        r = signature[0:32]
        s = signature[32:64]
    else:
        v, r, s = ecdsa_raw_sign(rawhash, key)
        r = u256be(r)
        s = u256be(s)
    return EcdsaSignature(v, r, s)