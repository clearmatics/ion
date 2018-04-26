import click

from .repl import repl
from ion.lithium.etheventrelay import etheventrelay


commands = click.Group('commands')
commands.add_command(etheventrelay, "etheventrelay")


@commands.command(help="REPL for Ion commands and scripting")
def shell():
    repl(click.get_current_context())


if __name__ == "__main__":
    commands.main()
