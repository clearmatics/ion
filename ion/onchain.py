import os
import glob
import json
import time

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


def _dispatch_cmd(meta, rpc, method, args, sig, wait=False, commit=False):
    result_types = [_['type'] for _ in method['outputs']]

    click.echo(method['name'] + " inputs:")
    for n, inp in enumerate(method['inputs']):
        click.echo("\t%s %s = %r" % (inp['type'], inp['name'], args[n]))

    if method['constant'] or not commit:
        result = rpc.call(meta['contract'], sig, args, result_types)
        if result:
            click.echo("\noutputs:")
            for idx, output in enumerate(method['outputs']):
                click.echo("\t%r %r = %r" % (output['type'], output['name'], result[idx]))
    else:
        tx = rpc.call_with_transaction(meta['account'], meta['contract'], sig, args, result_types)
        click.echo("Tx %s" % (tx,))
        tx_receipt = None
        if wait:
            first = True
            while True:
                tx_receipt = rpc.eth_getTransactionReceipt(tx)
                if tx_receipt is not None:
                    break
                if first:
                    click.echo("Waiting for transaction receipt")
                    first = False
                try:
                    time.sleep(1)
                except KeyboardInterrupt:
                    click.echo("Skipped waiting for transaction receipt...")
                    break
        if tx_receipt:
            click.echo("Receipt: %r" % (tx_receipt,))


def _make_abi_cmd(method):
    argsig = "%s" % (','.join([i['type'] for i in method['inputs']]))
    sig = bytes("%s(%s)" % (method['name'], argsig))
    sig_hash = keccak_256(sig).hexdigest()[:8]

    # TODO: if there are duplicate names, suffix with signature hash / fingerprint

    @click.command(method['name'], help=sig, short_help=argsig)
    @click.option('-w', 'wait', is_flag=True, help="Wait until mined")
    @click.option('-c', 'commit', is_flag=True, help="Commit transaction")
    @click.pass_context
    def cmd(ctx, wait, commit, **kwa):
        meta = ctx.meta
        rpc = ctx.obj
        assert isinstance(rpc, EthJsonRpc)
        args = [kwa[_['name']] for _ in method['inputs']]
        _dispatch_cmd(meta, rpc, method, args, sig, wait, commit)

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

        # parse size from type, e.g.  uint256
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
