import argparse

from .ethrpc import EthJsonRpc
from .utils import require, scan_bin


def bytes20(value):
    value = scan_bin(value)
    require( len(value) == 20, "20 bytes required" )
    return value


def bytes32(value):
    value = scan_bin(value)
    require( len(value) == 32, "32 bytes required" )
    return value


def posint256(value):
    value = int(value)
    require( value > 0 )
    require( value <= (1<<255) )
    return value

class BinAction(argparse.Action):
    def __init__(self, option_strings, dest, nargs=None, **kwa):
        require( nargs is None, 'nargs not allowed' )
        super(BinAction, self).__init__(option_strings, dest, **kwa)

    def __call__(self, parser, namespace, values, option_string=None):
        setattr(namespace, self.dest, scan_bin(values))


class PosInt256(argparse.Action):
    """Parse Ethereum address secret"""

    def __init__(self, option_strings, dest, nargs=None, **kwa):
        require( nargs is None, 'nargs not allowed' )
        super(PosInt256, self).__init__(option_strings, dest, **kwa)

    def __call__(self, parser, namespace, values, option_string=None):
        setattr(namespace, self.dest, posint256(values))


class Bytes32(argparse.Action):
    """Parse Ethereum bytes32"""

    def __init__(self, option_strings, dest, nargs=None, **kwa):
        require( nargs is None, 'nargs not allowed' )
        super(Bytes32, self).__init__(option_strings, dest, **kwa)

    def __call__(self, parser, namespace, values, option_string=None):
        setattr(namespace, self.dest, bytes32(values))


class Bytes20(argparse.Action):
    """Parse Ethereum address"""

    def __init__(self, option_strings, dest, **kwa):
        # require( nargs is None, 'nargs not allowed' )
        super(Bytes20, self).__init__(option_strings, dest, **kwa)

    def __call__(self, parser, namespace, values, option_string=None):
        if isinstance(values, list):
            values = map(bytes20, values)
        else:
            values = bytes20(values)
        setattr(namespace, self.dest, values)


class EthRpc(argparse.Action):
    """Parse RPC address, return EthJsonRpc handle"""

    def __init__(self, option_strings, dest, nargs=None, **kwa):
        require(nargs is None, 'nargs not allowed')
        super(EthRpc, self).__init__(option_strings, dest, **kwa)

    def __call__(self, parser, namespace, values, option_string=None):
        ip, port = values.split(':')
        port = int(port)
        setattr(namespace, self.dest, EthJsonRpc(ip, port))


class PaymentDependency(argparse.Action):
    def __init__(self, option_strings, dest, nargs=None, **kwa):
        require(nargs is None, 'nargs not allowed')
        super(PaymentDependency, self).__init__(option_strings, dest, **kwa)

    def __call__(self, parser, namespace, values, option_string=None):
        pay_to , pay_cur, pay_ref, pay_amt = values.split(',')
        pay_to = bytes20(pay_to)
        pay_cur = bytes20(pay_cur)
        pay_amt = posint256(pay_amt)
        pay_ref = bytes32(pay_ref)
        if len(scan_bin(values)) != 20:
            raise ValueError('Invalid address length')
        setattr(namespace, self.dest, (pay_to, pay_cur, pay_amt, pay_ref))
