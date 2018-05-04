import click

from Ion import commands as ion_commands
from ion.lithium.lithium import etheventrelay as lithium
from .repl import repl

commands = click.Group('commands')
commands.add_command(ion_commands, "ion")
commands.add_command(lithium, "lithium")


@commands.command(help="REPL for Ion commands and scripting")
def shell():
    repl(click.get_current_context())


if __name__ == "__main__":
    commands.main()
