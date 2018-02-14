import os
import glob
import click

from ..args import arg_ethrpc, arg_bytes20, arg_bytes32, arg_uint256, arg_bytes, make_uint_n, make_bytes_n
from ..crypto import keccak_256


@click.group("onchain", help="On-chain interfaces")
@click.option('--rpc', metavar="ip:port", callback=arg_ethrpc, help="Etherum JSON-RPC server", default='127.0.0.1:8545')
@click.option('--contract', metavar="0x...20", callback=arg_bytes20, help="Contract address")
@click.option('--account', metavar="0x...20", callback=arg_bytes20, help="Account to use for transactions")
@click.pass_context
def commands(ctx, rpc, contract, account=None):
    """
    :type ctx: click.Context
    """
    ctx.meta['rpc'] = rpc
    ctx.meta['contract'] = contract
    ctx.meta['account'] = account


def subcmd(abi_file):
    import json
    with open(abi_file, 'rb') as handle:
        data = json.load(handle)

        grou_name = os.path.basename(abi_file).split('.')[0]
        group = click.Group(grou_name )

        for method in data:
            if method['type'] != 'function':
                continue
            if method['name'][0] == '_':
                continue

            argsig = "%s" % (','.join([i['type'] for i in method['inputs']]))
            sig = "%s(%s)" % (method['name'], argsig)
            sig_hash = keccak_256(bytes(sig)).hexdigest()[:8]

            @click.command(method['name'], help=sig, short_help=argsig)
            def cmd():
                pass

            argtypes = {
                'bytes': arg_bytes,
                'address': arg_bytes20,
                'bytes32': arg_bytes32,
                'uint256': arg_uint256,
            }
            makers = {
                'uint': make_uint_n,
                'bytes': make_bytes_n,
            }
            for arg in method['inputs']:
                kwa = dict(
                    required=True,
                    metavar=arg['type'],
                )
                arg_type = arg['type']
                if arg_type.endswith("[]"):
                    arg_type = arg_type[:-2]
                    kwa['multiple'] = True

                without_ints = arg_type.rstrip('0123456789')
                if without_ints != arg_type:
                    n = int(arg_type[len(without_ints):])
                    kwa['type'] = makers[without_ints](n)
                else:
                    kwa['type'] = argtypes[arg_type]

                name = arg['name'].strip('_')
                cmd = click.option('--' + name, **kwa)(cmd)

            group.add_command(cmd)
        commands.add_command(group)


for filename in glob.glob('abi/*.abi'):
    subcmd(filename)

if __name__ == "__main__":
    commands()
