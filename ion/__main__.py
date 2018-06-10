import click

from .Ion import commands as ion_commands
from .lithium.lithium import etheventrelay as lithium
from .htlc.cli import commands as htlc_commands

commands = click.Group('commands')
commands.add_command(ion_commands, "ion")
commands.add_command(lithium, "lithium")
commands.add_command(htlc_commands, "htlc")

if __name__ == "__main__":
    commands.main()
