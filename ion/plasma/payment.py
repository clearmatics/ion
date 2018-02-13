from __future__ import print_function
import os
import sys
import random
from collections import defaultdict, namedtuple

import click
from ethereum.utils import privtoaddr

from ..args import arg_bytes20, arg_bytes32, arg_uint256
from ..utils import u256be, require, Marshalled
from ..crypto import EcdsaSignature, ecdsa_sign, keccak_256


_SignedPaymentStruct = namedtuple('SignedPayment', ('p', 's'))

_PaymentStruct = namedtuple('Payment', ('f', 't', 'c', 'v', 'r', 'd'))


class Payment(_PaymentStruct, Marshalled):
    def __init__(self, *args, **kwa):
        super(Payment, self).__init__(*args, **kwa)

    def seal(self, secret, prev_hash):
        # type: (bytes, bytes) -> SignedPayment
        require(len(secret) == 32)
        msg = self.hash(prev_hash)
        return SignedPayment(self, ecdsa_sign(msg, secret))

    def hash(self, prev_hash):
        return keccak_256(prev_hash + self.dump()).digest()

    def dependency_hash(self):
        return dependency_hash(self.t, self.c, self.v, self.r)

    def dump(self):
        require(all([len(x) == 20 for x in [self.f, self.t, self.c]]))
        require(len(self.r) == 32, "Bad payment reference size")
        require(isinstance(self.v, (int, long)), "Invalid paymentvalue type")
        require(self.d is None or isinstance(self.d, (tuple, list, str, bytes)))
        dependencies = ''
        if isinstance(self.d, (tuple, list)):
            require(all([len(x) == 32 for x in self.d]))
            dependencies = ''.join(self.d)
        return ''.join([self.f, self.t, self.c, u256be(self.v), self.r, dependencies])


class SignedPayment(_SignedPaymentStruct, Marshalled):
    def __init__(self, *args, **kwa):
        super(SignedPayment, self).__init__(*args, **kwa)

    def marshal(self):
        return SignedPayment(self.p.marshal(),
                             self.s.marshal())

    @classmethod
    def unmarshal(cls, args):
        assert isinstance(args, (tuple, list))
        assert len(args) == 2
        return cls(Payment.unmarshal(args[0]),
                   EcdsaSignature.unmarshal(args[1]))

    def verify(self, prev_hash):
        # type: (bytes) -> bool
        msg = self.p.hash(prev_hash)
        f = self.s.recover(msg)
        return f == self.p.f

    def open(self, prev_hash):
        # type: (bytes) -> Payment
        if not self.verify(prev_hash):
            raise ValueError("Invalid signature,")
        return self.p


# --------------------------------------------------------------------
# Payment functions


def dependency_hash(t, c, v, r):
    return keccak_256(''.join([t, c, u256be(v), r])).digest()


def payments_apply(payments):
    # type: (list[Payment]) -> (dict, dict)
    # Calculate resulting balance in each currency for every party
    balances = defaultdict(lambda: defaultdict(int))
    created = defaultdict(int)
    for p in payments:
        if p.v <= 0:
            raise RuntimeError("Payment value < 0: %r" % (p,))
        if p.f == p.c and p.f == p.t:
            # Self transaction mints your own currency
            created[p.f] += p.v
        else:
            balances[p.c][p.f] -= p.v
        balances[p.c][p.t] += p.v
    return balances, created


def payments_sum(payments):
    # type: (list[Payment]) -> list[Payment]
    outv = defaultdict(int)
    outr = defaultdict(list)
    outd = defaultdict(list)
    for p in payments:
        outv[(p.f,p.t,p.c)] += p.v
        outr[(p.f,p.t,p.c)] += p.r
        if p.d is not None:
            outd[(p.f,p.t,p.c)] += p.d
    return [ Payment(f, t, c, v, outr[(f,t)], outd[(f,t)] or None)
             for (f, t, c), v in outv.items() ]


# --------------------------------------------------------------------
# Functions for testing


def random_keypairs(n):
    keys = [os.urandom(32) for _ in range(0, n)]
    pubs = map(privtoaddr, keys)
    return dict(zip(pubs, keys))


