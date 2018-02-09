#!/usr/bin/env python
from __future__ import print_function
import random

from .utils import zpad, int_to_big_endian, bit_clear, bit_test, bit_set, bytes_to_int
from .crypto import keccak_256

def serialize(v):
    if isinstance(v, str):
        return v
    if isinstance(v, (int, long)):
        return zpad(int_to_big_endian(v), 32)
    raise NotImplementedError(v)


hashs = lambda *x: bytes_to_int(keccak_256(''.join(map(serialize, x))).digest())

merkle_hash = lambda *x: bit_clear(hashs(*x), 0xFF)


def merkle_tree(items):
    tree = [map(merkle_hash, items)]
    extra = merkle_hash("merkle-tree-extra")
    while True:
        level = tree[-1]
        # Ensure level has an even number of items, pad it with an 'extra item'
        if len(level) % 2 != 0:
            level.append( extra )
        # Hash each pair in the list to create the next level
        it = iter(level)
        tree.append([merkle_hash(item, next(it)) for item in it])
        if len(tree[-1]) == 1:
            break
    return tree, tree[-1][0]


def merkle_path(item, tree):
    """
    Create a merkle path for the item within the tree
    max length = (height*2) - 1
    min length = 1
    """
    item = merkle_hash(item)
    idx = tree[0].index(item)

    path = []
    for level in tree[:-1]:
        if (idx % 2) == 0:
            path.append(bit_set(level[idx+1], 255))
        else:
            path.append(level[idx-1])
        idx = idx // 2
    return path


def merkle_proof(leaf, path, root):
    """
    Verify merkle path for an item matches the root
    """
    node = merkle_hash(leaf)
    for item in path:
        if bit_test(item, 255):
            node = merkle_hash(node, bit_clear(item, 255))
        else:
            node = merkle_hash(item, node)
    return node == root


def main():
    for i in range(1, 100):
        items = range(0, i)
        tree, root = merkle_tree(items)
        random.shuffle(items)
        for item in items:
            proof = merkle_path(item, tree)
            assert merkle_proof(item, proof, root) is True


if __name__ == "__main__":
    main()
