import click

from Ion import commands as ion_commands
from ion.lithium.lithium import etheventrelay as lithium

commands = click.Group('commands')
commands.add_command(ion_commands, "ion")
commands.add_command(lithium, "lithium")

if __name__ == "__main__":
    commands.main()