def random_payment(balances, keypairs):
    # TODO: make code smaller, easier to read...
    while True:
        if len(balances) and random.choice([False, True]):
            c = random.choice(balances.keys())
            f, b = random.choice(balances[c].items())
            t = random.choice(keypairs.keys())
            if f == t:
                continue
        else:
            f = t = c = random.choice(keypairs.keys())
            b = 10000
        if b <= 0:
            continue
        v = random.randint(0, b)
        if v <= 0:
            continue
        if f != t or c != f:
            balances[c][f] -= v
        balances[c][t] += v
        return balances, Payment(f, t, c, v, os.urandom(32), None)


def random_payments(prev_hash, npayments=60, nparticipants=6):
    keypairs = random_keypairs(nparticipants)
    balances = defaultdict(lambda: defaultdict(int))
    out = []
    while len(out) < npayments:
        balances, p = random_payment(balances, keypairs)
        out.append( p.seal(keypairs[p.f], prev_hash) )
    return out


def payments_graphviz(payments, colours=dict()):
    # type: (list[Payment], dict) -> graphviz.Digraph

    #payments = payments_sum(payments)

    import graphviz
    d = graphviz.Digraph(engine='circo')

    # Used to determine size of line in relation to value
    maxvals = defaultdict(list)
    for p in payments:
        maxvals[p.c].append(p.v)
    maxvals = {c: max(v) for c, v in maxvals.items()}

    for p in payments:
        tx_name = p.dependency_hash()[:3].encode('hex')
        pen_width = str((float(abs(p.v)) / maxvals[p.c]) * 5)
        d.edge(p.f[:3].encode('hex'), 'tx ' + tx_name,
               weight=str(p.v), fontsize="8.0",
               label=' '.join([str(p.v), p.c[:3].encode('hex')]),
               color=colours.get(p.c, 'black'),
               penwidth=pen_width)
        d.edge('tx ' + tx_name, p.t[:3].encode('hex'),
               weight=str(p.v), fontsize="8.0",
               label=' '.join([str(p.v), p.c[:3].encode('hex')]),
               color=colours.get(p.c, 'black'),
               penwidth=pen_width)
        if p.d:
            for dep in p.d:
                d.edge('tx ' + tx_name, 'tx ' + dep[:3].encode('hex'),
                       color='grey', style="dotted")
    return d


# --------------------------------------------------------------------
# Program entry, demonstrates payment parameters and encoding


@click.command()
@click.option('--secret', '-s', metavar='BYTES32', callback=arg_bytes32)
@click.option('--input', '-i', metavar='FILE', type=click.File('r'), required=False)
@click.option('--from', '-f', 'source', metavar='ADDR', callback=arg_bytes20, required=False)
@click.option('--to', '-t', 'dest', metavar='ADDR', callback=arg_bytes20, required=False)
@click.option('--currency', '-c', metavar='ADDR', callback=arg_bytes20, required=False)
@click.option('--reference', '-r', metavar='BYTES32', callback=arg_bytes32, required=False)
@click.option('--value', '-v', metavar='UINT256', callback=arg_uint256)
@click.option('--block-hash', '-b', metavar='BYTES32', callback=arg_bytes32)
@click.option("--format", '-f', type=click.Choice(['meta', 'json']), default='meta')
def main(secret, input, source, dest, currency, reference, value, block_hash, format):
    import json

    currency = currency or os.urandom(20)
    reference = reference or os.urandom(32)
    block_hash = block_hash or os.urandom(32)
    dest = dest or os.urandom(20)
    secret = secret or os.urandom(32)
    value = value or random.randint(1, 1000)

    if secret:
        source = privtoaddr(secret)

    signed_payment = None

    if input:
        data = json.load(input)
        signed_payment = SignedPayment.unmarshal(data)
        payment = signed_payment.open(block_hash)
    else:
        payment = Payment(source, dest, currency, value, reference, [])

    if secret:
        signed_payment = payment.seal(secret, block_hash)

    assert payment is not None

    if format == 'meta':
        print("Block Hash:", block_hash.encode('hex'))
        print("Hash:", payment.hash(block_hash).encode('hex'))
        print("Dependency Hash:", payment.dependency_hash().encode('hex'))
        print("Target:", block_hash.encode('hex'))
        print("From:", payment.f.encode('hex'))
        print("To:", payment.t.encode('hex'))
        print("Currency:", payment.c.encode('hex'))
        print("Ref:", payment.r.encode('hex'))
        print("Value:", payment.v)
        if payment.d:
            for d in payment.d:
                print("\tDepends:", d.encode('hex'))

    elif format == 'json':
        print(signed_payment.tojson())

    return 0


if __name__ == "__main__":
    sys.exit(main())
