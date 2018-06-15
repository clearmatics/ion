import click

from .Ion import commands as ion_commands
from .lithium.lithium import etheventrelay as lithium
from .htlc.cli import COMMANDS as HTLC_COMMANDS

COMMANDS = click.Group('commands')
COMMANDS.add_command(ion_commands, "ion")
COMMANDS.add_command(lithium, "lithium")
COMMANDS.add_command(HTLC_COMMANDS, "htlc")

if __name__ == "__main__":
    COMMANDS.main()
