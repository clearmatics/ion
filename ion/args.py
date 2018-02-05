import argparse

from ethjsonrpc import EthJsonRpc
from ethereum.utils import scan_bin

from .utils import require


def arg_bytes20(value):
    value = scan_bin(value)
    require( len(value) == 20, "20 bytes required" )
    return value


def arg_bytes32(value):
    value = scan_bin(value)
    require( len(value) == 32, "32 bytes required" )
    return value


def arg_posint256(value):
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


class PosInt256Action(argparse.Action):
    """Parse Ethereum address secret"""

    def __init__(self, option_strings, dest, nargs=None, **kwa):
        require( nargs is None, 'nargs not allowed' )
        super(PosInt256Action, self).__init__(option_strings, dest, **kwa)

    def __call__(self, parser, namespace, values, option_string=None):
        setattr(namespace, self.dest, arg_posint256(values))


class Bytes32Action(argparse.Action):
    """Parse Ethereum bytes32"""

    def __init__(self, option_strings, dest, nargs=None, **kwa):
        require( nargs is None, 'nargs not allowed' )
        super(Bytes32Action, self).__init__(option_strings, dest, **kwa)

    def __call__(self, parser, namespace, values, option_string=None):
        setattr(namespace, self.dest, arg_bytes32(values))


class Bytes20Action(argparse.Action):
    """Parse Ethereum address"""

    def __init__(self, option_strings, dest, nargs=None, **kwa):
        require( nargs is None, 'nargs not allowed' )
        super(Bytes20Action, self).__init__(option_strings, dest, **kwa)

    def __call__(self, parser, namespace, values, option_string=None):
        require( len(scan_bin(values)) == 20, 'Invalid address length' )
        setattr(namespace, self.dest, arg_bytes20(values))


class EthRpcAction(argparse.Action):
    """Parse RPC address, return EthJsonRpc handle"""

    def __init__(self, option_strings, dest, nargs=None, **kwa):
        require(nargs is None, 'nargs not allowed')
        super(EthRpcAction, self).__init__(option_strings, dest, **kwa)

    def __call__(self, parser, namespace, values, option_string=None):
        ip, port = values.split(':')
        port = int(port)
        setattr(namespace, self.dest, EthJsonRpc(ip, port))


class PaymentDependencyAction(argparse.Action):
    def __init__(self, option_strings, dest, nargs=None, **kwa):
        require(nargs is None, 'nargs not allowed')
        super(PaymentDependencyAction, self).__init__(option_strings, dest, **kwa)

    def __call__(self, parser, namespace, values, option_string=None):
        pay_to , pay_cur, pay_ref, pay_amt = values.split(',')
        pay_to = arg_bytes20(pay_to)
        pay_cur = arg_bytes20(pay_cur)
        pay_amt = arg_posint256(pay_amt)
        pay_ref = arg_bytes32(pay_ref)
        if len(scan_bin(values)) != 20:
            raise ValueError('Invalid address length')
        setattr(namespace, self.dest, (pay_to, pay_cur, pay_amt, pay_ref))
