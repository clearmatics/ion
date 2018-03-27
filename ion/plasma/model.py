import json
from collections import namedtuple

from ..utils import Marshalled, require
from ..crypto import keccak_256


_BlockStruct= namedtuple('Block', ('prev', 'root'))


class Block(_BlockStruct, Marshalled):
    def __init__(self, *args, **kwa):
        _BlockStruct.__init__(*args, **kwa)

    @property
    def hash(self):
        require( len(self.prev) == 32 )
        require( len(self.root) == 32 )
        if long(self.prev.encode('hex'), 16) == 0:
            return keccak_256(self.root).digest()
        return keccak_256(self.prev + self.root).digest()

    def __str__(self):
        rootHex = '0x'+self.root.encode('hex')
        prevHex = '0x'+self.prev.encode('hex')
        return json.dumps(dict({'root': rootHex, 'prev': prevHex}), indent=2)
