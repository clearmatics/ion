
import os
import pyjsonrpc
from ethereum.utils import privtoaddr

from ..plasma.model import Block
from ..plasma.payment import dependency_hash, Payment, SignedPayment
from ..utils import unmarshal, marshal


class IonClient(object):
    rpc = None  # type: IonRpcServer

    def __init__(self, rpc_endpoint, secret, public=None):
        """
        :type rpc_endpoint: IonRpcServer | str
        :type secret: str | bytes
        """
        assert len(secret) == 32
        if isinstance(rpc_endpoint, (str, unicode, bytes)):
            assert rpc_endpoint[:4] == 'http'
            self.rpc = pyjsonrpc.HttpClient(rpc_endpoint)
        else:
            assert isinstance(rpc_endpoint, object)
            self.rpc = rpc_endpoint
        self.secret = secret

        if public:
            self.public = public
        else:
            self.public = privtoaddr(secret)
            print(self.public.encode('hex'))

    def __str__(self):
        return self.public.encode('hex')

    def block_get_latest(self):
        return Block.unmarshal(self.rpc.block_get_latest())

    def payment_sign(self, payment_args, prev_hash=None):
        if prev_hash is None:
            prev_hash = unmarshal(self.rpc.block_hash())
        payment = Payment(*payment_args)
        return payment.seal(self.secret, prev_hash)

    def payment_pending(self, ref):
        if isinstance(ref, Payment):
            ref = ref.r
        elif isinstance(ref, SignedPayment):
            ref = ref.p.r
        assert len(ref) == 32
        ref = marshal(ref)
        return self.rpc.payment_pending(ref)

    def graph(self):
        return self.rpc.graph()

    def balance(self, currency, holder=None):
        if holder is None:
            holder = self.public
        assert len(currency) == 20
        assert len(holder) == 20
        currency = marshal(currency)
        holder = marshal(holder)
        return self.rpc.balance(currency, holder)

    def mint(self, n_tokens):
        """Give myself N tokens of my own currency"""
        n_tokens = int(n_tokens)
        assert n_tokens > 0
        return self.pay(self.public, self.public, n_tokens)

    def commit(self):
        """Commit the transaction pool, sealing it into a new block"""
        return Block.unmarshal(self.rpc.block_commit())

    def _marshal_dependencies(self, dependencies):
        if dependencies is None:
            return None
        out = []
        for dep in dependencies:
            t, c, v, r = dep
            if isinstance(t, IonClient):
                t = t.public
            if isinstance(c, IonClient):
                c = c.public
            out.append( dependency_hash(t, c, v, r) )
        return out

    def pay(self, recipient, currency, value, r=None, deps=None):
        """
        :type recipient: IonClient | str
        :type currency: IonClient | str
        """
        if isinstance(recipient, IonClient):
            recipient = recipient.public
        if isinstance(currency, IonClient):
            currency = currency.public
        if r is None:
            r = os.urandom(32)
        deps = self._marshal_dependencies(deps)
        f = self.public
        sp = self.payment_sign((f, recipient, currency, value, r, deps))
        return self.payment_submit(sp)

    # TODO: payment_cancel, payment_update

    def payment_submit(self, sp):
        assert isinstance(sp, SignedPayment)
        return self.rpc.payment_submit(sp.marshal())
