import click

from .repl import repl
from .rpc.__main__ import commands as rpc_commands
from .plasma.__main__ import commands as plasma_commands
from .onchain import commands as onchain_commands


commands = click.Group('commands')
commands.add_command(rpc_commands, "rpc")
commands.add_command(plasma_commands, "plasma")
commands.add_command(onchain_commands, "onchain")


@commands.command(help="REPL for Ion commands and scripting")
def shell():
    repl(click.get_current_context())


if __name__ == "__main__":
    commands()
