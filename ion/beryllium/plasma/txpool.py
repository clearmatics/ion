import random

from .chain import chaindata_path
from .payment import payments_graphviz


class TxPool(object):
    def __init__(self, target):
        self._pool = dict()
        self._target = target

    @property
    def target(self):
        return self._target

    @property
    def payments(self):
        return [sp.p for sp in self._pool.values()]

    # TODO: does it balance?

    # TODO: are all balances above zero?
    # do { prune() } while ( ! balanced() );

    def add(self, lithium_root):
        """
        :type lithium_root: SignedPayment
        """
        if not signed_payment.verify(self._target):
            raise RuntimeError("Invalid payment signature")
        payment = signed_payment.p

        if payment.r in self._pool:
            # TODO: replacing payment must be authenticated
            #  e.g. include the hash of the previous payment used to replace it
            raise RuntimeError("Duplicate payment reference")

        self._pool[payment.r] = signed_payment

    def graph(self, filename='txpool'):
        g = payments_graphviz([sp.p for sp in self._pool.values()])
        g.render( chaindata_path(filename, 'graphviz') )
        return True

    def prune(self, balances):
        """Remove payments which don't have the necessary dependencies"""
        # TODO: remove payments which exceed available balance
        #     - prioritise?

        # Re-check dependencies after every removal until there is a stable set
        signed_payments = self._pool.values()
        # refs = set([sp.p.r for sp in signed_payments])

        while True:
            dependencies = [sp.p.dependency_hash() for sp in signed_payments]
            prev_len = len(signed_payments)
            random.shuffle(signed_payments)  # shuffle every time
            signed_payments = [sp for sp in signed_payments
                               if sp.p.d is None or any([d in dependencies for d in sp.p.d])]
            if prev_len == len(signed_payments):
                break

        # new_refs = set([sp.p.r for sp in signed_payments])
        # TODO: remove evict items from transaction pool
        # TODO: show diff of payments, which ones have been removed
        return signed_payments

    def pending(self, reference):
        return reference in self._pool

    def cancel(self, reference):
        if reference not in self._pool:
            return False
        del self._pool[reference]
        return True
