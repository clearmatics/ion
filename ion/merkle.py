#!/usr/bin/env python
## Copyright (c) 2016-2018 Clearmatics Technologies Ltd
## SPDX-License-Identifier: LGPL-3.0+

"""
Merkle:
Provides an interface to produce merkle trees, proofs, etc.
"""
from __future__ import print_function

import random

from .crypto import keccak_256
from .utils import zpad, int_to_big_endian, bit_clear, bit_test, bit_set, bytes_to_int


def serialize(v):
    """Convert to value to a hashable scalar"""
    if isinstance(v, str):
        return v
    if isinstance(v, (int, long)):
        return zpad(int_to_big_endian(v), 32)
    raise NotImplementedError(v)


hashs = lambda *x: bytes_to_int(keccak_256(''.join(map(serialize, x))).digest())

merkle_hash = lambda *x: bit_clear(hashs(*x), 0xFF)

def merkle_tree(items):
    """
    Hashes a list of items, then creates a Merkle tree where the items are
    hashed in pairs to form the next level of the tree until the level is
    only one item (the root).

    ```
    [ [ H(0), H(1), H(2), H(3) ]                # Level 0
      [ H(H(0)||H(1)), H(H(2)||H(3)) ]          # Level 1
      [ H(H(H(0)||H(1))||H(H(2)||H(3))) ] ]     # Level 2 (root)
    ```

    If a level has an odd number of items it is padded with an 'extra' item
    to keep the tree perfectly balanced.

    The first level of items is sorted.

    :type items: list
    :return: list, long
    """
    tree = [sorted(map(merkle_hash, items))]
    extra = merkle_hash("merkle-tree-extra")
    while True:
        level = tree[-1]
        # Ensure level has an even number of items, pad it with an 'extra item'
        if len(level) % 2 != 0:
            level.append(extra)
        # Hash each pair in the list to create the next level
        it = iter(level)
        tree.append([merkle_hash(item, next(it)) for item in it])
        if len(tree[-1]) == 1:
            break
    return tree, tree[-1][0]


def merkle_path(item, tree):
    """
    Given a tree and an item, return a path which can be used to
    verify the item exists within a root.

    The path for `x` is: L5, L11
    The root is: 17, or `H(11, H(5, H(D)))`

    ```
                   |
                   v
     a    b   c    x       <- items
     |    |   |    |
    L2   R3  [L5] H(x)     <- level 0
      \  /     \  /
       \/       \/
     [L11]  H(L5, H(x))    <- level 1
        \       /
          \   /
    H(L11, H(L5, H(x)))    <- level 2 (root)
    ```
    """
    item = merkle_hash(item)
    # TODO handle item passed not being in list more elegantly
    idx = tree[0].index(item)

    path = []
    for level in tree[:-1]:
        if (idx % 2) == 0:
            path.append(bit_set(level[idx+1], 0xFF))
        else:
            path.append(level[idx-1])
        idx = idx // 2
    return path


def merkle_proof(leaf, path, root):
    """
    Verify Merkle path for an item matches the root

    The most significant bit of every item in the path is used to
    determine if it's a 'left' or 'right' node, meaning:

        H(node, item) or H(item, node)
    """
    node = merkle_hash(leaf)
    for item in path:
        if bit_test(item, 0xFF):
            node = merkle_hash(node, bit_clear(item, 0xFF))
        else:
            node = merkle_hash(item, node)
    return node == root


def main():
    # Create 99 trees of 1..N items
    for i in range(1, 100):
        items = range(0, i)
        tree, root = merkle_tree(items)

        # Verify all items exist within the root
        random.shuffle(items)
        for item in items:
            proof = merkle_path(item, tree)
            assert merkle_proof(item, proof, root) is True


if __name__ == "__main__":
    main()
