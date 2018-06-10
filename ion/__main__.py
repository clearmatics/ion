import click

from .Ion import commands as ion_commands
from .lithium.lithium import etheventrelay as lithium
from .htlc.cli import COMMANDS as HTLC_COMMANDS

commands = click.Group('commands')
commands.add_command(ion_commands, "ion")
commands.add_command(lithium, "lithium")
commands.add_command(HTLC_COMMANDS, "htlc")

if __name__ == "__main__":
    commands.main()
