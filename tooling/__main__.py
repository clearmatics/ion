import click

from .encoder import commands as encoder_commands
from .trie import commands as trie_commands

commands = click.Group('commands')
commands.add_command(encoder_commands, 'encoder')
commands.add_command(trie_commands, 'trie')

if __name__ == "__main__":
    commands.main()