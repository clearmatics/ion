import os
import glob
import json

import click

from .ethrpc import EthJsonRpc
from .args import arg_ethrpc, arg_bytes20, arg_bytes, make_uint_n, make_bytes_n
from .crypto import keccak_256


@click.group("onchain", short_help="On-chain interfaces")
@click.option('--rpc', metavar="ip:port", callback=arg_ethrpc, help="Etherum JSON-RPC server", default='127.0.0.1:8545')
@click.option('--contract', metavar="0x...20", callback=arg_bytes20, help="Contract address")
@click.option('--account', metavar="0x...20", callback=arg_bytes20, help="Account to use for transactions")
@click.pass_context
def commands(ctx, rpc, contract, account=None):
    ctx.obj = rpc
    ctx.meta['contract'] = contract.encode('hex')
    ctx.meta['account'] = account.encode('hex')


def _make_abi_cmd(method):
    argsig = "%s" % (','.join([i['type'] for i in method['inputs']]))
    sig = bytes("%s(%s)" % (method['name'], argsig))
    sig_hash = keccak_256(sig).hexdigest()[:8]

    # TODO: if there are duplicate names, suffix with signature hash / fingerprint

    @click.command(method['name'], help=sig, short_help=argsig)
    @click.pass_context
    def cmd(ctx, **kwa):
        meta = ctx.meta
        rpc = ctx.obj
        assert isinstance(rpc, EthJsonRpc)

        args = [kwa[ _['name'] ] for _ in method['inputs']]
        result_types = [_['type'] for _ in method['outputs']]

        click.echo(sig)
        click.echo("\nInput:")
        for inp in method['inputs']:
            click.echo("\t%s %s" % (inp['name'], inp['type']))
        click.echo()

        result = rpc.call(meta['contract'], sig, args, result_types)
        for idx, output in enumerate(method['outputs']):
            print("\t%r %r = %r" % (output['type'], output['name'], result[idx]))

    argtypes = {
        'bytes': arg_bytes,
        'address': arg_bytes20,
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

        # is an array?
        arg_type = arg['type']
        if arg_type.endswith("[]"):
            arg_type = arg_type[:-2]
            kwa['multiple'] = True

        # parse size from type
        without_ints = arg_type.rstrip('0123456789')
        if without_ints != arg_type:
            n = int(arg_type[len(without_ints):])
            kwa['callback'] = makers[without_ints](n)
        else:
            kwa['callback'] = argtypes[arg_type]

        name = arg['name'].strip('_')
        cmd = click.option('--' + name, arg['name'], **kwa)(cmd)

    return cmd


def _make_abi_group(abi_file):
    with open(abi_file, 'rb') as handle:
        data = json.load(handle)

        group_name = os.path.basename(abi_file).split('.')[0]
        group = click.Group(group_name)

        for method in data:
            if method['type'] != 'function':
                continue
            if method['name'][0] == '_':
                continue
            group.add_command( _make_abi_cmd(method) )

        return group


for filename in glob.glob('abi/*.abi'):
    commands.add_command( _make_abi_group(filename) )


if __name__ == "__main__":
    commands()
