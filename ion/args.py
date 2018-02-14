
from .ethrpc import EthJsonRpc
from .utils import require, scan_bin


def arg_bytes(ctx, param, value):
    if value is None:
        return None
    value = scan_bin(value)
    return value

def make_bytes_n(n):
    def arg_bytes_n(ctx, param, value):
        value = arg_bytes(ctx, param, value)
        if value is None:
            return None
        require( len(value) == n, str(n) + " bytes required" )
        return value
    return arg_bytes_n


arg_bytes20 = make_bytes_n(20)
arg_bytes32 = make_bytes_n(32)


def make_uint_n(n):
    def arg_uint_n(ctx, param, value):
        if value is None:
            return None
        value = int(value)
        require( value >= 0 )
        require( value <= (1<<(n-1)) )
        return value
    return arg_uint_n

arg_uint256 = make_uint_n(256)


def arg_ethrpc(ctx, param, value):
    if value is None:
        return None
    ip, port = value.split(':')
    port = int(port)
    require( port > 0 )
    require( port < 0xFFFF )
    return EthJsonRpc(ip, port)


"""
class PaymentDependency(argparse.Action):
    def __init__(self, option_strings, dest, nargs=None, **kwa):
        require(nargs is None, 'nargs not allowed')
        super(PaymentDependency, self).__init__(option_strings, dest, **kwa)

    def __call__(self, parser, namespace, values, option_string=None):
        pay_to , pay_cur, pay_ref, pay_amt = values.split(',')
        pay_to = arg_bytes20(pay_to)
        pay_cur = arg_bytes20(pay_cur)
        pay_amt = arg_posint256(pay_amt)
        pay_ref = arg_bytes32(pay_ref)
        if len(scan_bin(values)) != 20:
            raise ValueError('Invalid address length')
        setattr(namespace, self.dest, (pay_to, pay_cur, pay_amt, pay_ref))
"""
